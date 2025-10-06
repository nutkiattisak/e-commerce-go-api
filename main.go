package main

import (
	"ecommerce-go-api/config"
	"ecommerce-go-api/entity"
	"ecommerce-go-api/utils"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func init() {
	config.InitialENV()
	config.ConnectDatabase()
	
	// Run database migration
	// migration.AutoMigrate()
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Welcome to E-commerce API",
			"status":  "running",
		})
	})

	e.GET("/health", func(c echo.Context) error {
		started := time.Now()
		duration := time.Since(started)

		if duration.Seconds() > 10 {
			return c.JSON(http.StatusInternalServerError, entity.ResponseError{})
		}

		return c.String(http.StatusOK, http.StatusText(http.StatusOK))
	})

	log.Printf("Server starting on port %s", config.HTTP_PORT)
	log.Fatal(e.Start(":" + config.HTTP_PORT))

	utils.ServeGracefulShutdown(e)
}