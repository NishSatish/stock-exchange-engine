package db

import (
	"context"
	"errors"
	"stock-exchange-simulator/pkg/models"
)

type IOrderRepository interface {
	CreateOrder(ctx context.Context, order models.Order) (models.Order, error)
	GetOrderByID(ctx context.Context, id string) (models.Order, error)
}

type OrderRepository struct {
}

var orders []models.Order

func NewOrderRepository() IOrderRepository {
	return &OrderRepository{}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, order models.Order) (models.Order, error) {
	orders = append(orders, order)
	return order, nil
}

func (r *OrderRepository) GetOrderByID(ctx context.Context, id string) (models.Order, error) {
	for _, order := range orders {
		if order.ID == id {
			return order, nil
		}
	}
	return models.Order{}, errors.New("order not found")
}
