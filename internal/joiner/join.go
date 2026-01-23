package joiner

import (
	"downloader-converter-pricelists/internal/currency"
	"downloader-converter-pricelists/internal/model"
)

func Join(
	items []model.DBFItem,
	offers map[string]model.XLSXOffer,
) []model.OutputItem {

	var result []model.OutputItem

	for _, it := range items {
		out := model.OutputItem{
			ItemID:       it.Code,
			Name:         it.Name,
			ProducerName: it.Producer,
			Quant:        it.Qty,
			Description:  it.History,
			PackQuant:    it.QntPack,
			ClassName:    it.ClassName,
			Munit:        "шт",
			Weight:       it.Weight,
			Moq:          it.MOQ,
			PriceBreaks:  it.Prices,
		}

		if offer, ok := offers[it.Name]; ok {
			out.ImageURL = offer.ImageURL
			for i := range out.PriceBreaks {
				out.PriceBreaks[i].Price =
					currency.ToUSD(out.PriceBreaks[i].Price, offer.Currency)
			}
		}

		result = append(result, out)
	}

	return result
}
