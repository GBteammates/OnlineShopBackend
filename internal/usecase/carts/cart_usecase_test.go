package usecase

import (
	"OnlineShopBackend/internal/models"
	mocks "OnlineShopBackend/internal/usecase/repo_mocks"
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var (
	testId         = uuid.New()
	testModelsCart = &models.Cart{
		Id:    testId,
		Items: testItems,
	}
	testItem = models.ItemWithQuantity{
		Quantity: 1,
	}
	testItems = []models.ItemWithQuantity{
		testItem,
	}
)

func TestGetCart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	cartRepo := mocks.NewMockCartStore(ctrl)
	usecase := NewCartUsecase(cartRepo, logger)
	ctx := context.Background()

	cartRepo.EXPECT().GetCart(ctx, testId).Return(nil, fmt.Errorf("error"))
	res, err := usecase.GetCart(ctx, testId)
	require.Error(t, err)
	require.Nil(t, res)

	cartRepo.EXPECT().GetCart(ctx, testId).Return(testModelsCart, nil)
	res, err = usecase.GetCart(ctx, testId)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, testModelsCart)
}

func TestGetCartByUserId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	cartRepo := mocks.NewMockCartStore(ctrl)
	usecase := NewCartUsecase(cartRepo, logger)
	ctx := context.Background()

	cartRepo.EXPECT().GetCartByUserId(ctx, testId).Return(nil, fmt.Errorf("error"))
	res, err := usecase.GetCartByUserId(ctx, testId)
	require.Error(t, err)
	require.Nil(t, res)

	cartRepo.EXPECT().GetCartByUserId(ctx, testId).Return(testModelsCart, nil)
	res, err = usecase.GetCartByUserId(ctx, testId)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, testModelsCart)
}

func TestDeleteItemFromCart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	cartRepo := mocks.NewMockCartStore(ctrl)
	usecase := NewCartUsecase(cartRepo, logger)
	ctx := context.Background()

	cartRepo.EXPECT().DeleteItemFromCart(ctx, testId, testId).Return(fmt.Errorf("error"))
	err := usecase.DeleteItemFromCart(ctx, testId, testId)
	require.Error(t, err)

	cartRepo.EXPECT().DeleteItemFromCart(ctx, testId, testId).Return(nil)
	err = usecase.DeleteItemFromCart(ctx, testId, testId)
	require.NoError(t, err)
}

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	cartRepo := mocks.NewMockCartStore(ctrl)
	usecase := NewCartUsecase(cartRepo, logger)
	ctx := context.Background()

	cartRepo.EXPECT().CreateCart(ctx, testId).Return(uuid.Nil, fmt.Errorf("error"))
	res, err := usecase.CreateCart(ctx, testId)
	require.Error(t, err)
	require.Equal(t, res, uuid.Nil)

	cartRepo.EXPECT().CreateCart(ctx, testId).Return(testId, nil)
	res, err = usecase.CreateCart(ctx, testId)
	require.NoError(t, err)
	require.Equal(t, res, testId)
}

func TestAddItemToCart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	cartRepo := mocks.NewMockCartStore(ctrl)
	usecase := NewCartUsecase(cartRepo, logger)
	ctx := context.Background()

	cartRepo.EXPECT().AddItemToCart(ctx, testId, testId).Return(fmt.Errorf("error"))
	err := usecase.AddItemToCart(ctx, testId, testId)
	require.Error(t, err)

	cartRepo.EXPECT().AddItemToCart(ctx, testId, testId).Return(nil)
	err = usecase.AddItemToCart(ctx, testId, testId)
	require.NoError(t, err)
}

func TestDeleteCart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	cartRepo := mocks.NewMockCartStore(ctrl)
	usecase := NewCartUsecase(cartRepo, logger)
	ctx := context.Background()

	cartRepo.EXPECT().DeleteCart(ctx, testId).Return(fmt.Errorf("error"))
	err := usecase.DeleteCart(ctx, testId)
	require.Error(t, err)

	cartRepo.EXPECT().DeleteCart(ctx, testId).Return(nil)
	err = usecase.DeleteCart(ctx, testId)
	require.NoError(t, err)
}
