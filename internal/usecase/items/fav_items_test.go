package usecase

import (
	"OnlineShopBackend/internal/models"
	mocks "OnlineShopBackend/internal/usecase/repo_mocks"
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestAddFavouriteItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cache := mocks.NewMockIItemsCache(ctrl)
	filestorage := mocks.NewMockFileStorager(ctrl)
	usecase := NewItemUsecase(itemRepo, cache, filestorage, logger)
	ctx := context.Background()

	itemRepo.EXPECT().AddFavouriteItem(ctx, testId, testItemId).Return(err)
	err := usecase.AddFavouriteItem(ctx, testId, testItemId)
	require.Error(t, err)

	itemRepo.EXPECT().AddFavouriteItem(ctx, testId, testItemId).Return(nil)
	cache.EXPECT().CheckCache(ctx, testId.String()+"nameasc").Return(false)
	cache.EXPECT().CheckCache(ctx, testId.String()+"namedesc").Return(false)
	cache.EXPECT().CheckCache(ctx, testId.String()+"priceasc").Return(false)
	cache.EXPECT().CheckCache(ctx, testId.String()+"pricedesc").Return(false)
	cache.EXPECT().CheckCache(ctx, testId.String()+"Fav").Return(true)
	cache.EXPECT().GetFavouriteItemsIdCache(ctx, testId.String()+"Fav").Return(nil, err)
	err = usecase.AddFavouriteItem(ctx, testId, testItemId)
	require.NoError(t, err)
}

func TestDeleteFavouriteItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cache := mocks.NewMockIItemsCache(ctrl)
	filestorage := mocks.NewMockFileStorager(ctrl)
	usecase := NewItemUsecase(itemRepo, cache, filestorage, logger)
	ctx := context.Background()

	itemRepo.EXPECT().DeleteFavouriteItem(ctx, testId, testItemId).Return(err)
	err := usecase.DeleteFavouriteItem(ctx, testId, testItemId)
	require.Error(t, err)

	itemRepo.EXPECT().DeleteFavouriteItem(ctx, testId, testItemId).Return(nil)
	cache.EXPECT().CheckCache(ctx, testId.String()+"nameasc").Return(false)
	cache.EXPECT().CheckCache(ctx, testId.String()+"namedesc").Return(false)
	cache.EXPECT().CheckCache(ctx, testId.String()+"priceasc").Return(false)
	cache.EXPECT().CheckCache(ctx, testId.String()+"pricedesc").Return(false)
	cache.EXPECT().CheckCache(ctx, testId.String()+"Fav").Return(true)
	cache.EXPECT().GetFavouriteItemsIdCache(ctx, testId.String()+"Fav").Return(nil, err)
	err = usecase.DeleteFavouriteItem(ctx, testId, testItemId)
	require.NoError(t, err)
}

func TestGetFavouriteItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cache := mocks.NewMockIItemsCache(ctrl)
	filestorage := mocks.NewMockFileStorager(ctrl)
	usecase := NewItemUsecase(itemRepo, cache, filestorage, logger)
	ctx := gomock.Any()
	param = testId.String()
	paramns := testId
	testItemChan := make(chan models.Item, 1)
	testItemChan <- testItemWithId
	close(testItemChan)

	cache.EXPECT().CheckCache(ctx, param+"nameasc").Return(false)
	itemRepo.EXPECT().GetFavouriteItems(ctx, paramns).Return(testItemChan, nil)
	cache.EXPECT().CreateItemsCache(ctx, items, param+"nameasc").Return(nil)
	cache.EXPECT().CreateItemsQuantityCache(ctx, 1, param+"Quantity").Return(nil)
	cache.EXPECT().GetItemsCache(ctx, param+"nameasc").Return(items, nil)
	res, err := usecase.GetFavouriteItems(context.Background(), paramns, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, items)

	testItemChan2 := make(chan models.Item, 1)
	testItemChan2 <- testItemWithId
	close(testItemChan2)
	cache.EXPECT().CheckCache(ctx, param+"nameasc").Return(false)
	itemRepo.EXPECT().GetFavouriteItems(ctx, paramns).Return(testItemChan2, nil)
	cache.EXPECT().CreateItemsCache(ctx, items, param+"nameasc").Return(fmt.Errorf("error"))
	cache.EXPECT().CreateItemsQuantityCache(ctx, 1, param+"Quantity").Return(fmt.Errorf("error"))
	cache.EXPECT().GetItemsCache(ctx, param+"nameasc").Return(nil, fmt.Errorf("error"))
	itemRepo.EXPECT().GetFavouriteItems(ctx, paramns).Return(testItemChan2, fmt.Errorf("error"))
	res, err = usecase.GetFavouriteItems(context.Background(), paramns, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.Error(t, err)
	require.Nil(t, res)

	cache.EXPECT().CheckCache(ctx, param+"nameasc").Return(true)
	cache.EXPECT().GetItemsCache(ctx, param+"nameasc").Return(items, nil)
	res, err = usecase.GetFavouriteItems(context.Background(), paramns, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, items)

	cache.EXPECT().CheckCache(ctx, param+"nameasc").Return(true)
	cache.EXPECT().GetItemsCache(ctx, param+"nameasc").Return(items2, nil)
	res, err = usecase.GetFavouriteItems(context.Background(), paramns, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, items)

	cache.EXPECT().CheckCache(ctx, param+"nameasc").Return(true)
	cache.EXPECT().GetItemsCache(ctx, param+"nameasc").Return(items, nil)
	res, err = usecase.GetFavouriteItems(context.Background(), paramns, testLimitOptionsItemsList2, testSortOptionsItemsList)
	require.Error(t, err)
	require.Nil(t, res)

	err = fmt.Errorf("error")
	cache.EXPECT().CheckCache(ctx, param+"nameasc").Return(false)
	itemRepo.EXPECT().GetFavouriteItems(ctx, paramns).Return(testItemChan, err)
	res, err = usecase.GetFavouriteItems(context.Background(), paramns, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.Error(t, err)
	require.Nil(t, res)

	testChan2 := make(chan models.Item, 1)
	testChan2 <- testItemWithId
	close(testChan2)

	cache.EXPECT().CheckCache(ctx, param+"nameasc").Return(false)
	itemRepo.EXPECT().GetFavouriteItems(ctx, paramns).Return(testChan2, nil)
	cache.EXPECT().CreateItemsCache(ctx, items, param+"nameasc").Return(fmt.Errorf("error"))
	cache.EXPECT().CreateItemsQuantityCache(ctx, 1, param+"Quantity").Return(fmt.Errorf("error"))
	cache.EXPECT().GetItemsCache(ctx, param+"nameasc").Return(nil, fmt.Errorf("error"))
	itemRepo.EXPECT().GetFavouriteItems(ctx, paramns).Return(testChan2, nil)
	res, err = usecase.GetFavouriteItems(context.Background(), paramns, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.NoError(t, err)
	require.NotNil(t, res)
}

func TestItemsQuantityInFavourite(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cache := mocks.NewMockIItemsCache(ctrl)
	filestorage := mocks.NewMockFileStorager(ctrl)
	usecase := NewItemUsecase(itemRepo, cache, filestorage, logger)
	ctx := gomock.Any()
	testFav := testId.String()
	key := testFav + "Quantity"

	cache.EXPECT().CheckCache(ctx, key).Return(false)
	itemRepo.EXPECT().ItemsInFavouriteQuantity(ctx, testId).Return(-1, fmt.Errorf("error"))
	res, err := usecase.ItemsQuantityInFavourite(context.Background(), testId)
	require.Error(t, err)
	require.Equal(t, res, -1)

	cache.EXPECT().CheckCache(ctx, key).Return(false)
	itemRepo.EXPECT().ItemsInFavouriteQuantity(ctx, testId).Return(1, nil)
	cache.EXPECT().CreateItemsQuantityCache(ctx, 1, key).Return(fmt.Errorf("error"))
	res, err = usecase.ItemsQuantityInFavourite(context.Background(), testId)
	require.NoError(t, err)
	require.Equal(t, res, 1)

	cache.EXPECT().CheckCache(ctx, key).Return(true)
	cache.EXPECT().GetItemsQuantityCache(ctx, key).Return(-1, fmt.Errorf("error"))
	itemRepo.EXPECT().ItemsInFavouriteQuantity(ctx, testId).Return(-1, fmt.Errorf("error"))
	res, err = usecase.ItemsQuantityInFavourite(context.Background(), testId)
	require.Error(t, err)
	require.Equal(t, res, -1)

	cache.EXPECT().CheckCache(ctx, key).Return(true)
	cache.EXPECT().GetItemsQuantityCache(ctx, key).Return(-1, fmt.Errorf("error"))
	itemRepo.EXPECT().ItemsInFavouriteQuantity(ctx, testId).Return(1, nil)
	res, err = usecase.ItemsQuantityInFavourite(context.Background(), testId)
	require.NoError(t, err)
	require.Equal(t, res, 1)

	cache.EXPECT().CheckCache(ctx, key).Return(true)
	cache.EXPECT().GetItemsQuantityCache(ctx, key).Return(1, nil)
	res, err = usecase.ItemsQuantityInFavourite(context.Background(), testId)
	require.NoError(t, err)
	require.Equal(t, res, 1)
}

func TestUpdateFavouriteItemsCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cache := mocks.NewMockIItemsCache(ctrl)
	filestorage := mocks.NewMockFileStorager(ctrl)
	usecase := NewItemUsecase(itemRepo, cache, filestorage, logger)
	ctx := context.Background()

	cacheResults := make([]models.Item, 0, 1)
	cacheResults = append(cacheResults, cacheItem)
	updateResults := make([]models.Item, 0, 2)
	updateResults = append(updateResults, *newItem)
	updateResults = append(updateResults, *newItem)

	cache.EXPECT().CheckCache(ctx, testId.String()+"nameasc").Return(true)
	cache.EXPECT().GetItemsCache(ctx, testId.String()+"nameasc").Return(nil, err)
	usecase.UpdateFavouriteItemsCache(ctx, testId, testItemId, "add")

	cache.EXPECT().CheckCache(ctx, testId.String()+"nameasc").Return(true)
	cache.EXPECT().GetItemsCache(ctx, testId.String()+"nameasc").Return(cacheResults, nil)
	itemRepo.EXPECT().GetItem(ctx, testItemId).Return(nil, err)
	usecase.UpdateFavouriteItemsCache(ctx, testId, testItemId, "add")

	cache.EXPECT().CheckCache(ctx, testId.String()+"nameasc").Return(true)
	cache.EXPECT().GetItemsCache(ctx, testId.String()+"nameasc").Return(cacheResults, nil)
	itemRepo.EXPECT().GetItem(ctx, testItemId).Return(nil, err)
	usecase.UpdateFavouriteItemsCache(ctx, testId, testItemId, "add")

	cache.EXPECT().CheckCache(ctx, testId.String()+"nameasc").Return(true)
	cache.EXPECT().GetItemsCache(ctx, testId.String()+"nameasc").Return(cacheResults, nil)
	itemRepo.EXPECT().GetItem(ctx, testItemId).Return(&testItem1, nil)
	cache.EXPECT().CreateItemsQuantityCache(ctx, 2, testId.String()+"Quantity").Return(err)
	usecase.UpdateFavouriteItemsCache(ctx, testId, testItemId, "add")

	cache.EXPECT().CheckCache(ctx, testId.String()+"nameasc").Return(true)
	cache.EXPECT().GetItemsCache(ctx, testId.String()+"nameasc").Return(cacheResults, nil)
	itemRepo.EXPECT().GetItem(ctx, testItemId).Return(newItem, nil)
	cache.EXPECT().CreateItemsQuantityCache(ctx, 2, testId.String()+"Quantity").Return(nil)
	cache.EXPECT().CreateItemsCache(ctx, updateResults, testId.String()+"nameasc").Return(err)
	usecase.UpdateFavouriteItemsCache(ctx, testId, testItemId, "add")

	cache.EXPECT().CheckCache(ctx, testId.String()+"nameasc").Return(true)
	cache.EXPECT().GetItemsCache(ctx, testId.String()+"nameasc").Return(cacheResults, nil)
	itemRepo.EXPECT().GetItem(ctx, testItemId).Return(newItem, nil)
	cache.EXPECT().CreateItemsQuantityCache(ctx, 2, testId.String()+"Quantity").Return(nil)
	cache.EXPECT().CreateItemsCache(ctx, updateResults, testId.String()+"nameasc").Return(nil)

	cache.EXPECT().GetItemsCache(ctx, testId.String()+"namedesc").Return(cacheResults, nil)
	itemRepo.EXPECT().GetItem(ctx, testItemId).Return(newItem, nil)
	cache.EXPECT().CreateItemsQuantityCache(ctx, 2, testId.String()+"Quantity").Return(nil)
	cache.EXPECT().CreateItemsCache(ctx, updateResults, testId.String()+"namedesc").Return(nil)

	cache.EXPECT().GetItemsCache(ctx, testId.String()+"priceasc").Return(cacheResults, nil)
	itemRepo.EXPECT().GetItem(ctx, testItemId).Return(newItem, nil)
	cache.EXPECT().CreateItemsQuantityCache(ctx, 2, testId.String()+"Quantity").Return(nil)
	cache.EXPECT().CreateItemsCache(ctx, updateResults, testId.String()+"priceasc").Return(nil)

	cache.EXPECT().GetItemsCache(ctx, testId.String()+"pricedesc").Return(cacheResults, nil)
	itemRepo.EXPECT().GetItem(ctx, testItemId).Return(newItem, nil)
	cache.EXPECT().CreateItemsQuantityCache(ctx, 2, testId.String()+"Quantity").Return(nil)
	cache.EXPECT().CreateItemsCache(ctx, updateResults, testId.String()+"pricedesc").Return(nil)
	usecase.UpdateFavouriteItemsCache(ctx, testId, testItemId, "add")

	cache.EXPECT().CheckCache(ctx, testId.String()+"nameasc").Return(true)
	cache.EXPECT().GetItemsCache(ctx, testId.String()+"nameasc").Return(updateResults, nil)
	cache.EXPECT().CreateItemsQuantityCache(ctx, 1, testId.String()+"Quantity").Return(nil)
	cache.EXPECT().CreateItemsCache(ctx, cacheResults, testId.String()+"nameasc").Return(err)
	usecase.UpdateFavouriteItemsCache(ctx, testId, testItemId, "delete")

	cache.EXPECT().CheckCache(ctx, testId.String()+"nameasc").Return(true)
	cache.EXPECT().GetItemsCache(ctx, testId.String()+"nameasc").Return(updateResults, nil)
	cache.EXPECT().CreateItemsQuantityCache(ctx, 1, testId.String()+"Quantity").Return(err)
	usecase.UpdateFavouriteItemsCache(ctx, testId, testItemId, "delete")
}

func TestGetFavouriteItemsId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cache := mocks.NewMockIItemsCache(ctrl)
	filestorage := mocks.NewMockFileStorager(ctrl)
	usecase := NewItemUsecase(itemRepo, cache, filestorage, logger)
	ctx := gomock.Any()

	cache.EXPECT().CheckCache(ctx, testId.String()+"Fav").Return(false)
	cache.EXPECT().CheckCache(ctx, testId.String()+"Quantity").Return(true)
	cache.EXPECT().GetItemsQuantityCache(ctx, testId.String()+"Quantity").Return(-1, err)
	itemRepo.EXPECT().ItemsInFavouriteQuantity(ctx, testId).Return(-1, fmt.Errorf("error"))
	res, err := usecase.GetFavouriteItemsId(context.Background(), testId)
	require.Error(t, err)
	require.Nil(t, res)

	cache.EXPECT().CheckCache(ctx, testId.String()+"Fav").Return(false)
	cache.EXPECT().CheckCache(ctx, testId.String()+"Quantity").Return(true)
	cache.EXPECT().GetItemsQuantityCache(ctx, testId.String()+"Quantity").Return(0, nil)
	res, err = usecase.GetFavouriteItemsId(context.Background(), testId)
	require.Error(t, err)
	require.ErrorIs(t, err, models.ErrorNotFound{})
	require.Nil(t, res)

	cache.EXPECT().CheckCache(ctx, testId.String()+"Fav").Return(false)
	cache.EXPECT().CheckCache(ctx, testId.String()+"Quantity").Return(true)
	cache.EXPECT().GetItemsQuantityCache(ctx, testId.String()+"Quantity").Return(1, nil)
	itemRepo.EXPECT().GetFavouriteItemsId(ctx, testId).Return(nil, models.ErrorNotFound{})
	res, err = usecase.GetFavouriteItemsId(context.Background(), testId)
	require.Error(t, err)
	require.ErrorIs(t, err, models.ErrorNotFound{})
	require.Nil(t, res)

	cache.EXPECT().CheckCache(ctx, testId.String()+"Fav").Return(false)
	cache.EXPECT().CheckCache(ctx, testId.String()+"Quantity").Return(true)
	cache.EXPECT().GetItemsQuantityCache(ctx, testId.String()+"Quantity").Return(1, nil)
	itemRepo.EXPECT().GetFavouriteItemsId(ctx, testId).Return(nil, err)
	res, err = usecase.GetFavouriteItemsId(context.Background(), testId)
	require.Error(t, err)
	require.Nil(t, res)

	cache.EXPECT().CheckCache(ctx, testId.String()+"Fav").Return(false)
	cache.EXPECT().CheckCache(ctx, testId.String()+"Quantity").Return(true)
	cache.EXPECT().GetItemsQuantityCache(ctx, testId.String()+"Quantity").Return(1, nil)
	itemRepo.EXPECT().GetFavouriteItemsId(ctx, testId).Return(&testFavUids, nil)
	cache.EXPECT().CreateFavouriteItemsIdCache(ctx, testFavUids, testId.String()+"Fav").Return(fmt.Errorf("error"))
	cache.EXPECT().GetFavouriteItemsIdCache(ctx, testId.String()+"Fav").Return(nil, fmt.Errorf("error"))
	itemRepo.EXPECT().GetFavouriteItemsId(ctx, testId).Return(&testFavUids, nil)
	res, err = usecase.GetFavouriteItemsId(context.Background(), testId)
	require.NoError(t, err)
	require.Equal(t, res, &testFavUids)

	cache.EXPECT().CheckCache(ctx, testId.String()+"Fav").Return(false)
	cache.EXPECT().CheckCache(ctx, testId.String()+"Quantity").Return(true)
	cache.EXPECT().GetItemsQuantityCache(ctx, testId.String()+"Quantity").Return(1, nil)
	itemRepo.EXPECT().GetFavouriteItemsId(ctx, testId).Return(&testFavUids, nil)
	cache.EXPECT().CreateFavouriteItemsIdCache(ctx, testFavUids, testId.String()+"Fav").Return(fmt.Errorf("error"))
	cache.EXPECT().GetFavouriteItemsIdCache(ctx, testId.String()+"Fav").Return(nil, fmt.Errorf("error"))
	itemRepo.EXPECT().GetFavouriteItemsId(ctx, testId).Return(nil, fmt.Errorf("error"))
	res, err = usecase.GetFavouriteItemsId(context.Background(), testId)
	require.Error(t, err)
	require.Nil(t, res)

	cache.EXPECT().CheckCache(ctx, testId.String()+"Fav").Return(false)
	cache.EXPECT().CheckCache(ctx, testId.String()+"Quantity").Return(true)
	cache.EXPECT().GetItemsQuantityCache(ctx, testId.String()+"Quantity").Return(1, nil)
	itemRepo.EXPECT().GetFavouriteItemsId(ctx, testId).Return(nil, fmt.Errorf("error"))
	res, err = usecase.GetFavouriteItemsId(context.Background(), testId)
	require.Error(t, err)
	require.Nil(t, res)

	cache.EXPECT().CheckCache(ctx, testId.String()+"Fav").Return(false)
	cache.EXPECT().CheckCache(ctx, testId.String()+"Quantity").Return(true)
	cache.EXPECT().GetItemsQuantityCache(ctx, testId.String()+"Quantity").Return(1, nil)
	itemRepo.EXPECT().GetFavouriteItemsId(ctx, testId).Return(&testFavUids, nil)
	cache.EXPECT().CreateFavouriteItemsIdCache(ctx, testFavUids, testId.String()+"Fav").Return(nil)
	cache.EXPECT().GetFavouriteItemsIdCache(ctx, testId.String()+"Fav").Return(nil, fmt.Errorf("error"))
	itemRepo.EXPECT().GetFavouriteItemsId(ctx, testId).Return(nil, models.ErrorNotFound{})
	res, err = usecase.GetFavouriteItemsId(context.Background(), testId)
	require.Error(t, err)
	require.ErrorIs(t, err, models.ErrorNotFound{})
	require.Nil(t, res)
}

func TestUpdateFavIdsCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cache := mocks.NewMockIItemsCache(ctrl)
	filestorage := mocks.NewMockFileStorager(ctrl)
	usecase := NewItemUsecase(itemRepo, cache, filestorage, logger)
	ctx := context.Background()

	cache.EXPECT().CheckCache(ctx, testId.String()+"Fav").Return(false)
	cache.EXPECT().CreateFavouriteItemsIdCache(ctx, testFavUids, testId.String()+"Fav").Return(err)
	usecase.UpdateFavIdsCache(ctx, testId, testItemId, "add")

	cache.EXPECT().CheckCache(ctx, testId.String()+"Fav").Return(false)
	cache.EXPECT().CreateFavouriteItemsIdCache(ctx, testFavUids, testId.String()+"Fav").Return(nil)
	usecase.UpdateFavIdsCache(ctx, testId, testItemId, "add")

	cache.EXPECT().CheckCache(ctx, testId.String()+"Fav").Return(true)
	cache.EXPECT().GetFavouriteItemsIdCache(ctx, testId.String()+"Fav").Return(nil, err)
	usecase.UpdateFavIdsCache(ctx, testId, testItemId, "add")

	cache.EXPECT().CheckCache(ctx, testId.String()+"Fav").Return(true)
	cache.EXPECT().GetFavouriteItemsIdCache(ctx, testId.String()+"Fav").Return(&testFavUids, nil)
	cache.EXPECT().CreateFavouriteItemsIdCache(ctx, testFavUids, testId.String()+"Fav").Return(err)
	usecase.UpdateFavIdsCache(ctx, testId, testItemId, "add")

	cache.EXPECT().CheckCache(ctx, testId.String()+"Fav").Return(true)
	cache.EXPECT().GetFavouriteItemsIdCache(ctx, testId.String()+"Fav").Return(&testFavUids, nil)
	cache.EXPECT().CreateFavouriteItemsIdCache(ctx, testFavUids, testId.String()+"Fav").Return(nil)
	usecase.UpdateFavIdsCache(ctx, testId, testItemId, "delete")
}