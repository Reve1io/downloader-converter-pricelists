package utils

import (
	"strings"

	"github.com/xuri/excelize/v2"
)

func GetCellByIndex(f *excelize.File, sheet string, row int, col int) string {
	cell, _ := excelize.CoordinatesToCellName(col+1, row)
	val, _ := f.GetCellValue(sheet, cell)
	return strings.TrimSpace(val)
}
