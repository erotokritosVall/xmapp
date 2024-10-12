package kafka

import (
	"context"
	"os"

	"github.com/rs/zerolog/log"
	kf "github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader    *kf.Reader
	stopChann chan os.Signal
}

func NewConsumer(cfg *KafkaConfig) *Consumer {
	return &Consumer{
		reader: kf.NewReader(kf.ReaderConfig{
			Brokers: []string{cfg.Address},
			Topic:   cfg.Topic,
		}),
		stopChann: make(chan os.Signal),
	}
}

func (c *Consumer) Start(receiver chan<- *Message) {
	go func() {
		for {
			select {
			case <-c.stopChann:
				log.Debug().Msgf("consumer stopping")
				c.reader.Close()
				return

			default:
				m, err := c.reader.FetchMessage(context.Background())
				if err != nil {
					log.Err(err).Msg("failed to consume kafka message")
					break
				}

				if len(m.Headers) == 0 || len(m.Value) == 0 {
					break
				}

				headers := make(map[string][]byte, len(m.Headers))
				for _, h := range m.Headers {
					headers[h.Key] = h.Value
				}

				message := &Message{
					Headers: headers,
					Body:    m.Value,
				}

				receiver <- message
			}
		}
	}()
}

func (c *Consumer) Stop(s os.Signal) {
	c.stopChann <- s
}
