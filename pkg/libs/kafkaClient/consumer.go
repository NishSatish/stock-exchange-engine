package kafkaClient

import (
	"context"
	"github.com/segmentio/kafka-go"
)

type Reader struct {
	r *kafka.Reader
}

func NewReader(kafkaURL, topic, groupID string) *Reader {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{kafkaURL},
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	return &Reader{r: reader}
}

func (kr *Reader) ReadMessage(ctx context.Context) (kafka.Message, error) {
	return kr.r.ReadMessage(ctx)
}

// Close closes the Kafka reader.
func (kr *Reader) Close() error {
	return kr.r.Close()
}
