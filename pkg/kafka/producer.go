package kafka

import (
	"context"

	"github.com/rs/zerolog/log"
	kafkago "github.com/segmentio/kafka-go"
)

// Producer wraps a Kafka writer for publishing messages.
type Producer struct {
	writer *kafkago.Writer
}

// NewProducer creates a new Kafka producer for the given topic.
func NewProducer(brokers []string, topic string) *Producer {
	w := &kafkago.Writer{
		Addr:     kafkago.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafkago.LeastBytes{},
	}
	return &Producer{writer: w}
}

// Publish sends a message to Kafka.
func (p *Producer) Publish(ctx context.Context, key, value []byte) error {
	err := p.writer.WriteMessages(ctx, kafkago.Message{
		Key:   key,
		Value: value,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to publish kafka message")
		return err
	}
	return nil
}

// Close closes the producer connection.
func (p *Producer) Close() error {
	return p.writer.Close()
}
