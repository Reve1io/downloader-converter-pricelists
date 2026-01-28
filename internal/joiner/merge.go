package joiner

import (
	"downloader-converter-pricelists/internal/model"
	"strings"
)

func Merge(all ...[]model.DBFItem) []model.DBFItem {
	index := make(map[string]*model.DBFItem)

	for _, src := range all {
		for _, item := range src {
			key := strings.ToUpper(item.Name)

			if existing, ok := index[key]; ok {
				existing.Prices = append(existing.Prices, item.Prices...)
				if existing.ImageURL == "" {
					existing.ImageURL = item.ImageURL
				}
			} else {
				cp := item
				index[key] = &cp
			}
		}
	}

	var result []model.DBFItem
	for _, v := range index {
		if len(v.Prices) == 0 {
			v.Prices = []model.PriceBreak{{0, 0}}
		}
		result = append(result, *v)
	}

	return result
}
