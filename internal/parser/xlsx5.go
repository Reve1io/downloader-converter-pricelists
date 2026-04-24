package parser

import (
	"downloader-converter-pricelists/internal/model"
	"downloader-converter-pricelists/internal/utils"
	"sort"
	"strings"

	"github.com/xuri/excelize/v2"
)

func ParseXLSX5(path string, out chan<- model.DBFItem) error {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return err
	}
	defer f.Close()

	sheetName := "Лист_1"

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return err
	}

	priceMap := map[int]int{
		4:  1,
		5:  10,
		6:  25,
		7:  50,
		8:  100,
		9:  500,
		10: 1000,
	}

	for i, row := range rows {

		if i == 0 || len(row) < 4 {
			continue
		}

		totalQty := 0
		for colIdx := 10; colIdx < len(row); colIdx++ {
			val := strings.TrimSpace(row[colIdx])
			if val != "" {
				cleanQty := strings.ReplaceAll(val, " ", "")
				parsed := utils.ParseInt(cleanQty)
				if parsed > 0 {
					totalQty = parsed
					break
				}
			}
		}

		item := model.DBFItem{
			Code:     utils.Cell(row, 1),
			Name:     utils.Cell(row, 2),
			Producer: utils.Cell(row, 3),
			Qty:      totalQty,
			Currency: "RUB",
			Supplier: "voltbricks",
		}

		for colIdx, quant := range priceMap {
			rawPrice := utils.Cell(row, colIdx)
			if rawPrice == "" {
				continue
			}

			cleanPrice := strings.Map(func(r rune) rune {
				if (r >= '0' && r <= '9') || r == '.' || r == ',' {
					return r
				}
				return -1
			}, rawPrice)

			lastMark := strings.LastIndexAny(cleanPrice, ".,")

			var finalPriceStr string
			if lastMark == -1 {
				finalPriceStr = cleanPrice
			} else {
				integerPart := strings.NewReplacer(".", "", ",", "").Replace(cleanPrice[:lastMark])
				fractionalPart := cleanPrice[lastMark+1:]
				finalPriceStr = integerPart + "." + fractionalPart
			}

			price := utils.ParseFloatSafe(finalPriceStr)

			if price > 0 {
				item.Prices = append(item.Prices, model.PriceBreak{
					Price: price,
					Quant: quant,
				})
			}
		}

		sort.Slice(item.Prices, func(i, j int) bool {
			return item.Prices[i].Quant < item.Prices[j].Quant
		})

		out <- item
	}
	return nil
}
