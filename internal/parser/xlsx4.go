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

	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 2 {
			continue
		}

		name := strings.TrimSpace(row[0])
		if name == "" {
			continue
		}

		item := model.DBFItem{
			Code:      utils.Cell(row, 4),
			Name:      utils.Cell(row, 3),
			Producer:  utils.Cell(row, 5),
			ClassName: utils.Cell(row, 3),
			Supplier:  "radioelementy",
		}

		for i := 11; i <= 15; i += 2 {
			qty := utils.AtoiSafe(utils.Cell(row, i+1))
			price := utils.ParseFloatSafe(utils.Cell(row, i))

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
