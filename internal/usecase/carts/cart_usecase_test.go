package usecase

import (
	"OnlineShopBackend/internal/models"
	mocks "OnlineShopBackend/internal/usecase/repo_mocks"
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L().Sugar()
	store := mocks.NewMockCartStore(ctrl)
	usecase := NewCartUsecase(store, logger)
	ctx := context.Background()

	t.Run("error get cart", func(t *testing.T) {
		id := uuid.New()
		testErr := errors.New("error")
		store.EXPECT().Get(ctx, id).Return(nil, testErr)
		res, err := usecase.Get(ctx, id)
		require.Error(t, err)
		require.Nil(t, res)
	})
	t.Run("success get cart", func(t *testing.T) {
		id := uuid.New()
		store.EXPECT().Get(ctx, id).Return(&models.Cart{}, nil)
		res, err := usecase.Get(ctx, id)
		require.NoError(t, err)
		require.NotNil(t, res)
	})
}

func TestCartByUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L().Sugar()
	store := mocks.NewMockCartStore(ctrl)
	usecase := NewCartUsecase(store, logger)
	ctx := context.Background()
	t.Run("error get cart by user id", func(t *testing.T) {
		id := uuid.New()
		testErr := errors.New("error")
		store.EXPECT().CartByUserId(ctx, id).Return(nil, testErr)
		res, err := usecase.CartByUserId(ctx, id)
		require.Error(t, err)
		require.Nil(t, res)
	})
	t.Run("success get cart by user id", func(t *testing.T) {
		id := uuid.New()
		store.EXPECT().CartByUserId(ctx, id).Return(&models.Cart{}, nil)
		res, err := usecase.CartByUserId(ctx, id)
		require.NoError(t, err)
		require.NotNil(t, res)
	})
}

func TestDeleteItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L().Sugar()
	store := mocks.NewMockCartStore(ctrl)
	usecase := NewCartUsecase(store, logger)
	ctx := context.Background()
	t.Run("error delete item", func(t *testing.T) {
		cartId := uuid.New()
		userId := uuid.New()
		testErr := errors.New("error")
		store.EXPECT().DeleteItem(ctx, cartId, userId).Return(testErr)
		err := usecase.DeleteItem(ctx, cartId, userId)
		require.Error(t, err)
	})
	t.Run("success delete item", func(t *testing.T) {
		cartId := uuid.New()
		userId := uuid.New()
		store.EXPECT().DeleteItem(ctx, cartId, userId).Return(nil)
		err := usecase.DeleteItem(ctx, cartId, userId)
		require.NoError(t, err)
	})
}

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L().Sugar()
	store := mocks.NewMockCartStore(ctrl)
	usecase := NewCartUsecase(store, logger)
	ctx := context.Background()
	t.Run("error create cart", func(t *testing.T) {
		userId := uuid.New()
		testErr := errors.New("error")
		store.EXPECT().Create(ctx, userId).Return(uuid.Nil, testErr)
		res, err := usecase.Create(ctx, userId)
		require.Error(t, err)
		require.Equal(t, uuid.Nil, res)
	})
	t.Run("success create cart", func(t *testing.T) {
		userId := uuid.New()
		cartId := uuid.New()
		store.EXPECT().Create(ctx, userId).Return(cartId, nil)
		res, err := usecase.Create(ctx, userId)
		require.NoError(t, err)
		require.Equal(t, cartId, res)
	})
}

func TestAddItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L().Sugar()
	store := mocks.NewMockCartStore(ctrl)
	usecase := NewCartUsecase(store, logger)
	ctx := context.Background()
	t.Run("error add item", func(t *testing.T) {
		cartId := uuid.New()
		userId := uuid.New()
		testErr := errors.New("error")
		store.EXPECT().AddItem(ctx, cartId, userId).Return(testErr)
		err := usecase.AddItem(ctx, cartId, userId)
		require.Error(t, err)
	})
	t.Run("success add item", func(t *testing.T) {
		cartId := uuid.New()
		userId := uuid.New()
		store.EXPECT().AddItem(ctx, cartId, userId).Return(nil)
		err := usecase.AddItem(ctx, cartId, userId)
		require.NoError(t, err)
	})
}
