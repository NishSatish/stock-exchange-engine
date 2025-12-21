package service

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"stock-exchange-simulator/pkg/db"
	"stock-exchange-simulator/pkg/libs"
	"stock-exchange-simulator/pkg/libs/kafkaClient"
	kafkaDto "stock-exchange-simulator/pkg/libs/kafkaClient/dto"
	"stock-exchange-simulator/pkg/libs/taskqueue/dto"
	"stock-exchange-simulator/pkg/models"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

type TradeMasterService struct {
	db            *db.RepositoryFactory
	libsService   *libs.LibsFactory
	kafkaProducer *kafkaClient.Writer
}

type ITradeMasterServiceInterface interface {
	ExecuteTrade(trade *models.Trade) error
	OrderProcessor(ctx context.Context, enqueuedOrder *asynq.Task) error
}

func NewTradeMasterService(db *db.RepositoryFactory, libsService *libs.LibsFactory) *TradeMasterService {
	kafkaProducer := libsService.KafkaFactory.NewProducer(kafkaDto.LtpGenerationTopic)
	return &TradeMasterService{
		db,
		libsService,
		kafkaProducer,
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
		allSellOrdersForStock, err := this.libsService.RedisClient.ZRangeWithScores(ctx, sellOrderBookKey, 0, int64(orderDto.Price)).Result()
		if err != nil && err != redis.Nil {
			return fmt.Errorf("failed to get best sell order: %w", err)
		}

		if len(allSellOrdersForStock) > 0 {
			orderMatchIndex := slices.IndexFunc(allSellOrdersForStock, func(o redis.Z) bool {
				return o.Score == orderDto.Price
			})

			if orderMatchIndex != -1 {
				// Match found
				matchingOrder := allSellOrdersForStock[orderMatchIndex]
				fmt.Printf("Match found for BUY order %s: selling at %f\n", orderDto.OrderID, matchingOrder.Score)

				sellOrderID := matchingOrder.Member.(string)
				trade := &models.Trade{
					StockID:     orderDto.StockID,
					BuyOrderID:  orderDto.OrderID,
					SellOrderID: sellOrderID,
					Quantity:    orderDto.Quantity, // Assuming full quantity match TODO
					Price:       matchingOrder.Score,
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
		allBuyOrdersForStock, err := this.libsService.RedisClient.ZRevRangeWithScores(ctx, buyOrderBookKey, 0, int64(orderDto.Price)).Result()
		fmt.Println("AHA Here are all the buy orders that i needed", allBuyOrdersForStock)
		if err != nil && err != redis.Nil {
			return fmt.Errorf("failed to get best buy order: %w", err)
		}

		if len(allBuyOrdersForStock) > 0 {
			orderMatchIndex := slices.IndexFunc(allBuyOrdersForStock, func(o redis.Z) bool {
				return o.Score == orderDto.Price
			})
			if orderMatchIndex != -1 {
				// Match found
				matchingOrder := allBuyOrdersForStock[orderMatchIndex]
				fmt.Printf("Match found for SELL order %s: buying at %f\n", orderDto.OrderID, matchingOrder.Score)

				// For simplicity, assume full match for now
				// Create trade
				buyOrderID := matchingOrder.Member.(string)
				trade := &models.Trade{
					StockID:     orderDto.StockID,
					BuyOrderID:  buyOrderID,
					SellOrderID: orderDto.OrderID,
					Quantity:    orderDto.Quantity, // Assuming full quantity match
					Price:       matchingOrder.Score,
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
	// Difficult to revert a db operation, so do it at the end of whatever other stuff you want to do coz they can be reversed
	// this is just the kafka dispatch part, let kafka handle the failures once the message is received
	err := this.dispatchTradeToGenerateLTP(trade)
	if err != nil {
		return fmt.Errorf("failed to dispatch to kafka ltp generation queue: %w", err)
	}

	// TODO: Transactionalise this whole thing (trade creation and order completetion)
	createdTrade, err := this.db.TradeRepo.CreateTrade(context.Background(), *trade)
	if err != nil {
		return fmt.Errorf("failed to create trade: %w", err)
	}

	orderIDs := []string{createdTrade.BuyOrderID, createdTrade.SellOrderID}
	err = this.db.OrderRepo.MarkOrdersCompleteInBulk(context.Background(), orderIDs, createdTrade.ID)
	if err != nil {
		return fmt.Errorf("failed to mark orders as complete: %w", err)
	}

	fmt.Printf("Trade executed successfully: %+v\n", createdTrade)
	return nil
}

func (this *TradeMasterService) dispatchTradeToGenerateLTP(trade *models.Trade) error {
	kafkaKeyString := trade.StockID + "_" + trade.Timestamp.String()
	ltpGenerationKafkaPayload, err := this.kafkaProducer.PrepareKafkaPayload(kafkaKeyString, kafkaDto.LtpGenerationKafkaDTO{
		StockID:   trade.StockID,
		Price:     trade.Price,
		Timestamp: trade.Timestamp,
	})

	if err != nil {
		fmt.Errorf("Failed to prepare kafkaPayload %s", kafkaKeyString)
	}
	
	err = this.kafkaProducer.WriteMessages(context.Background(), ltpGenerationKafkaPayload)
	if err != nil {
		fmt.Println("Failed to dispatch trade for LTP generation:", err)
		return err
	}

	fmt.Println("Dispatched trade for LTP generation:", trade.ID)
	return nil
}

func constructOrderBookKey(stockID string, orderType string) string {
	return "orderbook:" + stockID + ":" + orderType
}
