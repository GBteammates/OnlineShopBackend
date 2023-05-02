package usecase

import (
	"OnlineShopBackend/internal/models"
	mocks "OnlineShopBackend/internal/usecase/repo_mocks"
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

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

func TestGetOrdersByUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	orderRepo := mocks.NewMockOrderStore(ctrl)
	logger := zap.L()
	ctx := context.Background()
	ctxWithCancel, cancel := context.WithCancel(ctx)
	cancel()
	testChan := make(chan models.Order, 1)
	testOrder := models.Order{}
	testChan <- testOrder
	testRes := make([]models.Order, 0, 10)
	testRes = append(testRes, testOrder)
	close(testChan)

	testCases := []struct {
		name   string
		store  *mocks.MockOrderStore
		logger *zap.Logger
		user   *models.User
		res    []models.Order
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
			name:   "error on get orders by user",
			store:  orderRepo,
			logger: logger,
			user:   &models.User{},
			expect: func(mos *mocks.MockOrderStore) {
				mos.EXPECT().GetOrdersByUser(ctx, &models.User{}).Return(nil, fmt.Errorf("error"))
			},
			ctx: ctx,
		},
		{
			name:   "success get orders by user",
			store:  orderRepo,
			logger: logger,
			user:   &models.User{},
			res:    testRes,
			expect: func(mos *mocks.MockOrderStore) {
				mos.EXPECT().GetOrdersByUser(ctx, &models.User{}).Return(testChan, nil)
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
			res, err := usecase.GetOrdersByUser(tc.ctx, tc.user)
			if strings.Contains(tc.name, "error") {
				require.Error(t, err)
				require.Nil(t, res)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.res, res)
			}
		})
	}
}

func TestDeleteOrder(t *testing.T) {
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
			name:   "error on delete order",
			store:  orderRepo,
			logger: logger,
			order:  &models.Order{},
			expect: func(mos *mocks.MockOrderStore) {
				mos.EXPECT().DeleteOrder(ctx, &models.Order{}).Return(fmt.Errorf("error"))
			},
			ctx: ctx,
		},
		{
			name:   "success delete order",
			store:  orderRepo,
			logger: logger,
			order:  &models.Order{},
			expect: func(mos *mocks.MockOrderStore) {
				mos.EXPECT().DeleteOrder(ctx, &models.Order{}).Return(nil)
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
			err := usecase.DeleteOrder(tc.ctx, tc.order)
			if strings.Contains(tc.name, "error") {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	orderRepo := mocks.NewMockOrderStore(ctrl)
	logger := zap.L()
	ctx := context.Background()
	ctxWithCancel, cancel := context.WithCancel(ctx)
	cancel()

	var (
		id = uuid.New()
	)

	testCases := []struct {
		name   string
		store  *mocks.MockOrderStore
		logger *zap.Logger
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
			name:   "error on get order",
			store:  orderRepo,
			logger: logger,
			expect: func(mos *mocks.MockOrderStore) {
				mos.EXPECT().GetOrderById(ctx, id).Return(models.Order{}, fmt.Errorf("error"))
			},
			ctx: ctx,
		},
		{
			name:   "success get order",
			store:  orderRepo,
			logger: logger,
			expect: func(mos *mocks.MockOrderStore) {
				mos.EXPECT().GetOrderById(ctx, id).Return(models.Order{}, nil)
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
			res, err := usecase.GetOrder(tc.ctx, id)
			if strings.Contains(tc.name, "error") {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, res)
			}
		})
	}
}
