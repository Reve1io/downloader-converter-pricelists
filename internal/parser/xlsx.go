package parser

import (
	"strings"

	"downloader-converter-pricelists/internal/model"

	"github.com/xuri/excelize/v2"
)

func ParseXLSX(path string) (map[string]model.XLSXOffer, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}

	rows, err := f.Rows("Sheet1")
	if err != nil {
		return nil, err
	}

	offers := make(map[string]model.XLSXOffer)
	rows.Next() // header

	for rows.Next() {
		row, _ := rows.Columns()
		name := strings.TrimSpace(row[0])

		if name == "" {
			continue
		}

		offers[name] = model.XLSXOffer{
			Name:     name,
			Price:    parseFloat(row[1]),
			Currency: row[2],
			ImageURL: row[3],
		}
	}

	return offers, nil
}
