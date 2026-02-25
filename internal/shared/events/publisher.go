package shared_events

import (
	"context"
	"fmt"

	infra_kafka "github.com/patato8984/Shop/internal/infra/kafka"
	"github.com/segmentio/kafka-go"
)

func RunEventListener(ctx context.Context, reader *infra_kafka.Consumer, handler func(ctx context.Context, msg kafka.Message) error) {
	go func() {
		for {
			msg, err := reader.FetchMessage(ctx)
			if err != nil {
				return
			}
			if err := handler(ctx, msg); err == nil {
				if err := reader.Commit(ctx, msg); err != nil {
					fmt.Print("error commit")
				}
			} else {
				fmt.Printf("Consumer error: %v\n", err)
				if err := reader.Commit(ctx, msg); err != nil {
					fmt.Print("error commit")
				}
			}
		}
	}()
}
