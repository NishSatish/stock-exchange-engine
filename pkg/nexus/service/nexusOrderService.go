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
	"time"

	"github.com/hibiken/asynq"
)

/*
  NEXUS is the  client facing module of the system
  for accepting order requests from users.
*/

type INexusServiceInterface interface {
	PlaceOrder(order *models.Order) (string, error)
	DumpOrderBookToRedis()
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
		OrderID:   createdOrder.ID,
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

func (s *NexusService) DumpOrderBookToRedis() {
	/*
		Use this function in case of a case in which redis just gets empty.
		We currently don't have a read-thru mechanism for redis, and the whole order matching engine works on the redis in memroy order book.
		So if redis is empty the whole system tanks

		TODO: Currently have this wired up to an endpoint to run on demand, will later have an automated redis health script or something.
	*/
	fmt.Println("RUNNING ORDERBOOK HYDRATION JOB")
	startTime := time.Now()
	pendingOrders, err := s.dbService.OrderRepo.FetchOrdersByStatus(context.Background(), models.Pending)
	if err != nil {
		_ = fmt.Errorf("failed to fetch pending orders from DB: %w", err)
	}

	for _, order := range pendingOrders {
		orderPayload, marshallErr := json.Marshal(dto.EnqueueOrderPlacedDTO{
			OrderID:   order.ID,
			StockID:   order.StockID,
			Price:     order.Price,
			OrderType: order.Type,
			Quantity:  order.Quantity,
		})
		if marshallErr != nil {
			_ = fmt.Errorf("failed to marshal order for enqueueing: %w", marshallErr)
		}

		task := asynq.NewTask(dto.EnqueueOrderPlacedTopic, orderPayload)
		// Good to enqueue one by one so that the processor activates and can also match orders during the hydration. I mean don't manually write sorted set members.
		_, enqueueErr := s.libsService.TaskQueueClient.Enqueue(task)
		if enqueueErr != nil {
			_ = fmt.Errorf("failed to enqueue order for processing: %w", enqueueErr)
		}
	}

	fmt.Println("HYDRATION JOB DONE IN", time.Since(startTime))
}
