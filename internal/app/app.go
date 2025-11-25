package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sample_project/internal/handler"
	"sample_project/internal/middleware"
	"sample_project/internal/repository"
	"sample_project/internal/service"
	"syscall"
	"time"

	_ "sample_project/internal/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type App struct {
	router  *gin.Engine
	server  *http.Server
	repo    repository.Repository
	monitor *service.MonitorService
	service *service.Service
}

func New(repo repository.Repository, monitor *service.MonitorService, svc *service.Service) *App {
	router := gin.Default()

	h := handler.NewHandler(svc, monitor)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/health", h.HealthCheck)

	api := router.Group("/api")
	{
		api.POST("/sign_in", h.SignIn)

		services := api.Group("/services")
		{
			services.POST("", middleware.CheckAuth(), h.CreateService)
			services.GET("", h.GetServices)
			services.GET("/status", h.GetServicesStatus)
			services.DELETE("/:id", middleware.CheckAuth(), h.DeleteService)
			services.GET("/:id/results", h.GetServiceResults)
			services.GET("/:id/results/filter", h.GetServiceResultsByStatus)
		}
	}

	return &App{
		router: router,
		server: &http.Server{
			Addr:    ":7070",
			Handler: router,
		},
		repo:    repo,
		monitor: monitor,
		service: svc,
	}
}

func (a *App) Start() {
	go func() {
		log.Printf("Server starting on %s", ":7070")
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()
}

func (a *App) Stop() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	a.monitor.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
