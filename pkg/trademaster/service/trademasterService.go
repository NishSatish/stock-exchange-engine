package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
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
		return fmt.Errorf("failed to unmarshal order: %w", err)
	}

	buyOrderBookKey := constructOrderBookKey(orderDto.StockID, "buy")
	sellOrderBookKey := constructOrderBookKey(orderDto.StockID, "sell")

	// TODO: Think about what happens when 2 matching buy sells are picked up concurrently
	if orderDto.OrderType == models.Buy {
		// Look for a matching sell order
		results, err := this.libsService.RedisClient.ZRangeWithScores(ctx, sellOrderBookKey, 0, 0).Result()
		if err != nil && err != redis.Nil {
			return fmt.Errorf("failed to get best sell order: %w", err)
		}

		if len(results) > 0 {
			bestSellOrder := results[0]
			if orderDto.Price >= bestSellOrder.Score {
				// Match found
				fmt.Printf("Match found for BUY order %s: selling at %f\n", orderDto.OrderID, bestSellOrder.Score)

				// For simplicity, assume full match for now
				// Create trade
				sellOrderID := bestSellOrder.Member.(string)
				trade := &models.Trade{
					StockID:     orderDto.StockID,
					BuyOrderID:  orderDto.OrderID,
					SellOrderID: sellOrderID,
					Quantity:    orderDto.Quantity, // Assuming full quantity match
					Price:       bestSellOrder.Score,
				}

				if err := this.ExecuteTrade(trade); err != nil {
					return err
				}

				// Remove matched orders
				this.libsService.RedisClient.ZRem(ctx, sellOrderBookKey, sellOrderID)

				return nil
			}
		}

		// No match found, add to buy order book
		this.libsService.RedisClient.ZAdd(ctx, buyOrderBookKey, redis.Z{
			Score:  orderDto.Price,
			Member: orderDto.OrderID,
		})

		fmt.Println("Buy order added to order book %w", buyOrderBookKey)

	} else if orderDto.OrderType == models.Sell {
		// Look for a matching buy order
		results, err := this.libsService.RedisClient.ZRevRangeWithScores(ctx, buyOrderBookKey, 0, 0).Result()
		if err != nil && err != redis.Nil {
			return fmt.Errorf("failed to get best buy order: %w", err)
		}

		if len(results) > 0 {
			bestBuyOrder := results[0]
			if orderDto.Price <= bestBuyOrder.Score {
				// Match found
				fmt.Printf("Match found for SELL order %s: buying at %f\n", orderDto.OrderID, bestBuyOrder.Score)

				// For simplicity, assume full match for now
				// Create trade
				buyOrderID := bestBuyOrder.Member.(string)
				trade := &models.Trade{
					StockID:     orderDto.StockID,
					BuyOrderID:  buyOrderID,
					SellOrderID: orderDto.OrderID,
					Quantity:    orderDto.Quantity, // Assuming full quantity match
					Price:       bestBuyOrder.Score,
				}

				if err := this.ExecuteTrade(trade); err != nil {
					fmt.Println("Error saving trade in db", err)
					return err
				}

				// Remove matched orders
				fmt.Println("REmoving buy order from order book:", buyOrderBookKey, buyOrderID)
				this.libsService.RedisClient.ZRem(ctx, buyOrderBookKey, buyOrderID)
				return nil
			}
		}
		// No match found, add to sell order book
		this.libsService.RedisClient.ZAdd(ctx, sellOrderBookKey, redis.Z{
			Score:  orderDto.Price,
			Member: orderDto.OrderID,
		})
	}

	return nil
}

func (this *TradeMasterService) ExecuteTrade(trade *models.Trade) error {
	_, err := this.db.TradeRepo.CreateTrade(context.Background(), *trade)
	if err != nil {
		return fmt.Errorf("failed to create trade: %w", err)
	}
	fmt.Printf("Trade executed successfully: %+v\n", trade)
	return nil
}

func constructOrderBookKey(stockID string, orderType string) string {
	return "orderbook:" + stockID + ":" + orderType
}
