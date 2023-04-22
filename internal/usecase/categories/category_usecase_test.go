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

func TestUploadCategoryImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	categoryRepo := mocks.NewMockCategoryStore(ctrl)
	cache := mocks.NewMockICategoriesCache(ctrl)
	filestorage := mocks.NewMockFileStorager(ctrl)
	ctx := context.Background()
	var (
		id       = uuid.New()
		filename = "testName"
		file     = []byte{0xff, 0xd8, 0xff, 0xe0, 0x0, 0x10, 0x4a, 0x46, 0x49, 0x46, 0x0, 0x1, 0x1, 0x1, 0x0, 0x48, 0x0, 0x48, 0x0, 0x0, 0xff, 0xe1, 0x0, 0x22, 0x45, 0x78, 0x69, 0x66, 0x0, 0x0, 0x4d, 0x4d, 0x0, 0x2a, 0x0, 0x0, 0x0, 0x8, 0x0, 0x1, 0x1, 0x12, 0x0, 0x3, 0x0, 0x0, 0x0, 0x1, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xff, 0xfe, 0x0, 0xd, 0x53, 0x65, 0x63, 0x6c, 0x75, 0x62, 0x2e, 0x6f, 0x72, 0x67, 0x0, 0xff, 0xdb, 0x0, 0x43, 0x0, 0x2, 0x1, 0x1, 0x2, 0x1, 0x1, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x2, 0x3, 0x5, 0x3, 0x3, 0x3, 0x3, 0x3, 0x6, 0x4, 0x4, 0x3, 0x5, 0x7, 0x6, 0x7, 0x7, 0x7, 0x6, 0x7, 0x7, 0x8, 0x9, 0xb, 0x9, 0x8, 0x8, 0xa, 0x8, 0x7, 0x7, 0xa, 0xd, 0xa, 0xa, 0xb, 0xc, 0xc, 0xc, 0xc, 0x7, 0x9, 0xe, 0xf, 0xd, 0xc, 0xe, 0xb, 0xc, 0xc, 0xc, 0xff, 0xdb, 0x0, 0x43, 0x1, 0x2, 0x2, 0x2, 0x3, 0x3, 0x3, 0x6, 0x3, 0x3, 0x6, 0xc, 0x8, 0x7, 0x8, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xc, 0xff, 0xc0, 0x0, 0x11, 0x8, 0x0, 0x1, 0x0, 0x1, 0x3, 0x1, 0x22, 0x0, 0x2, 0x11, 0x1, 0x3, 0x11, 0x1, 0xff, 0xc4, 0x0, 0x1f, 0x0, 0x0, 0x1, 0x5, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, 0xb, 0xff, 0xc4, 0x0, 0xb5, 0x10, 0x0, 0x2, 0x1, 0x3, 0x3, 0x2, 0x4, 0x3, 0x5, 0x5, 0x4, 0x4, 0x0, 0x0, 0x1, 0x7d, 0x1, 0x2, 0x3, 0x0, 0x4, 0x11, 0x5, 0x12, 0x21, 0x31, 0x41, 0x6, 0x13, 0x51, 0x61, 0x7, 0x22, 0x71, 0x14, 0x32, 0x81, 0x91, 0xa1, 0x8, 0x23, 0x42, 0xb1, 0xc1, 0x15, 0x52, 0xd1, 0xf0, 0x24, 0x33, 0x62, 0x72, 0x82, 0x9, 0xa, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5a, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6a, 0x73, 0x74, 0x75, 0x76, 0x77, 0x78, 0x79, 0x7a, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88, 0x89, 0x8a, 0x92,
			0x93, 0x94, 0x95, 0x96, 0x97, 0x98, 0x99, 0x9a, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6, 0xa7, 0xa8, 0xa9, 0xaa, 0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7, 0xb8, 0xb9, 0xba, 0xc2, 0xc3, 0xc4, 0xc5, 0xc6, 0xc7, 0xc8, 0xc9, 0xca, 0xd2, 0xd3, 0xd4, 0xd5, 0xd6, 0xd7, 0xd8, 0xd9, 0xda, 0xe1, 0xe2, 0xe3, 0xe4, 0xe5, 0xe6, 0xe7, 0xe8, 0xe9, 0xea, 0xf1, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6, 0xf7, 0xf8, 0xf9, 0xfa, 0xff, 0xc4, 0x0, 0x1f, 0x1, 0x0, 0x3, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, 0xb, 0xff, 0xc4, 0x0, 0xb5, 0x11, 0x0, 0x2, 0x1, 0x2, 0x4, 0x4, 0x3, 0x4, 0x7, 0x5, 0x4, 0x4, 0x0, 0x1, 0x2, 0x77, 0x0, 0x1, 0x2, 0x3, 0x11, 0x4, 0x5, 0x21, 0x31, 0x6, 0x12, 0x41, 0x51, 0x7, 0x61, 0x71, 0x13, 0x22, 0x32, 0x81, 0x8, 0x14, 0x42, 0x91, 0xa1, 0xb1, 0xc1, 0x9, 0x23, 0x33, 0x52, 0xf0, 0x15, 0x62, 0x72, 0xd1, 0xa, 0x16, 0x24, 0x34, 0xe1, 0x25, 0xf1, 0x17, 0x18, 0x19, 0x1a, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5a, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6a, 0x73, 0x74, 0x75, 0x76, 0x77, 0x78, 0x79, 0x7a, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88, 0x89, 0x8a, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97, 0x98, 0x99, 0x9a, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6, 0xa7, 0xa8, 0xa9, 0xaa, 0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7, 0xb8, 0xb9, 0xba, 0xc2, 0xc3, 0xc4, 0xc5, 0xc6, 0xc7, 0xc8, 0xc9, 0xca, 0xd2, 0xd3, 0xd4, 0xd5, 0xd6, 0xd7, 0xd8, 0xd9, 0xda, 0xe2, 0xe3, 0xe4, 0xe5, 0xe6, 0xe7, 0xe8, 0xe9, 0xea, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6, 0xf7, 0xf8, 0xf9, 0xfa, 0xff, 0xda, 0x0, 0xc, 0x3, 0x1, 0x0, 0x2, 0x11, 0x3, 0x11, 0x0, 0x3f, 0x0, 0xfc, 0x8b, 0xa2, 0x8a, 0x2b, 0xf3, 0xb3, 0xf6, 0x83, 0xff, 0xd9}
		category = &models.Category{
			Id:          id,
			Name:        "testName",
			Description: "testDescription",
			Image:       "testImage",
		}
		path = "testPath"
	)

	testCases := []struct {
		name           string
		store          *mocks.MockCategoryStore
		fs             *mocks.MockFileStorager
		cache          *mocks.MockICategoriesCache
		logger         *zap.Logger
		id             uuid.UUID
		filename       string
		file           []byte
		path           string
		category       *models.Category
		storeGetExpect func(*mocks.MockCategoryStore)
		fsExpect       func(*mocks.MockFileStorager)
		storeUpdExpect func(*mocks.MockCategoryStore)
		cacheExpect    func(*mocks.MockICategoriesCache)
	}{
		{
			name:           "error on get category",
			store:          categoryRepo,
			fs:             filestorage,
			cache:          cache,
			logger:         logger,
			id:             id,
			filename:       filename,
			file:           file,
			storeGetExpect: func(mis *mocks.MockCategoryStore) { mis.EXPECT().GetCategory(ctx, id).Return(nil, fmt.Errorf("error")) },
		},
		{
			name:           "error on filestorage put category image",
			store:          categoryRepo,
			fs:             filestorage,
			cache:          cache,
			logger:         logger,
			id:             id,
			filename:       filename,
			file:           file,
			category:       category,
			storeGetExpect: func(mis *mocks.MockCategoryStore) { mis.EXPECT().GetCategory(ctx, id).Return(category, nil) },
			fsExpect: func(mfs *mocks.MockFileStorager) {
				mfs.EXPECT().PutCategoryImage(id.String(), filename, file).Return("", fmt.Errorf("error"))
			},
		},
		{
			name:           "error on update category",
			store:          categoryRepo,
			fs:             filestorage,
			cache:          cache,
			logger:         logger,
			id:             id,
			filename:       filename,
			file:           file,
			category:       category,
			storeGetExpect: func(mis *mocks.MockCategoryStore) { mis.EXPECT().GetCategory(ctx, id).Return(category, nil) },
			fsExpect: func(mfs *mocks.MockFileStorager) {
				mfs.EXPECT().PutCategoryImage(id.String(), filename, file).Return(path, nil)
			},
			storeUpdExpect: func(mis *mocks.MockCategoryStore) {
				mis.EXPECT().UpdateCategory(ctx, category).Return(fmt.Errorf("error"))
			},
		},
		{
			name:           "success upload",
			store:          categoryRepo,
			fs:             filestorage,
			cache:          cache,
			logger:         logger,
			id:             id,
			filename:       filename,
			file:           file,
			category:       category,
			storeGetExpect: func(mis *mocks.MockCategoryStore) { mis.EXPECT().GetCategory(ctx, id).Return(category, nil) },
			fsExpect: func(mfs *mocks.MockFileStorager) {
				mfs.EXPECT().PutCategoryImage(id.String(), filename, file).Return(path, nil)
			},
			storeUpdExpect: func(mis *mocks.MockCategoryStore) {
				mis.EXPECT().UpdateCategory(ctx, category).Return(nil)
			},
			cacheExpect: func(mic *mocks.MockICategoriesCache) { mic.EXPECT().CheckCache(ctx, categoriesListKey).Return(false) },
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			usecase := NewCategoryUsecase(tc.store, tc.cache, tc.fs, tc.logger)
			tc.storeGetExpect(tc.store)
			if tc.fsExpect != nil {
				tc.fsExpect(tc.fs)
			}
			if tc.storeUpdExpect != nil {
				tc.storeUpdExpect(tc.store)
			}
			if tc.cacheExpect != nil {
				tc.cacheExpect(tc.cache)
			}
			err := usecase.UploadCategoryImage(ctx, tc.id, tc.filename, tc.file)
			if strings.Contains(tc.name, "error") {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDeleteCategoryImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	categoryRepo := mocks.NewMockCategoryStore(ctrl)
	cache := mocks.NewMockICategoriesCache(ctrl)
	filestorage := mocks.NewMockFileStorager(ctrl)
	ctx := context.Background()
	var (
		id       = uuid.New()
		filename = "testName"
		category = &models.Category{
			Id:          id,
			Name:        "testName",
			Description: "testDescription",
			Image:       "testImage",
		}
	)

	testCases := []struct {
		name           string
		store          *mocks.MockCategoryStore
		fs             *mocks.MockFileStorager
		cache          *mocks.MockICategoriesCache
		logger         *zap.Logger
		id             uuid.UUID
		filename       string
		file           []byte
		path           string
		category       *models.Category
		storeGetExpect func(*mocks.MockCategoryStore)
		fsExpect       func(*mocks.MockFileStorager)
		storeUpdExpect func(*mocks.MockCategoryStore)
		cacheExpect    func(*mocks.MockICategoriesCache)
	}{
		{
			name:           "error on get category",
			store:          categoryRepo,
			fs:             filestorage,
			cache:          cache,
			logger:         logger,
			id:             id,
			filename:       filename,
			storeGetExpect: func(mis *mocks.MockCategoryStore) { mis.EXPECT().GetCategory(ctx, id).Return(nil, fmt.Errorf("error")) },
		},
		{
			name:           "error on filestorage delete category image",
			store:          categoryRepo,
			fs:             filestorage,
			cache:          cache,
			logger:         logger,
			id:             id,
			filename:       filename,
			category:       category,
			storeGetExpect: func(mis *mocks.MockCategoryStore) { mis.EXPECT().GetCategory(ctx, id).Return(category, nil) },
			fsExpect: func(mfs *mocks.MockFileStorager) {
				mfs.EXPECT().DeleteCategoryImage(id.String(), filename).Return(fmt.Errorf("error"))
			},
		},
		{
			name:           "error on update category",
			store:          categoryRepo,
			fs:             filestorage,
			cache:          cache,
			logger:         logger,
			id:             id,
			filename:       filename,
			category:       category,
			storeGetExpect: func(mis *mocks.MockCategoryStore) { mis.EXPECT().GetCategory(ctx, id).Return(category, nil) },
			fsExpect: func(mfs *mocks.MockFileStorager) {
				mfs.EXPECT().DeleteCategoryImage(id.String(), filename).Return(nil)
			},
			storeUpdExpect: func(mis *mocks.MockCategoryStore) {
				mis.EXPECT().UpdateCategory(ctx, category).Return(fmt.Errorf("error"))
			},
		},
		{
			name:           "success deleted",
			store:          categoryRepo,
			fs:             filestorage,
			cache:          cache,
			logger:         logger,
			id:             id,
			filename:       filename,
			category:       category,
			storeGetExpect: func(mis *mocks.MockCategoryStore) { mis.EXPECT().GetCategory(ctx, id).Return(category, nil) },
			fsExpect: func(mfs *mocks.MockFileStorager) {
				mfs.EXPECT().DeleteCategoryImage(id.String(), filename).Return(nil)
			},
			storeUpdExpect: func(mis *mocks.MockCategoryStore) {
				mis.EXPECT().UpdateCategory(ctx, category).Return(nil)
			},
			cacheExpect: func(mic *mocks.MockICategoriesCache) { mic.EXPECT().CheckCache(ctx, categoriesListKey).Return(false) },
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			usecase := NewCategoryUsecase(tc.store, tc.cache, tc.fs, tc.logger)
			tc.storeGetExpect(tc.store)
			if tc.fsExpect != nil {
				tc.fsExpect(tc.fs)
			}
			if tc.storeUpdExpect != nil {
				tc.storeUpdExpect(tc.store)
			}
			if tc.cacheExpect != nil {
				tc.cacheExpect(tc.cache)
			}
			err := usecase.DeleteCategoryImage(ctx, tc.id, tc.filename)
			if strings.Contains(tc.name, "error") {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetCategoryImagesList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	categoryRepo := mocks.NewMockCategoryStore(ctrl)
	cache := mocks.NewMockICategoriesCache(ctrl)
	filestorage := mocks.NewMockFileStorager(ctrl)
	ctx := context.Background()
	var (
		fileInfos = []*models.FileInfo{
			{
				Name:       "testName",
				Path:       "testPath",
				CreateDate: "testCreateDate",
				ModifyDate: "testModifyDate",
			},
		}
	)

	testCases := []struct {
		name      string
		store     *mocks.MockCategoryStore
		fs        *mocks.MockFileStorager
		cache     *mocks.MockICategoriesCache
		logger    *zap.Logger
		fileInfos []*models.FileInfo
		fsExpect  func(*mocks.MockFileStorager)
	}{
		{
			name:   "error on filestorage get categories images list",
			store:  categoryRepo,
			fs:     filestorage,
			cache:  cache,
			logger: logger,
			fsExpect: func(mfs *mocks.MockFileStorager) {
				mfs.EXPECT().GetCategoriesImagesList().Return(nil, fmt.Errorf("error"))
			},
		},
		{
			name:   "success get categories images list",
			store:  categoryRepo,
			fs:     filestorage,
			cache:  cache,
			logger: logger,
			fsExpect: func(mfs *mocks.MockFileStorager) {
				mfs.EXPECT().GetCategoriesImagesList().Return(fileInfos, nil)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			usecase := NewCategoryUsecase(tc.store, tc.cache, tc.fs, tc.logger)
			tc.fsExpect(tc.fs)

			res, err := usecase.GetCategoriesImagesList(ctx)
			if strings.Contains(tc.name, "error") {
				require.Error(t, err)
				require.Nil(t, res)
			} else {
				require.NoError(t, err)
				require.Equal(t, res, fileInfos)
			}
		})
	}
}

func TestDeleteCategoryImageById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	categoryRepo := mocks.NewMockCategoryStore(ctrl)
	cache := mocks.NewMockICategoriesCache(ctrl)
	filestorage := mocks.NewMockFileStorager(ctrl)
	ctx := context.Background()
	var (
		id = uuid.New()
	)

	testCases := []struct {
		name     string
		store    *mocks.MockCategoryStore
		fs       *mocks.MockFileStorager
		cache    *mocks.MockICategoriesCache
		logger   *zap.Logger
		id       uuid.UUID
		fsExpect func(*mocks.MockFileStorager)
	}{
		{
			name:   "error on filestorage delete category image",
			store:  categoryRepo,
			fs:     filestorage,
			cache:  cache,
			logger: logger,
			fsExpect: func(mfs *mocks.MockFileStorager) {
				mfs.EXPECT().DeleteCategoryImageById(id.String()).Return(fmt.Errorf("error"))
			},
		},
		{
			name:   "success delete category image by id",
			store:  categoryRepo,
			fs:     filestorage,
			cache:  cache,
			logger: logger,
			fsExpect: func(mfs *mocks.MockFileStorager) {
				mfs.EXPECT().DeleteCategoryImageById(id.String()).Return(nil)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			usecase := NewCategoryUsecase(tc.store, tc.cache, tc.fs, tc.logger)
			tc.fsExpect(tc.fs)

			err := usecase.DeleteCategoryImageById(ctx, id)
			if strings.Contains(tc.name, "error") {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
