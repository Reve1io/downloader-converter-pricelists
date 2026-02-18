package database

import (
	"context"
	"log"
	"time"

	"downloader-converter-pricelists/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const batchSize = 1000
const progressStep = 100000

type PGWriter struct {
	pool *pgxpool.Pool

	products  map[string]int64
	producers map[string]int32
	suppliers map[string]int32
}

func NewPGWriter(pool *pgxpool.Pool) *PGWriter {
	return &PGWriter{
		pool:      pool,
		products:  make(map[string]int64),
		producers: make(map[string]int32),
		suppliers: make(map[string]int32),
	}
}

func (w *PGWriter) WriteStream(ctx context.Context, in <-chan model.DBFItem) (err error) {

	tx, err := w.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			log.Println("ROLLBACK")
			_ = tx.Rollback(ctx)
		}
	}()

	start := time.Now()
	var priceCounter int
	var batch pgx.Batch

	for item := range in {

		productID, err := w.upsertProduct(ctx, tx, item)
		if err != nil {
			return err
		}

		supplierID, err := w.getSupplierID(ctx, tx, item.Supplier)
		if err != nil {
			return err
		}

		for _, p := range item.Prices {

			batch.Queue(
				`INSERT INTO current_prices
				(product_id, supplier_id, quant, price, currency, qty_available, moq, qnt_pack, updated_at)
				VALUES ($1,$2,$3,$4,$5,$6,$7,$8,now())
				ON CONFLICT (product_id, supplier_id, quant)
				DO UPDATE SET
					price = EXCLUDED.price,
					qty_available = EXCLUDED.qty_available,
					moq = EXCLUDED.moq,
					qnt_pack = EXCLUDED.qnt_pack,
					updated_at = now()`,
				productID,
				supplierID,
				p.Quant,
				p.Price,
				item.Currency,
				item.Qty,
				item.MOQ,
				item.QntPack,
			)

			priceCounter++

			if priceCounter%batchSize == 0 {
				if err := w.flushBatch(ctx, tx, &batch); err != nil {
					return err
				}
				batch = pgx.Batch{}
			}

			if priceCounter%progressStep == 0 {
				log.Printf("PG progress: %d rows (%.2fs)",
					priceCounter,
					time.Since(start).Seconds(),
				)
			}
		}
	}

	if batch.Len() > 0 {
		if err := w.flushBatch(ctx, tx, &batch); err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	log.Printf("PG DONE: %d rows in %.2fs",
		priceCounter,
		time.Since(start).Seconds(),
	)

	return nil
}

func (w *PGWriter) flushBatch(ctx context.Context, tx pgx.Tx, batch *pgx.Batch) error {
	br := tx.SendBatch(ctx, batch)
	return br.Close()
}

func (w *PGWriter) upsertProduct(
	ctx context.Context,
	tx pgx.Tx,
	item model.DBFItem,
) (int64, error) {

	// ---- cache check ----
	if id, ok := w.products[item.Code]; ok {
		return id, nil
	}

	producerID, err := w.getProducerID(ctx, tx, item.Producer)
	if err != nil {
		return 0, err
	}

	var id int64

	err = tx.QueryRow(ctx,
		`INSERT INTO products
		(code,name,producer_id,category_id,history,image_url,weight,updated_at)
		VALUES ($1,$2,$3,NULL,$4,$5,$6,now())
		ON CONFLICT (code)
		DO UPDATE SET
			name = EXCLUDED.name,
			producer_id = EXCLUDED.producer_id,
			history = EXCLUDED.history,
			image_url = EXCLUDED.image_url,
			weight = EXCLUDED.weight,
			updated_at = now()
		RETURNING id`,
		item.Code,
		item.Name,
		producerID,
		item.History,
		item.ImageURL,
		item.Weight,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	w.products[item.Code] = id

	return id, nil
}

func (w *PGWriter) getSupplierID(
	ctx context.Context,
	tx pgx.Tx,
	name string,
) (int32, error) {

	if name == "" {
		return 0, nil
	}

	if id, ok := w.suppliers[name]; ok {
		return id, nil
	}

	var id int32

	err := tx.QueryRow(ctx,
		`INSERT INTO suppliers(name)
		 VALUES ($1)
		 ON CONFLICT (name)
		 DO UPDATE SET name = EXCLUDED.name
		 RETURNING id`,
		name,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	w.suppliers[name] = id

	return id, nil
}

func (w *PGWriter) getProducerID(
	ctx context.Context,
	tx pgx.Tx,
	name string,
) (*int32, error) {

	if name == "" {
		return nil, nil
	}

	if id, ok := w.producers[name]; ok {
		return &id, nil
	}

	var id int32

	err := tx.QueryRow(ctx,
		`INSERT INTO producers(name)
		 VALUES ($1)
		 ON CONFLICT (name)
		 DO UPDATE SET name = EXCLUDED.name
		 RETURNING id`,
		name,
	).Scan(&id)

	if err != nil {
		return nil, err
	}

	w.producers[name] = id

	return &id, nil
}
