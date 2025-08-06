package main

import (
	"context"
	"fmt"
	"grupos-cmd/internal/config"
	"grupos-cmd/internal/event"
	"grupos-cmd/internal/handler"
	"grupos-cmd/internal/repository"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "grupos-cmd/docs" // This is for swagger docs

	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Grupos Command API
// @version 1.0
// @description This is the command API for the Grupos service.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
func main() {
	// --- Configuration ---
	cfg := config.New()

	// --- Logger ---
	logger := log.New(os.Stdout, "", log.LstdFlags)
	logger.Println("Starting service...")

	// --- Context for initialization ---
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// --- Repository ---
	repo, err := repository.NewFirestoreRepository(ctx, cfg)
	if err != nil {
		logger.Fatalf("Failed to create repository: %v", err)
	}
	logger.Println("Firestore repository connected.")

	// --- Publisher ---
	publisher, err := event.NewWatermillPublisher(cfg)
	if err != nil {
		logger.Fatalf("Failed to create publisher: %v", err)
	}
	defer publisher.Close()
	logger.Println("Watermill publisher connected.")

	// --- Echo ---
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// --- Handler ---
	h := handler.NewHandler(repo, publisher)
	h.RegisterRoutes(e)

	// --- Swagger ---
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// --- Start Server ---
	go func() {
		if err := e.Start(":" + cfg.Port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()
	logger.Printf("Server started on port %s", cfg.Port)

	// --- Graceful Shutdown ---
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Println("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		e.Logger.Fatal(err)
	}

	logger.Println("Server gracefully stopped.")
}
