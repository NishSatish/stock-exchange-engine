package db

import (
	"context"
	"errors"
	"stock-exchange-simulator/pkg/models"
)

type ITradeRepository interface {
	CreateTrade(ctx context.Context, trade models.Trade) (models.Trade, error)
	GetTradeByID(ctx context.Context, id string) (models.Trade, error)
}

type TradeRepository struct {
}

var trades []models.Trade

func NewTradeRepository() ITradeRepository {
	return &TradeRepository{}
}

func (r *TradeRepository) CreateTrade(ctx context.Context, trade models.Trade) (models.Trade, error) {
	trades = append(trades, trade)
	return trade, nil
}

func (r *TradeRepository) GetTradeByID(ctx context.Context, id string) (models.Trade, error) {
	for _, trade := range trades {
		if trade.ID == id {
			return trade, nil
		}
	}
	return models.Trade{}, errors.New("trade not found")
}
