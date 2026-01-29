package utils

import "strconv"

func ParseInt(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}
