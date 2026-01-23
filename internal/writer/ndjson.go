package writer

import (
	"encoding/json"
	"os"

	"downloader-converter-pricelists/internal/model"
)

func WriteNDJSON(path string, items []model.OutputItem) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	for _, it := range items {
		if err := enc.Encode(it); err != nil {
			return err
		}
	}
	return nil
}
