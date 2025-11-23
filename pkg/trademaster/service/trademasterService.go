package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"stock-exchange-simulator/pkg/db"
	"stock-exchange-simulator/pkg/libs"
	"stock-exchange-simulator/pkg/libs/taskqueue/dto"
	"stock-exchange-simulator/pkg/models"
)

type TradeMasterService struct {
	db          *db.RepositoryFactory
	libsService *libs.LibsFactory
}

type ITradeMasterServiceInterface interface {
	ExecuteTrade(trade *models.Trade) error
	OrderProcessor(ctx context.Context, enqueuedOrder *asynq.Task) error
}

func NewTradeMasterService(db *db.RepositoryFactory, libsService *libs.LibsFactory) *TradeMasterService {
	return &TradeMasterService{
		db,
		libsService,
	}
}

func (this *TradeMasterService) OrderProcessor(ctx context.Context, enqueuedOrder *asynq.Task) error {
	var orderDto dto.EnqueueOrderPlacedDTO
	if err := json.Unmarshal(enqueuedOrder.Payload(), &orderDto); err != nil {
		fmt.Println("ERROOOOOOOOOO")
	}
	fmt.Println("Your order is here", orderDto)
	return nil
}

func (this *TradeMasterService) ExecuteTrade(trade *models.Trade) error {

	return nil
}
