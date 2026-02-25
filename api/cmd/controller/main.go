package main

import (
	"log"
	"net/http"
	"time"

	"github.com/adityawiryaa/api/internal/config"
	delivery "github.com/adityawiryaa/api/internal/delivery/http/controller"
	"github.com/adityawiryaa/api/internal/repository"
	"github.com/adityawiryaa/api/internal/repository/commands"
	"github.com/adityawiryaa/api/internal/repository/queries"
	controlleruc "github.com/adityawiryaa/api/internal/usecases/controller"
	"github.com/adityawiryaa/api/pkg/shutdown"
)

func main() {
	cfg := config.LoadControllerConfig()

	db, err := config.NewDB(cfg.DBPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	if err := repository.Migrate(db); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	agentCmd := commands.NewAgentCommand(db)
	configCmd := commands.NewConfigCommand(db)
	configQuery := queries.NewConfigQuery(db)

	commandUC := controlleruc.NewCommandUsecase(agentCmd, configCmd, configQuery)
	queryUC := controlleruc.NewQueryUsecase(configQuery)

	handler := delivery.NewHandler(commandUC, queryUC)
	router := delivery.SetupRouter(handler, cfg.APIKey)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		log.Printf("controller starting on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	shutdown.GracefulShutdown(srv, 5*time.Second)
}
