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

var _ usecase.IOrderUsecase = (*orderUsecase)(nil)

type orderUsecase struct {
	store  usecase.OrderStore
	logger *zap.SugaredLogger
}

func NewOrderUsecase(store usecase.OrderStore, logger *zap.SugaredLogger) *orderUsecase {
	return &orderUsecase{
		store:  store,
		logger: logger,
	}
}

func (u *orderUsecase) Create(ctx context.Context, cart *models.Cart, user models.User, address models.UserAddress) (*models.Order, error) {
	select {
	case <-ctx.Done():
		u.logger.Error("context closed")
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
		res, err := u.store.Create(ctx, &ordr)
		if err != nil {
			u.logger.Errorf("can't add order to db %s", err)
			return nil, fmt.Errorf("can't place order to db : %w", err)
		}
		u.logger.Debugf("order %s created", res.Id.String())
		return res, nil
	}
}

func (u *orderUsecase) ChangeStatus(ctx context.Context, order *models.Order) error {
	select {
	case <-ctx.Done():
		u.logger.Error("context closed")
		return fmt.Errorf("context closed")
	default:
		if err := u.store.ChangeStatus(ctx, order); err != nil {
			u.logger.Errorf("can't change status of order: %s", err)
			return fmt.Errorf("can't change status of order: %w", err)
		}
	}
	return nil
}
func (u *orderUsecase) List(ctx context.Context, user *models.User) ([]models.Order, error) {
	select {
	case <-ctx.Done():
		u.logger.Error("context closed")
		return nil, fmt.Errorf("context closed")
	default:
		result := make([]models.Order, 0, 10)
		resChan, err := u.store.List(ctx, user)
		if err != nil {
			u.logger.Errorf("can't get orders for user %s: %s", user.Id.String(), err)
			return nil, fmt.Errorf("can't get orders for user %s: %w", user.Id.String(), err)
		}
		for ordr := range resChan {
			result = append(result, ordr)
		}
		return result, nil
	}
}
func (u *orderUsecase) Delete(ctx context.Context, order *models.Order) error {
	select {
	case <-ctx.Done():
		u.logger.Error("context closed")
		return fmt.Errorf("context closed")
	default:
		if err := u.store.Delete(ctx, order); err != nil {
			u.logger.Error("can't delete order %s", err)
			return fmt.Errorf("can't delete order %w", err)
		}
		return nil
	}
}

func (u *orderUsecase) ChangeAddress(ctx context.Context, order *models.Order) error {
	select {
	case <-ctx.Done():
		u.logger.Error("context closed")
		return fmt.Errorf("context closed")
	default:
		if err := u.store.ChangeAddress(ctx, order); err != nil {
			u.logger.Errorf("can't change address %s: ", err)
			return fmt.Errorf("can't change address %w: ", err)
		}
		return nil
	}
}

func (u *orderUsecase) Get(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	select {
	case <-ctx.Done():
		u.logger.Error("context closed")
		return nil, fmt.Errorf("context closed")
	default:
		res, err := u.store.Get(ctx, id)
		if err != nil {
			u.logger.Errorf("can't get order: %s", err)
			return nil, fmt.Errorf("can't get order: %w", err)
		}
		return &res, nil
	}
}
