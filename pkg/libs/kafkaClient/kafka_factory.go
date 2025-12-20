package kafkaClient

// Adding this factory file because i want to be ableto spin up topics and their producers and consumers on demand
// only the kafka infra level config will be in the libs

type KafkaFactory struct {
	url string
}

func NewKafkaFactory(kafkaUrl string) *KafkaFactory {
	return &KafkaFactory{
		url: kafkaUrl,
	}
}

func (f *KafkaFactory) NewProducer(topic string) *Writer {
	return NewWriter(f.url, topic)
}

func (f *KafkaFactory) NewConsumer(topic, groupID string) *Reader {
	return NewReader(f.url, topic, groupID)
}
