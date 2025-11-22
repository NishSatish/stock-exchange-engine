package db

import (
	"context"
	"stock-exchange-simulator/pkg/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type IStockRepository interface {
	CreateStock(ctx context.Context, stock models.Stock) (models.Stock, error)
	GetStockByTicker(ctx context.Context, ticker string) (models.Stock, error)
	GetStockByIsin(ctx context.Context, isin string) (models.Stock, error)
}

type StockRepository struct {
	db *pgxpool.Pool
}

func NewStockRepository(db *pgxpool.Pool) IStockRepository {
	return &StockRepository{
		db: db,
	}
}

func (r *StockRepository) CreateStock(ctx context.Context, stock models.Stock) (models.Stock, error) {
	err := r.db.QueryRow(ctx, "INSERT INTO stocks (ticker, name, ltp, isin) VALUES ($1, $2, $3, $4) RETURNING id", stock.Ticker, stock.Name, stock.LTP, stock.Isin).Scan(&stock.ID)
	if err != nil {
		return models.Stock{}, err
	}
	return stock, nil
}

func (r *StockRepository) GetStockByTicker(ctx context.Context, ticker string) (models.Stock, error) {
	var stock models.Stock
	err := r.db.QueryRow(ctx, "SELECT id, ticker, name, ltp, isin FROM stocks WHERE ticker = $1", ticker).Scan(&stock.ID, &stock.Ticker, &stock.Name, &stock.LTP, &stock.Isin)
	if err != nil {
		return models.Stock{}, err
	}
	return stock, nil
}

func (r *StockRepository) GetStockByIsin(ctx context.Context, isin string) (models.Stock, error) {
	var stock models.Stock
	err := r.db.QueryRow(ctx, "SELECT id, ticker, name, ltp, isin FROM stocks WHERE isin = $1", isin).Scan(&stock.ID, &stock.Ticker, &stock.Name, &stock.LTP, &stock.Isin)
	if err != nil {
		return models.Stock{}, err
	}
	return stock, nil
}
