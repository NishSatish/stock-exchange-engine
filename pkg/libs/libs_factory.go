package libs

import (
	"github.com/redis/go-redis/v9"
	"stock-exchange-simulator/pkg/libs/taskqueue"

	"github.com/hibiken/asynq"
)

// LibsFactory holds clients for all external libraries.
type LibsFactory struct {
	TaskQueueClient *taskqueue.TaskClient
	RedisClient     *redis.Client
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

	return &LibsFactory{
		TaskQueueClient: taskClient,
		RedisClient:     redisClient,
	}
}
