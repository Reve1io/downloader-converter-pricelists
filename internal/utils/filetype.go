package utils

import (
	"os"
)

func IsZipFile(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	buf := make([]byte, 4)
	if _, err := f.Read(buf); err != nil {
		return false
	}

	// ZIP signature: PK\x03\x04
	return buf[0] == 0x50 && buf[1] == 0x4B
}
