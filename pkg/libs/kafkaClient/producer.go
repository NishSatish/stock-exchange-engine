package kafkaClient

import (
	"context"
	"github.com/segmentio/kafka-go"
)

type Writer struct {
	w *kafka.Writer
}

func NewWriter(kafkaURL, topic string) *Writer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(kafkaURL),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
		Async:        true,
		Completion: func(messages []kafka.Message, err error) {
			if err != nil {
				// In a real app, use a structured logger
				println("kafka write failed:", err.Error())
			}
		},
	}
	return &Writer{w: writer}
}

func (kw *Writer) WriteMessages(ctx context.Context, messages ...kafka.Message) error {
	return kw.w.WriteMessages(ctx, messages...)
}

// Close closes the Kafka writer.
func (kw *Writer) Close() error {
	return kw.w.Close()
}
