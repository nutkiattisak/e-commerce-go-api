package main

import (
	"context"
	"ecommerce-go-api/config"
	resp "ecommerce-go-api/internal/response"
	"ecommerce-go-api/utils"
	"log"
	"net/http"
	"os"
	"time"

	intvalidator "ecommerce-go-api/internal/validator"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "ecommerce-go-api/docs"
	authDelivery "ecommerce-go-api/feature/auth/delivery"
	cartDelivery "ecommerce-go-api/feature/cart/delivery"
	courierDelivery "ecommerce-go-api/feature/courier/delivery"
	locationDelivery "ecommerce-go-api/feature/location/delivery"
	orderDelivery "ecommerce-go-api/feature/order/delivery"
	productDelivery "ecommerce-go-api/feature/product/delivery"
	refundDelivery "ecommerce-go-api/feature/refund/delivery"
	shopDelivery "ecommerce-go-api/feature/shop/delivery"
	userDelivery "ecommerce-go-api/feature/user/delivery"

	orderRepo "ecommerce-go-api/feature/order/repository"
	productRepo "ecommerce-go-api/feature/product/repository"
	"ecommerce-go-api/internal/cron"

	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title						E-commerce API
// @version					1.0.0
// @description				This is E-commerce API documentation.
// @BasePath					/
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description				Type "Bearer" followed by a space and JWT token.
func init() {
	config.InitialENV()
	config.ConnectDatabase()
}

func main() {
	e := echo.New()

	e.Validator = intvalidator.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/", func(c echo.Context) error {
		return resp.Success(c, http.StatusOK, "Welcome to E-commerce API", map[string]string{
			"status":  "running",
			"version": "v1.0.0",
		})
	})

	appEnv := os.Getenv("APP_ENV")

	if appEnv != "production" {
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	e.GET("/health", func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
		defer cancel()

		db := config.DB
		sqlDB, err := db.DB()
		if err != nil {
			return resp.Error(c, http.StatusInternalServerError, "service unhealthy (db error)")
		}

		if err := sqlDB.PingContext(ctx); err != nil {
			return resp.Error(c, http.StatusInternalServerError, "service unhealthy (db ping failed)")
		}

		return resp.Success(c, http.StatusOK, "healthy", map[string]interface{}{
			"database": "connected",
		})
	})

	db := config.DB

	oRepo := orderRepo.NewOrderRepository(db)
	pRepo := productRepo.NewProductRepository(db)
	scheduler, err := cron.NewScheduler(oRepo, pRepo)
	if err != nil {
		log.Fatalf("Failed to create scheduler: %v", err)
	}
	if err := scheduler.Start(); err != nil {
		log.Fatalf("Failed to start scheduler: %v", err)
	}
	defer scheduler.Stop()

	api := e.Group("/api")
	{
		authDelivery.RegisterAuthHandler(api, db)
		userDelivery.RegisterUserHandler(api, db)
		locationDelivery.RegisterLocationHandler(api, db)
		productDelivery.RegisterProductHandler(api, db)
		shopDelivery.RegisterShopHandler(api, db)
		cartDelivery.RegisterCartHandler(api, db)
		orderDelivery.RegisterOrderHandler(api, db)
		courierDelivery.RegisterCourierHandler(api, db)
		refundDelivery.RegisterRefundHandler(api, db)
	}

	utils.ServeGracefulShutdown(e)
}
