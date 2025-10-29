package main

import (
	"ecommerce-go-api/config"
	resp "ecommerce-go-api/internal/response"
	"ecommerce-go-api/utils"
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
	shopDelivery "ecommerce-go-api/feature/shop/delivery"
	userDelivery "ecommerce-go-api/feature/user/delivery"

	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title						E-commerce API
// @version						1.0.0
// @description					This is E-commerce API documentation.
//
// @host						localhost:8080
// @BasePath					/
//
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description					Type "Bearer" followed by a space and JWT token.
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
		started := time.Now()
		duration := time.Since(started)

		if duration.Seconds() > 10 {
			return resp.Error(c, http.StatusInternalServerError, "service unhealthy")
		}

		return resp.Success(c, http.StatusOK, "healthy", map[string]interface{}{
			"duration": duration.String(),
			"database": "connected",
		})
	})

	db := config.DB

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
	}

	utils.ServeGracefulShutdown(e)
}
