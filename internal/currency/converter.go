package currency

func ToUSD(value float64, currency string) float64 {
	switch currency {
	case "USD":
		return value
	case "EUR":
		return value * 1.08
	case "RUB":
		return value / 90
	default:
		return value
	}
}
