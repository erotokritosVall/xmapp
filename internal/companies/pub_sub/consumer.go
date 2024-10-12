package pubsub

import (
	"encoding/json"
	"os"

	"github.com/erotokritosVall/xmapp/internal/events"
	"github.com/erotokritosVall/xmapp/pkg/kafka"
	"github.com/rs/zerolog/log"
)

type ConsumerManager struct {
	consumer  kafka.Consumer
	stopChann chan os.Signal
}

func NewConsumerManager(cfg *kafka.KafkaConfig) *ConsumerManager {
	return &ConsumerManager{
		consumer:  *kafka.NewConsumer(cfg),
		stopChann: make(chan os.Signal),
	}
}

func (c *ConsumerManager) Start() {
	ch := make(chan *kafka.Message)

	c.consumer.Start(ch)

	go func() {
		defer close(ch)

		for {
			select {
			case <-c.stopChann:
				log.Debug().Msgf("consumer manager stopping")
				return

			case incomingMessage := <-ch:
				domainEventTypeStr := incomingMessage.Headers[domainEventTypeHeader]

				domainEventType, err := events.DomainEventTypeString(string(domainEventTypeStr))
				if err != nil {
					log.Err(err).Msg("received invalid domain event type header")
				}

				var e events.DomainEvent

				switch domainEventType {
				case events.EventTypeCompanyCreated:
					c := &events.CompanyCreated{}
					if err := json.Unmarshal(incomingMessage.Body, c); err != nil {
						log.Err(err).Msg("failed to unmarshal CompanyCreated event")
					}

					e = c

				case events.EventTypeCompanyUpdated:
					c := &events.CompanyUpdated{}
					if err := json.Unmarshal(incomingMessage.Body, c); err != nil {
						log.Err(err).Msg("failed to unmarshal CompanyUpdated event")
					}

					e = c

				case events.EventTypeCompanyDeleted:
					c := &events.CompanyDeleted{}
					if err := json.Unmarshal(incomingMessage.Body, c); err != nil {
						log.Err(err).Msg("failed to unmarshal CompanyDeleted event")
					}

					e = c
				}

				if e != nil {
					log.Debug().Msgf("received domain event %s", e.UniqueId())
				}
			}
		}
	}()
}

func (c *ConsumerManager) Stop(s os.Signal) {
	c.consumer.Stop(s)
	c.stopChann <- s
}
