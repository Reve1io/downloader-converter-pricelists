package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"downloader-converter-pricelists/internal/config"
	"downloader-converter-pricelists/internal/model"
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

	os.MkdirAll(cfg.InputDir, 0755)
	utils.CleanDir(cfg.InputDir)

	httpDl := &source.HTTPDownloader{
		Timeout: time.Second * time.Duration(cfg.HTTP.TimeoutSeconds),
		Retries: cfg.HTTP.Retries,
	}

	for _, src := range cfg.HTTP.Sources {
		dest := filepath.Join(cfg.InputDir, src.Filename)
		log.Println("Downloading:", src.URL)

		if err := httpDl.Download(src.URL, dest, src.Type); err != nil {
			log.Fatal(err)
		}
	}

	ftpDl := &source.FTPDownloader{
		Addr:     cfg.FTP.Addr,
		User:     cfg.FTP.User,
		Password: cfg.FTP.Password,
		Timeout:  time.Second * time.Duration(cfg.FTP.TimeoutSeconds),
	}

	for _, src := range cfg.FTP.Sources {
		dest := filepath.Join(cfg.InputDir, src.Filename)
		log.Println("Downloading FTP:", src.Remote)

		if err := ftpDl.Download(src.Remote, dest); err != nil {
			log.Fatal(err)
		}
	}

	os.MkdirAll("output", 0755)
	ts := time.Now().Format("2006-01-02_15-04-05")
	jsonPath := "output/out_" + ts + ".ndjson"

	items := make(chan model.DBFItem, 2000)

	go func() {
		defer close(items)

		log.Println("Parsing DBF...")
		if err := parser.ParseDBF("input/COMPELDISTI2.dbf", items); err != nil {
			log.Fatal(err)
		}

		log.Println("Parsing XLSX1...")
		if err := parser.ParseXLSX1("input/x1.xlsx", items); err != nil {
			log.Fatal(err)
		}

		log.Println("Parsing CSV Available...")
		if err := parser.ParseCSVAvailable("input/ITECS_price_available.csv", items); err != nil {
			log.Fatal(err)
		}

		log.Println("Parsing CSV Stock...")
		if err := parser.ParseCSVStock("input/ITECS_price_stock.csv", items); err != nil {
			log.Fatal(err)
		}

		log.Println("Parsing XLSX4...")
		if err := parser.ParseXLSX4("input/x4.xlsx", items); err != nil {
			log.Fatal(err)
		}

		log.Println("All parsers finished")
	}()

	log.Println("Writing NDJSON stream...")
	if err := writer.WriteNDJSONStream(jsonPath, items); err != nil {
		log.Fatal(err)
	}

	log.Println("✔ Done:", jsonPath)
}
