package dto

import "stock-exchange-simulator/pkg/models"

const EnqueueOrderPlacedTopic = "asynq.orders.placed"

type EnqueueOrderPlacedDTO struct {
	OrderID   string           `json:"order_id"`
	StockID   string           `json:"stock_id"`
	Price     float64          `json:"price"`
	OrderType models.OrderType `json:"order_type"`
	Quantity  int              `json:"quantity"`
}
