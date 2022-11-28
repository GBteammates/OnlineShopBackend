package handlers

import (
	"OnlineShopBackend/internal/models"
	"OnlineShopBackend/internal/usecase"
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
	usecase := usecase.NewStorage(itemRepo, categoryRepo, cash, logger)
	handlers := NewHandlers(usecase, logger)
	testCategory := Category{
		Name:        "TestName",
		Description: "TestDescription",
	}
	testModelCategory := &models.Category{
		Name:        testCategory.Name,
		Description: testCategory.Description,
	}
	expect, _ := uuid.Parse("feb77bbc-1b8a-4739-bd68-d3b052af9a80")
	categoryRepo.EXPECT().CreateCategory(ctx, testModelCategory).Return(expect, nil)
	res, err := handlers.CreateCategory(ctx, testCategory)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, expect)
	err = fmt.Errorf("error on create category")
	categoryRepo.EXPECT().CreateCategory(ctx, testModelCategory).Return(uuid.Nil, err)
	res, err = handlers.CreateCategory(ctx, testCategory)
	require.Error(t, err)
	require.Equal(t, res, uuid.Nil)
}

func TestGetCategoryList(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	categoryRepo := mocks.NewMockCategoryStore(ctrl)
	itemRepo := mocks.NewMockItemStore(ctrl)
	cash := mocks.NewMockCash(ctrl)
	usecase := usecase.NewStorage(itemRepo, categoryRepo, cash, logger)
	handlers := NewHandlers(usecase, logger)
	id, _ := uuid.Parse("feb77bbc-1b8a-4739-bd68-d3b052af9a80")
	testModelCategory := models.Category{
		Id:          id,
		Name:        "TestName",
		Description: "TestDescription",
	}
	testChan := make(chan models.Category, 2)
	testChan <- testModelCategory
	close(testChan)
	testChan2 := make(chan models.Category, 1)

	expect := make([]Category, 0, 100)
	expect = append(expect, Category{
		Id:          testModelCategory.Id.String(),
		Name:        testModelCategory.Name,
		Description: testModelCategory.Description,
	})
	categoryRepo.EXPECT().GetCategoryList(ctx).Return(testChan, nil)
	res, err := handlers.GetCategoryList(ctx)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, expect)

	categoryRepo.EXPECT().GetCategoryList(ctx).Return(testChan2, fmt.Errorf("error on categories list query context"))
	res, err = handlers.GetCategoryList(ctx)
	expect2 := make([]Category, 0, 100)
	require.Error(t, err)
	require.Equal(t, res, expect2)
}
