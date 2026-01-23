package writer

import (
	"downloader-converter-pricelists/internal/model"
	"strconv"

	"github.com/xuri/excelize/v2"
)

func WriteXLSX(path string, items []model.OutputItem) error {
	f := excelize.NewFile()
	sw, _ := f.NewStreamWriter("Sheet1")

	row := []interface{}{"ID", "Name", "Producer", "Weight"}
	sw.SetRow("A1", row)

	for i, it := range items {
		sw.SetRow(
			"A"+strconv.Itoa(i+2),
			[]interface{}{it.ItemID, it.Name, it.ProducerName, it.Weight},
		)
	}

	sw.Flush()
	return f.SaveAs(path)
}
