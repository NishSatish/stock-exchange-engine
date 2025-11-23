package taskqueue

import "github.com/hibiken/asynq"

// TaskClient is a wrapper for the asynq client.
type TaskClient struct {
	*asynq.Client
}

func NewTaskClient(redisOpt asynq.RedisClientOpt) *TaskClient {
	return &TaskClient{
		Client: asynq.NewClient(redisOpt),
	}
}
