package shared_events

import (
	"context"
	"encoding/json"
)

type KafkaProducer interface {
	SendMessage(ctx context.Context, topic string, key, value []byte) error
}
type EventPublisher struct {
	kp KafkaProducer
}

func NewEventPublisher(kp KafkaProducer) *EventPublisher {
	return &EventPublisher{kp: kp}
}
func (p *EventPublisher) Publisher(ctx context.Context, topic string, key string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return p.kp.SendMessage(ctx, topic, []byte(key), data)
}
