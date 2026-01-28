package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"downloader-converter-pricelists/internal/config"
	"downloader-converter-pricelists/internal/joiner"
	"downloader-converter-pricelists/internal/parser"
	"downloader-converter-pricelists/internal/source"
	"downloader-converter-pricelists/internal/utils"
	"downloader-converter-pricelists/internal/writer"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// input/
	os.MkdirAll(cfg.InputDir, 0755)
	utils.CleanDir(cfg.InputDir)

	// downloader
	httpDl := &source.HTTPDownloader{
		Timeout: time.Second * time.Duration(cfg.HTTP.TimeoutSeconds),
		Retries: cfg.HTTP.Retries,
	}

	for _, src := range cfg.HTTP.Sources {
		dest := filepath.Join(cfg.InputDir, src.Filename)
		log.Println("Downloading:", src.URL)

		if err := httpDl.Download(src.URL, dest); err != nil {
			log.Fatal(err)
		}
	}

	// parsing → ВСЁ В DBFItem
	itemsDBF, err := parser.ParseDBF("input/COMPELDISTI2.dbf")
	if err != nil {
		log.Fatal(err)
	}

	x1, err := parser.ParseXLSX1("input/x1.xlsx")
	if err != nil {
		log.Fatal(err)
	}

	x2, err := parser.ParseXLSX2("input/x2.xlsx")
	if err != nil {
		log.Fatal(err)
	}

	x3, err := parser.ParseXLSX3("input/x3.xlsx")
	if err != nil {
		log.Fatal(err)
	}

	x4, err := parser.ParseXLSX4("input/x4.xlsx")
	if err != nil {
		log.Fatal(err)
	}

	merged := joiner.Merge(itemsDBF, x1, x2, x3, x4)

	os.MkdirAll("output", 0755)

	ts := time.Now().Format("2006-01-02_15-04-05")
	jsonPath := "output/out_" + ts + ".ndjson"
	xlsxPath := "output/out_" + ts + ".xlsx"

	if err := writer.WriteNDJSON(jsonPath, merged); err != nil {
		log.Fatal(err)
	}

	if err := writer.WriteXLSX(xlsxPath, merged); err != nil {
		log.Fatal(err)
	}

	log.Println("✔ Done:", jsonPath, xlsxPath)
}
