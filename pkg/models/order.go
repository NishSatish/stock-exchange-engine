package models

import "time"

// OrderType defines the type of an order (buy or sell).
type OrderType string
type OrderStatus string

const (
	Pending   OrderStatus = "PENDING"
	Completed OrderStatus = "COMPLETED"
	Cancelled OrderStatus = "CANCELLED"
)

const (
	Buy  OrderType = "BUY"
	Sell OrderType = "SELL"
)

// Order represents a buy or sell order.
type Order struct {
	ID        string      `json:"id"`
	StockID   string      `json:"stock_id"`
	Type      OrderType   `json:"type"`
	Quantity  int         `json:"quantity"`
	Price     float64     `json:"price"`
	Timestamp time.Time   `json:"timestamp"`
	Status    OrderStatus `json:"status"`
	TradeID   string      `json:"trade_id,omitempty"`
}
