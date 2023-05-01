package usecase

import (
	"OnlineShopBackend/internal/models"
	mocks "OnlineShopBackend/internal/usecase/repo_mocks"
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

/*var (
	testUser = models.User{
		Firstname: "TestName",
		Lastname:  "TestLastName",
		Password:  "123",
		Email:     "testeamil@123.ru",
		Address: models.UserAddress{
			Zipcode: "123",
			Country: "Israel",
			City:    "Haifa",
			Street:  "דניל 4",
		},
		Rights: models.Rights{
			Name:  "pacan",
			Rules: []string{"buy"},
		},
	}

	testCategory = models.Category{
		Name:        "Electronics",
		Description: "Electric stuff",
		Image:       "image.url",
	}

	testItem11 = models.Item{
		Title:       "testItem11",
		Description: "Awesome chinese item",
		Price:       300,
		Category:    testCategory,
		Vendor:      "chinese factory",
		Images:      []string{},
	}
	testItem2 = models.Item{
		Title:       "testItem2",
		Description: "Awesome chinese item",
		Price:       500,
		Category:    testCategory,
		Vendor:      "russian factory",
		Images:      []string{},
	}

	testOrder = models.Order{
		ShipmentTime: time.Now().Add(models.StandardShipmentPeriod),
		User:         testUser,
		Address:      testUser.Address,
		Status:       models.StatusCreated,
		Items: []models.ItemWithQuantity{
			{Item: testItem11, Quantity: 2}, {Item: testItem2, Quantity: 1},
		},
	}

	lgr = zap.NewExample().Sugar()
)

type orderRepoMock struct {
	err error
}

var _ usecase.OrderStore = (*orderRepoMock)(nil)

func (orMock *orderRepoMock) CreateOrder(ctx context.Context, order *models.Order) (*models.Order, error) {
	order.Id, _ = uuid.NewRandom()
	return order, orMock.err
}
func (orMock *orderRepoMock) DeleteOrder(ctx context.Context, order *models.Order) error {
	return orMock.err
}
func (orMock *orderRepoMock) ChangeAddress(ctx context.Context, order *models.Order, address models.UserAddress) error {
	order.Address = address
	return orMock.err
}
func (orMock *orderRepoMock) ChangeStatus(ctx context.Context, order *models.Order, status models.Status) error {
	order.Status = status
	return orMock.err
}

func (orMock *orderRepoMock) CreateCart(ctx context.Context, userId uuid.UUID) (uuid.UUID, error) {
	return uuid.Nil, nil
}

func (orMock *orderRepoMock) DeleteCart(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (orMock *orderRepoMock) GetOrderById(ctx context.Context, id uuid.UUID) (models.Order, error) {
	userID, _ := uuid.NewRandom()
	itemID1, _ := uuid.NewRandom()
	itemID2, _ := uuid.NewRandom()
	order := testOrder
	order.Id = id
	order.User.Id = userID
	order.Items[0].Id = itemID1
	order.Items[1].Id = itemID2
	return order, orMock.err
}

func (orMock *orderRepoMock) GetOrdersByUser(ctx context.Context, user *models.User) (chan models.Order, error) {
	res := make(chan models.Order, 1)
	return res, orMock.err
}

func TestPlaceOrder(t *testing.T) {
	uscs := NewOrderUsecase(&orderRepoMock{}, lgr)
	cartID, _ := uuid.NewRandom()
	userID, _ := uuid.NewRandom()
	cart := models.Cart{
		Id:     cartID,
		UserId: userID,
		Items: []models.ItemWithQuantity{
			{Item: testItem11, Quantity: 1},
			{Item: testItem2, Quantity: 1},
		},
		ExpireAt: time.Now().Add(2 * time.Hour),
	}
	res, err := uscs.CreateOrder(context.Background(), &cart, testUser, testOrder.Address)
	require.NoError(t, err)
	assert.Equal(t, testUser.Address, res.Address)
	assert.Equal(t, cart.Items, res.Items)
}

func TestPlaceOrderDBError(t *testing.T) {
	uscs := NewOrderUsecase(&orderRepoMock{err: fmt.Errorf("test error")}, lgr)
	cartID, _ := uuid.NewRandom()
	userID, _ := uuid.NewRandom()
	cart := models.Cart{
		Id:     cartID,
		UserId: userID,
		Items: []models.ItemWithQuantity{
			{Item: testItem11, Quantity: 2},
			{Item: testItem2, Quantity: 1},
		},
		ExpireAt: time.Now().Add(2 * time.Hour),
	}
	res, err := uscs.CreateOrder(context.Background(), &cart, testUser, testOrder.Address)
	require.Error(t, err)
	assert.Nil(t, res)
}

func TestChangeStatus(t *testing.T) {
	uscs := NewOrderUsecase(&orderRepoMock{}, lgr)
	err := uscs.ChangeStatus(context.Background(), &testOrder, models.StatusProcessed)
	defer func() {
		testOrder.Status = models.StatusCreated
	}()
	require.NoError(t, err)
	assert.Equal(t, models.StatusProcessed, testOrder.Status)

}

func TestChangeStatusError(t *testing.T) {
	uscs := NewOrderUsecase(&orderRepoMock{err: fmt.Errorf("test error")}, lgr)
	err := uscs.ChangeStatus(context.Background(), &testOrder, models.StatusProcessed)
	defer func() {
		testOrder.Status = models.StatusCreated
	}()
	require.Error(t, err)
}

func TestChangeAddress(t *testing.T) {
	uscs := NewOrderUsecase(&orderRepoMock{}, lgr)
	oldAddress := testOrder.Address
	err := uscs.ChangeAddress(context.Background(), &testOrder, models.UserAddress{
		Street:  "הלל 49",
		City:    "חיפה",
		Zipcode: "313455",
		Country: "Israel",
	})
	defer func() {
		testOrder.Address = oldAddress
	}()
	require.NoError(t, err)
}

func TestChangeAddressError(t *testing.T) {
	uscs := NewOrderUsecase(&orderRepoMock{err: fmt.Errorf("test error")}, lgr)
	oldAddress := testOrder.Address
	err := uscs.ChangeAddress(context.Background(), &testOrder, models.UserAddress{
		Street:  "הלל 49",
		City:    "חיפה",
		Zipcode: "313455",
		Country: "Israel",
	})
	defer func() {
		testOrder.Address = oldAddress
	}()
	require.Error(t, err)
}

func TestDeleteOrder(t *testing.T) {
	uscs := NewOrderUsecase(&orderRepoMock{}, lgr)
	err := uscs.DeleteOrder(context.Background(), &testOrder)
	require.NoError(t, err)
}

func TestGetOrder(t *testing.T) {
	id, _ := uuid.NewRandom()
	uscs := NewOrderUsecase(&orderRepoMock{}, lgr)
	order, err := uscs.GetOrder(context.Background(), id)
	require.NoError(t, err)
	assert.Equal(t, testOrder.User.Firstname, order.User.Firstname)
	assert.Equal(t, testOrder.ShipmentTime, order.ShipmentTime)
}*/

func TestCreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	orderRepo := mocks.NewMockOrderStore(ctrl)
	logger := zap.L()
	ctx := context.Background()
	ctxWithCancel, cancel := context.WithCancel(ctx)
	cancel()

	testCases := []struct {
		name    string
		store   *mocks.MockOrderStore
		logger  *zap.Logger
		cart    *models.Cart
		user    models.User
		address models.UserAddress
		order   *models.Order
		expect  func(*mocks.MockOrderStore)
		ctx     context.Context
	}{
		{
			name:    "error context closed",
			store:   orderRepo,
			logger:  logger,
			cart:    &models.Cart{},
			user:    models.User{},
			address: models.UserAddress{},
			ctx:     ctxWithCancel,
		},
		{
			name:    "error on create order",
			store:   orderRepo,
			logger:  logger,
			cart:    &models.Cart{},
			user:    models.User{},
			address: models.UserAddress{},
			order:   &models.Order{},
			expect: func(mos *mocks.MockOrderStore) {
				mos.EXPECT().CreateOrder(ctx, gomock.Any()).Return(nil, fmt.Errorf("error"))
			},
			ctx: ctx,
		},
		{
			name:    "success create order",
			store:   orderRepo,
			logger:  logger,
			cart:    &models.Cart{},
			user:    models.User{},
			address: models.UserAddress{},
			order:   &models.Order{},
			expect: func(mos *mocks.MockOrderStore) {
				mos.EXPECT().CreateOrder(ctx, gomock.Any()).Return(&models.Order{}, nil)
			},
			ctx: ctx,
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			usecase := NewOrderUsecase(orderRepo, logger.Sugar())
			if tc.expect != nil {
				tc.expect(tc.store)
			}
			res, err := usecase.CreateOrder(tc.ctx, tc.cart, tc.user, tc.address)
			if strings.Contains(tc.name, "error") {
				require.Error(t, err)
				require.Nil(t, res)
			} else {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, res, &models.Order{})
			}
		})
	}
}

func TestChangeStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	orderRepo := mocks.NewMockOrderStore(ctrl)
	logger := zap.L()
	ctx := context.Background()
	ctxWithCancel, cancel := context.WithCancel(ctx)
	cancel()

	testCases := []struct {
		name   string
		store  *mocks.MockOrderStore
		logger *zap.Logger
		order  *models.Order
		expect func(*mocks.MockOrderStore)
		ctx    context.Context
	}{
		{
			name:   "error context closed",
			store:  orderRepo,
			logger: logger,
			ctx:    ctxWithCancel,
		},
		{
			name:   "error on change order status",
			store:  orderRepo,
			logger: logger,
			order:  &models.Order{},
			expect: func(mos *mocks.MockOrderStore) {
				mos.EXPECT().ChangeStatus(ctx, &models.Order{}).Return(fmt.Errorf("error"))
			},
			ctx: ctx,
		},
		{
			name:   "success change order status",
			store:  orderRepo,
			logger: logger,
			order:  &models.Order{},
			expect: func(mos *mocks.MockOrderStore) {
				mos.EXPECT().ChangeStatus(ctx, &models.Order{}).Return(nil)
			},
			ctx: ctx,
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			usecase := NewOrderUsecase(orderRepo, logger.Sugar())
			if tc.expect != nil {
				tc.expect(tc.store)
			}
			err := usecase.ChangeStatus(tc.ctx, tc.order)
			if strings.Contains(tc.name, "error") {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestChangeAddress(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	orderRepo := mocks.NewMockOrderStore(ctrl)
	logger := zap.L()
	ctx := context.Background()
	ctxWithCancel, cancel := context.WithCancel(ctx)
	cancel()

	testCases := []struct {
		name   string
		store  *mocks.MockOrderStore
		logger *zap.Logger
		order  *models.Order
		expect func(*mocks.MockOrderStore)
		ctx    context.Context
	}{
		{
			name:   "error context closed",
			store:  orderRepo,
			logger: logger,
			ctx:    ctxWithCancel,
		},
		{
			name:   "error on change order address",
			store:  orderRepo,
			logger: logger,
			order:  &models.Order{},
			expect: func(mos *mocks.MockOrderStore) {
				mos.EXPECT().ChangeAddress(ctx, &models.Order{}).Return(fmt.Errorf("error"))
			},
			ctx: ctx,
		},
		{
			name:   "success change order address",
			store:  orderRepo,
			logger: logger,
			order:  &models.Order{},
			expect: func(mos *mocks.MockOrderStore) {
				mos.EXPECT().ChangeAddress(ctx, &models.Order{}).Return(nil)
			},
			ctx: ctx,
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			usecase := NewOrderUsecase(orderRepo, logger.Sugar())
			if tc.expect != nil {
				tc.expect(tc.store)
			}
			err := usecase.ChangeAddress(tc.ctx, tc.order)
			if strings.Contains(tc.name, "error") {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}