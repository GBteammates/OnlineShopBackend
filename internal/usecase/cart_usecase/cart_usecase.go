package cart_usecase

import (
	"OnlineShopBackend/internal/models"
	"OnlineShopBackend/internal/usecase"
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var _ usecase.ICartUsecase = &CartUseCase{}

type CartUseCase struct {
	store  usecase.CartStore
	logger *zap.Logger
}

func NewCartUseCase(store usecase.CartStore, logger *zap.Logger) usecase.ICartUsecase {
	logger.Debug("Enter in usecase NewCartUseCase()")
	cart := &CartUseCase{store: store, logger: logger}
	return cart
}

// GetCart creates request in db and returns cart or error
func (c *CartUseCase) GetCart(ctx context.Context, cartId uuid.UUID) (*models.Cart, error) {
	c.logger.Sugar().Debugf("Enter in usecase GetCart() with args: ctx, cartId: %v", cartId)
	cart, err := c.store.GetCart(ctx, cartId)
	if err != nil {
		return nil, err
	}
	return cart, nil
}

// GetCartByUserId creates request in db and returns cart or error
func (c *CartUseCase) GetCartByUserId(ctx context.Context, userId uuid.UUID) (*models.Cart, error) {
	c.logger.Sugar().Debugf("Enter in usecase GetCart() with args: ctx, userId: %v", userId)
	cart, err := c.store.GetCartByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	return cart, nil
}

// DeleteItemFromCart delete item from cart
func (c *CartUseCase) DeleteItemFromCart(ctx context.Context, cartId uuid.UUID, itemId uuid.UUID) error {
	c.logger.Sugar().Debugf("Enter in usecase DeleteItemFromCart() with args: ctx, cartId: %v, itemId: %v", cartId, itemId)
	err := c.store.DeleteItemFromCart(ctx, cartId, itemId)
	if err != nil {
		return err
	}
	return nil
}

// Create create new cart
func (c *CartUseCase) Create(ctx context.Context, userId uuid.UUID) (uuid.UUID, error) {
	c.logger.Sugar().Debugf("Enter in usecase cart Create() with args: ctx, userId: %v", userId)
	cartId, err := c.store.Create(ctx, userId)
	if err != nil {
		return cartId, err
	}
	return cartId, nil
}

// AddItemToCart add item to cart
func (c *CartUseCase) AddItemToCart(ctx context.Context, cartId uuid.UUID, itemId uuid.UUID) error {
	c.logger.Sugar().Debugf("Enter in usecase AddItemToCart() with args: ctx, cartId: %v, itemId: %v", cartId, itemId)
	err := c.store.AddItemToCart(ctx, cartId, itemId)
	if err != nil {
		return err
	}
	return nil
}

// DeleteCart delete cart from db
func (c *CartUseCase) DeleteCart(ctx context.Context, cartId uuid.UUID) error {
	c.logger.Sugar().Debugf("Enter in usecase DeleteCart() with args: ctx, cartId: %v", cartId)
	err := c.store.DeleteCart(ctx, cartId)
	if err != nil {
		return err
	}
	return nil
}
