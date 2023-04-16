package usecase

import (
	"OnlineShopBackend/internal/models"
	"OnlineShopBackend/internal/repository/mocks"
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var (
	testItemId    = uuid.New()
	testModelItem = models.Item{
		Title:       "test",
		Description: "test",
		Category:    models.Category{},
	}
	testItemWithId = models.Item{
		Id:          testItemId,
		Title:       "test",
		Description: "test",
		Category:    models.Category{},
	}
	testItemWithId2 = models.Item{
		Id:       testItemId,
		Category: models.Category{},
	}
	emptyItem = models.Item{}
	param     = "test"
	items     = []models.Item{testItemWithId}
	items2    = []models.Item{testItemWithId, testItemWithId2}
	newItem   = &models.Item{
		Id:          testItemId,
		Title:       "test",
		Description: "test",
		Category:    models.Category{},
		Price:       0,
		Vendor:      "test",
	}
	cacheItem = models.Item{
		Id:          testItemId,
		Title:       "test",
		Description: "test",
		Category:    models.Category{},
		Price:       0,
		Vendor:      "test",
	}
	testCategoryName          = "testName"
	testSearch                = "testSearch"
	err                       = errors.New("error")
	testLimitOptionsItemsList = map[string]int{
		"offset": 0,
		"limit":  1,
	}
	testLimitOptionsItemsList2 = map[string]int{
		"offset": 2,
		"limit":  1,
	}
	testSortOptionsItemsList = map[string]string{
		"sortType":  "name",
		"sortOrder": "asc",
	}
	testFavUids = map[uuid.UUID]uuid.UUID{
		testItemId: testId,
	}
)

func TestCreateItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cache := mocks.NewMockIItemsCache(ctrl)
	filestorage := mocks.NewMockFileStorager(ctrl)
	usecase := NewItemUsecase(itemRepo, cache, filestorage, logger)
	ctx := context.Background()

	itemRepo.EXPECT().CreateItem(ctx, &testModelItem).Return(uuid.Nil, err)
	res, err := usecase.CreateItem(ctx, &testModelItem)
	require.Error(t, err)
	require.Equal(t, res, uuid.Nil)

	itemRepo.EXPECT().CreateItem(ctx, &testModelItem).Return(testId, nil)
	cache.EXPECT().CheckCache(ctx, itemsListKeyNameAsc).Return(false)
	cache.EXPECT().CheckCache(ctx, itemsListKeyNameDesc).Return(false)
	cache.EXPECT().CheckCache(ctx, itemsListKeyPriceAsc).Return(false)
	cache.EXPECT().CheckCache(ctx, itemsListKeyPriceDesc).Return(false)
	res, err = usecase.CreateItem(ctx, &testModelItem)
	require.NoError(t, err)
	require.Equal(t, res, testId)
}

func TestUpdateItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cache := mocks.NewMockIItemsCache(ctrl)
	filestorage := mocks.NewMockFileStorager(ctrl)
	usecase := NewItemUsecase(itemRepo, cache, filestorage, logger)
	ctx := context.Background()

	itemRepo.EXPECT().UpdateItem(ctx, &testModelItem).Return(err)
	err := usecase.UpdateItem(ctx, &testModelItem)
	require.Error(t, err)

	itemRepo.EXPECT().UpdateItem(ctx, &testModelItem).Return(nil)
	cache.EXPECT().CheckCache(ctx, itemsListKeyNameAsc).Return(false)
	cache.EXPECT().CheckCache(ctx, itemsListKeyNameDesc).Return(false)
	cache.EXPECT().CheckCache(ctx, itemsListKeyPriceAsc).Return(false)
	cache.EXPECT().CheckCache(ctx, itemsListKeyPriceDesc).Return(false)
	err = usecase.UpdateItem(ctx, &testModelItem)
	require.NoError(t, err)
}

func TestGetItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cache := mocks.NewMockIItemsCache(ctrl)
	filestorage := mocks.NewMockFileStorager(ctrl)
	usecase := NewItemUsecase(itemRepo, cache, filestorage, logger)
	ctx := context.Background()

	itemRepo.EXPECT().GetItem(ctx, testItemId).Return(&testItemWithId, nil)
	res, err := usecase.GetItem(ctx, testItemId)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, &testItemWithId)

	err = fmt.Errorf("error on get item")
	itemRepo.EXPECT().GetItem(ctx, testItemId).Return(&emptyItem, err)
	res, err = usecase.GetItem(ctx, testItemId)
	require.Error(t, err)
	require.Nil(t, res)
}

