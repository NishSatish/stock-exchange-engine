package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"stock-exchange-simulator/pkg/libs"
	"stock-exchange-simulator/pkg/libs/kafkaClient"
	kafkaDto "stock-exchange-simulator/pkg/libs/kafkaClient/dto"
	"stock-exchange-simulator/pkg/tickbus/service/dto"
	"time"

	"github.com/redis/go-redis/v9"
)

/*
 TickBus is the module for watching trades get matched
 and to generate "ticks" and candles on a stream.
*/

type ITickBusServiceInterface interface {
	Start()
}

type TickBusService struct {
	redisClient *redis.Client
	kafkaReader *kafkaClient.Reader
	ctx         context.Context
}

func NewTickBusService(ctx context.Context, libs *libs.LibsFactory) *TickBusService {
	kafkaReader := libs.KafkaFactory.NewConsumer(kafkaDto.LtpGenerationTopic, "tickbus-group")
	return &TickBusService{
		redisClient: libs.RedisClient,
		kafkaReader: kafkaReader,
		ctx:         ctx,
	}
}

func (s *TickBusService) Start() {
	fmt.Println("kafka listener started")
	go func() {
		for {
			m, err := s.kafkaReader.ReadMessage(s.ctx)
			if err != nil {
				log.Println("error reading message from kafka", err)
				break
			}
			fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))

			var ltpDto kafkaDto.LtpGenerationKafkaDTO 
			if err := json.Unmarshal(m.Value, &ltpDto); err != nil {
				log.Println("error unmarshalling trade:", err)
				continue
			}

			log.Println("unmarshalled trade", ltpDto)

			s.processTrade(ltpDto)
		}
	}()
}

func (s *TickBusService) processTrade(ltpDto kafkaDto.LtpGenerationKafkaDTO) {
	// Get previous LTP data from Redis,
	// for calculating change.
	redisHashKey := "ltp:data"
	prevLTPData, err := s.redisClient.HGet(s.ctx, redisHashKey, ltpDto.StockID).Result()

	var ltpData dto.LTPData
	if err == redis.Nil {
		// successful not-found response
		// If no previous data, create a new one
		fmt.Printf("creating new ltp hash for %s \n", ltpDto.StockID)
		ltpData = dto.LTPData{}
	} else if err != nil {
		log.Println("error getting prev ltp data from redis", err)
		return
	} else {
		if err := json.Unmarshal([]byte(prevLTPData), &ltpData); err != nil {
			log.Println("error unmarshalling prev ltp data", err)
			return
		}
	}

	//Calculate new LTP, Change, Volume
	newLTP := dto.LTPData{
		StockID: ltpDto.StockID,
		LTP:         ltpDto.Price,
		Change:      ltpDto.Price - ltpData.LTP,
		LastUpdated: time.Now(),
	}

	// Update Redis Hash (HSET ltp:data <stock_symbol> <ltp_data>)
	newLTPJSON, err := json.Marshal(newLTP)
	if err != nil {
		log.Println("error marshalling new ltp data", err)
		return
	}

	if err := s.redisClient.HSet(s.ctx, redisHashKey, newLTP.StockID, newLTPJSON).Err(); err != nil {
		log.Println("error setting new ltp data in redis", err)
		return
	}
}
