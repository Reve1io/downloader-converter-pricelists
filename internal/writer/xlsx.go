package writer

import (
	"strconv"

	"downloader-converter-pricelists/internal/model"

	"github.com/xuri/excelize/v2"
)

func WriteXLSXStream(path string, in <-chan model.DBFItem) error {
	f := excelize.NewFile()
	sheet := "Sheet1"

	sw, err := f.NewStreamWriter(sheet)
	if err != nil {
		return err
	}

	headers := []interface{}{
		"Code",
		"Name",
		"Producer",
		"Class",
		"Quantity",
		"MOQ",
		"Pack",
		"Weight",
		"Prices",
		"Supplier",
	}

	cell, _ := excelize.CoordinatesToCellName(1, 1)
	if err := sw.SetRow(cell, headers); err != nil {
		return err
	}

	rowID := 2

	for it := range in {
		var prices string
		for _, p := range it.Prices {
			prices += strconv.Itoa(p.Quant) + ":" +
				strconv.FormatFloat(p.Price, 'f', 2, 64) + " "
		}

		row := []interface{}{
			it.Code,
			it.Name,
			it.Producer,
			it.ClassName,
			it.Qty,
			it.MOQ,
			it.QntPack,
			it.Weight,
			prices,
			it.Supplier,
		}

		cell, _ := excelize.CoordinatesToCellName(1, rowID)
		if err := sw.SetRow(cell, row); err != nil {
			return err
		}

		rowID++
	}

	if err := sw.Flush(); err != nil {
		return err
	}

	return f.SaveAs(path)
}
