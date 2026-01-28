package model

type PriceBreak struct {
	Quant int     `json:"quant"`
	Price float64 `json:"price"`
}

type DBFItem struct {
	Code      string
	Name      string
	Producer  string
	QntPack   int
	MOQ       int
	Qty       int
	Prices    []PriceBreak
	ClassName string
	History   string
	ImageURL  string
	Weight    float64
	Currency  string
}
