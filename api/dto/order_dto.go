package dto

import "stock-exchange-simulator/pkg/models"

// PlaceOrderRequest defines the request body for placing a new order.
type PlaceOrderRequest struct {
	StockID  string           `json:"stock_id" binding:"required"`
	Type     models.OrderType `json:"type" binding:"required"`
	Quantity int              `json:"quantity" binding:"required,min=1"`
	Price    float64          `json:"price" binding:"required,min=0"`
}
