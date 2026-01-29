package parser

import (
	"strings"

	"downloader-converter-pricelists/internal/model"

	"github.com/xuri/excelize/v2"
)

func ParseXLSX3(path string) ([]model.DBFItem, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := f.GetRows("ITECS_price_stock")
	if err != nil {
		return nil, err
	}

	var items []model.DBFItem

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

		items = append(items, model.DBFItem{
			Name:     name,
			Supplier: "Dip8",
		})
	}

	return items, nil
}
