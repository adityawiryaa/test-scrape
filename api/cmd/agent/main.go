package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/adityawiryaa/api/domain/entity"
	"github.com/adityawiryaa/api/internal/config"
	"github.com/adityawiryaa/api/internal/repository/memory"
	agentuc "github.com/adityawiryaa/api/internal/usecases/agent"
	"github.com/adityawiryaa/api/pkg/backoff"
	controllerclient "github.com/adityawiryaa/api/pkg/controller"
	workerclient "github.com/adityawiryaa/api/pkg/worker"
)

func main() {
	cfg := config.LoadAgentConfig()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	controllerClient := controllerclient.NewClient(cfg.ControllerURL, cfg.APIKey, cfg.RequestTimeout)
	workerClient := workerclient.NewClient(cfg.WorkerURL, cfg.RequestTimeout)

	store := memory.NewConfigStore()

	commandUC := agentuc.NewCommandUsecase(controllerClient, workerClient, store, backoff.DefaultConfig())
	queryUC := agentuc.NewQueryUsecase(controllerClient, store)

	resp, err := commandUC.RegisterWithController(ctx, &entity.RegistrationRequest{
		Hostname:  cfg.Hostname,
		IPAddress: cfg.IPAddress,
		Port:      cfg.Port,
	})
	if err != nil {
		log.Fatalf("failed to register: %v", err)
	}
	log.Printf("registered as agent %s", resp.AgentID)

	log.Printf("starting config polling (interval: %s)", cfg.PollInterval)
	queryUC.StartPolling(ctx, cfg.PollInterval, commandUC.ForwardConfigToWorker)

	log.Println("agent stopped")
	os.Exit(0)
}
