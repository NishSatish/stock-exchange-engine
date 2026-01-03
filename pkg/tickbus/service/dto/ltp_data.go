package dto

import "time"

// LTPData represents the Last Traded Price data for a stock.
type LTPData struct {
	StockID string    `json:"stock_id"`
	LTP         float64   `json:"ltp"`
	Change      float64   `json:"change"` // Change fromg last recorded LTP
	LastUpdated time.Time `json:"last_updated"`
}
