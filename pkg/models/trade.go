package models

import "time"

// Trade represents a completed trade.
type Trade struct {
	ID          string    `json:"id"`
	StockID     string    `json:"stock_id"`
	BuyOrderID  string    `json:"buy_order_id"`
	SellOrderID string    `json:"sell_order_id"`
	Quantity    int       `json:"quantity"`
	Price       float64   `json:"price"`
	Timestamp   time.Time `json:"timestamp"`
}
