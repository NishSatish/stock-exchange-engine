package models

// Stock represents a stock in the stock market.
type Stock struct {
	ID     string  `json:"id"`
	Ticker string  `json:"ticker"`
	Name   string  `json:"name"`
	LTP    float64 `json:"ltp"`
	Isin   string  `json:"isin"`
}
