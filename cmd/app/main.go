package main

import (
	"log"
	"time"

	"downloader-converter-pricelists/internal/joiner"
	"downloader-converter-pricelists/internal/parser"
	"downloader-converter-pricelists/internal/writer"
)

func main() {
	dbfItems, err := parser.ParseDBF("input/input.dbf")
	if err != nil {
		log.Fatal(err)
	}

	offers, err := parser.ParseXLSX("input/offers.xlsx")
	if err != nil {
		log.Fatal(err)
	}

	result := joiner.Join(dbfItems, offers)

	ts := time.Now().Format("2006-01-02_15-04-05")
	path := "output/"
	jsonPath := path + "out_" + ts + ".ndjson"
	xlsxPath := path + "out_" + ts + ".xlsx"

	if err := writer.WriteNDJSON(jsonPath, result); err != nil {
		log.Fatal(err)
	}

	if err := writer.WriteXLSX(xlsxPath, result); err != nil {
		log.Fatal(err)
	}
}
