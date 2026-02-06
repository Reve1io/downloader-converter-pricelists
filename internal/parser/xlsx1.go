package parser

import (
	"time"

	"downloader-converter-pricelists/internal/model"
	"downloader-converter-pricelists/internal/utils"

	"github.com/xuri/excelize/v2"
)

func ParseXLSX1(path string, out chan<- model.DBFItem) error {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	sheetName := "Прайс-лист-" + time.Now().Format("02-01.2006")

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil
	}

	for i, row := range rows {
		if i == 0 {
			continue
		}

		qty := utils.AtoiSafe(utils.Cell(row, 11))
		price := utils.ParseFloatSafe(utils.Cell(row, 12))

		item := model.DBFItem{
			Code:     utils.Cell(row, 0),
			Name:     utils.Cell(row, 5),
			Producer: utils.Cell(row, 3),
			Qty:      utils.ParseInt(utils.Cell(row, 3)),
			Supplier: "ruelectronics",
		}

		item.Prices = append(item.Prices, model.PriceBreak{
			Quant: qty,
			Price: price,
		})

		out <- item
	}
	return nil
}
