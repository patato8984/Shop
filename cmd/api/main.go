package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/patato8984/Shop/internal/app"
	infra_kafka "github.com/patato8984/Shop/internal/infra/kafka"
	"github.com/patato8984/Shop/internal/infra/postgres"
	"github.com/patato8984/Shop/internal/infra/redis"
	shared_events "github.com/patato8984/Shop/internal/shared/events"
	"github.com/patato8984/Shop/internal/shared/logger"
	"go.uber.org/zap"
)

func main() {
	logger.Init()
	log := logger.L()
	config, err := app.LoadApiConfig()
	if err != nil {
		log.Fatal("error load config",
			zap.Error(err),
		)
	}
	err = postgres.WaitForPostgres(config.DbPatch, time.Second*50)
	if err != nil {
		log.Fatal("error wait for postgres",
			zap.Error(err),
		)
	}
	database, err := postgres.NewPostgresConnection(config.DbPatch)
	if err != nil {
		log.Fatal("error db",
			zap.Error(err),
		)
	}
	redis, err := redis.InitRedis(config.Addr, config.RedisPassword)
	if err != nil {
		log.Fatal("error redis",
			zap.Error(err),
		)
	}
	kafka := infra_kafka.NewProducer([]string{"kafka:29092"}, true, log)
	kp := shared_events.NewEventPublisher(kafka)
	DiAndRouterMethods, standaloneHandlers := app.DependencyInitiationApi(database, redis, *config, log, kp)
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := standaloneHandlers.SeedAdmins.SeedAdmins(ctx); err != nil {
		log.Error("error add admin",
			zap.Error(err),
		)
		return
	}
	mux := http.NewServeMux()
	mux = app.NewApp(DiAndRouterMethods).RegisterAuthRoutes(mux)
	mux.Handle("POST /api/v1/webhook/bank", http.HandlerFunc(standaloneHandlers.WebHookBank.WebHookPayment))
	ctxGF, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go func() {
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				ctxWorker, cancel := context.WithTimeout(ctxGF, 10*time.Second)
				times := time.Now()
				err := standaloneHandlers.WorkerCartService.UpdateWorkerPrice(ctxWorker)
				duration := time.Since(times)
				if err != nil {
					log.Error("worker",
						zap.Error(err),
						zap.Duration("time", duration),
					)
				} else {
					log.Info("worker",
						zap.String("message", "the prices in the carts have been updated"),
						zap.Duration("time", duration),
					)
				}
				cancel()
			case <-ctxGF.Done():
				log.Info("worker stopping")
				return
			}
		}
	}()
	go func() {
		standaloneHandlers.WorkerOutbox.Start(ctxGF)
	}()
	srv := &http.Server{
		Addr:    config.Port,
		Handler: mux,
	}
	go func() {
		log.Info("server started on " + config.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server forced to shutdown",
				zap.Error(err),
			)
		}
	}()
	<-ctxGF.Done()
	log.Info("completion of work...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Info("server forced to shutdown", zap.Error(err))
	}
	log.Info("server stopped")
	if err := database.Close(); err != nil {
		log.Info("database forced to shutdown", zap.Error(err))
	}
	log.Info("database stopped")
	if err := kafka.Close(); err != nil {
		log.Info("kafka forced to shutdown", zap.Error(err))
	}
	log.Info("kafka stopped")
	if err := redis.Close(); err != nil {
		log.Info("redis forced to shutdown", zap.Error(err))
	}
	log.Info("redis stopped")
	log.Info("server api exited bb)")
}
