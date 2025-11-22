package db

import "github.com/jackc/pgx/v5/pgxpool"

type RepositoryFactory struct {
	OrderRepo IOrderRepository
	StockRepo IStockRepository
	TradeRepo ITradeRepository
}

func NewRepositoryFactory(db *pgxpool.Pool) *RepositoryFactory {
	return &RepositoryFactory{
		OrderRepo: NewOrderRepository(db),
		StockRepo: NewStockRepository(db),
		TradeRepo: NewTradeRepository(db),
	}
}
