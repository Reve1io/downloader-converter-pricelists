package parser

import (
	"strconv"
	"strings"
)

func parseFloat(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}

	// DBF / XLSX часто используют запятую как десятичный разделитель
	s = strings.ReplaceAll(s, ",", ".")

	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return v
}
