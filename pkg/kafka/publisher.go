package kafka

import (
	"context"

	kf "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/protocol"
)

type Publisher struct {
	writer *kf.Writer
}

func NewPublisher(cfg *KafkaConfig) *Publisher {
	p := &Publisher{
		writer: kf.NewWriter(kf.WriterConfig{
			Brokers: []string{cfg.Address},
			Topic:   cfg.Topic,
		}),
	}

	return p
}

func (p *Publisher) Publish(ctx context.Context, messages []Message) error {
	toPublish := make([]kf.Message, len(messages))

	for _, m := range messages {
		kfMsg := kf.Message{
			Value:   m.Body,
			Headers: make([]protocol.Header, len(m.Headers)),
		}

		for k, v := range m.Headers {
			kfMsg.Headers = append(kfMsg.Headers, kf.Header{
				Key:   k,
				Value: v,
			})
		}

		toPublish = append(toPublish, kfMsg)
	}

	return p.writer.WriteMessages(ctx, toPublish...)
}
