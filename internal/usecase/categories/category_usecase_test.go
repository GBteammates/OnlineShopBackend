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

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L().Sugar()
	store := mocks.NewMockCategoryStore(ctrl)
	cache := mocks.NewMockICategoriesCache(ctrl)
	filestorage := mocks.NewMockFilestorage(ctrl)
	ctx := context.Background()
	usecase := New(store, cache, filestorage, logger)
	id := uuid.New()
	category := &models.Category{
		Name:        "testName",
		Description: "testDescription",
		Image:       "testImage",
	}
	testErr := errors.New("error")
	t.Run("error store create category", func(t *testing.T) {
		store.EXPECT().Create(ctx, category).Return(uuid.Nil, testErr)
		res, err := usecase.Create(ctx, category)
		require.Error(t, err)
		require.Equal(t, uuid.Nil, res)
	})
	t.Run("success create category", func(t *testing.T) {
		store.EXPECT().Create(ctx, category).Return(id, nil)
		category.Id = id
		cache.EXPECT().UpdateCache(ctx, category, models.CreateOp).Return(nil)
		res, err := usecase.Create(ctx, category)
		require.NoError(t, err)
		require.Equal(t, id, res)
	})
}

func TestUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L().Sugar()
	store := mocks.NewMockCategoryStore(ctrl)
	cache := mocks.NewMockICategoriesCache(ctrl)
	filestorage := mocks.NewMockFilestorage(ctrl)
	ctx := context.Background()
	usecase := New(store, cache, filestorage, logger)
	id := uuid.New()
	category := &models.Category{
		Id:          id,
		Name:        "testName",
		Description: "testDescription",
		Image:       "testImage",
	}
	testErr := errors.New("error")
	t.Run("error store update category", func(t *testing.T) {
		store.EXPECT().Update(ctx, category).Return(testErr)
		err := usecase.Update(ctx, category)
		require.Error(t, err)
	})
	t.Run("success update category", func(t *testing.T) {
		store.EXPECT().Update(ctx, category).Return(nil)
		cache.EXPECT().UpdateCache(ctx, category, models.UpdateOp).Return(nil)
		err := usecase.Update(ctx, category)
		require.NoError(t, err)
	})
}

func TestGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L().Sugar()
	store := mocks.NewMockCategoryStore(ctrl)
	cache := mocks.NewMockICategoriesCache(ctrl)
	filestorage := mocks.NewMockFilestorage(ctrl)
	ctx := context.Background()
	usecase := New(store, cache, filestorage, logger)
	id := uuid.New()
	category := &models.Category{
		Id:          id,
		Name:        "testName",
		Description: "testDescription",
		Image:       "testImage",
	}
	testErr := errors.New("error")
	t.Run("error get category", func(t *testing.T) {
		store.EXPECT().Get(ctx, id).Return(nil, testErr)
		res, err := usecase.Get(ctx, id)
		require.Error(t, err)
		require.Equal(t, &models.Category{}, res)
	})
	t.Run("success get category", func(t *testing.T) {
		store.EXPECT().Get(ctx, id).Return(category, nil)
		res, err := usecase.Get(ctx, id)
		require.NoError(t, err)
		require.Equal(t, category, res)
	})
}

func TestCategoryByName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L().Sugar()
	store := mocks.NewMockCategoryStore(ctrl)
	cache := mocks.NewMockICategoriesCache(ctrl)
	filestorage := mocks.NewMockFilestorage(ctrl)
	ctx := context.Background()
	usecase := New(store, cache, filestorage, logger)
	name := "name"
	id := uuid.New()
	category := &models.Category{
		Id:          id,
		Name:        name,
		Description: "testDescription",
		Image:       "testImage",
	}
	testErr := errors.New("error")
	t.Run("error get category by name", func(t *testing.T) {
		store.EXPECT().CategoryByName(ctx, name).Return(nil, testErr)
		res, err := usecase.CategoryByName(ctx, name)
		require.Error(t, err)
		require.Equal(t, &models.Category{}, res)
	})
	t.Run("success get category", func(t *testing.T) {
		store.EXPECT().CategoryByName(ctx, name).Return(category, nil)
		res, err := usecase.CategoryByName(ctx, name)
		require.NoError(t, err)
		require.Equal(t, category, res)
	})
}

