package delivery

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"ecommerce-go-api/feature/order/repository"
	refundRepo "ecommerce-go-api/feature/refund/repository"
	"ecommerce-go-api/feature/refund/usecase"
	shopRepo "ecommerce-go-api/feature/shop/repository"
	"ecommerce-go-api/middleware"
)

func RegisterRoutes(group *echo.Group, handler *RefundHandler) {
	shopRefunds := group.Group("/shop/refunds", middleware.JWTAuth(), middleware.ShopOwnerOnly())
	shopRefunds.POST("", handler.CreateRefund)
	shopRefunds.PUT("/:refundId/approve", handler.ApproveRefund)

	userRefunds := group.Group("/refunds", middleware.JWTAuth(), middleware.UserOnly())
	userRefunds.PUT("/:refundId/bank-account", handler.SubmitRefundBankAccount)
}

func RegisterRefundHandler(group *echo.Group, db *gorm.DB) {
	refundRepository := refundRepo.NewRefundRepository(db)
	orderRepository := repository.NewOrderRepository(db)
	shopRepository := shopRepo.NewShopRepository(db)
	refundUsecase := usecase.NewRefundUsecase(refundRepository, orderRepository, shopRepository)
	handler := NewRefundHandler(refundUsecase)
	RegisterRoutes(group, handler)
}
