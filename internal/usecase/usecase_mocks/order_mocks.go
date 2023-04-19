package mocks

import (
	"OnlineShopBackend/internal/models"
	usecase "OnlineShopBackend/internal/usecase/interfaces"
	context "context"
	"time"

	"github.com/google/uuid"
)

type OrderUsecaseMock struct {
	Err error
}

var _ usecase.IOrderUsecase = (*OrderUsecaseMock)(nil)

func (o *OrderUsecaseMock) PlaceOrder(ctx context.Context, cart *models.Cart, user models.User, address models.UserAddress) (*models.Order, error) {
	return &models.Order{
		Id: uuid.New(),
	}, o.Err
}
func (o *OrderUsecaseMock) ChangeStatus(ctx context.Context, order *models.Order, newStatus models.Status) error {
	return o.Err
}
func (o *OrderUsecaseMock) GetOrdersByUser(ctx context.Context, user *models.User) ([]models.Order, error) {
	return []models.Order{
		{
			Id:           uuid.New(),
			ShipmentTime: time.Now().Add(time.Duration(models.StandardShipmentPeriod.Hours())),
			User:         *user,
			Status:       models.StatusCourier,
			Items: []models.ItemWithQuantity{
				{
					Item: models.Item{
						Id:    uuid.New(),
						Title: "test1",
					},
					Quantity: 1,
				},
				{
					Item: models.Item{
						Id:    uuid.New(),
						Title: "test2",
					},
					Quantity: 2,
				},
			},
		}, {
			Id:           uuid.New(),
			ShipmentTime: time.Now().Add(time.Duration(models.StandardShipmentPeriod.Hours())),
			User:         *user,
			Status:       models.StatusCourier,
			Items: []models.ItemWithQuantity{
				{
					Item: models.Item{
						Id:    uuid.New(),
						Title: "test2",
					},
					Quantity: 3,
				},
				{
					Item: models.Item{
						Id:    uuid.New(),
						Title: "test2",
					},
					Quantity: 4,
				},
			},
		},
	}, o.Err
}

func (o *OrderUsecaseMock) DeleteOrder(ctx context.Context, order *models.Order) error {
	return o.Err
}
func (o *OrderUsecaseMock) ChangeAddress(ctx context.Context, order *models.Order, newAddress models.UserAddress) error {
	return o.Err
}
func (o *OrderUsecaseMock) GetOrder(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	return &models.Order{
		Id:           uuid.New(),
		ShipmentTime: time.Now().Add(time.Duration(models.StandardShipmentPeriod.Hours())),
		User: models.User{
			Id: uuid.New(),
		},
		Status: models.StatusCourier,
		Items: []models.ItemWithQuantity{
			{
				Item: models.Item{
					Id:    uuid.New(),
					Title: "test1",
				},
				Quantity: 1,
			},
			{
				Item: models.Item{
					Id:    uuid.New(),
					Title: "test2",
				},
				Quantity: 2,
			},
		},
	}, o.Err
}
