package main

import (
	"ecommerce-go-api/config"
	"ecommerce-go-api/entity"
	"ecommerce-go-api/utils"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "ecommerce-go-api/docs"
	authDelivery "ecommerce-go-api/feature/auth/delivery"
	locationDelivery "ecommerce-go-api/feature/location/delivery"
	userDelivery "ecommerce-go-api/feature/user/delivery"

	echoSwagger "github.com/swaggo/echo-swagger"
)

func init() {
	config.InitialENV()
	config.ConnectDatabase()
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
			"version": "v1.0.0",
		})
	})

	e.GET("/health", func(c echo.Context) error {
		started := time.Now()
		duration := time.Since(started)

		if duration.Seconds() > 10 {
			return c.JSON(http.StatusInternalServerError, entity.ResponseError{})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":   "healthy",
			"duration": duration.String(),
			"database": "connected",
		})
	})

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	db := config.DB
	api := e.Group("/api")
	{
		authDelivery.RegisterAuthHandler(api, db)
		userDelivery.RegisterUserHandler(api, db)
		locationDelivery.RegisterLocationHandler(api, db)
	}

	utils.ServeGracefulShutdown(e)
}
