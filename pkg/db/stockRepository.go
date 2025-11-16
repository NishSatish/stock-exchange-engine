package db

import (
	"context"
	"errors"
	"stock-exchange-simulator/pkg/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type IStockRepository interface {
	CreateStock(ctx context.Context, stock models.Stock) bool
	GetStockByTicker(ctx context.Context, ticker string) (models.Stock, error)
	GetStockByIsin(ctx context.Context, isin string) (models.Stock, error)
}

type StockRepository struct {
	db *pgxpool.Pool
}

var stocks []models.Stock

func NewStockRepository(db *pgxpool.Pool) IStockRepository {
	return &StockRepository{
		db,
	}
}

func (r *StockRepository) CreateStock(ctx context.Context, stock models.Stock) bool {
	stocks = append(stocks, stock)
	return true
}

func (r *StockRepository) GetStockByTicker(ctx context.Context, ticker string) (models.Stock, error) {
	for _, stock := range stocks {
		if stock.Ticker == ticker {
			return stock, nil
		}
	}
	return models.Stock{}, errors.New("stock no found moinseur")
}

func (r *StockRepository) GetStockByIsin(ctx context.Context, isin string) (models.Stock, error) {
	for _, stock := range stocks {
		if stock.Ticker == isin {
			return stock, nil
		}
	}
	return models.Stock{}, errors.New("stock no found moinseur")
}
