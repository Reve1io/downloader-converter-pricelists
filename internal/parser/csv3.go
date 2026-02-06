package parser

import (
	"encoding/csv"
	"io"
	"os"

	"downloader-converter-pricelists/internal/model"
	"downloader-converter-pricelists/internal/utils"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func ParseCSVStock(path string, out chan<- model.DBFItem) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := transform.NewReader(f, charmap.Windows1251.NewDecoder())

	r := csv.NewReader(decoder)
	r.Comma = ';'

	_, _ = r.Read()

	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		item := model.DBFItem{
			Code:      utils.Cell(row, 0),
			Name:      utils.Cell(row, 1),
			Producer:  utils.Cell(row, 2),
			MOQ:       utils.AtoiSafe(utils.Cell(row, 5)),
			Qty:       utils.AtoiSafe(utils.Cell(row, 21)),
			ClassName: utils.Cell(row, 3),
			Supplier:  "Dip8",
		}

		for i := 7; i <= 20; i += 2 {
			qty := utils.AtoiSafe(utils.Cell(row, i))
			price := utils.ParseFloatSafe(utils.Cell(row, i+1))

			if qty > 0 && price > 0 {
				item.Prices = append(item.Prices, model.PriceBreak{
					Quant: qty,
					Price: price,
				})
			}
		}

		out <- item
	}

	return nil
}
