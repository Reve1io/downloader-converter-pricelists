package parser

import (
	"strings"

	"github.com/LindsayBradford/go-dbf/godbf"
)

func fieldIndex(dbf *godbf.DbfTable, name string) int {
	for i, field := range dbf.FieldNames() {
		if strings.EqualFold(field, name) {
			return i
		}
	}
	return -1
}