func TestList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L().Sugar()
	store := mocks.NewMockCategoryStore(ctrl)
	cache := mocks.NewMockICategoriesCache(ctrl)
	filestorage := mocks.NewMockFilestorage(ctrl)
	ctx := context.Background()
	usecase := New(store, cache, filestorage, logger)
	categories := []models.Category{}
	ch := make(chan models.Category, 1)
	testErr := errors.New("error")
	t.Run("error on list categories", func(t *testing.T) {
		cache.EXPECT().CategoriesFromCache(gomock.Any(), models.CategoriesList).Return(nil, testErr)
		store.EXPECT().List(ctx).Return(ch, testErr)
		res, err := usecase.List(ctx)
		require.Error(t, err)
		require.Nil(t, res)
	})
	t.Run("success list categories", func(t *testing.T) {
		cache.EXPECT().CategoriesFromCache(gomock.Any(), models.CategoriesList).Return(categories, nil)
		res, err := usecase.List(ctx)
		require.NoError(t, err)
		require.Equal(t, categories, res)
	})
}

func TestDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L().Sugar()
	store := mocks.NewMockCategoryStore(ctrl)
	cache := mocks.NewMockICategoriesCache(ctrl)
	filestorage := mocks.NewMockFilestorage(ctrl)
	ctx := context.Background()
	usecase := New(store, cache, filestorage, logger)
	id := uuid.New()
	category := &models.Category{Id: id}
	testErr := errors.New("error")
	t.Run("error delete category", func(t *testing.T) {
		store.EXPECT().Delete(ctx, id).Return(testErr)
		err := usecase.Delete(ctx, id)
		require.Error(t, err)
	})
	t.Run("success delete category", func(t *testing.T) {
		store.EXPECT().Delete(ctx, id).Return(nil)
		cache.EXPECT().UpdateCache(ctx, category, models.DeleteOp).Return(nil)
		err := usecase.Delete(ctx, id)
		require.NoError(t, err)
	})
}

func Test_getCategories(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L().Sugar()
	store := mocks.NewMockCategoryStore(ctrl)
	cache := mocks.NewMockICategoriesCache(ctrl)
	filestorage := mocks.NewMockFilestorage(ctrl)
	ctx := context.Background()
	usecase := New(store, cache, filestorage, logger)
	categories := []models.Category{}
	ch := make(chan models.Category, 1)
	testErr := errors.New("error")
	t.Run("success get categories from cache", func(t *testing.T) {
		cache.EXPECT().CategoriesFromCache(gomock.Any(), models.CategoriesList).Return(categories, nil)
		res, err := usecase.getCategories(ctx)
		require.NoError(t, err)
		require.Equal(t, categories, res)
	})
	t.Run("error on get categories from store", func(t *testing.T) {
		cache.EXPECT().CategoriesFromCache(gomock.Any(), models.CategoriesList).Return(nil, testErr)
		store.EXPECT().List(ctx).Return(ch, testErr)
		res, err := usecase.getCategories(ctx)
		require.Error(t, err)
		require.Nil(t, res)
	})
	t.Run("success get empty categories list", func(t *testing.T) {
		cache.EXPECT().CategoriesFromCache(gomock.Any(), models.CategoriesList).Return(nil, testErr)
		store.EXPECT().List(ctx).Return(ch, nil)
		res, err := usecase.getCategories(ctx)
		require.NoError(t, err)
		require.Equal(t, categories, res)
	})
	t.Run("success get not empty categories list", func(t *testing.T) {
		category := models.Category{Name: "name"}
		categories := make([]models.Category, 1)
		categories = append(categories, category)
		ch <- category
		close(ch)
		cache.EXPECT().CategoriesFromCache(gomock.Any(), models.CategoriesList).Return(nil, testErr)
		store.EXPECT().List(ctx).Return(ch, nil)
		cache.EXPECT().CategoriesToCache(ctx, categories).Return(nil)
		res, err := usecase.getCategories(ctx)
		require.NoError(t, err)
		require.Equal(t, categories, res)
	})
}
