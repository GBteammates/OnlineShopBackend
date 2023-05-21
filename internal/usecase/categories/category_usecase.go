package usecase

import (
	"OnlineShopBackend/internal/models"
	usecase "OnlineShopBackend/internal/usecase/interfaces"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var _ usecase.ICategoryUsecase = (*categoryUsecase)(nil)

const (
	categoriesListKey = "CategoriesList"
	createOp          = "create"
	updateOp          = "update"
	deleteOp          = "delete"
)

type categoryUsecase struct {
	categoryStore   usecase.CategoryStore
	categoriesCache usecase.ICategoriesCache
	filestorage     usecase.FileStorager
	logger          *zap.Logger
}

func NewCategoryUsecase(store usecase.CategoryStore, cache usecase.ICategoriesCache, filestorage usecase.FileStorager, logger *zap.Logger) *categoryUsecase {
	logger.Debug("Enter in usecase NewcategoryUsecase()")
	return &categoryUsecase{
		categoryStore:   store,
		categoriesCache: cache,
		filestorage:     filestorage,
		logger:          logger}
}

// CreateCategory call database method and returns id of created category or error
func (usecase *categoryUsecase) CreateCategory(ctx context.Context, category *models.Category) (uuid.UUID, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase CreateCategory() with args: ctx, category: %v", category)
	id, err := usecase.categoryStore.CreateCategory(ctx, category)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error on create category: %w", err)
	}
	category.Id = id
	err = usecase.categoriesCache.UpdateCategoryCache(ctx, category, createOp)
	if err != nil {
		usecase.logger.Error(fmt.Sprintf("error on update cache: %v", err))
	} else {
		usecase.logger.Info("Update cache success")
	}
	return id, nil
}

// UpdateCategory call database method to update category and returns error or nil
func (usecase *categoryUsecase) UpdateCategory(ctx context.Context, category *models.Category) error {
	usecase.logger.Sugar().Debugf("Enter in usecase UpdateCategory() with args: ctx, category: %v", category)
	err := usecase.categoryStore.UpdateCategory(ctx, category)
	if err != nil {
		return fmt.Errorf("error on update category: %w", err)
	}
	err = usecase.categoriesCache.UpdateCategoryCache(ctx, category, updateOp)
	if err != nil {
		usecase.logger.Error(fmt.Sprintf("error on update cache: %v", err))
	} else {
		usecase.logger.Info("Update cache success")
	}
	return nil
}

// GetCategory call database and returns *models.Category with given id or returns error
func (usecase *categoryUsecase) GetCategory(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase GetCategory() with args: ctx, id: %v", id)
	category, err := usecase.categoryStore.GetCategory(ctx, id)
	if err != nil {
		return &models.Category{}, fmt.Errorf("error on get category: %w", err)
	}
	return category, nil
}

// GetCategoryList call database method and returns chan with all models.Category or error
func (usecase *categoryUsecase) GetCategoryList(ctx context.Context) ([]models.Category, error) {
	usecase.logger.Debug("Enter in usecase GetCategoryList() with args: ctx")

	// Context with timeout so as not to wait for an answer from the cache for too long
	ctxT, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	// Сheck whether there is a cache with a list of categories
	if ok := usecase.categoriesCache.CheckCache(ctxT, categoriesListKey); !ok {
		// If cache does not exist, request a list of categories from the database
		categoryIncomingChan, err := usecase.categoryStore.GetCategoryList(ctx)
		if err != nil {
			return nil, err
		}
		categories := make([]models.Category, 0, 100)
		for category := range categoryIncomingChan {
			categories = append(categories, category)
		}
		// Create a cache with a list of categories
		err = usecase.categoriesCache.CreateCategoriesListСache(ctxT, categories, categoriesListKey)
		if err != nil {
			usecase.logger.Sugar().Warnf("error on create categories list cache with key: %s, error: %v", categoriesListKey, err)
		} else {
			usecase.logger.Sugar().Infof("Create categories list cache with key: %s success", categoriesListKey)
		}
	}

	// Get a list of categories from cache
	categories, err := usecase.categoriesCache.GetCategoriesListCache(ctxT, categoriesListKey)
	if err != nil {
		usecase.logger.Sugar().Warnf("error on get cache with key: %s, err: %v", categoriesListKey, err)
		// If error on get cache, request a list of categories from the database
		categoryIncomingChan, err := usecase.categoryStore.GetCategoryList(ctx)
		if err != nil {
			return nil, err
		}
		dbCategories := make([]models.Category, 0, 100)
		for category := range categoryIncomingChan {
			dbCategories = append(dbCategories, category)
			categories = dbCategories
		}
		usecase.logger.Info("Get category list from db success")
		return categories, nil
	}
	usecase.logger.Info("Get category list from cache success")
	return categories, nil
}

// DeleteCategory call database method for deleting category
func (usecase *categoryUsecase) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	usecase.logger.Sugar().Debugf("Enter in usecase DeleteCategory() with args: ctx, id: %v", id)
	err := usecase.categoryStore.DeleteCategory(ctx, id)
	if err != nil {
		return err
	}
	err = usecase.categoriesCache.UpdateCategoryCache(ctx, &models.Category{Id: id}, deleteOp)
	if err != nil {
		usecase.logger.Error(fmt.Sprintf("error on update cache: %v", err))
	}
	usecase.logger.Info("Delete category success")
	return nil
}

// GetCategoryByName call database method for get category by name
func (usecase *categoryUsecase) GetCategoryByName(ctx context.Context, name string) (*models.Category, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase GetCategoryByName() with args: ctx, name: %s", name)
	category, err := usecase.categoryStore.GetCategoryByName(ctx, name)
	if err != nil {
		return nil, err
	}
	usecase.logger.Info("Get category by name success")
	return category, nil
}
