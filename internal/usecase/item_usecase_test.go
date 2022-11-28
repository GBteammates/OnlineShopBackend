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

func TestCreateItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	categoryRepo := mocks.NewMockCategoryStore(ctrl)
	itemRepo := mocks.NewMockItemStore(ctrl)
	cash := mocks.NewMockCash(ctrl)
	usecase := NewStorage(itemRepo, categoryRepo, cash, logger)
	testCategoryId, _ := uuid.Parse("feb77bbc-1b8a-4739-bd68-d3b052af9a80")
	testModelItem := &models.Item{
		Title:       "TestTitle",
		Description: "TestDescription",
		Category:    testCategoryId,
		Price:       1,
		Vendor:      "TestVendor",
	}
	expect, _ := uuid.Parse("13574b3b-0c44-4864-89de-a086ad68ec4b")
	itemRepo.EXPECT().CreateItem(ctx, testModelItem).Return(expect, nil)
	res, err := usecase.CreateItem(ctx, testModelItem)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, expect)
	err = fmt.Errorf("error on create item")
	itemRepo.EXPECT().CreateItem(ctx, testModelItem).Return(uuid.Nil, err)
	res, err = usecase.CreateItem(ctx, testModelItem)
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
	usecase := NewStorage(itemRepo, categoryRepo, cash, logger)

	itemId, _ := uuid.Parse("feb77bbc-1b8a-4739-bd68-d3b052af9a80")
	testCategoryId, _ := uuid.Parse("b02c1542-dba1-46d2-ac71-e770c13d0d50")
	testModelItem := &models.Item{
		Id:          itemId,
		Title:       "TestTitle",
		Description: "TestDescription",
		Category:    testCategoryId,
		Price:       1,
		Vendor:      "TestVendor",
	}
	itemRepo.EXPECT().UpdateItem(ctx, testModelItem).Return(nil)
	err := usecase.UpdateItem(ctx, testModelItem)
	require.NoError(t, err)

	err = fmt.Errorf("error on update item")
	itemRepo.EXPECT().UpdateItem(ctx, testModelItem).Return(err)
	err = usecase.UpdateItem(ctx, testModelItem)
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
	usecase := NewStorage(itemRepo, categoryRepo, cash, logger)
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
	itemRepo.EXPECT().GetItem(ctx, uid).Return(testModelItem, nil)
	res, err := usecase.GetItem(ctx, uid)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, testModelItem)

	err = fmt.Errorf("error on get item")
	itemRepo.EXPECT().GetItem(ctx, uid).Return(&models.Item{}, err)
	res, err = usecase.GetItem(ctx, uid)
	require.Error(t, err)
	require.Equal(t, res, &models.Item{})
}

func TestItemsList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	categoryRepo := mocks.NewMockCategoryStore(ctrl)
	itemRepo := mocks.NewMockItemStore(ctrl)
	cash := mocks.NewMockCash(ctrl)
	usecase := NewStorage(itemRepo, categoryRepo, cash, logger)
	
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
	testKey := "ItemsList"
	testChan := make(chan models.Item, 1)
	testChan <- testModelItem
	close(testChan)
	expect := make([]models.Item, 1)
	expect = append(expect, testModelItem)
	
	cash.EXPECT().CheckCash(testKey).Return(false)
	itemRepo.EXPECT().ItemsList(ctx).Return(testChan, nil)
	cash.EXPECT().CreateCash(ctx, testChan, testKey).Return(nil)
	cash.EXPECT().GetCash(testKey).Return(expect, nil)
	res, err := usecase.ItemsList(ctx)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, expect)

	cash.EXPECT().CheckCash(testKey).Return(true)
	cash.EXPECT().GetCash(testKey).Return(expect, nil)
	res, err = usecase.ItemsList(ctx)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, expect)

	err = fmt.Errorf("error on itemslist()")
	cash.EXPECT().CheckCash(testKey).Return(false)
	itemRepo.EXPECT().ItemsList(ctx).Return(testChan, err)
	res, err = usecase.ItemsList(ctx)
	require.Error(t, err)
	require.Nil(t, res)

	cash.EXPECT().CheckCash(testKey).Return(false)
	itemRepo.EXPECT().ItemsList(ctx).Return(testChan, nil)
	cash.EXPECT().CreateCash(ctx, testChan, testKey).Return(err)
	res, err = usecase.ItemsList(ctx)
	require.Error(t, err)
	require.Nil(t, res)

	cash.EXPECT().CheckCash(testKey).Return(false)
	itemRepo.EXPECT().ItemsList(ctx).Return(testChan, nil)
	cash.EXPECT().CreateCash(ctx, testChan, testKey).Return(nil)
	cash.EXPECT().GetCash(testKey).Return(nil, err)
	res, err = usecase.ItemsList(ctx)
	require.Error(t, err)
	require.Nil(t, res)
}

func TestSearchLine(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	categoryRepo := mocks.NewMockCategoryStore(ctrl)
	itemRepo := mocks.NewMockItemStore(ctrl)
	cash := mocks.NewMockCash(ctrl)
	usecase := NewStorage(itemRepo, categoryRepo, cash, logger)
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
	param := "est"
	testChan := make(chan models.Item, 1)
	testChan <- testModelItem
	close(testChan)
	itemRepo.EXPECT().SearchLine(ctx, param).Return(testChan, nil)
	res, err := usecase.SearchLine(ctx, param)
	require.NoError(t, err)
	require.NotNil(t, res)

	err = fmt.Errorf("error on search line()")
	itemRepo.EXPECT().SearchLine(ctx, param).Return(testChan, err)
	res, err = usecase.SearchLine(ctx, param)
	require.Error(t, err)
	require.Nil(t, res)
}
