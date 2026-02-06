package writer

import (
	"encoding/json"
	"os"

	"downloader-converter-pricelists/internal/model"
)

func WriteNDJSONStream(path string, ch <-chan model.DBFItem) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)

	for it := range ch {
		if err := enc.Encode(it); err != nil {
			return err
		}
	}
	return nil
}
