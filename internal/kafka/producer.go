package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   topic,
	})
	return &Producer{writer: writer}
}

func (p *Producer) Publish(message string) error {
	return p.writer.WriteMessages(
		context.Background(),
		kafka.Message{Value: []byte(message)},
	)
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
