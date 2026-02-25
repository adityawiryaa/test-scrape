package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hibiken/asynq"

	"github.com/adityawiryaa/api/internal/config"
	delivery "github.com/adityawiryaa/api/internal/delivery/http/worker"
	"github.com/adityawiryaa/api/internal/repository/memory"
	workeruc "github.com/adityawiryaa/api/internal/usecases/worker"
	"github.com/adityawiryaa/api/pkg/cache"
	hitqueue "github.com/adityawiryaa/api/pkg/hit/queue"
)

func main() {
	cfg := config.LoadWorkerConfig()

	log.Printf("[init] loading config: port=%s redis=%s", cfg.Port, cfg.Redis.Addr())

	store := memory.NewConfigStore()
	executor := workeruc.NewHTTPExecutor(cfg.RequestTimeout)

	log.Printf("[init] connecting to redis: addr=%s db=%d asynq_db=%d", cfg.Redis.Addr(), cfg.Redis.DB, cfg.Redis.AsynqDB)
	rdb := cache.NewRedisClient(cfg.Redis.Addr(), cfg.Redis.DB)

	queueClient := hitqueue.NewClient(cfg.Redis.Addr(), cfg.Redis.AsynqDB)
	resultStore := hitqueue.NewResultStore(rdb)

	commandUC := workeruc.NewCommandUsecase(executor, store, queueClient)
	queryUC := workeruc.NewQueryUsecase(store, resultStore)

	handler := delivery.NewHandler(commandUC, queryUC)
	router := delivery.SetupRouter(handler, cfg.APIKey)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	processor := hitqueue.NewProcessor(executor, resultStore)
	asynqSrv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: cfg.Redis.Addr(), DB: cfg.Redis.AsynqDB},
		asynq.Config{Concurrency: 10},
	)

	mux := asynq.NewServeMux()
	processor.RegisterHandlers(mux)

	go func() {
		log.Printf("[http] server starting on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[http] server error: %v", err)
		}
	}()

	go func() {
		log.Printf("[asynq] worker starting: concurrency=10 redis=%s db=%d", cfg.Redis.Addr(), cfg.Redis.AsynqDB)
		if err := asynqSrv.Run(mux); err != nil {
			log.Fatalf("[asynq] server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[shutdown] signal received, shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("[shutdown] stopping http server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("[shutdown] http server error: %v", err)
	}

	log.Println("[shutdown] stopping asynq worker...")
	asynqSrv.Shutdown()

	log.Println("[shutdown] closing redis connections...")
	_ = queueClient.Close()
	_ = rdb.Close()

	log.Println("[shutdown] worker exited")
}
