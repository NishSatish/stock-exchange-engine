package libs

import (
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"stock-exchange
	"stock-exchange-simulator/pkg/libs/kafkaClient"
	"stock-exchange-simulator/pkg/libs/taskqueue"
)

// LibsFactory holds clients for all external libraries.
type LibsFactory struct {
	TaskQueueClient *taskqueue.TaskClient
	RedisClient     *redis.Client
	KafkaFactory    *kafkaClient.KafkaFactory
}

// NewLibsFactory creates and configures all library clients.
func NewLibsFactory() *LibsFactory {
	redisOpt := asynq.RedisClientOpt{
		Addr: "localhost:6379",
	}

	taskClient := taskqueue.NewTaskClient(redisOpt)

	redisClient := redis.NewClient(&redis.Options{
		Addr: redisOpt.Addr,
	})

	kafkaFactory := kafkaClient.NewKafkaFactory()

	return &LibsFactory{
		TaskQueueClient: taskClient,
		RedisClient:     redisClient,
		KafkaFactory:    kafkaFactory,
	}
}
