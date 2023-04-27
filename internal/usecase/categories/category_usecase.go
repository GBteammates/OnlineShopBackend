package usecase

import (
	"OnlineShopBackend/internal/models"
	usecase "OnlineShopBackend/internal/usecase/interfaces"
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var _ usecase.ICategoryUsecase = (*categoryUsecase)(nil)

var (
	categoriesListKey = "CategoriesList"
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
	err = usecase.UpdateCache(ctx, id, "create")
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
	err = usecase.UpdateCache(ctx, category.Id, "update")
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
	err = usecase.UpdateCache(ctx, id, "delete")
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

// UpdateCache updating cache when creating or updating category
func (usecase *categoryUsecase) UpdateCache(ctx context.Context, id uuid.UUID, op string) error {
	usecase.logger.Sugar().Debugf("Enter in usecase UpdateCache() with args: ctx, id: %v, op: %s", id, op)
	// If the cache with such a key does not exist, we return the error, there is nothing to update
	if !usecase.categoriesCache.CheckCache(ctx, categoriesListKey) {
		return fmt.Errorf("cache is not exists")
	}

	// Get a category from the database for updating in the cache
	newCategory, err := usecase.categoryStore.GetCategory(ctx, id)
	if err != nil {
		// If the error returned and the cache is updated in connection
		// with the removal of the category, we use an empty category with the Id
		if op == "delete" {
			newCategory = &models.Category{Id: id}
		} else {
			return fmt.Errorf("error on get category: %w", err)
		}
	}
	// Get a list of categories from cache
	categories, err := usecase.categoriesCache.GetCategoriesListCache(ctx, categoriesListKey)
	if err != nil {
		return fmt.Errorf("error on get cache: %w", err)
	}
	// Change list of categories for update the cache
	if op == "update" {
		for i, category := range categories {
			if category.Id == id {
				categories[i] = *newCategory
				break
			}
		}
	}
	if op == "create" {
		categories = append(categories, *newCategory)
	}
	if op == "delete" {
		for i, category := range categories {
			if category.Id == id {
				categories = append(categories[:i], categories[i+1:]...)
				break
			}
		}
	}
	// Sort list of categories by name in alphabetical order
	sort.Slice(categories, func(i, j int) bool { return categories[i].Name < categories[j].Name })
	// Create new cache with list of categories
	err = usecase.categoriesCache.CreateCategoriesListСache(ctx, categories, categoriesListKey)
	if err != nil {
		return err
	}
	usecase.logger.Info("Category cache update success")
	return nil
}

// DeleteCategoryCache deleted cache after deleting categories
func (usecase *categoryUsecase) DeleteCategoryCache(ctx context.Context, name string) error {
	usecase.logger.Debug(fmt.Sprintf("Enter in usecase DeleteCategoryCache() with args: ctx, name: %s", name))
	// keys is a list of cache keys with items in deleted category sorting by name and price
	keys := []string{name + "nameasc", name + "namedesc", name + "priceasc", name + "pricedesc"}
	for _, key := range keys {
		// For each key from list delete cache
		err := usecase.categoriesCache.DeleteCache(ctx, key)
		if err != nil {
			usecase.logger.Error(fmt.Sprintf("error on delete cache with key: %s, error is %v", key, err))
			return err
		}
	}
	// Delete cache with quantity of items in deleted category
	err := usecase.categoriesCache.DeleteCache(ctx, name+"Quantity")
	if err != nil {
		usecase.logger.Error(fmt.Sprintf("error on delete cache with key: %s, error is %v", name, err))
		return err
	}
	usecase.logger.Info("Category cache deleted success")
	return nil
}

func (usecase *categoryUsecase) UploadCategoryImage(ctx context.Context, id uuid.UUID, name string, file []byte) error {
	usecase.logger.Sugar().Debugf("Enter in usecase UploadCategoryImage() with args: ctx, id: %v, name: %s, file", id, name)

	category, err := usecase.categoryStore.GetCategory(ctx, id)
	if err != nil {
		return fmt.Errorf("error on get category: %w", err)
	}

	path, err := usecase.filestorage.PutCategoryImage(id.String(), name, file)
	if err != nil {
		return fmt.Errorf("error on put image to filestorage: %w", err)
	}

	category.Image = path

	err = usecase.categoryStore.UpdateCategory(ctx, category)
	if err != nil {
		return fmt.Errorf("error on update category: %w", err)
	}

	err = usecase.UpdateCache(ctx, category.Id, "update")
	if err != nil {
		usecase.logger.Error(fmt.Sprintf("error on update cache: %v", err))
	} else {
		usecase.logger.Info("Update cache success")
	}
	return nil
}

func (usecase *categoryUsecase) DeleteCategoryImage(ctx context.Context, id uuid.UUID, name string) error {
	usecase.logger.Sugar().Debugf("Enter in usecase DeleteCategoryImage() with args: ctx, id: %v, name: %s", id, name)

	category, err := usecase.categoryStore.GetCategory(ctx, id)
	if err != nil {
		return fmt.Errorf("error on get category: %w", err)
	}
	if strings.Contains(category.Image, name) {
		category.Image = ""
	}

	err = usecase.filestorage.DeleteCategoryImage(id.String(), name)
	if err != nil {
		return fmt.Errorf("error on delete category image: %w", err)
	}

	err = usecase.categoryStore.UpdateCategory(ctx, category)
	if err != nil {
		return fmt.Errorf("error on update category: %w", err)
	}
	err = usecase.UpdateCache(ctx, category.Id, "update")
	if err != nil {
		usecase.logger.Error(fmt.Sprintf("error on update cache: %v", err))
	} else {
		usecase.logger.Info("Update cache success")
	}
	return nil
}

func (usecase *categoryUsecase) GetCategoriesImagesList(ctx context.Context) ([]*models.FileInfo, error) {
	usecase.logger.Debug("Enter in usecase GetCategoriesImagesList()")

	result, err := usecase.filestorage.GetCategoriesImagesList()
	if err != nil {
		return nil, fmt.Errorf("error on get categories images list: %w", err)
	}
	return result, nil
}

func (usecase *categoryUsecase) DeleteCategoryImageById(ctx context.Context, id uuid.UUID) error {
	usecase.logger.Sugar().Debugf("Enter in usecase DeleteCategoryImageById() with arts: ctx, id: %v", id)

	err := usecase.filestorage.DeleteCategoryImageById(id.String())
	if err != nil {
		return fmt.Errorf("error on delete category image by id: %w", err)
	}
	return nil
}