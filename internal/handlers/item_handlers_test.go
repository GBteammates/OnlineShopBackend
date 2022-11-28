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

func TestCreateItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	categoryRepo := mocks.NewMockCategoryStore(ctrl)
	itemRepo := mocks.NewMockItemStore(ctrl)
	cash := mocks.NewMockCash(ctrl)
	usecase := usecase.NewStorage(itemRepo, categoryRepo, cash, logger)
	handlers := NewHandlers(usecase, logger)
	testKey := "ItemsList"
	testItem := Item{
		Title:       "TestTitle",
		Description: "TestDescription",
		Category:    "b02c1542-dba1-46d2-ac71-e770c13d0d50",
		Price:       1,
		Vendor:      "TestVendor",
	}
	testCategoryId, _ := uuid.Parse(testItem.Category)
	testModelItem := &models.Item{
		Title:       testItem.Title,
		Description: testItem.Description,
		Category:    testCategoryId,
		Price:       testItem.Price,
		Vendor:      testItem.Vendor,
	}
	expect, _ := uuid.Parse("feb77bbc-1b8a-4739-bd68-d3b052af9a80")
	itemRepo.EXPECT().CreateItem(ctx, testModelItem).Return(expect, nil)
	cash.EXPECT().CheckCash(testKey).Return(false)
	res, err := handlers.CreateItem(ctx, testItem)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, expect)

	err = fmt.Errorf("error on create item")
	itemRepo.EXPECT().CreateItem(ctx, testModelItem).Return(uuid.Nil, err)
	res, err = handlers.CreateItem(ctx, testItem)
	require.Error(t, err)
	require.Equal(t, res, uuid.Nil)
}

func TestUpdateItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	categoryRepo := mocks.NewMockCategoryStore(ctrl)
	itemRepo := mocks.NewMockItemStore(ctrl)
	cash := mocks.NewMockCash(ctrl)
	usecase := usecase.NewStorage(itemRepo, categoryRepo, cash, logger)
	handlers := NewHandlers(usecase, logger)
	testKey := "ItemsList"
	testItem := Item{
		Id:          "feb77bbc-1b8a-4739-bd68-d3b052af9a80",
		Title:       "TestTitle",
		Description: "TestDescription",
		Category:    "b02c1542-dba1-46d2-ac71-e770c13d0d50",
		Price:       1,
		Vendor:      "TestVendor",
	}
	itemId, _ := uuid.Parse(testItem.Id)
	testCategoryId, _ := uuid.Parse(testItem.Category)
	testModelItem := &models.Item{
		Id:          itemId,
		Title:       testItem.Title,
		Description: testItem.Description,
		Category:    testCategoryId,
		Price:       testItem.Price,
		Vendor:      testItem.Vendor,
	}
	itemRepo.EXPECT().UpdateItem(ctx, testModelItem).Return(nil)
	cash.EXPECT().CheckCash(testKey).Return(false)
	err := handlers.UpdateItem(ctx, testItem)
	require.NoError(t, err)

	err = fmt.Errorf("error on update item")
	itemRepo.EXPECT().UpdateItem(ctx, testModelItem).Return(err)
	err = handlers.UpdateItem(ctx, testItem)
	require.Error(t, err)
}

func TestGetItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	categoryRepo := mocks.NewMockCategoryStore(ctrl)
	itemRepo := mocks.NewMockItemStore(ctrl)
	cash := mocks.NewMockCash(ctrl)
	usecase := usecase.NewStorage(itemRepo, categoryRepo, cash, logger)
	handlers := NewHandlers(usecase, logger)
	id := "feb77bbc-1b8a-4739-bd68-d3b052af9a80"
	uid, _ := uuid.Parse(id)
	testModelItem := &models.Item{
		Id:          uid,
		Title:       "TestTitle",
		Description: "TestDescription",
		Category:    uuid.New(),
		Price:       1,
		Vendor:      "TestVendor",
	}
	testItem := Item{
		Id:          id,
		Title:       testModelItem.Title,
		Description: testModelItem.Description,
		Category:    testModelItem.Category.String(),
		Price:       testModelItem.Price,
		Vendor:      testModelItem.Vendor,
	}
	itemRepo.EXPECT().GetItem(ctx, uid).Return(testModelItem, nil)
	res, err := handlers.GetItem(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, testItem)

	res, err = handlers.GetItem(ctx, "invalidId")
	require.Error(t, err)
	require.Equal(t, res, Item{})

	err = fmt.Errorf("error on get item")
	itemRepo.EXPECT().GetItem(ctx, uid).Return(&models.Item{}, err)
	res, err = handlers.GetItem(ctx, id)
	require.Error(t, err)
	require.Equal(t, res, Item{})
}

