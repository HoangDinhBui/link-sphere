package kafka

import (
	"context"

	"github.com/rs/zerolog/log"
	kafkago "github.com/segmentio/kafka-go"
)

// Consumer wraps a Kafka reader for consuming messages.
type Consumer struct {
	reader *kafkago.Reader
}

// NewConsumer creates a new Kafka consumer for the given topic and group.
func NewConsumer(brokers []string, topic, groupID string) *Consumer {
	r := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})
	return &Consumer{reader: r}
}

// MessageHandler is a function that processes a Kafka message.
type MessageHandler func(key, value []byte) error

// Consume starts consuming messages and calls handler for each one.
func (c *Consumer) Consume(ctx context.Context, handler MessageHandler) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				log.Error().Err(err).Msg("failed to read kafka message")
				continue
			}
			if err := handler(msg.Key, msg.Value); err != nil {
				log.Error().Err(err).Msg("failed to handle kafka message")
			}
		}
	}
}

// Close closes the consumer connection.
func (c *Consumer) Close() error {
	return c.reader.Close()
}
