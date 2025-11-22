package models

import "time"

// Trade represents a completed trade.
type Trade struct {
	ID        string    `json:"id"`
	ISIN      string    `json:"stock_id"`
	OrderID   string    `json:"buy_order_id"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
}
