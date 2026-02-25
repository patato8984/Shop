package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/patato8984/Shop/internal/app"
	infra_kafka "github.com/patato8984/Shop/internal/infra/kafka"
	catalog_model "github.com/patato8984/Shop/internal/modules/catalog/model"
	shared_events "github.com/patato8984/Shop/internal/shared/events"
	"github.com/patato8984/Shop/internal/shared/logger"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Envelope struct {
	Type     string          `json:"type"`
	Payload  json.RawMessage `json:"payload"`
	MetaDate MetaDate        `json:"meta_date"`
}
type MetaDate struct {
	IDUser   int    `json:"id_user"`
	RoleUser string `json:"role"`
}

func main() {
	logger.Init()
	log := logger.L()
	config := app.LoadCDCConfig()
	kp := infra_kafka.NewProducer(config.Broker, true, log)
	readerCatalog := infra_kafka.NewConsumer(config.Broker, "catalog_event", "change_data_capture")
	ctxGF, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	shared_events.RunEventListener(ctxGF, readerCatalog, func(ctx context.Context, msg kafka.Message) error {
		var event Envelope
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			return err
		}
		switch event.Type {
		case "product_create":
			var product catalog_model.Product
			if err := json.Unmarshal(event.Payload, &product); err != nil {
				if er := kp.SendMessage(ctxGF, "topic error", []byte(err.Error()), msg.Value); er != nil {
					log.Error("event", zap.Error(er))
				}
				return err
			}
			log.Info(event.Type,
				zap.Int("id_user", event.MetaDate.IDUser),
				zap.String("role", event.MetaDate.RoleUser),
				zap.Int("id_product", product.Id),
				zap.String("name_product", product.Name),
			)
			return nil
		case "product_delete":
			var product catalog_model.Product
			if err := json.Unmarshal(event.Payload, &product); err != nil {
				if er := kp.SendMessage(ctxGF, "topic error", []byte(err.Error()), msg.Value); er != nil {
					log.Error("event", zap.Error(er))
				}
				return err
			}
			log.Info(event.Type,
				zap.Int("id_user", event.MetaDate.IDUser),
				zap.String("role", event.MetaDate.RoleUser),
				zap.Int("id_product", product.Id),
			)
			return nil
		case "skus_created":
			var product catalog_model.Product
			if err := json.Unmarshal(event.Payload, &product); err != nil {
				if er := kp.SendMessage(ctxGF, "topic error", []byte(err.Error()), msg.Value); er != nil {
					log.Error("event", zap.Error(er))
				}
				return err
			}
			if len(product.SKUs) == 0 {
				return errors.New("SKUs not found")
			}
			log.Info(event.Type,
				zap.Int("id_user", event.MetaDate.IDUser),
				zap.String("role", event.MetaDate.RoleUser),
				zap.Int("id_product", product.Id),
				zap.Int("id_sku", product.SKUs[0].Id),
				zap.String("name_product", product.Name),
			)
			return nil
		case "skus_addStock":
			var product catalog_model.StockUpdatedLoad
			if err := json.Unmarshal(event.Payload, &product); err != nil {
				if er := kp.SendMessage(ctxGF, "topic error", []byte(err.Error()), msg.Value); er != nil {
					log.Error("event", zap.Error(er))
				}
				return err
			}
			log.Info(event.Type,
				zap.Int("id_user", event.MetaDate.IDUser),
				zap.String("role", event.MetaDate.RoleUser),
				zap.Int("id_skus", product.SkusID),
				zap.Int("new_stock", product.NewStock),
			)
			return nil
		case "catalog_event":
			var product catalog_model.PriceUpdatedLoad
			if err := json.Unmarshal(event.Payload, &product); err != nil {
				if er := kp.SendMessage(ctxGF, "topic error", []byte(err.Error()), msg.Value); er != nil {
					log.Error("event", zap.Error(er))
				}
				return err
			}
			log.Info(event.Type,
				zap.Int("id_user", event.MetaDate.IDUser),
				zap.String("role", event.MetaDate.RoleUser),
				zap.Int("id_skus", product.SkusID),
				zap.Float64("new_price", product.NewPrice),
			)
			return nil
		case "skus_deleted":
			var product catalog_model.SkusDeleteLoad
			if err := json.Unmarshal(event.Payload, &product); err != nil {
				if er := kp.SendMessage(ctxGF, "topic error", []byte(err.Error()), msg.Value); er != nil {
					log.Error("event", zap.Error(er))
				}
				return err
			}
			log.Info(event.Type,
				zap.Int("id_user", event.MetaDate.IDUser),
				zap.String("role", event.MetaDate.RoleUser),
				zap.Int("id_skus", product.SkusID),
				zap.Int("id_product", product.ProductID),
			)
			return nil
		}
		fmt.Printf("skipping message: %s", event.Type)
		return nil
	})
	<-ctxGF.Done()
	log.Info("completion of work...")
	if err := readerCatalog.Close(); err != nil {
		log.Info("readerCatalog kafka forced to shutdown", zap.Error(err))
	}
	log.Info("kafka stopped")
}
