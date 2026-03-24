package parser

import (
	"strings"

	"downloader-converter-pricelists/internal/model"
	"downloader-converter-pricelists/internal/utils"

	"github.com/xuri/excelize/v2"
)

func ParseXLSX4(path string, out chan<- model.DBFItem) error {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return err
	}
	defer f.Close()

	rows, err := f.GetRows("Worksheet")
	if err != nil {
		return err
	}

	agg := make(map[string]*model.DBFItem)

	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 8 {
			continue
		}

		name := strings.TrimSpace(utils.Cell(row, 3))
		if name == "" {
			continue
		}

		producer := utils.Cell(row, 5)
		code := utils.Cell(row, 4)
		qty := utils.ParseInt(row[7])

		key := producer + "|" + name

		if _, ok := agg[key]; !ok {

			item := &model.DBFItem{
				Code:      code,
				Name:      name,
				Producer:  producer,
				ClassName: name,
				Supplier:  "radioelementy",
			}

			for j := 11; j <= 15; j += 2 {
				price := utils.ParseFloatSafe(utils.Cell(row, j))
				q := utils.AtoiSafe(utils.Cell(row, j+1))

				if q > 0 && price > 0 {
					item.Prices = append(item.Prices, model.PriceBreak{
						Quant: q,
						Price: price,
					})
				}
			}

			agg[key] = item
		}

		agg[key].Qty += qty
	}

	for _, item := range agg {
		out <- *item
	}

	return nil
}
