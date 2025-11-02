package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/kubilayrn/ChronoGo/internal/database"
	"github.com/kubilayrn/ChronoGo/internal/handler"
	"github.com/kubilayrn/ChronoGo/internal/queue"
	"github.com/kubilayrn/ChronoGo/internal/repository"
	"github.com/kubilayrn/ChronoGo/internal/sender"

	_ "github.com/kubilayrn/ChronoGo/docs"
)

// @title           ChronoGo API
// @version         1.0
// @description     Automatic message sending system API
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@chronogo.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @schemes http https
func main() {

	ctx := context.Background()
	dbConfig := database.LoadConfigFromEnv()
	if err := database.Connect(ctx, dbConfig); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()
	log.Println("Database connection established")

	messageRepo := repository.NewMessageRepository()
	webhookSender := sender.NewWebhookSender()
	scheduler := queue.NewScheduler(messageRepo, webhookSender)
	h := handler.NewHandler(messageRepo, scheduler)

	if err := scheduler.Start(); err != nil {
		log.Printf("Failed to start scheduler automatically: %v", err)
	} else {
		log.Println("Scheduler started automatically on deployment")
	}

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		ctx := context.Background()
		if err := database.DB.Ping(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "error",
				"message": "Database connection failed",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Database connection is healthy",
		})
	})

	api := r.Group("/api")
	{
		api.GET("/messages/sent", h.ListSentMessages)
		api.POST("/scheduler/toggle", h.ToggleScheduler)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Println("Server started on :8080")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
