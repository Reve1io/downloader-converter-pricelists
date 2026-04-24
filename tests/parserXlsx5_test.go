package tests

import (
	"downloader-converter-pricelists/internal/model"
	"downloader-converter-pricelists/internal/parser"
	"downloader-converter-pricelists/internal/writer"
	"testing"
	"time"
)

func TestParseXLSX5(t *testing.T) {

	items := make(chan model.DBFItem, 2000)
	localPath := "test_input/unload_all.xlsx"

	ts := time.Now().Format("2006-01-02_15-04-05")
	jsonPath := "output/out_" + ts + ".ndjson"

	t.Log("Parsing XLSX5 started")

	ndjson := make(chan error)

	go func() {
		ndjson <- writer.WriteNDJSONStream(jsonPath, items)
	}()

	go func() {
		defer close(items)
		if err := parser.ParseXLSX5(localPath, items); err != nil {
			t.Error(err)
		}
	}()

	if err := <-ndjson; err != nil {
		t.Errorf("Error during writing ndjson: %v", err)
	}

	t.Log("Parsing and writing ndjson successfully!")
}
