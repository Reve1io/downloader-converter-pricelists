package writer

import (
	"strconv"

	"downloader-converter-pricelists/internal/model"

	"github.com/xuri/excelize/v2"
)

func WriteXLSX(path string, items []model.DBFItem) error {
	f := excelize.NewFile()
	sheet := "Sheet1"
	f.SetSheetName("Sheet1", sheet)

	// Header
	headers := []string{
		"Code",
		"Name",
		"Producer",
		"Class",
		"Quantity",
		"MOQ",
		"Pack",
		"Weight",
		"Prices",
	}

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// Rows
	for r, it := range items {
		row := r + 2

		f.SetCellValue(sheet, "A"+strconv.Itoa(row), it.Code)
		f.SetCellValue(sheet, "B"+strconv.Itoa(row), it.Name)
		f.SetCellValue(sheet, "C"+strconv.Itoa(row), it.Producer)
		f.SetCellValue(sheet, "D"+strconv.Itoa(row), it.ClassName)
		f.SetCellValue(sheet, "E"+strconv.Itoa(row), it.Qty)
		f.SetCellValue(sheet, "F"+strconv.Itoa(row), it.MOQ)
		f.SetCellValue(sheet, "G"+strconv.Itoa(row), it.QntPack)
		f.SetCellValue(sheet, "H"+strconv.Itoa(row), it.Weight)

		// price breaks as text
		var prices string
		for _, p := range it.Prices {
			prices += strconv.Itoa(p.Quant) + ":" +
				strconv.FormatFloat(p.Price, 'f', 2, 64) + " "
		}
		f.SetCellValue(sheet, "I"+strconv.Itoa(row), prices)
	}

	return f.SaveAs(path)
}
