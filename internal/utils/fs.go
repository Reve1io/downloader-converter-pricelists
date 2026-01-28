package utils

import (
	"os"
	"path/filepath"
)

func CleanDir(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if !e.IsDir() {
			os.Remove(filepath.Join(dir, e.Name()))
		}
	}
	return nil
}
