package cmd

import (
	"log"
	"sample_project/internal/app"
	"sample_project/internal/conf"
	"sample_project/internal/migrations"
	"sample_project/internal/repository"
	"sample_project/internal/service"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Запуск HTTP API и мониторинга сервисом",
	Long:  "Запускает веб-сервер на порту :7070 и начинает мониторинг всех сервисов из БД",
	Run: func(cmd *cobra.Command, args []string) {
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
	defer repo.Close()

	if err := migrations.Migrate(config.DatabaseURL); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	svc := service.NewService(repo)
	monitorService := service.NewMonitorService(repo)

	if err := monitorService.Start(); err != nil {
		log.Fatalf("Failed to start monitor service: %v", err)
	}

	application = app.New(repo, monitorService, svc)

	application.Start()
	defer application.Stop()

}
