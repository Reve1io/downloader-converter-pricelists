package utils

import (
	"strconv"
	"strings"
)

func AtoiSafe(s string) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	v, _ := strconv.Atoi(s)
	return v
}

func AtofSafe(s string) float64 {
	s = strings.TrimSpace(strings.ReplaceAll(s, ",", "."))
	if s == "" {
		return 0
	}
	v, _ := strconv.ParseFloat(s, 64)
	return v
}
