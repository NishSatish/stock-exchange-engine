package service

import (
	"stock-exchange-simulator/pkg/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ITradeMasterServiceInterface interface {
	ExecuteTrade(trade *models.Trade) error
}

type TradeMasterService struct {
	db *pgxpool.Pool
}

func NewTradeMasterService(db *pgxpool.Pool) *TradeMasterService {
	return &TradeMasterService{
		db,
	}
}

func (s *TradeMasterService) ExecuteTrade(trade *models.Trade) error {
	return nil
}
