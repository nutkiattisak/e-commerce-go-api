package delivery

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"ecommerce-go-api/feature/courier/repository"
	"ecommerce-go-api/feature/courier/usecase"
	"ecommerce-go-api/middleware"
)

func RegisterCourierHandler(group *echo.Group, db *gorm.DB) {
	repo := repository.NewCourierRepository(db)
	uc := usecase.NewCourierUsecase(repo)
	handler := NewCourierHandler(uc)

	couriers := group.Group("/couriers")
	couriers.Use(middleware.JWTAuth(), middleware.ShopOwnerOnly())

	couriers.GET("", handler.ListCouriers)
}
