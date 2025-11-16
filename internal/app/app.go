package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	Item "sample_project/internal/handler"
	MW "sample_project/internal/middleware"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	_ "sample_project/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type App struct {
	router *gin.Engine
	server *http.Server
}

func New() *App {
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")

	//api.POST("/create_user", Item.New().Register())
	api.POST("/sign_in", Item.SignIn())

	service := api.Group("/service")
	{
		service.POST("/item", MW.CheckAuth(), Item.CreateService)
		service.GET("/items", Item.GetService)
		service.GET("/item/:id", Item.SearchServiceId)
		service.PUT("/item/:id", MW.CheckAuth(), Item.ChangeService)
		service.DELETE("/item/:id", MW.CheckAuth(), Item.DeleteService)
	}

	result := api.Group("/result")
	{
		result.POST("/item", MW.CheckAuth(), Item.CreateResult)
		result.GET("/items", Item.GetResult)
		result.GET("/item/:id", Item.SearchResultId)
		result.PUT("/item/:id", MW.CheckAuth(), Item.ChangeResult)
		result.DELETE("/item/:id", MW.CheckAuth(), Item.DeleteResult)
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