func TestItemsList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cache := mocks.NewMockIItemsCache(ctrl)
	filestorage := mocks.NewMockFileStorager(ctrl)
	usecase := NewItemUsecase(itemRepo, cache, filestorage, logger)
	ctx := gomock.Any()
	testItemChan := make(chan models.Item, 1)
	testItemChan <- testItemWithId
	close(testItemChan)

	cache.EXPECT().CheckCache(ctx, itemsListKey+"nameasc").Return(false)
	itemRepo.EXPECT().ItemsList(ctx).Return(testItemChan, nil)
	cache.EXPECT().CreateItemsCache(ctx, items, itemsListKey+"nameasc").Return(nil)
	cache.EXPECT().CreateItemsQuantityCache(ctx, len(items), itemsQuantityKey).Return(nil)
	cache.EXPECT().GetItemsCache(ctx, itemsListKey+"nameasc").Return(items, nil)

	res, err := usecase.ItemsList(context.Background(), testLimitOptionsItemsList, testSortOptionsItemsList)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, items)

	cache.EXPECT().CheckCache(ctx, itemsListKey+"nameasc").Return(true)
	cache.EXPECT().GetItemsCache(ctx, itemsListKey+"nameasc").Return(items, nil)
	res, err = usecase.ItemsList(context.Background(), testLimitOptionsItemsList, testSortOptionsItemsList)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, items)

	cache.EXPECT().CheckCache(ctx, itemsListKey+"nameasc").Return(true)
	cache.EXPECT().GetItemsCache(ctx, itemsListKey+"nameasc").Return(items2, nil)
	res, err = usecase.ItemsList(context.Background(), testLimitOptionsItemsList, testSortOptionsItemsList)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, items)

	cache.EXPECT().CheckCache(ctx, itemsListKey+"nameasc").Return(true)
	cache.EXPECT().GetItemsCache(ctx, itemsListKey+"nameasc").Return(items, nil)
	res, err = usecase.ItemsList(context.Background(), testLimitOptionsItemsList2, testSortOptionsItemsList)
	require.Error(t, err)
	require.Nil(t, res)

	err = fmt.Errorf("error on itemslist()")
	cache.EXPECT().CheckCache(ctx, itemsListKey+"nameasc").Return(false)
	itemRepo.EXPECT().ItemsList(ctx).Return(testItemChan, err)
	res, err = usecase.ItemsList(context.Background(), testLimitOptionsItemsList, testSortOptionsItemsList)
	require.Error(t, err)
	require.Nil(t, res)

	testChan2 := make(chan models.Item, 1)
	testChan2 <- testItemWithId
	close(testChan2)

	cache.EXPECT().CheckCache(ctx, itemsListKey+"nameasc").Return(false)
	itemRepo.EXPECT().ItemsList(ctx).Return(testChan2, nil)
	cache.EXPECT().CreateItemsCache(ctx, items, itemsListKey+"nameasc").Return(err)
	cache.EXPECT().CreateItemsQuantityCache(ctx, len(items), itemsQuantityKey).Return(err)
	cache.EXPECT().GetItemsCache(ctx, itemsListKeyNameAsc).Return(items, nil)
	res, err = usecase.ItemsList(context.Background(), testLimitOptionsItemsList, testSortOptionsItemsList)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, items)

	testChan3 := make(chan models.Item, 1)
	testChan3 <- testItemWithId
	close(testChan3)

	cache.EXPECT().CheckCache(ctx, itemsListKey+"nameasc").Return(false)
	itemRepo.EXPECT().ItemsList(ctx).Return(testChan3, nil)
	cache.EXPECT().CreateItemsCache(ctx, items, itemsListKey+"nameasc").Return(err)
	cache.EXPECT().CreateItemsQuantityCache(ctx, len(items), itemsQuantityKey).Return(err)
	cache.EXPECT().GetItemsCache(ctx, itemsListKeyNameAsc).Return(nil, fmt.Errorf("error"))
	itemRepo.EXPECT().ItemsList(ctx).Return(testChan3, fmt.Errorf("error"))
	res, err = usecase.ItemsList(context.Background(), testLimitOptionsItemsList, testSortOptionsItemsList)
	require.Error(t, err)
	require.Nil(t, res)

	testChan4 := make(chan models.Item, 1)
	testChan4 <- testItemWithId
	close(testChan4)

	cache.EXPECT().CheckCache(ctx, itemsListKey+"nameasc").Return(false)
	itemRepo.EXPECT().ItemsList(ctx).Return(testChan4, nil)
	cache.EXPECT().CreateItemsCache(ctx, items, itemsListKey+"nameasc").Return(err)
	cache.EXPECT().CreateItemsQuantityCache(ctx, len(items), itemsQuantityKey).Return(err)
	cache.EXPECT().GetItemsCache(ctx, itemsListKeyNameAsc).Return(nil, fmt.Errorf("error"))
	itemRepo.EXPECT().ItemsList(ctx).Return(testChan4, nil)
	res, err = usecase.ItemsList(context.Background(), testLimitOptionsItemsList, testSortOptionsItemsList)
	require.NoError(t, err)
	require.NotNil(t, res)
}

