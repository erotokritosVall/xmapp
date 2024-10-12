package pubsub

import (
	"context"
	"encoding/json"

	"github.com/erotokritosVall/xmapp/internal/events"
	"github.com/erotokritosVall/xmapp/pkg/kafka"
)

type PublisherManager struct {
	publisher *kafka.Publisher
}

func NewPublisher(cfg *kafka.KafkaConfig) *PublisherManager {
	return &PublisherManager{
		publisher: kafka.NewPublisher(cfg),
	}
}

func (p *PublisherManager) PublishDomainEvents(ctx context.Context, evts []events.DomainEvent) error {
	toPublish := make([]kafka.Message, len(evts))

	for _, e := range evts {
		body, err := json.Marshal(e)
		if err != nil {
			return err
		}

		headers := map[string][]byte{
			domainEventTypeHeader: []byte(e.Type().String()),
			domainEventIdHeader:   []byte(e.UniqueId()),
		}

		m := kafka.Message{
			Headers: headers,
			Body:    body,
		}

		toPublish = append(toPublish, m)
	}

	return p.publisher.Publish(ctx, toPublish)
}
