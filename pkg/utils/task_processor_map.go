package utils

import (
	"github.com/hibiken/asynq"
)

func MapEventsToProcessors(mux *asynq.ServeMux) *asynq.ServeMux {
	// Just take the mux object, assign handlers and send it back

	/*
	 * ORDER HANDLERS
	 */
	//mux.HandleFunc(dto.EnqueueOrderPlacedTopic)

	return mux
}
