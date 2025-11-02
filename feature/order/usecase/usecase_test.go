package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"ecommerce-go-api/domain/mock"
	"ecommerce-go-api/entity"
	"ecommerce-go-api/internal/errmap"
)

func TestCreateOrderFromCart_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup mocks
	mockOrderRepo := mock.NewMockOrderRepository(ctrl)
	mockShopRepo := mock.NewMockShopRepository(ctrl)
	mockProductRepo := mock.NewMockProductRepository(ctrl)
	mockUserRepo := mock.NewMockUserRepository(ctrl)

	uc := NewOrderUsecase(mockOrderRepo, mockShopRepo, mockProductRepo, mockUserRepo)

	// Test data
	ctx := context.Background()
	userID := uuid.New()
	addressID := 123
	shopID := uuid.New()
	productID := 1

	cart := &entity.Cart{
		ID:     1,
		UserID: userID,
	}

	cartItems := []*entity.CartItem{
		{
			ID:        1,
			CartID:    cart.ID,
			ProductID: productID,
			Qty:       2,
			Product: entity.Product{
				ID:       productID,
				Name:     "Test Product",
				Price:    100.0,
				StockQty: 10,
				ShopID:   shopID,
				IsActive: true,
			},
		},
	}

	address := &entity.Address{
		ID:          addressID,
		UserID:      userID,
		Name:        "John Doe",
		PhoneNumber: "0812345678",
		Line1:       "123 Main St",
		Line2:       "Apt 4B",
		Zipcode:     10110,
		SubDistrict: entity.SubDistrict{NameTH: "ปทุมวัน"},
		District:    entity.District{NameTH: "ปทุมวัน"},
		Province:    entity.Province{NameTH: "กรุงเทพมหานคร"},
	}

	shopCouriers := []*entity.ShopCourier{
		{
			ID:        1,
			ShopID:    shopID,
			CourierID: 1,
			Rate:      50.0,
		},
	}

	req := entity.CreateOrderRequest{
		AddressID:       addressID,
		PaymentMethodID: 1,
	}

	// Setup expectations
	mockOrderRepo.EXPECT().
		GetCartByUserID(ctx, userID).
		Return(cart, nil).
		Times(1)

	mockOrderRepo.EXPECT().
		ListCartItems(ctx, cart.ID).
		Return(cartItems, nil).
		Times(1)

	mockUserRepo.EXPECT().
		GetAddressByID(ctx, addressID).
		Return(address, nil).
		Times(1)

	mockShopRepo.EXPECT().
		ListShopCouriersByShopIDs(ctx, gomock.Any()).
		Return(shopCouriers, nil).
		Times(1)

	mockOrderRepo.EXPECT().
		CreateFullOrder(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), cart.ID, userID).
		DoAndReturn(func(ctx context.Context, order *entity.Order, shopOrders []*entity.ShopOrder, orderItemsByShop map[string][]*entity.OrderItem, payment *entity.Payment, cartID int, uid uuid.UUID) error {
			order.ID = uuid.New()
			for _, so := range shopOrders {
				so.ID = uuid.New()
				so.OrderID = order.ID
			}
			return nil
		}).
		Times(1)

	mockOrderRepo.EXPECT().
		GetOrderLogsByOrderID(ctx, gomock.Any()).
		Return([]*entity.OrderLog{}, nil).
		AnyTimes()

	mockOrderRepo.EXPECT().
		GetOrderLogsByShopOrderID(ctx, gomock.Any()).
		Return([]*entity.OrderLog{}, nil).
		AnyTimes()

	mockOrderRepo.EXPECT().
		GetOrderByID(ctx, gomock.Any()).
		DoAndReturn(func(ctx context.Context, id uuid.UUID) (*entity.Order, error) {
			return &entity.Order{
				ID:                  id,
				UserID:              userID,
				AddressID:           addressID,
				GrandTotal:          250.0, // 200 (2*100) + 50 (shipping)
				PaymentMethodID:     1,
				ShippingName:        "John Doe",
				ShippingPhone:       "0812345678",
				ShippingLine1:       "123 Main St",
				ShippingLine2:       "Apt 4B",
				ShippingSubDistrict: "ปทุมวัน",
				ShippingDistrict:    "ปทุมวัน",
				ShippingProvince:    "กรุงเทพมหานคร",
				ShippingZipcode:     "10110",
				ShopOrders:          []entity.ShopOrder{},
			}, nil
		}).
		Times(1)

	// Execute
	result, err := uc.CreateOrderFromCart(ctx, userID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 250.0, result.GrandTotal)
	assert.Equal(t, "John Doe", result.ShippingName)
	assert.Equal(t, "0812345678", result.ShippingPhone)
	assert.Equal(t, "ปทุมวัน", result.ShippingSubDistrict)
}

