package parser

import (
	"time"

	"downloader-converter-pricelists/internal/model"
	"downloader-converter-pricelists/internal/utils"

	"github.com/xuri/excelize/v2"
)

func ParseXLSX1(path string) ([]model.DBFItem, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sheetName := "Прайс-лист-" + time.Now().Format("02-01.2006")

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	var items []model.DBFItem

	for i, row := range rows {
		if i == 0 {
			continue
		}

		item := model.DBFItem{
			Code:     utils.Cell(row, 0),
			Name:     utils.Cell(row, 1),
			Producer: utils.Cell(row, 2),
			Qty:      utils.ParseInt(utils.Cell(row, 3)),
			Supplier: "ruelectronics",
		}

		items = append(items, item)
	}

	return items, nil
}
