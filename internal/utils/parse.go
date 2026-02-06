package utils

import (
	"strconv"
	"strings"
)

func ParseInt(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}

func ParseFloatSafe(s string) float64 {
	s = strings.ReplaceAll(s, ",", ".")
	v, _ := strconv.ParseFloat(s, 64)
	return v
}
