package kafkaClient

// Adding this factory file because i want to be ableto spin up topics and their producers and consumers on demand
// only the kafka infra level config will be in the libs

const (
	KafkaURL = "localhost:9092"
)

type KafkaFactory struct {
}

func NewKafkaFactory() *KafkaFactory {
	return &KafkaFactory{}
}

func (f *KafkaFactory) NewProducer(topic string) *Writer {
	return NewWriter(KafkaURL, topic)
}

func (f *KafkaFactory) NewConsumer(topic, groupID string) *Reader {
	return NewReader(KafkaURL, topic, groupID)
}
