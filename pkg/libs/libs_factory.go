package libs

import (
	"stock-exchange-simulator/pkg/libs/taskqueue"

	"github.com/hibiken/asynq"
)

// LibsFactory holds clients for all external libraries.
type LibsFactory struct {
	TaskQueueClient *taskqueue.TaskClient
}

// NewLibsFactory creates and configures all library clients.
func NewLibsFactory() *LibsFactory {
	redisOpt := asynq.RedisClientOpt{
		Addr: "localhost:6379",
	}

	taskClient := taskqueue.NewTaskClient(redisOpt)

	return &LibsFactory{
		TaskQueueClient: taskClient,
	}
}
