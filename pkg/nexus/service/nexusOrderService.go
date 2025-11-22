package service

import (
	"context"
	"fmt"
	"stock-exchange-simulator/pkg/db"
	"stock-exchange-simulator/pkg/models"
	trademaster "stock-exchange-simulator/pkg/trademaster/service"
)

/*
  NEXUS is the  client facing module of the system
  for accepting order requests from users.
*/

type INexusServiceInterface interface {
	PlaceOrder(order *models.Order) (string, error)
}

type NexusService struct {
	trademasterService trademaster.ITradeMasterServiceInterface
	dbService          *db.RepositoryFactory
}

func NewNexusService(
	trademasterService trademaster.ITradeMasterServiceInterface,
	dbService *db.RepositoryFactory,
) *NexusService {
	return &NexusService{
		trademasterService: trademasterService,
		dbService:          dbService,
	}
}

func (s *NexusService) PlaceOrder(order *models.Order) (string, error) {
	createdOrder, err := s.dbService.OrderRepo.CreateOrder(context.Background(), *order)
	if err != nil {
		return "", fmt.Errorf("failed to create order: %w", err)
	}
	return createdOrder.ID, nil
}
