package utils

import (
	"context"
	"ecommerce-go-api/config"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func ServeGracefulShutdown(e *echo.Echo) {
	go func() {
		port := ":" + config.HTTP_PORT
		log.Infof("Starting server on port %s", port)

		if err := e.Start(port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	log.Info("Server started. Press Ctrl+C to shutdown.")
	<-ctx.Done()

	log.Info("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Info("Server shutdown complete.")
}