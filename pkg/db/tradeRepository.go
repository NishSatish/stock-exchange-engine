package db

import (
	"context"
	"fmt"
	"stock-exchange-simulator/pkg/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ITradeRepository interface {
	CreateTrade(ctx context.Context, trade models.Trade) (models.Trade, error)
	GetTradeByID(ctx context.Context, id string) (models.Trade, error)
}

type TradeRepository struct {
	db *pgxpool.Pool
}

func NewTradeRepository(db *pgxpool.Pool) ITradeRepository {
	return &TradeRepository{
		db: db,
	}
}

func (r *TradeRepository) CreateTrade(ctx context.Context, trade models.Trade) (models.Trade, error) {
	err := r.db.QueryRow(ctx, "INSERT INTO trades (isin, order_id, quantity, price, status) VALUES ($1, $2, $3, $4, $5) RETURNING id, timestamp", trade.ISIN, trade.OrderID, trade.Quantity, trade.Price, "open").Scan(&trade.ID, &trade.Timestamp)
	if err != nil {
		return models.Trade{}, err
	}
	return trade, nil
}

func (r *TradeRepository) GetTradeByID(ctx context.Context, id string) (models.Trade, error) {
	var trade models.Trade
	err := r.db.QueryRow(ctx, "SELECT id, quantity, price, timestamp FROM trades WHERE id = $1", id).Scan(&trade.ID, &trade.Quantity, &trade.Price, &trade.Timestamp)
	if err != nil {
		return models.Trade{}, err
	}
	return trade, nil
}

func (r *TradeRepository) GetTradesByOrderID(ctx context.Context, order_id string) ([]models.Trade, error) {
	var trades []models.Trade

	rows, err := r.db.Query(ctx, "SELECT id, isin, order_id, quantity, price, timestamp, status FROM trades WHERE order_id = $1", order_id)
	if err != nil {
		fmt.Print("error getting tradesn by order_id")
		return nil, err
	}

	for rows.Next() {
		var t models.Trade
		errLocal := rows.Scan(
			&t.ID,
			&t.ISIN,
			&t.OrderID,
			&t.Quantity,
			&t.Price,
			&t.Timestamp,
			&t.Status,
		)
		if err != nil {
			fmt.Print("Error scanning a row", errLocal)
			continue
		}

		trades = append(trades, t)
	}

	return trades, err
}
