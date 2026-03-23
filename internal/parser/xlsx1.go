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
		return err
	}
	defer f.Close()

	sheetName := "Прайс-лист-" + time.Now().Format("02-01.2006")

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return err
	}

	for i, row := range rows {
		if i < 5 {
			continue
		}

		item := model.DBFItem{
			Code:      utils.Cell(row, 0),
			Name:      utils.Cell(row, 5),
			Producer:  utils.Cell(row, 3),
			Qty:       utils.AtoiSafe(utils.Cell(row, 9)),
			ClassName: utils.Cell(row, 6),
			Supplier:  "ruelectronics",
		}

		item.Prices = append(item.Prices, model.PriceBreak{
			Quant: utils.AtoiSafe(utils.Cell(row, 9)),
			Price: utils.ParseFloatSafe(utils.Cell(row, 10)),
		})

		out <- item
	}
	return nil
}