func TestSearchLine(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cache := mocks.NewMockIItemsCache(ctrl)
	filestorage := mocks.NewMockFileStorager(ctrl)
	usecase := NewItemUsecase(itemRepo, cache, filestorage, logger)
	ctx := gomock.Any()

	testItemChan := make(chan models.Item, 1)
	testItemChan <- testItemWithId
	close(testItemChan)

	cache.EXPECT().CheckCache(ctx, param+"nameasc").Return(false)
	itemRepo.EXPECT().SearchLine(ctx, param).Return(testItemChan, nil)
	cache.EXPECT().CreateItemsQuantityCache(ctx, 1, param+"Quantity").Return(nil)
	cache.EXPECT().CreateItemsCache(ctx, items, param+"nameasc").Return(nil)
	cache.EXPECT().GetItemsCache(ctx, param+"nameasc").Return(items, nil)
	res, err := usecase.SearchLine(context.Background(), param, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, items)

	cache.EXPECT().CheckCache(ctx, param+"nameasc").Return(true)
	cache.EXPECT().GetItemsCache(ctx, param+"nameasc").Return(items, nil)
	res, err = usecase.SearchLine(context.Background(), param, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, items)

	cache.EXPECT().CheckCache(ctx, param+"nameasc").Return(true)
	cache.EXPECT().GetItemsCache(ctx, param+"nameasc").Return(items2, nil)
	res, err = usecase.SearchLine(context.Background(), param, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, items)

	cache.EXPECT().CheckCache(ctx, param+"nameasc").Return(true)
	cache.EXPECT().GetItemsCache(ctx, param+"nameasc").Return(items, nil)
	res, err = usecase.SearchLine(context.Background(), param, testLimitOptionsItemsList2, testSortOptionsItemsList)
	require.Error(t, err)
	require.Nil(t, res)

	err = fmt.Errorf("error on search()")
	cache.EXPECT().CheckCache(ctx, param+"nameasc").Return(false)
	itemRepo.EXPECT().SearchLine(ctx, param).Return(testItemChan, err)
	res, err = usecase.SearchLine(context.Background(), param, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.Error(t, err)
	require.Nil(t, res)

	testChan2 := make(chan models.Item, 1)
	testChan2 <- testItemWithId
	close(testChan2)

	cache.EXPECT().CheckCache(ctx, param+"nameasc").Return(false)
	itemRepo.EXPECT().SearchLine(ctx, param).Return(testChan2, nil)
	cache.EXPECT().CreateItemsQuantityCache(ctx, 1, param+"Quantity").Return(err)
	cache.EXPECT().CreateItemsCache(ctx, items, param+"nameasc").Return(err)
	cache.EXPECT().GetItemsCache(ctx, param+"nameasc").Return(items, nil)
	res, err = usecase.SearchLine(context.Background(), param, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, items)

	testChan3 := make(chan models.Item, 1)
	testChan3 <- testItemWithId
	close(testChan3)

	cache.EXPECT().CheckCache(ctx, param+"nameasc").Return(false)
	itemRepo.EXPECT().SearchLine(ctx, param).Return(testChan3, nil)
	cache.EXPECT().CreateItemsQuantityCache(ctx, 1, param+"Quantity").Return(nil)
	cache.EXPECT().CreateItemsCache(ctx, items, param+"nameasc").Return(nil)
	cache.EXPECT().GetItemsCache(ctx, param+"nameasc").Return(nil, fmt.Errorf("error"))
	itemRepo.EXPECT().SearchLine(ctx, param).Return(testChan3, fmt.Errorf("error"))
	res, err = usecase.SearchLine(context.Background(), param, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.Error(t, err)
	require.Nil(t, res)

	testChan4 := make(chan models.Item, 1)
	testChan4 <- testItemWithId
	close(testChan4)

	cache.EXPECT().CheckCache(ctx, param+"nameasc").Return(false)
	itemRepo.EXPECT().SearchLine(ctx, param).Return(testChan4, nil)
	cache.EXPECT().CreateItemsQuantityCache(ctx, 1, param+"Quantity").Return(nil)
	cache.EXPECT().CreateItemsCache(ctx, items, param+"nameasc").Return(nil)
	cache.EXPECT().GetItemsCache(ctx, param+"nameasc").Return(nil, fmt.Errorf("error"))
	itemRepo.EXPECT().SearchLine(ctx, param).Return(testChan4, nil)
	res, err = usecase.SearchLine(context.Background(), param, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.NoError(t, err)
	require.NotNil(t, res)
}

func TestGetItemsByCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cache := mocks.NewMockIItemsCache(ctrl)
	filestorage := mocks.NewMockFileStorager(ctrl)
	usecase := NewItemUsecase(itemRepo, cache, filestorage, logger)
	ctx := gomock.Any()

	testItemChan := make(chan models.Item, 1)
	testItemChan <- testItemWithId
	close(testItemChan)

	cache.EXPECT().CheckCache(ctx, param+"nameasc").Return(false)
	itemRepo.EXPECT().GetItemsByCategory(ctx, param).Return(testItemChan, nil)
	cache.EXPECT().CreateItemsCache(ctx, items, param+"nameasc").Return(nil)
	cache.EXPECT().CreateItemsQuantityCache(ctx, 1, param+"Quantity")
	cache.EXPECT().GetItemsCache(ctx, param+"nameasc").Return(items, nil)
	res, err := usecase.GetItemsByCategory(context.Background(), param, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, items)

	testItemChan = make(chan models.Item, 1)
	testItemChan <- testItemWithId
	close(testItemChan)
	cache.EXPECT().CheckCache(ctx, param+"nameasc").Return(false)
	itemRepo.EXPECT().GetItemsByCategory(ctx, param).Return(testItemChan, nil)
	cache.EXPECT().CreateItemsCache(ctx, items, param+"nameasc").Return(nil)
	cache.EXPECT().CreateItemsQuantityCache(ctx, 1, param+"Quantity").Return(fmt.Errorf("error"))
	cache.EXPECT().GetItemsCache(ctx, param+"nameasc").Return(nil, fmt.Errorf("error"))
	itemRepo.EXPECT().GetItemsByCategory(ctx, param).Return(nil, fmt.Errorf("error"))
	res, err = usecase.GetItemsByCategory(context.Background(), param, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.Error(t, err)
	require.Nil(t, res)

	cache.EXPECT().CheckCache(ctx, param+"nameasc").Return(true)
	cache.EXPECT().GetItemsCache(ctx, param+"nameasc").Return(items, nil)
	res, err = usecase.GetItemsByCategory(context.Background(), param, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, items)

	cache.EXPECT().CheckCache(ctx, param+"nameasc").Return(true)
	cache.EXPECT().GetItemsCache(ctx, param+"nameasc").Return(items2, nil)
	res, err = usecase.GetItemsByCategory(context.Background(), param, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, items)

	cache.EXPECT().CheckCache(ctx, param+"nameasc").Return(true)
	cache.EXPECT().GetItemsCache(ctx, param+"nameasc").Return(items, nil)
	res, err = usecase.GetItemsByCategory(context.Background(), param, testLimitOptionsItemsList2, testSortOptionsItemsList)
	require.Error(t, err)
	require.Nil(t, res)

	err = fmt.Errorf("error")
	cache.EXPECT().CheckCache(ctx, param+"nameasc").Return(false)
	itemRepo.EXPECT().GetItemsByCategory(ctx, param).Return(testItemChan, err)
	res, err = usecase.GetItemsByCategory(context.Background(), param, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.Error(t, err)
	require.Nil(t, res)

	testChan2 := make(chan models.Item, 1)
	testChan2 <- testItemWithId
	close(testChan2)

	cache.EXPECT().CheckCache(ctx, param+"nameasc").Return(false)
	itemRepo.EXPECT().GetItemsByCategory(ctx, param).Return(testChan2, nil)
	cache.EXPECT().CreateItemsCache(ctx, items, param+"nameasc").Return(nil)
	cache.EXPECT().CreateItemsQuantityCache(ctx, 1, param+"Quantity").Return(fmt.Errorf("error"))
	cache.EXPECT().GetItemsCache(ctx, param+"nameasc").Return(nil, fmt.Errorf("error"))
	itemRepo.EXPECT().GetItemsByCategory(ctx, param).Return(testChan2, nil)
	res, err = usecase.GetItemsByCategory(context.Background(), param, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.NoError(t, err)
	require.NotNil(t, res)

	testChan3 := make(chan models.Item, 1)
	testChan3 <- testItemWithId
	close(testChan3)

	cache.EXPECT().CheckCache(ctx, param+"nameasc").Return(false)
	itemRepo.EXPECT().GetItemsByCategory(ctx, param).Return(testChan3, nil)
	cache.EXPECT().CreateItemsCache(ctx, items, param+"nameasc").Return(fmt.Errorf("error"))
	cache.EXPECT().CreateItemsQuantityCache(ctx, 1, param+"Quantity").Return(fmt.Errorf("error"))
	cache.EXPECT().GetItemsCache(ctx, param+"nameasc").Return(nil, fmt.Errorf("error"))
	itemRepo.EXPECT().GetItemsByCategory(ctx, param).Return(testChan3, nil)
	res, err = usecase.GetItemsByCategory(context.Background(), param, testLimitOptionsItemsList, testSortOptionsItemsList)
	require.NoError(t, err)
	require.NotNil(t, res)

}

func TestItemsQuantity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cache := mocks.NewMockIItemsCache(ctrl)
	filestorage := mocks.NewMockFileStorager(ctrl)
	usecase := NewItemUsecase(itemRepo, cache, filestorage, logger)
	ctx := gomock.Any()

	cache.EXPECT().CheckCache(ctx, itemsQuantityKey).Return(false)
	itemRepo.EXPECT().ItemsListQuantity(ctx).Return(-1, fmt.Errorf("error"))
	res, err := usecase.ItemsQuantity(context.Background())
	require.Error(t, err)
	require.Equal(t, res, -1)

	cache.EXPECT().CheckCache(ctx, itemsQuantityKey).Return(false)
	itemRepo.EXPECT().ItemsListQuantity(ctx).Return(1, nil)
	cache.EXPECT().CreateItemsQuantityCache(ctx, 1, itemsQuantityKey).Return(fmt.Errorf("error"))
	res, err = usecase.ItemsQuantity(context.Background())
	require.NoError(t, err)
	require.Equal(t, res, 1)

	cache.EXPECT().CheckCache(ctx, itemsQuantityKey).Return(true)
	cache.EXPECT().GetItemsQuantityCache(ctx, itemsQuantityKey).Return(-1, fmt.Errorf("error"))
	itemRepo.EXPECT().ItemsListQuantity(ctx).Return(-1, fmt.Errorf("error"))
	res, err = usecase.ItemsQuantity(context.Background())
	require.Error(t, err)
	require.Equal(t, res, -1)

	cache.EXPECT().CheckCache(ctx, itemsQuantityKey).Return(true)
	cache.EXPECT().GetItemsQuantityCache(ctx, itemsQuantityKey).Return(-1, fmt.Errorf("error"))
	itemRepo.EXPECT().ItemsListQuantity(ctx).Return(1, nil)
	res, err = usecase.ItemsQuantity(context.Background())
	require.NoError(t, err)
	require.Equal(t, res, 1)

	cache.EXPECT().CheckCache(ctx, itemsQuantityKey).Return(true)
	cache.EXPECT().GetItemsQuantityCache(ctx, itemsQuantityKey).Return(1, nil)
	res, err = usecase.ItemsQuantity(context.Background())
	require.NoError(t, err)
	require.Equal(t, res, 1)
}

func TestItemsQuantityInCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cache := mocks.NewMockIItemsCache(ctrl)
	filestorage := mocks.NewMockFileStorager(ctrl)
	usecase := NewItemUsecase(itemRepo, cache, filestorage, logger)
	ctx := gomock.Any()
	key := testCategoryName + "Quantity"

	cache.EXPECT().CheckCache(ctx, key).Return(false)
	itemRepo.EXPECT().ItemsByCategoryQuantity(ctx, testCategoryName).Return(-1, fmt.Errorf("error"))
	res, err := usecase.ItemsQuantityInCategory(context.Background(), testCategoryName)
	require.Error(t, err)
	require.Equal(t, res, -1)

	cache.EXPECT().CheckCache(ctx, key).Return(false)
	itemRepo.EXPECT().ItemsByCategoryQuantity(ctx, testCategoryName).Return(1, nil)
	cache.EXPECT().CreateItemsQuantityCache(ctx, 1, key).Return(fmt.Errorf("error"))
	res, err = usecase.ItemsQuantityInCategory(context.Background(), testCategoryName)
	require.NoError(t, err)
	require.Equal(t, res, 1)

	cache.EXPECT().CheckCache(ctx, key).Return(true)
	cache.EXPECT().GetItemsQuantityCache(ctx, key).Return(-1, fmt.Errorf("error"))
	itemRepo.EXPECT().ItemsByCategoryQuantity(ctx, testCategoryName).Return(-1, fmt.Errorf("error"))
	res, err = usecase.ItemsQuantityInCategory(context.Background(), testCategoryName)
	require.Error(t, err)
	require.Equal(t, res, -1)

	cache.EXPECT().CheckCache(ctx, key).Return(true)
	cache.EXPECT().GetItemsQuantityCache(ctx, key).Return(-1, fmt.Errorf("error"))
	itemRepo.EXPECT().ItemsByCategoryQuantity(ctx, testCategoryName).Return(1, nil)
	res, err = usecase.ItemsQuantityInCategory(context.Background(), testCategoryName)
	require.NoError(t, err)
	require.Equal(t, res, 1)

	cache.EXPECT().CheckCache(ctx, key).Return(true)
	cache.EXPECT().GetItemsQuantityCache(ctx, key).Return(1, nil)
	res, err = usecase.ItemsQuantityInCategory(context.Background(), testCategoryName)
	require.NoError(t, err)
	require.Equal(t, res, 1)
}

func TestItemsQuantityInSearch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cache := mocks.NewMockIItemsCache(ctrl)
	filestorage := mocks.NewMockFileStorager(ctrl)
	usecase := NewItemUsecase(itemRepo, cache, filestorage, logger)
	ctx := gomock.Any()
	key := testSearch + "Quantity"

	cache.EXPECT().CheckCache(ctx, key).Return(false)
	itemRepo.EXPECT().ItemsInSearchQuantity(ctx, testSearch).Return(-1, fmt.Errorf("error"))
	res, err := usecase.ItemsQuantityInSearch(context.Background(), testSearch)
	require.Error(t, err)
	require.Equal(t, res, -1)

	cache.EXPECT().CheckCache(ctx, key).Return(false)
	itemRepo.EXPECT().ItemsInSearchQuantity(ctx, testSearch).Return(1, nil)
	cache.EXPECT().CreateItemsQuantityCache(ctx, 1, key).Return(fmt.Errorf("error"))
	res, err = usecase.ItemsQuantityInSearch(context.Background(), testSearch)
	require.NoError(t, err)
	require.Equal(t, res, 1)

	cache.EXPECT().CheckCache(ctx, key).Return(true)
	cache.EXPECT().GetItemsQuantityCache(ctx, key).Return(-1, fmt.Errorf("error"))
	itemRepo.EXPECT().ItemsInSearchQuantity(ctx, testSearch).Return(-1, fmt.Errorf("error"))
	res, err = usecase.ItemsQuantityInSearch(context.Background(), testSearch)
	require.Error(t, err)
	require.Equal(t, res, -1)

	cache.EXPECT().CheckCache(ctx, key).Return(true)
	cache.EXPECT().GetItemsQuantityCache(ctx, key).Return(-1, fmt.Errorf("error"))
	itemRepo.EXPECT().ItemsInSearchQuantity(ctx, testSearch).Return(1, nil)
	res, err = usecase.ItemsQuantityInSearch(context.Background(), testSearch)
	require.NoError(t, err)
	require.Equal(t, res, 1)

	cache.EXPECT().CheckCache(ctx, key).Return(true)
	cache.EXPECT().GetItemsQuantityCache(ctx, key).Return(1, nil)
	res, err = usecase.ItemsQuantityInSearch(context.Background(), testSearch)
	require.NoError(t, err)
	require.Equal(t, res, 1)
}

func TestUpdateCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cache := mocks.NewMockIItemsCache(ctrl)
	filestorage := mocks.NewMockFileStorager(ctrl)
	usecase := NewItemUsecase(itemRepo, cache, filestorage, logger)
	ctx := context.Background()

	cache.EXPECT().CheckCache(ctx, "ItemsList"+"nameasc").Return(false)
	cache.EXPECT().CheckCache(ctx, "ItemsList"+"namedesc").Return(false)
	cache.EXPECT().CheckCache(ctx, "ItemsList"+"priceasc").Return(false)
	cache.EXPECT().CheckCache(ctx, "ItemsList"+"pricedesc").Return(false)
	err := usecase.UpdateCache(ctx, uuid.New(), "create")
	require.Error(t, err)

	cache.EXPECT().CheckCache(ctx, itemsListKeyNameAsc).Return(true)
	cache.EXPECT().GetItemsCache(ctx, itemsListKeyNameAsc).Return(nil, err)
	err = usecase.UpdateCache(ctx, testId, "create")
	require.Error(t, err)

	cacheResults := make([]models.Item, 0, 1)
	cacheResults = append(cacheResults, cacheItem)
	cache.EXPECT().CheckCache(ctx, itemsListKeyNameAsc).Return(true)
	cache.EXPECT().GetItemsCache(ctx, itemsListKeyNameAsc).Return(cacheResults, nil)
	itemRepo.EXPECT().GetItem(ctx, testId).Return(nil, err)
	err = usecase.UpdateCache(ctx, testId, "create")
	require.Error(t, err)

	cacheResults = make([]models.Item, 0, 1)
	cacheResults = append(cacheResults, cacheItem)
	cache.EXPECT().CheckCache(ctx, itemsListKeyNameAsc).Return(true)
	cache.EXPECT().GetItemsCache(ctx, itemsListKeyNameAsc).Return(cacheResults, nil)
	itemRepo.EXPECT().GetItem(ctx, testId).Return(&testItemWithId, nil)
	cache.EXPECT().CreateItemsCache(ctx, cacheResults, itemsListKeyNameAsc).Return(err)
	err = usecase.UpdateCache(ctx, testId, "update")
	require.Error(t, err)

	cacheResults = make([]models.Item, 0, 1)
	cacheResults = append(cacheResults, cacheItem)
	updateResults := append(cacheResults, cacheItem)
	cache.EXPECT().CheckCache(ctx, itemsListKeyNameAsc).Return(true)
	cache.EXPECT().GetItemsCache(ctx, itemsListKeyNameAsc).Return(cacheResults, nil)
	itemRepo.EXPECT().GetItem(ctx, testId).Return(&testItemWithId, nil)
	cache.EXPECT().CreateItemsQuantityCache(ctx, len(updateResults), itemsQuantityKey).Return(err)
	err = usecase.UpdateCache(ctx, testId, "create")
	require.Error(t, err)
}

func TestUpdateItemsInCategoryCache(t *testing.T) {
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

	cache.EXPECT().CheckCache(ctx, newItem.Category.Name+"nameasc").Return(false)
	cache.EXPECT().CheckCache(ctx, newItem.Category.Name+"namedesc").Return(false)
	cache.EXPECT().CheckCache(ctx, newItem.Category.Name+"priceasc").Return(false)
	cache.EXPECT().CheckCache(ctx, newItem.Category.Name+"pricedesc").Return(false)
	err := usecase.UpdateItemsInCategoryCache(ctx, newItem, "create")
	require.Error(t, err)

	cache.EXPECT().CheckCache(ctx, newItem.Category.Name+"nameasc").Return(true)
	cache.EXPECT().GetItemsCache(ctx, newItem.Category.Name+"nameasc").Return(nil, fmt.Errorf("error"))
	err = usecase.UpdateItemsInCategoryCache(ctx, newItem, "create")
	require.Error(t, err)

	cache.EXPECT().CheckCache(ctx, newItem.Category.Name+"nameasc").Return(true)
	cache.EXPECT().GetItemsCache(ctx, newItem.Category.Name+"nameasc").Return(cacheResults, nil)
	cache.EXPECT().CreateItemsQuantityCache(ctx, len(updateResults), newItem.Category.Name+"Quantity").Return(fmt.Errorf("error"))
	err = usecase.UpdateItemsInCategoryCache(ctx, newItem, "create")
	require.Error(t, err)

	cache.EXPECT().CheckCache(ctx, newItem.Category.Name+"nameasc").Return(true)
	cache.EXPECT().GetItemsCache(ctx, newItem.Category.Name+"nameasc").Return(cacheResults, nil)
	cache.EXPECT().CreateItemsQuantityCache(ctx, len(updateResults), newItem.Category.Name+"Quantity").Return(nil)
	cache.EXPECT().CreateItemsCache(ctx, updateResults, newItem.Category.Name+"nameasc").Return(fmt.Errorf("error"))
	err = usecase.UpdateItemsInCategoryCache(ctx, newItem, "create")
	require.Error(t, err)

	cache.EXPECT().CheckCache(ctx, newItem.Category.Name+"nameasc").Return(true)
	cache.EXPECT().GetItemsCache(ctx, newItem.Category.Name+"nameasc").Return(cacheResults, nil)
	cache.EXPECT().CreateItemsCache(ctx, cacheResults, newItem.Category.Name+"nameasc").Return(nil)

	cache.EXPECT().GetItemsCache(ctx, newItem.Category.Name+"namedesc").Return(cacheResults, nil)
	cache.EXPECT().CreateItemsCache(ctx, cacheResults, newItem.Category.Name+"namedesc").Return(nil)

	cache.EXPECT().GetItemsCache(ctx, newItem.Category.Name+"priceasc").Return(cacheResults, nil)
	cache.EXPECT().CreateItemsCache(ctx, cacheResults, newItem.Category.Name+"priceasc").Return(nil)

	cache.EXPECT().GetItemsCache(ctx, newItem.Category.Name+"pricedesc").Return(cacheResults, nil)
	cache.EXPECT().CreateItemsCache(ctx, cacheResults, newItem.Category.Name+"pricedesc").Return(nil)

	err = usecase.UpdateItemsInCategoryCache(ctx, newItem, "update")
	require.NoError(t, err)

	deletedResults := []models.Item{testItemWithId}
	deleteResults := []models.Item{}
	cache.EXPECT().CheckCache(ctx, newItem.Category.Name+"nameasc").Return(true)
	cache.EXPECT().GetItemsCache(ctx, newItem.Category.Name+"nameasc").Return(deletedResults, nil)
	cache.EXPECT().CreateItemsQuantityCache(ctx, len(deleteResults), testItemWithId.Category.Name+"Quantity").Return(err)
	cache.EXPECT().CreateItemsCache(ctx, deleteResults, newItem.Category.Name+"nameasc").Return(nil)

	cache.EXPECT().GetItemsCache(ctx, newItem.Category.Name+"namedesc").Return(deletedResults, nil)
	cache.EXPECT().CreateItemsQuantityCache(ctx, len(deleteResults), testItemWithId.Category.Name+"Quantity").Return(err)
	cache.EXPECT().CreateItemsCache(ctx, deleteResults, newItem.Category.Name+"namedesc").Return(nil)

	cache.EXPECT().GetItemsCache(ctx, newItem.Category.Name+"priceasc").Return(deletedResults, nil)
	cache.EXPECT().CreateItemsQuantityCache(ctx, len(deleteResults), testItemWithId.Category.Name+"Quantity").Return(err)
	cache.EXPECT().CreateItemsCache(ctx, deleteResults, newItem.Category.Name+"priceasc").Return(nil)

	cache.EXPECT().GetItemsCache(ctx, newItem.Category.Name+"pricedesc").Return(deletedResults, nil)
	cache.EXPECT().CreateItemsQuantityCache(ctx, len(deleteResults), testItemWithId.Category.Name+"Quantity").Return(err)
	cache.EXPECT().CreateItemsCache(ctx, deleteResults, newItem.Category.Name+"pricedesc").Return(nil)
	err = usecase.UpdateItemsInCategoryCache(ctx, newItem, "delete")
	require.NoError(t, err)
}

func TestDeleteItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cache := mocks.NewMockIItemsCache(ctrl)
	filestorage := mocks.NewMockFileStorager(ctrl)
	usecase := NewItemUsecase(itemRepo, cache, filestorage, logger)
	ctx := context.Background()

	itemRepo.EXPECT().DeleteItem(ctx, testId).Return(err)
	err := usecase.DeleteItem(ctx, testId)
	require.Error(t, err)

	itemRepo.EXPECT().DeleteItem(ctx, testId).Return(nil)
	cache.EXPECT().CheckCache(ctx, itemsListKeyNameAsc).Return(false)
	cache.EXPECT().CheckCache(ctx, itemsListKeyNameDesc).Return(false)
	cache.EXPECT().CheckCache(ctx, itemsListKeyPriceAsc).Return(false)
	cache.EXPECT().CheckCache(ctx, itemsListKeyPriceDesc).Return(false)
	err = usecase.DeleteItem(ctx, testId)
	require.NoError(t, err)
}

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

func TestSortItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	itemRepo := mocks.NewMockItemStore(ctrl)
	cache := mocks.NewMockIItemsCache(ctrl)
	filestorage := mocks.NewMockFileStorager(ctrl)
	usecase := NewItemUsecase(itemRepo, cache, filestorage, logger)

	testItems := []models.Item{
		{Title: "A"},
		{Title: "C"},
		{Title: "B"},
	}
	testItems2 := []models.Item{
		{Price: 10},
		{Price: 30},
		{Price: 20},
	}

	usecase.SortItems(testItems, "name", "asc")
	require.Equal(t, testItems, []models.Item{
		{Title: "A"},
		{Title: "B"},
		{Title: "C"},
	})
	usecase.SortItems(testItems, "name", "desc")
	require.Equal(t, testItems, []models.Item{
		{Title: "C"},
		{Title: "B"},
		{Title: "A"},
	})
	usecase.SortItems(testItems2, "price", "asc")
	require.Equal(t, testItems2, []models.Item{
		{Price: 10},
		{Price: 20},
		{Price: 30},
	})
	usecase.SortItems(testItems2, "price", "desc")
	require.Equal(t, testItems2, []models.Item{
		{Price: 30},
		{Price: 20},
		{Price: 10},
	})
	usecase.SortItems(testItems, "pricee", "desc")
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
