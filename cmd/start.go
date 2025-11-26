package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"sample_project/internal/app"
	"sample_project/internal/conf"
	"sample_project/internal/migrations"
	"sample_project/internal/repository"
	"sample_project/internal/service"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Запуск HTTP API и мониторинга сервисом",
	Long:  "Запускает веб-сервер на порту :7070 и начинает мониторинг всех сервисов из БД",
	Run: func(_ *cobra.Command, _ []string) {
		runServer()
	},
}

var application *app.App

func init() {
	rootCmd.AddCommand(startCmd)
}

func runServer() {
	config := conf.Load()

	repo, err := repository.NewPostgresRepository(config.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := migrations.Migrate(config.DatabaseURL); err != nil {
		repo.Close()
		log.Fatalf("Failed to run migrations: %v", err)
	}

	svc := service.NewService(repo)
	monitorService := service.NewMonitorService(repo)

	if err := monitorService.Start(); err != nil {
		repo.Close()
		log.Fatalf("Failed to start monitor service: %v", err)
	}

	application = app.New(repo, monitorService, svc)
	application.Start()

	waitForShutdown()

	application.Stop()
	if err := repo.Close(); err != nil {
		log.Printf("Error closing database connection: %v", err)
	}
}

func waitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown signal received...")
}
