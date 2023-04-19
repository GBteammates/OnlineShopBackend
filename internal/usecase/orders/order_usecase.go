package usecase

import (
	"OnlineShopBackend/internal/models"
	usecase "OnlineShopBackend/internal/usecase/interfaces"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type order struct {
	orderStore usecase.OrderStore
	logger     *zap.SugaredLogger
}

var _ usecase.IOrderUsecase = (*order)(nil)

func NewOrderUsecase(orderStore usecase.OrderStore, logger *zap.SugaredLogger) usecase.IOrderUsecase {
	return &order{
		orderStore: orderStore,
		logger:     logger,
	}
}

func (o *order) PlaceOrder(ctx context.Context, cart *models.Cart, user models.User, address models.UserAddress) (*models.Order, error) {
	select {
	case <-ctx.Done():
		o.logger.Error("context closed")
		return nil, fmt.Errorf("context closed")
	default:
		ordr := models.Order{
			User:         user,
			Address:      address,
			Status:       models.StatusCreated,
			CreatedAt:    time.Now(),
			ShipmentTime: time.Now().Add(models.ProlongedShipmentPeriod),
			Items:        append([]models.ItemWithQuantity{}[:0:0], cart.Items...),
		}
		res, err := o.orderStore.CreateOrder(ctx, &ordr)
		if err != nil {
			o.logger.Errorf("can't add order to db %s", err)
			return nil, fmt.Errorf("can't place order to db : %w", err)
		}
		o.logger.Debugf("order %s created", res.ID.String())
		return res, nil
	}
}

func (o *order) ChangeStatus(ctx context.Context, order *models.Order, newStatus models.Status) error {
	select {
	case <-ctx.Done():
		o.logger.Error("context closed")
		return fmt.Errorf("context closed")
	default:
		if newStatus == order.Status {
			return nil
		}
		if err := o.orderStore.ChangeStatus(ctx, order, newStatus); err != nil {
			o.logger.Errorf("can't change status of order: %s", err)
			return fmt.Errorf("can't change status of order: %w", err)
		}
	}
	return nil
}
func (o *order) GetOrdersForUser(ctx context.Context, user *models.User) ([]models.Order, error) {
	select {
	case <-ctx.Done():
		o.logger.Error("context closed")
		return nil, fmt.Errorf("context closed")
	default:
		result := make([]models.Order, 0, 10)
		resChan, err := o.orderStore.GetOrdersForUser(ctx, user)
		if err != nil {
			o.logger.Errorf("can't get orders for user %s: %s", user.ID.String(), err)
			return nil, fmt.Errorf("can't get orders for user %s: %w", user.ID.String(), err)
		}
		for ordr := range resChan {
			result = append(result, ordr)
		}
		return result, nil
	}
}
func (o *order) DeleteOrder(ctx context.Context, order *models.Order) error {
	select {
	case <-ctx.Done():
		o.logger.Error("context closed")
		return fmt.Errorf("context closed")
	default:
		if err := o.orderStore.DeleteOrder(ctx, order); err != nil {
			o.logger.Error("can't delete order %s", err)
			return fmt.Errorf("can't delete order %w", err)
		}
		return nil
	}
}

func (o *order) ChangeAddress(ctx context.Context, order *models.Order, newAddress models.UserAddress) error {
	select {
	case <-ctx.Done():
		o.logger.Error("context closed")
		return fmt.Errorf("context closed")
	default:
		if newAddress == order.Address {
			return nil
		}
		if err := o.orderStore.ChangeAddress(ctx, order, newAddress); err != nil {
			o.logger.Errorf("can't change address %s: ", err)
			return fmt.Errorf("can't change address %w: ", err)
		}
		return nil
	}
}

func (o *order) GetOrder(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	select {
	case <-ctx.Done():
		o.logger.Error("context closed")
		return nil, fmt.Errorf("context closed")
	default:
		res, err := o.orderStore.GetOrderByID(ctx, id)
		if err != nil {
			o.logger.Errorf("can't get order: %s", err)
			return nil, fmt.Errorf("can't get order: %w", err)
		}
		return &res, nil
	}
}
