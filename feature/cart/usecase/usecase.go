package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"gorm.io/gorm"

	"ecommerce-go-api/domain"
	"ecommerce-go-api/entity"
	"ecommerce-go-api/internal/errmap"
)

type cartUsecase struct {
	cartRepo    domain.CartRepository
	productRepo domain.ProductRepository
	shopRepo    domain.ShopRepository
	validate    *validator.Validate
}

func NewCartUsecase(cartRepo domain.CartRepository, productRepo domain.ProductRepository, shopRepo domain.ShopRepository) domain.CartUsecase {
	return &cartUsecase{
		cartRepo:    cartRepo,
		productRepo: productRepo,
		shopRepo:    shopRepo,
		validate:    validator.New(),
	}
}

func (u *cartUsecase) EstimateShipping(ctx context.Context, userID uuid.UUID, cartItemIDs []int) (*entity.CartShippingEstimateResponse, error) {
	cartItems, err := u.cartRepo.GetCartItemsByIDs(ctx, cartItemIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart items: %w", err)
	}

	shopMap := make(map[string][]entity.CartItemShop)
	shopSubtotals := make(map[string]float64)
	shopUUIDs := make([]uuid.UUID, 0)
	shopSeen := make(map[string]bool)

	for _, it := range cartItems {
		shopID := it.Product.ShopID.String()
		if !shopSeen[shopID] {
			shopSeen[shopID] = true
			shopUUIDs = append(shopUUIDs, it.Product.ShopID)
		}
		cartItemShop := entity.CartItemShop{
			CartItemID: it.ID,
			ProductID:  it.ProductID,
			Qty:        it.Qty,
			UnitPrice:  it.Product.Price,
			Subtotal:   float64(it.Qty) * it.Product.Price,
		}
		shopMap[shopID] = append(shopMap[shopID], cartItemShop)
		shopSubtotals[shopID] += float64(it.Qty) * it.Product.Price
	}

	shopCouriers, err := u.shopRepo.ListShopCouriersByShopIDs(ctx, shopUUIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get shop couriers: %w", err)
	}

	scMap := make(map[string][]*entity.ShopCourier)
	for _, shopCourier := range shopCouriers {
		scMap[shopCourier.ShopID.String()] = append(scMap[shopCourier.ShopID.String()], shopCourier)
	}

	var resp entity.CartShippingEstimateResponse
	var grandTotal float64

	for shopID, items := range shopMap {
		subtotal := shopSubtotals[shopID]

		var courierOpt entity.CourierOption
		bestPrice := 0.0
		if scsForShop, ok := scMap[shopID]; ok && len(scsForShop) > 0 {
			for i, sc := range scsForShop {
				price := sc.Rate
				var courierName string
				if sc.Courier != nil {
					courierName = sc.Courier.Name
				}
				if i == 0 || price < bestPrice {
					bestPrice = price
					courierOpt = entity.CourierOption{
						CourierID: sc.CourierID,
						Name:      courierName,
						Price:     price,
					}
				}
			}
		} else {
			return nil, errmap.ErrNoShippingOptions
		}

		shopEstimate := entity.CartShopEstimate{
			ShopID:   shopID,
			Items:    items,
			Subtotal: subtotal,
			Courier:  courierOpt,
		}
		resp.Shop = append(resp.Shop, shopEstimate)
		grandTotal += subtotal + courierOpt.Price
	}

	resp.GrandTotal = grandTotal
	return &resp, nil
}

func (u *cartUsecase) AddItem(ctx context.Context, userID uuid.UUID, productID int, qty int) (*entity.CartItem, bool, error) {
	if qty <= 0 {
		return nil, false, errmap.ErrQuantityMustBeGreaterThanZero
	}

	product, err := u.productRepo.GetProductByID(ctx, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, errmap.ErrProductNotFound
		}
		return nil, false, fmt.Errorf("product lookup error: %w", err)
	}
	if !product.IsActive {
		return nil, false, errmap.ErrProductInactive
	}

	cart, err := u.cartRepo.EnsureCartForUser(ctx, userID)
	if err != nil {
		return nil, false, fmt.Errorf("failed to ensure cart: %w", err)
	}

	existing, _ := u.cartRepo.GetCartItemByUserAndProduct(ctx, userID, productID)
	existingQty := 0
	if existing != nil {
		existingQty = existing.Qty
	}

	if existingQty+qty > product.StockQty {
		return nil, false, errmap.ErrInsufficientStock
	}

	now := time.Now()
	item := &entity.CartItem{
		CartID:    cart.ID,
		ProductID: productID,
		Qty:       qty,
		CreatedAt: now,
		UpdatedAt: now,
	}

	res, created, err := u.cartRepo.UpsertCartItem(ctx, item)
	if err != nil {
		return nil, false, fmt.Errorf("failed to upsert cart item: %w", err)
	}
	return res, created, nil
}

func (u *cartUsecase) GetCart(ctx context.Context, userID uuid.UUID) (*entity.Cart, []*entity.CartItem, *entity.CartSummary, error) {
	cart, err := u.cartRepo.GetCartByUserID(ctx, userID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get cart: %w", err)
	}
	items, err := u.cartRepo.ListCartItems(ctx, cart.ID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to list cart items: %w", err)
	}

	var totalItems int
	var totalQty int
	var subtotal float64
	for _, it := range items {
		totalItems++
		totalQty += it.Qty
		price := it.Product.Price
		subtotal += float64(it.Qty) * price
	}

	summary := &entity.CartSummary{
		TotalItems: totalItems,
		TotalQty:   totalQty,
		Subtotal:   subtotal,
	}
	return cart, items, summary, nil
}

func (u *cartUsecase) UpdateItem(ctx context.Context, userID uuid.UUID, itemID int, qty int) (*entity.CartItem, error) {
	if qty <= 0 {
		return nil, errmap.ErrQuantityMustBeGreaterThanZero
	}

	ci, err := u.cartRepo.GetCartItemByID(ctx, itemID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errmap.ErrCartItemNotFound
		}
		return nil, fmt.Errorf("failed to get cart item: %w", err)
	}

	product, err := u.productRepo.GetProductByID(ctx, ci.ProductID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errmap.ErrProductNotFound
		}
		return nil, fmt.Errorf("product lookup error: %w", err)
	}
	if qty > product.StockQty {
		return nil, errmap.ErrInsufficientStock
	}
	cart, err := u.cartRepo.GetCartByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user's cart: %w", err)
	}
	if ci.CartID != cart.ID {
		return nil, errmap.ErrUnauthorized
	}

	ci.Qty = qty
	ci.UpdatedAt = time.Now()

	if err := u.cartRepo.UpdateCartItem(ctx, ci); err != nil {
		return nil, fmt.Errorf("failed to update cart item: %w", err)
	}

	return ci, nil
}

func (u *cartUsecase) DeleteItem(ctx context.Context, userID uuid.UUID, itemID int) error {
	ci, err := u.cartRepo.GetCartItemByID(ctx, itemID)
	if err != nil {
		return fmt.Errorf("failed to get cart item: %w", err)
	}

	cart, err := u.cartRepo.GetCartByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user's cart: %w", err)
	}
	if ci.CartID != cart.ID {
		return errmap.ErrUnauthorized
	}

	if err := u.cartRepo.DeleteCartItem(ctx, itemID); err != nil {
		return fmt.Errorf("failed to delete cart item: %w", err)
	}

	return nil
}
