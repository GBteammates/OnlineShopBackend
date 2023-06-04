package usecase

import (
	"OnlineShopBackend/internal/models"
	usecase "OnlineShopBackend/internal/usecase/interfaces"
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var _ usecase.ICartUsecase = (*cartUsecase)(nil)

type cartUsecase struct {
	store  usecase.CartStore
	logger *zap.SugaredLogger
}

func NewCartUsecase(store usecase.CartStore, logger *zap.SugaredLogger) *cartUsecase {
	logger.Debug("Enter in usecase NewCartUsecase()")
	cart := &cartUsecase{store: store, logger: logger}
	return cart
}

// GetCart creates request in db and returns cart or error
func (c *cartUsecase) Get(ctx context.Context, cartId uuid.UUID) (*models.Cart, error) {
	c.logger.Debugf("Enter in usecase Get() with args: ctx, cartId: %v", cartId)
	cart, err := c.store.Get(ctx, cartId)
	if err != nil {
		return nil, err
	}
	return cart, nil
}

// GetCartByUserId creates request in db and returns cart or error
func (c *cartUsecase) CartByUserId(ctx context.Context, userId uuid.UUID) (*models.Cart, error) {
	c.logger.Debugf("Enter in usecase CartByUserId() with args: ctx, userId: %v", userId)
	cart, err := c.store.CartByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	return cart, nil
}

// DeleteItemFromCart delete item from cart
func (c *cartUsecase) DeleteItem(ctx context.Context, cartId uuid.UUID, itemId uuid.UUID) error {
	c.logger.Debugf("Enter in usecase DeleteItem() with args: ctx, cartId: %v, itemId: %v", cartId, itemId)
	err := c.store.DeleteItem(ctx, cartId, itemId)
	if err != nil {
		return err
	}
	return nil
}

// Create create new cart
func (c *cartUsecase) Create(ctx context.Context, userId uuid.UUID) (uuid.UUID, error) {
	c.logger.Debugf("Enter in usecase cart Create() with args: ctx, userId: %v", userId)
	cartId, err := c.store.Create(ctx, userId)
	if err != nil {
		return cartId, err
	}
	return cartId, nil
}

// AddItem add item to cart
func (c *cartUsecase) AddItem(ctx context.Context, cartId uuid.UUID, itemId uuid.UUID) error {
	c.logger.Debugf("Enter in usecase AddItem() with args: ctx, cartId: %v, itemId: %v", cartId, itemId)
	err := c.store.AddItem(ctx, cartId, itemId)
	if err != nil {
		return err
	}
	return nil
}

// DeleteCart delete cart from db
func (c *cartUsecase) Delete(ctx context.Context, cartId uuid.UUID) error {
	c.logger.Debugf("Enter in usecase Delete() with args: ctx, cartId: %v", cartId)
	err := c.store.Delete(ctx, cartId)
	if err != nil {
		return err
	}
	return nil
}
