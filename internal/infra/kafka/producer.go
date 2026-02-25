package infra_kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(broker []string, async bool, log *zap.Logger) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:        kafka.TCP(broker...),
			Balancer:    &kafka.LeastBytes{},
			MaxAttempts: 3,
			Async:       async,
			Completion: func(messages []kafka.Message, err error) {
				if err != nil {
					log.Error("Kafka async write failed",
						zap.Error(err),
						zap.String("topic", messages[0].Topic),
					)
				}
			},
		},
	}
}
func (p *Producer) SendMessage(ctx context.Context, topic string, key, value []byte) error {
	return p.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   key,
		Value: value,
	})
}
func (p *Producer) Close() error {
	return p.writer.Close()
}
