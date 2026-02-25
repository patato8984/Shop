package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/patato8984/Shop/internal/app"
	infra_kafka "github.com/patato8984/Shop/internal/infra/kafka"
	"github.com/patato8984/Shop/internal/infra/redis"
	shared_events "github.com/patato8984/Shop/internal/shared/events"
	"github.com/patato8984/Shop/internal/shared/logger"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Envelope struct {
	Type     string          `json:"type"`
	Payload  json.RawMessage `json:"payload"`
	MetaDate json.RawMessage `json:"meta_date"`
}

func main() {
	logger.Init()
	log := logger.L()
	config := app.LoadCacheConfig()
	redis, err := redis.InitRedis(config.Addr, config.RedisPassword)
	if err != nil {
		log.Fatal("redis",
			zap.Error(err),
		)
	}
	manager := app.DependencyInitiationCacheInvalidator(redis, log)
	kp := infra_kafka.NewProducer(config.Broker, false, log)
	readerCatalog := infra_kafka.NewConsumer(config.Broker, "catalog_event", "clear_cache")
	readerCart := infra_kafka.NewConsumer(config.Broker, "cart_event", "clear_cache")
	ctxGF, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	handlers := func(ctxGF context.Context, msg kafka.Message) error {
		log.Info("kafka", zap.String("message", fmt.Sprintf("%s -> value: %s", msg.Key, string(msg.Value))))
		var event Envelope
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			return err
		}
		if handler, ok := manager.Handler[event.Type]; ok {
			return handler(context.Background(), event.Payload)
		}
		errMsg := fmt.Sprintf("unsupported event type: %s", event.Type)
		log.Warn("skipping event", zap.String("type", event.Type))
		if er := kp.SendMessage(ctxGF, "topic error", []byte(errMsg), msg.Value); er != nil {
			log.Error("event", zap.Error(er))
		}
		return nil
	}
	shared_events.RunEventListener(ctxGF, readerCatalog, handlers)
	shared_events.RunEventListener(ctxGF, readerCart, handlers)
	<-ctxGF.Done()
	log.Info("completion of work...")
	if err := readerCatalog.Close(); err != nil {
		log.Info("readerCatalog kafka forced to shutdown", zap.Error(err))
	}
	log.Info("reader catalog stopped")
	if err := readerCart.Close(); err != nil {
		log.Info("readerCart kafka forced to shutdown", zap.Error(err))
	}
	log.Info("reader cart stopped")
	if err := redis.Close(); err != nil {
		log.Info("redis forced to shutdown", zap.Error(err))
	}
	log.Info("redis stopped")
	log.Info("kafka stopped")
}
