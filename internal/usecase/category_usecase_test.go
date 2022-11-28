package usecase

import (
	"OnlineShopBackend/internal/models"
	"OnlineShopBackend/internal/usecase/mocks"
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestCreateCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	categoryRepo := mocks.NewMockCategoryStore(ctrl)
	itemRepo := mocks.NewMockItemStore(ctrl)
	cash := mocks.NewMockCash(ctrl)
	usecase := NewStorage(itemRepo, categoryRepo, cash, logger)

	testModelCategory := &models.Category{
		Name:        "test name",
		Description: "test description",
	}
	expect, _ := uuid.Parse("feb77bbc-1b8a-4739-bd68-d3b052af9a80")
	categoryRepo.EXPECT().CreateCategory(ctx, testModelCategory).Return(expect, nil)
	res, err := usecase.CreateCategory(ctx, testModelCategory)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, expect)
	err = fmt.Errorf("error on create category")
	categoryRepo.EXPECT().CreateCategory(ctx, testModelCategory).Return(uuid.Nil, err)
	res, err = usecase.CreateCategory(ctx, testModelCategory)
	require.Error(t, err)
	require.Equal(t, res, uuid.Nil)
}

func TestGetCategoryList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	categoryRepo := mocks.NewMockCategoryStore(ctrl)
	itemRepo := mocks.NewMockItemStore(ctrl)
	cash := mocks.NewMockCash(ctrl)
	usecase := NewStorage(itemRepo, categoryRepo, cash, logger)
	id, _ := uuid.Parse("feb77bbc-1b8a-4739-bd68-d3b052af9a80")
	testModelCategory := models.Category{
		Id:          id,
		Name:        "TestName",
		Description: "TestDescription",
	}
	testChan := make(chan models.Category, 2)
	testChan <- testModelCategory
	close(testChan)
	categoryRepo.EXPECT().GetCategoryList(ctx).Return(testChan, nil)
	res, err := usecase.GetCategoryList(ctx)
	require.NoError(t, err)
	require.NotNil(t, res)

	categoryRepo.EXPECT().GetCategoryList(ctx).Return(nil, fmt.Errorf("error on categories list query context"))
	res, err = usecase.GetCategoryList(ctx)
	require.Error(t, err)
	require.Nil(t, res)
}
