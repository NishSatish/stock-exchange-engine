package service

import (
	"context"
	"encoding/json"
	"fmt"
	"stock-exchange-simulator/pkg/db"
	"stock-exchange-simulator/pkg/libs"
	"stock-exchange-simulator/pkg/libs/taskqueue/dto"
	"stock-exchange-simulator/pkg/models"
	trademaster "stock-exchange-simulator/pkg/trademaster/service"

	"github.com/hibiken/asynq"
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
	libsService        *libs.LibsFactory
}

func NewNexusService(
	trademasterService trademaster.ITradeMasterServiceInterface,
	dbService *db.RepositoryFactory,
	libsService *libs.LibsFactory,
) *NexusService {
	return &NexusService{
		trademasterService,
		dbService,
		libsService,
	}
}

func (s *NexusService) PlaceOrder(order *models.Order) (string, error) {
	createdOrder, err := s.dbService.OrderRepo.CreateOrder(context.Background(), *order)
	if err != nil {
		// TODO: find a way to transaction-alise the db save and enqueuing, because the order matching engine
		// won't work without the queue item
		return "", fmt.Errorf("saving placed order to DB failed: %w", err)
	}
	fmt.Println("Order successfully created in db", createdOrder.ID)

	queuePayload, marshallErr := json.Marshal(dto.EnqueueOrderPlacedDTO{
		StockID:   order.StockID,
		Price:     order.Price,
		OrderType: order.Type,
	})
	if marshallErr != nil {
		return "", fmt.Errorf("enqueueing placed order failed: %w", marshallErr)
	}

	task := asynq.NewTask(dto.EnqueueOrderPlacedTopic, queuePayload)
	taskInfo, enqueueErr := s.libsService.TaskQueueClient.Enqueue(task)
	if enqueueErr != nil {
		return "", fmt.Errorf("enqueueing placed order failed: %w", enqueueErr)
	}
	fmt.Println("open order sent for processing", taskInfo.ID)
	return createdOrder.ID, nil
}
