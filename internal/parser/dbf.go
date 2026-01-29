package parser

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"downloader-converter-pricelists/internal/model"

	"github.com/LindsayBradford/go-dbf/godbf"

	enc "downloader-converter-pricelists/internal/utils"
)

func ParseDBF(path string) ([]model.DBFItem, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	header := make([]byte, 1)
	if _, err := f.Read(header); err != nil {
		return nil, err
	}

	log.Printf("DBF signature: 0x%X\n", header[0])

	switch header[0] {
	case 0x03, 0x83, 0x8B:
		// OK
	default:
		return nil, fmt.Errorf(
			"file %s is not DBF (signature 0x%X)",
			path,
			header[0],
		)
	}

	dbf, err := godbf.NewFromFile(path, "CP866")
	if err != nil {
		return nil, err
	}

	idx := map[string]int{
		"CODE":       fieldIndex(dbf, "CODE"),
		"NAME":       fieldIndex(dbf, "NAME"),
		"PRODUCER":   fieldIndex(dbf, "PRODUCER"),
		"QNT_PACK":   fieldIndex(dbf, "QNT_PACK"),
		"MOQ":        fieldIndex(dbf, "MOQ"),
		"QTY":        fieldIndex(dbf, "QTY"),
		"WEIGHT":     fieldIndex(dbf, "WEIGHT"),
		"CLASS_NAME": fieldIndex(dbf, "CLASS_NAME"),
		"HISTORY":    fieldIndex(dbf, "HISTORY"),
		"SUPPLIER":   fieldIndex(dbf, "SUPPLIER"),
	}

	if idx["NAME"] == -1 {
		return nil, fmt.Errorf("required field NAME not found in DBF")
	}

	var items []model.DBFItem

	for i := 0; i < dbf.NumberOfRecords(); i++ {
		name := strings.TrimSpace(enc.DecodeDBF(dbf.FieldValue(i, idx["NAME"])))
		if name == "" {
			continue
		}

		item := model.DBFItem{
			Name:      name,
			Producer:  enc.DecodeDBF(dbf.FieldValue(i, idx["PRODUCER"])),
			ClassName: enc.DecodeDBF(dbf.FieldValue(i, idx["CLASS_NAME"])),
			History:   enc.DecodeDBF(dbf.FieldValue(i, idx["HISTORY"])),
			Supplier:  "compel",
		}

		item.Code = strings.TrimSpace(dbf.FieldValue(i, idx["CODE"]))
		item.QntPack, _ = strconv.Atoi(dbf.FieldValue(i, idx["QNT_PACK"]))
		item.MOQ, _ = strconv.Atoi(dbf.FieldValue(i, idx["MOQ"]))
		item.Qty, _ = strconv.Atoi(dbf.FieldValue(i, idx["QTY"]))
		item.Weight, _ = strconv.ParseFloat(
			strings.ReplaceAll(dbf.FieldValue(i, idx["WEIGHT"]), ",", "."),
			64,
		)

		for n := 1; n <= 8; n++ {
			qtyIdx := fieldIndex(dbf, "QTY_"+strconv.Itoa(n))
			priceIdx := fieldIndex(dbf, "PRICE_"+strconv.Itoa(n))

			if qtyIdx == -1 || priceIdx == -1 {
				continue
			}

			qty, _ := strconv.Atoi(dbf.FieldValue(i, qtyIdx))
			price, _ := strconv.ParseFloat(
				strings.ReplaceAll(dbf.FieldValue(i, priceIdx), ",", "."),
				64,
			)

			if qty > 0 {
				item.Prices = append(item.Prices, model.PriceBreak{
					Quant: qty,
					Price: price,
				})
			}
		}

		if len(item.Prices) == 0 {
			item.Prices = []model.PriceBreak{{Quant: 0, Price: 0}}
		}

		items = append(items, item)
	}

	return items, nil
}
