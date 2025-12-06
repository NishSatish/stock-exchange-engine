package db

import (
	"context"
	"stock-exchange-simulator/pkg/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type IOrderRepository interface {
	CreateOrder(ctx context.Context, order models.Order) (models.Order, error)
	GetOrderByID(ctx context.Context, id string) (models.Order, error)
	MarkOrdersCompleteInBulk(ctx context.Context, orderIDs []string, tradeId string) error
}

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) IOrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, order models.Order) (models.Order, error) {
	err := r.db.QueryRow(ctx, "INSERT INTO orders (stock_id, type, quantity, price, status) VALUES ($1, $2, $3, $4, $5) RETURNING id, timestamp", order.StockID, order.Type, order.Quantity, order.Price, order.Status).Scan(&order.ID, &order.Timestamp)
	if err != nil {
		return models.Order{}, err
	}
	return order, nil
}

func (r *OrderRepository) GetOrderByID(ctx context.Context, id string) (models.Order, error) {
	var order models.Order
	err := r.db.QueryRow(ctx, "SELECT id, stock_id, type, quantity, price, timestamp, status FROM orders WHERE id = $1", id).Scan(&order.ID, &order.StockID, &order.Type, &order.Quantity, &order.Price, &order.Timestamp, &order.Status)
	if err != nil {
		return models.Order{}, err
	}
	return order, nil
}

func (r *OrderRepository) MarkOrdersCompleteInBulk(ctx context.Context, orderIDs []string, tradeID string) error {
	_, err := r.db.Exec(ctx, "UPDATE orders SET status = $2, trade_id = $3 WHERE id = ANY($1)", orderIDs, models.Completed, tradeID)
	return err
}
