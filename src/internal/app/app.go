package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	Item "sample_project/internal/handler"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type App struct {
	router *gin.Engine
	server *http.Server
}

func New() *App {
	router := gin.Default()

	api := router.Group("/api")

	service := api.Group("/service")
	{
		service.POST("/item", Item.CreateService)
		service.GET("/items", Item.GetService)
		service.GET("/item/:id", Item.SearchServiceId)
		service.PUT("/item/:id", Item.ChangeService)
		service.DELETE("/item/:id", Item.DeleteService)
	}

	result := api.Group("/result")
	{
		result.POST("/item", Item.CreateResult)
		result.GET("/items", Item.GetResult)
		result.GET("/item/:id", Item.SearchResultId)
		result.PUT("/item/:id", Item.ChangeResult)
		result.DELETE("/item/:id", Item.DeleteResult)
	}

	return &App{
		router: router,
		server: &http.Server{
			Addr:    ":7070",
			Handler: router,
		},
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
