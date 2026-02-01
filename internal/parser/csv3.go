package parser

import (
	"encoding/csv"
	"os"

	"downloader-converter-pricelists/internal/model"
	"downloader-converter-pricelists/internal/utils"
)

func ParseCSVStock(path string) ([]model.DBFItem, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.Comma = ';'

	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	var items []model.DBFItem

	for i, row := range rows {
		if i == 0 {
			continue // header
		}

		item := model.DBFItem{
			Code:     row[0],
			Name:     row[1],
			Qty:      utils.AtoiSafe(row[2]),
			Supplier: "Dip8",
		}

		items = append(items, item)
	}

	return items, nil
}
