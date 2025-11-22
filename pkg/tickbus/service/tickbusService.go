package service

import "github.com/jackc/pgx/v5/pgxpool"

/*
 TickBus is the module for watching trades get matched
 and to generate "ticks" and candles on a stream.
*/

type ITickBusServiceInterface interface {
	PublishTick(symbol string, price float64)
}

type TickBusService struct {
	db *pgxpool.Pool
}

func NewTickBusService(db *pgxpool.Pool) *TickBusService {
	return &TickBusService{
		db,
	}
}

func (s *TickBusService) PublishTick(symbol string, price float64) {
}