func TestItemsList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	categoryRepo := mocks.NewMockCategoryStore(ctrl)
	itemRepo := mocks.NewMockItemStore(ctrl)
	cash := mocks.NewMockCash(ctrl)
	usecase := usecase.NewStorage(itemRepo, categoryRepo, cash, logger)
	handlers := NewHandlers(usecase, logger)
	testKey := "ItemsList"
	id := "feb77bbc-1b8a-4739-bd68-d3b052af9a80"
	uid, _ := uuid.Parse(id)
	testModelItem := models.Item{
		Id:          uid,
		Title:       "TestTitle",
		Description: "TestDescription",
		Category:    uuid.New(),
		Price:       1,
		Vendor:      "TestVendor",
	}
	testModelSlice := []models.Item{testModelItem}
	testItem := Item{
		Id:          id,
		Title:       testModelItem.Title,
		Description: testModelItem.Description,
		Category:    testModelItem.Category.String(),
		Price:       testModelItem.Price,
		Vendor:      testModelItem.Vendor,
	}
	testSlice := make([]Item, 0, 100)
	testSlice = append(testSlice, testItem)

	cash.EXPECT().CheckCash(testKey).Return(true)
	cash.EXPECT().GetCash(testKey).Return(testModelSlice, nil)
	res, err := handlers.ItemsList(ctx)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, testSlice)

	err = fmt.Errorf("error on itemslist()")
	testSlice2 := make([]Item, 0, 100)
	cash.EXPECT().CheckCash(testKey).Return(true)
	cash.EXPECT().GetCash(testKey).Return(nil, err)
	res, err = handlers.ItemsList(ctx)
	require.Error(t, err)
	require.Equal(t, res, testSlice2)
}

func TestSearchLine(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	categoryRepo := mocks.NewMockCategoryStore(ctrl)
	itemRepo := mocks.NewMockItemStore(ctrl)
	cash := mocks.NewMockCash(ctrl)
	usecase := usecase.NewStorage(itemRepo, categoryRepo, cash, logger)
	handlers := NewHandlers(usecase, logger)
	id := "feb77bbc-1b8a-4739-bd68-d3b052af9a80"
	uid, _ := uuid.Parse(id)
	testModelItem := models.Item{
		Id:          uid,
		Title:       "TestTitle",
		Description: "TestDescription",
		Category:    uuid.New(),
		Price:       1,
		Vendor:      "TestVendor",
	}
	testItem := Item{
		Id:          id,
		Title:       testModelItem.Title,
		Description: testModelItem.Description,
		Category:    testModelItem.Category.String(),
		Price:       testModelItem.Price,
		Vendor:      testModelItem.Vendor,
	}
	param := "est"
	testChan := make(chan models.Item, 1)
	testChan <- testModelItem
	close(testChan)
	testSlice := make([]Item, 0, 100)
	testSlice = append(testSlice, testItem)
	itemRepo.EXPECT().SearchLine(ctx, param).Return(testChan, nil)
	res, err := handlers.SearchLine(ctx, param)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, testSlice)

	err = fmt.Errorf("error on search line()")
	testSlice2 := make([]Item, 0, 100)
	itemRepo.EXPECT().SearchLine(ctx, param).Return(testChan, err)
	res, err = handlers.SearchLine(ctx, param)
	require.Error(t, err)
	require.Equal(t, res, testSlice2)
}
