package model

type OutputItem struct {
	ItemID       int          `json:"item_id"`
	Name         string       `json:"name"`
	ProducerName string       `json:"producer_name"`
	Quant        int          `json:"quant"`
	Description  string       `json:"description"`
	PackQuant    int          `json:"pack_quant"`
	ClassName    string       `json:"class_name"`
	Munit        string       `json:"munit"`
	Weight       float64      `json:"weight"`
	Moq          int          `json:"moq"`
	ImageURL     string       `json:"image_url,omitempty"`
	PriceBreaks  []PriceBreak `json:"pricebreaks"`
	Supplier     string       `json:"supplier"`
}
