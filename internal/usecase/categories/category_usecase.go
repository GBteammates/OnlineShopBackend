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

type categoryUsecase struct {
	timeout     time.Duration
	store       usecase.CategoryStore
	cache       usecase.ICategoriesCache
	filestorage usecase.Filestorage
	logger      *zap.SugaredLogger
}

func New(
	store usecase.CategoryStore,
	cache usecase.ICategoriesCache,
	filestorage usecase.Filestorage,
	logger *zap.SugaredLogger,
) *categoryUsecase {
	logger.Debug("Enter in usecase NewcategoryUsecase()")
	return &categoryUsecase{
		store:       store,
		cache:       cache,
		filestorage: filestorage,
		logger:      logger}
}

// Create call database method and returns id of created category or error
func (u *categoryUsecase) Create(ctx context.Context, category *models.Category) (uuid.UUID, error) {
	u.logger.Debugf("Enter in usecase Create() with args: ctx, category: %v", category)
	id, err := u.store.Create(ctx, category)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error on create category: %w", err)
	}
	category.Id = id
	err = u.cache.UpdateCache(ctx, category, models.CreateOp)
	if err != nil {
		u.logger.Error(fmt.Sprintf("error on update cache: %v", err))
	} else {
		u.logger.Info("Update cache success")
	}
	return id, nil
}

// Update call database method to update category and returns error or nil
func (u *categoryUsecase) Update(ctx context.Context, category *models.Category) error {
	u.logger.Debugf("Enter in usecase Update() with args: ctx, category: %v", category)
	err := u.store.Update(ctx, category)
	if err != nil {
		return fmt.Errorf("error on update category: %w", err)
	}
	err = u.cache.UpdateCache(ctx, category, models.UpdateOp)
	if err != nil {
		u.logger.Error(fmt.Sprintf("error on update cache: %v", err))
	} else {
		u.logger.Info("Update cache success")
	}
	return nil
}

// GetCategory call database and returns *models.Category with given id or returns error
func (u *categoryUsecase) Get(ctx context.Context, param string) (*models.Category, error) {
	u.logger.Debugf("Enter in usecase Get() with args: ctx, param", param)
	category, err := u.store.Get(ctx, param)
	if err != nil {
		return &models.Category{}, fmt.Errorf("error on get category: %w", err)
	}
	return category, nil
}

// GetCategoryList call database method and returns chan with all models.Category or error
func (u *categoryUsecase) List(ctx context.Context) ([]models.Category, error) {
	u.logger.Debug("Enter in usecase List()")

	categories, err := u.getCategories(ctx)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// DeleteCategory call database method for deleting category
func (u *categoryUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	u.logger.Debugf("Enter in usecase Delete() with args: ctx, id: %v", id)
	err := u.store.Delete(ctx, id)
	if err != nil {
		return err
	}
	err = u.cache.UpdateCache(ctx, &models.Category{Id: id}, models.DeleteOp)
	if err != nil {
		u.logger.Error(fmt.Sprintf("error on update cache: %v", err))
	}
	u.logger.Info("Delete category success")
	return nil
}

func (u *categoryUsecase) getCategories(ctx context.Context) ([]models.Category, error) {
	u.logger.Debugf("Enter in usecasesecae getCategories()")

	// Context with timeout so as not to wait for an answer from the cache for too long
	ctxT, cancel := context.WithTimeout(ctx, u.timeout*time.Millisecond)
	defer cancel()

	categories, err := u.cache.CategoriesFromCache(ctxT, models.CategoriesList)
	if err == nil {
		return categories, nil
	}
	if err != nil {
		u.logger.Warnf("error on categories from cache: %v", err)
	}
	categoriesChan, err := u.store.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("error on get categories list from store: %w", err)
	}
	categories = make([]models.Category, 100)
	for category := range categoriesChan {
		categories = append(categories, category)
	}
	err = u.cache.CategoriesToCache(ctx, categories)
	if err != nil {
		u.logger.Warnf("error on categories to cache: %v", err)
	}
	return categories, nil
}
