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
	testModelCategory = &models.Category{
		Name: "test name",
	}
	testModelCategoryWithId = &models.Category{
		Id:   testId,
		Name: "test name",
	}
	emptyCategory = &models.Category{}
	testId        = uuid.New()
	testChan      = make(chan models.Category, 2)

	categories = []models.Category{*testModelCategoryWithId}
)

func TestCreateCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	categoryRepo := mocks.NewMockCategoryStore(ctrl)
	cache := mocks.NewMockICategoriesCache(ctrl)
	filestorage := mocks.NewMockFileStorager(ctrl)
	usecase := NewCategoryUsecase(categoryRepo, cache, filestorage, logger)

	categoryRepo.EXPECT().CreateCategory(ctx, testModelCategory).Return(testId, nil)
	cache.EXPECT().CheckCache(ctx, categoriesListKey).Return(false)
	res, err := usecase.CreateCategory(ctx, testModelCategory)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, testId)

	emptyCategories := make([]models.Category, 0)
	categories := []models.Category{*testModelCategoryWithId}
	categoryRepo.EXPECT().CreateCategory(ctx, testModelCategory).Return(testId, nil)
	cache.EXPECT().CheckCache(ctx, categoriesListKey).Return(true)
	categoryRepo.EXPECT().GetCategory(ctx, testId).Return(testModelCategoryWithId, nil)
	cache.EXPECT().GetCategoriesListCache(ctx, categoriesListKey).Return(emptyCategories, nil)
	cache.EXPECT().CreateCategoriesListСache(ctx, categories, categoriesListKey).Return(nil)
	res, err = usecase.CreateCategory(ctx, testModelCategory)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, testId)

	err = fmt.Errorf("error on create category")
	categoryRepo.EXPECT().CreateCategory(ctx, testModelCategory).Return(uuid.Nil, err)
	res, err = usecase.CreateCategory(ctx, testModelCategory)
	require.Error(t, err)
	require.Equal(t, res, uuid.Nil)
}

func TestUpdateCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	categoryRepo := mocks.NewMockCategoryStore(ctrl)
	cache := mocks.NewMockICategoriesCache(ctrl)
	filestorage := mocks.NewMockFileStorager(ctrl)
	usecase := NewCategoryUsecase(categoryRepo, cache, filestorage, logger)

	categoryRepo.EXPECT().UpdateCategory(ctx, testModelCategoryWithId).Return(nil)
	cache.EXPECT().CheckCache(ctx, categoriesListKey).Return(false)
	err := usecase.UpdateCategory(ctx, testModelCategoryWithId)
	require.NoError(t, err)

	categories := []models.Category{*testModelCategoryWithId}
	categoryRepo.EXPECT().UpdateCategory(ctx, testModelCategoryWithId).Return(nil)
	cache.EXPECT().CheckCache(ctx, categoriesListKey).Return(true)
	categoryRepo.EXPECT().GetCategory(ctx, testId).Return(testModelCategoryWithId, nil)
	cache.EXPECT().GetCategoriesListCache(ctx, categoriesListKey).Return(categories, nil)
	cache.EXPECT().CreateCategoriesListСache(ctx, categories, categoriesListKey).Return(nil)
	err = usecase.UpdateCategory(ctx, testModelCategoryWithId)
	require.NoError(t, err)

	categoryRepo.EXPECT().UpdateCategory(ctx, testModelCategoryWithId).Return(fmt.Errorf("error on update"))
	err = usecase.UpdateCategory(ctx, testModelCategoryWithId)
	require.Error(t, err)
}

func TestGetCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	categoryRepo := mocks.NewMockCategoryStore(ctrl)
	cache := mocks.NewMockICategoriesCache(ctrl)
	filestorage := mocks.NewMockFileStorager(ctrl)
	usecase := NewCategoryUsecase(categoryRepo, cache, filestorage, logger)

	categoryRepo.EXPECT().GetCategory(ctx, testId).Return(testModelCategoryWithId, nil)
	res, err := usecase.GetCategory(ctx, testId)
	require.NoError(t, err)
	require.Equal(t, res, testModelCategoryWithId)

	categoryRepo.EXPECT().GetCategory(ctx, testId).Return(emptyCategory, fmt.Errorf("error on get category"))
	res, err = usecase.GetCategory(ctx, testId)
	require.Error(t, err)
	require.Equal(t, res, emptyCategory)
}

