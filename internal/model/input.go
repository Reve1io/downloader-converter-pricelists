package model

type DBFItem struct {
	Code      int
	Name      string
	Producer  string
	QntPack   int
	MOQ       int
	Qty       int
	Prices    []PriceBreak
	ClassName string
	History   string
	Weight    float64
}

type XLSXOffer struct {
	Name     string
	Price    float64
	Currency string
	ImageURL string
	Stock    int
}

type PriceBreak struct {
	Quant int     `json:"quant"`
	Price float64 `json:"price"`
}
