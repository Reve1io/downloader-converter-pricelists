package utils

import (
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"
)

func DecodeDBF(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	// ✅ уже UTF-8 → ничего не делаем
	if utf8.ValidString(s) {
		return s
	}

	// пробуем CP866
	if out, err := charmap.CodePage866.NewDecoder().String(s); err == nil {
		if utf8.ValidString(out) {
			return strings.TrimSpace(out)
		}
	}

	// пробуем Windows-1251
	if out, err := charmap.Windows1251.NewDecoder().String(s); err == nil {
		if utf8.ValidString(out) {
			return strings.TrimSpace(out)
		}
	}

	// fallback
	return s
}