func TestGetCategoryList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := gomock.Any()
	logger := zap.L()
	categoryRepo := mocks.NewMockCategoryStore(ctrl)
	cache := mocks.NewMockICategoriesCache(ctrl)
	filestorage := mocks.NewMockFileStorager(ctrl)
	usecase := NewCategoryUsecase(categoryRepo, cache, filestorage, logger)

	cache.EXPECT().CheckCache(ctx, categoriesListKey).Return(true)
	cache.EXPECT().GetCategoriesListCache(ctx, categoriesListKey).Return(categories, nil)
	res, err := usecase.GetCategoryList(context.Background())
	require.NoError(t, err)
	require.Equal(t, res, categories)

	testChan0 := make(chan models.Category, 1)
	testChan0 <- *testModelCategoryWithId
	close(testChan0)
	cache.EXPECT().CheckCache(ctx, categoriesListKey).Return(true)
	cache.EXPECT().GetCategoriesListCache(ctx, categoriesListKey).Return(nil, fmt.Errorf("error"))
	categoryRepo.EXPECT().GetCategoryList(ctx).Return(testChan0, fmt.Errorf("error"))
	res, err = usecase.GetCategoryList(context.Background())
	require.Error(t, err)
	require.Nil(t, res)

	cache.EXPECT().CheckCache(ctx, categoriesListKey).Return(true)
	cache.EXPECT().GetCategoriesListCache(ctx, categoriesListKey).Return(nil, fmt.Errorf("error"))
	categoryRepo.EXPECT().GetCategoryList(ctx).Return(testChan0, nil)
	res, err = usecase.GetCategoryList(context.Background())
	require.NoError(t, err)
	require.Equal(t, res, categories)

	testChan <- *testModelCategoryWithId
	close(testChan)
	cache.EXPECT().CheckCache(ctx, categoriesListKey).Return(false)
	categoryRepo.EXPECT().GetCategoryList(ctx).Return(testChan, nil)
	cache.EXPECT().CreateCategoriesListСache(ctx, categories, categoriesListKey).Return(nil)
	cache.EXPECT().GetCategoriesListCache(ctx, categoriesListKey).Return(categories, nil)
	res, err = usecase.GetCategoryList(context.Background())
	require.NoError(t, err)
	require.Equal(t, res, categories)

	cache.EXPECT().CheckCache(ctx, categoriesListKey).Return(false)
	categoryRepo.EXPECT().GetCategoryList(ctx).Return(testChan, fmt.Errorf("error"))
	res, err = usecase.GetCategoryList(context.Background())
	require.Error(t, err)
	require.Nil(t, res)

	testChan2 := make(chan models.Category, 1)
	testChan2 <- *testModelCategoryWithId
	close(testChan2)
	cache.EXPECT().CheckCache(ctx, categoriesListKey).Return(false)
	categoryRepo.EXPECT().GetCategoryList(ctx).Return(testChan2, nil)
	cache.EXPECT().CreateCategoriesListСache(ctx, categories, categoriesListKey).Return(fmt.Errorf("error"))
	cache.EXPECT().GetCategoriesListCache(ctx, categoriesListKey).Return(categories, nil)
	res, err = usecase.GetCategoryList(context.Background())
	require.NoError(t, err)
	require.NotNil(t, res)
}

func TestUpdateCategoryCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	categoryRepo := mocks.NewMockCategoryStore(ctrl)
	cache := mocks.NewMockICategoriesCache(ctrl)
	filestorage := mocks.NewMockFileStorager(ctrl)
	usecase := NewCategoryUsecase(categoryRepo, cache, filestorage, logger)

	cache.EXPECT().CheckCache(ctx, categoriesListKey).Return(false)
	err := usecase.UpdateCache(ctx, testId, "create")
	require.Error(t, err)

	cache.EXPECT().CheckCache(ctx, categoriesListKey).Return(true)
	categoryRepo.EXPECT().GetCategory(ctx, testId).Return(nil, fmt.Errorf("error"))
	err = usecase.UpdateCache(ctx, testId, "create")
	require.Error(t, err)

	cache.EXPECT().CheckCache(ctx, categoriesListKey).Return(true)
	categoryRepo.EXPECT().GetCategory(ctx, testId).Return(testModelCategoryWithId, nil)
	cache.EXPECT().GetCategoriesListCache(ctx, categoriesListKey).Return(nil, fmt.Errorf("error"))
	err = usecase.UpdateCache(ctx, testId, "create")
	require.Error(t, err)

	cache.EXPECT().CheckCache(ctx, categoriesListKey).Return(true)
	categoryRepo.EXPECT().GetCategory(ctx, testId).Return(testModelCategoryWithId, nil)
	cache.EXPECT().GetCategoriesListCache(ctx, categoriesListKey).Return(categories, nil)
	cache.EXPECT().CreateCategoriesListСache(ctx, categories, categoriesListKey).Return(fmt.Errorf("error"))
	err = usecase.UpdateCache(ctx, testId, "update")
	require.Error(t, err)

	cache.EXPECT().CheckCache(ctx, categoriesListKey).Return(true)
	categoryRepo.EXPECT().GetCategory(ctx, testId).Return(testModelCategoryWithId, nil)
	cache.EXPECT().GetCategoriesListCache(ctx, categoriesListKey).Return(categories, nil)
	cache.EXPECT().CreateCategoriesListСache(ctx, categories, categoriesListKey).Return(nil)
	err = usecase.UpdateCache(ctx, testId, "update")
	require.NoError(t, err)
}

func TestDeleteCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	categoryRepo := mocks.NewMockCategoryStore(ctrl)
	cache := mocks.NewMockICategoriesCache(ctrl)
	filestorage := mocks.NewMockFileStorager(ctrl)
	usecase := NewCategoryUsecase(categoryRepo, cache, filestorage, logger)

	categoryRepo.EXPECT().DeleteCategory(ctx, testId).Return(fmt.Errorf("error"))
	err := usecase.DeleteCategory(ctx, testId)
	require.Error(t, err)

	categoryRepo.EXPECT().DeleteCategory(ctx, testId).Return(nil)
	cache.EXPECT().CheckCache(ctx, categoriesListKey).Return(true)
	categoryRepo.EXPECT().GetCategory(ctx, testId).Return(testModelCategoryWithId, nil)
	cache.EXPECT().GetCategoriesListCache(ctx, categoriesListKey).Return(categories, nil)
	cache.EXPECT().CreateCategoriesListСache(ctx, []models.Category{}, categoriesListKey).Return(nil)
	err = usecase.DeleteCategory(ctx, testId)
	require.NoError(t, err)

	categoryRepo.EXPECT().DeleteCategory(ctx, testId).Return(nil)
	cache.EXPECT().CheckCache(ctx, categoriesListKey).Return(false)
	err = usecase.DeleteCategory(ctx, testId)
	require.NoError(t, err)
}

func TestGetCategoryByName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	categoryRepo := mocks.NewMockCategoryStore(ctrl)
	cache := mocks.NewMockICategoriesCache(ctrl)
	filestorage := mocks.NewMockFileStorager(ctrl)
	usecase := NewCategoryUsecase(categoryRepo, cache, filestorage, logger)

	categoryRepo.EXPECT().GetCategoryByName(ctx, testModelCategoryWithId.Name).Return(nil, fmt.Errorf("error"))
	res, err := usecase.GetCategoryByName(ctx, testModelCategoryWithId.Name)
	require.Error(t, err)
	require.Nil(t, res)

	categoryRepo.EXPECT().GetCategoryByName(ctx, testModelCategoryWithId.Name).Return(testModelCategoryWithId, nil)
	res, err = usecase.GetCategoryByName(ctx, testModelCategoryWithId.Name)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, res, testModelCategoryWithId)
}

func TestDeleteCategoryCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()
	logger := zap.L()
	categoryRepo := mocks.NewMockCategoryStore(ctrl)
	cache := mocks.NewMockICategoriesCache(ctrl)
	filestorage := mocks.NewMockFileStorager(ctrl)
	usecase := NewCategoryUsecase(categoryRepo, cache, filestorage, logger)

	cache.EXPECT().DeleteCache(ctx, "testNamenameasc").Return(fmt.Errorf("error"))
	err := usecase.DeleteCategoryCache(ctx, "testName")
	require.Error(t, err)

	cache.EXPECT().DeleteCache(ctx, "testNamenameasc").Return(nil)
	cache.EXPECT().DeleteCache(ctx, "testNamenamedesc").Return(nil)
	cache.EXPECT().DeleteCache(ctx, "testNamepriceasc").Return(nil)
	cache.EXPECT().DeleteCache(ctx, "testNamepricedesc").Return(nil)
	cache.EXPECT().DeleteCache(ctx, "testNameQuantity").Return(err)
	err = usecase.DeleteCategoryCache(ctx, "testName")
	require.Error(t, err)

	cache.EXPECT().DeleteCache(ctx, "testNamenameasc").Return(nil)
	cache.EXPECT().DeleteCache(ctx, "testNamenamedesc").Return(nil)
	cache.EXPECT().DeleteCache(ctx, "testNamepriceasc").Return(nil)
	cache.EXPECT().DeleteCache(ctx, "testNamepricedesc").Return(nil)
	cache.EXPECT().DeleteCache(ctx, "testNameQuantity").Return(nil)
	err = usecase.DeleteCategoryCache(ctx, "testName")
	require.NoError(t, err)
}
