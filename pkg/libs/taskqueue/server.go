package taskqueue

import (
	"log"
	"stock-exchange-simulator/pkg/utils"

	"github.com/hibiken/asynq"
)

type TaskServer struct {
	Server *asynq.Server
}

func NewTaskServer(redisOpt asynq.RedisClientOpt) *TaskServer {
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: 10, // TODO: Check flexible concurrency per queue type
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)
	return &TaskServer{Server: server}
}

func (s *TaskServer) Run() error {
	var mux = asynq.NewServeMux()

	// only shitty part about Asynq, you have to register all processors to event types in one place
	mux = utils.MapEventsToProcessors(mux)

	log.Println("Asynq worker server started...")
	return s.Server.Start(mux)
}