func TestCreateOrderFromCart_CartNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := mock.NewMockOrderRepository(ctrl)
	mockShopRepo := mock.NewMockShopRepository(ctrl)
	mockProductRepo := mock.NewMockProductRepository(ctrl)
	mockUserRepo := mock.NewMockUserRepository(ctrl)

	uc := NewOrderUsecase(mockOrderRepo, mockShopRepo, mockProductRepo, mockUserRepo)

	ctx := context.Background()
	userID := uuid.New()

	req := entity.CreateOrderRequest{
		AddressID:       123,
		PaymentMethodID: 1,
	}

	// Setup expectations
	mockOrderRepo.EXPECT().
		GetCartByUserID(ctx, userID).
		Return(nil, gorm.ErrRecordNotFound).
		Times(1)

	// Execute
	result, err := uc.CreateOrderFromCart(ctx, userID, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, errmap.ErrCartIsEmpty, err)
}

func TestCreateOrderFromCart_EmptyCart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := mock.NewMockOrderRepository(ctrl)
	mockShopRepo := mock.NewMockShopRepository(ctrl)
	mockProductRepo := mock.NewMockProductRepository(ctrl)
	mockUserRepo := mock.NewMockUserRepository(ctrl)

	uc := NewOrderUsecase(mockOrderRepo, mockShopRepo, mockProductRepo, mockUserRepo)

	ctx := context.Background()
	userID := uuid.New()

	cart := &entity.Cart{
		ID:     1,
		UserID: userID,
	}

	req := entity.CreateOrderRequest{
		AddressID:       123,
		PaymentMethodID: 1,
	}

	// Setup expectations
	mockOrderRepo.EXPECT().
		GetCartByUserID(ctx, userID).
		Return(cart, nil).
		Times(1)

	mockOrderRepo.EXPECT().
		ListCartItems(ctx, cart.ID).
		Return([]*entity.CartItem{}, nil). // Empty cart items
		Times(1)

	// Execute
	result, err := uc.CreateOrderFromCart(ctx, userID, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, errmap.ErrCartIsEmpty, err)
}

func TestCreateOrderFromCart_AddressIDRequired(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := mock.NewMockOrderRepository(ctrl)
	mockShopRepo := mock.NewMockShopRepository(ctrl)
	mockProductRepo := mock.NewMockProductRepository(ctrl)
	mockUserRepo := mock.NewMockUserRepository(ctrl)

	uc := NewOrderUsecase(mockOrderRepo, mockShopRepo, mockProductRepo, mockUserRepo)

	ctx := context.Background()
	userID := uuid.New()
	shopID := uuid.New()

	cart := &entity.Cart{
		ID:     1,
		UserID: userID,
	}

	cartItems := []*entity.CartItem{
		{
			ID:        1,
			CartID:    cart.ID,
			ProductID: 1,
			Qty:       2,
			Product: entity.Product{
				ID:       1,
				Name:     "Test Product",
				Price:    100.0,
				ShopID:   shopID,
				IsActive: true,
			},
		},
	}

	req := entity.CreateOrderRequest{
		AddressID:       0, // Missing address ID
		PaymentMethodID: 1,
	}

	// Setup expectations
	mockOrderRepo.EXPECT().
		GetCartByUserID(ctx, userID).
		Return(cart, nil).
		Times(1)

	mockOrderRepo.EXPECT().
		ListCartItems(ctx, cart.ID).
		Return(cartItems, nil).
		Times(1)

	// Execute
	result, err := uc.CreateOrderFromCart(ctx, userID, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, errmap.ErrAddressIDRequired, err)
}

