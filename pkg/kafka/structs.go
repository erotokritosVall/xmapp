package kafka

type KafkaConfig struct {
	Address string `envconfig:"KAFKA_ADDRESS"`
	Topic   string `envconfig:"KAFKA_TOPIC"`
}

type Message struct {
	Headers map[string][]byte
	Body    []byte
}
