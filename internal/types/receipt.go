package types

import "encoding/json"

type Receipt struct {
	IsReceipt               bool        `json:"isReceipt"`
	Date                    string      `json:"date,omitempty"`
	TotalAmount             json.Number `json:"totalAmount,omitempty"`
	TotalAmountIncludingTax json.Number `json:"totalAmountIncludingTax,omitempty"`
	CurrencySymbol          string      `json:"currencySymbol,omitempty"`
	Items                   []Item      `json:"items"`
}

type Item struct {
	Name   string      `json:"name,omitempty"`
	Amount json.Number `json:"Amount,omitempty"`
	Count  json.Number `json:"count,omitempty"`
}