func TestCreateOrderFromCart_NoShippingOptions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := mock.NewMockOrderRepository(ctrl)
	mockShopRepo := mock.NewMockShopRepository(ctrl)
	mockProductRepo := mock.NewMockProductRepository(ctrl)
	mockUserRepo := mock.NewMockUserRepository(ctrl)

	uc := NewOrderUsecase(mockOrderRepo, mockShopRepo, mockProductRepo, mockUserRepo)

	ctx := context.Background()
	userID := uuid.New()
	addressID := 123
	shopID := uuid.New()

	cart := &entity.Cart{
		ID:     1,
		UserID: userID,
	}

	cartItems := []*entity.CartItem{
		{
			ID:        1,
			CartID:    cart.ID,
			ProductID: 1,
			Qty:       2,
			Product: entity.Product{
				ID:       1,
				Name:     "Test Product",
				Price:    100.0,
				ShopID:   shopID,
				IsActive: true,
			},
		},
	}

	address := &entity.Address{
		ID:          addressID,
		UserID:      userID,
		Name:        "John Doe",
		PhoneNumber: "0812345678",
		Line1:       "123 Main St",
	}

	req := entity.CreateOrderRequest{
		AddressID:       addressID,
		PaymentMethodID: 1,
	}

	// Setup expectations
	mockOrderRepo.EXPECT().
		GetCartByUserID(ctx, userID).
		Return(cart, nil).
		Times(1)

	mockOrderRepo.EXPECT().
		ListCartItems(ctx, cart.ID).
		Return(cartItems, nil).
		Times(1)

	mockUserRepo.EXPECT().
		GetAddressByID(ctx, addressID).
		Return(address, nil).
		Times(1)

	mockShopRepo.EXPECT().
		ListShopCouriersByShopIDs(ctx, gomock.Any()).
		Return([]*entity.ShopCourier{}, nil). // No couriers available
		Times(1)

	// Execute
	result, err := uc.CreateOrderFromCart(ctx, userID, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, errmap.ErrNoShippingOptions, err)
}

func TestCreateOrderFromCart_CreateFullOrderError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := mock.NewMockOrderRepository(ctrl)
	mockShopRepo := mock.NewMockShopRepository(ctrl)
	mockProductRepo := mock.NewMockProductRepository(ctrl)
	mockUserRepo := mock.NewMockUserRepository(ctrl)

	uc := NewOrderUsecase(mockOrderRepo, mockShopRepo, mockProductRepo, mockUserRepo)

	ctx := context.Background()
	userID := uuid.New()
	addressID := 123
	shopID := uuid.New()

	cart := &entity.Cart{
		ID:     1,
		UserID: userID,
	}

	cartItems := []*entity.CartItem{
		{
			ID:        1,
			CartID:    cart.ID,
			ProductID: 1,
			Qty:       2,
			Product: entity.Product{
				ID:       1,
				Name:     "Test Product",
				Price:    100.0,
				ShopID:   shopID,
				IsActive: true,
			},
		},
	}

	address := &entity.Address{
		ID:      addressID,
		UserID:  userID,
		Name:    "John Doe",
		Zipcode: 10110,
	}

	shopCouriers := []*entity.ShopCourier{
		{
			ID:        1,
			ShopID:    shopID,
			CourierID: 1,
			Rate:      50.0,
		},
	}

	req := entity.CreateOrderRequest{
		AddressID:       addressID,
		PaymentMethodID: 1,
	}

	dbError := errors.New("database error")

	// Setup expectations
	mockOrderRepo.EXPECT().
		GetCartByUserID(ctx, userID).
		Return(cart, nil).
		Times(1)

	mockOrderRepo.EXPECT().
		ListCartItems(ctx, cart.ID).
		Return(cartItems, nil).
		Times(1)

	mockUserRepo.EXPECT().
		GetAddressByID(ctx, addressID).
		Return(address, nil).
		Times(1)

	mockShopRepo.EXPECT().
		ListShopCouriersByShopIDs(ctx, gomock.Any()).
		Return(shopCouriers, nil).
		Times(1)

	mockOrderRepo.EXPECT().
		CreateFullOrder(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), cart.ID, userID).
		Return(dbError).
		Times(1)

	// Execute
	result, err := uc.CreateOrderFromCart(ctx, userID, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, dbError, err)
}
