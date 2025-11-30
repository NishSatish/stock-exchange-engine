package db

import (
	"context"
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
	err := r.db.QueryRow(ctx, "INSERT INTO trades (stock_id, buy_order_id, sell_order_id, quantity, price) VALUES ($1, $2, $3, $4, $5) RETURNING id, timestamp", trade.StockID, trade.BuyOrderID, trade.SellOrderID, trade.Quantity, trade.Price).Scan(&trade.ID, &trade.Timestamp)
	if err != nil {
		return models.Trade{}, err
	}
	return trade, nil
}

func (r *TradeRepository) GetTradeByID(ctx context.Context, id string) (models.Trade, error) {
	var trade models.Trade
	err := r.db.QueryRow(ctx, "SELECT id, stock_id, buy_order_id, sell_order_id, quantity, price, timestamp FROM trades WHERE id = $1", id).Scan(&trade.ID, &trade.StockID, &trade.BuyOrderID, &trade.SellOrderID, &trade.Quantity, &trade.Price, &trade.Timestamp)
	if err != nil {
		return models.Trade{}, err
	}
	return trade, nil
}
