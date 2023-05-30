package redis

import (
	"OnlineShopBackend/internal/models"
	usecase "OnlineShopBackend/internal/usecase/interfaces"
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

const (
	categoriesListKey = "CategoriesList"
)

var _ usecase.ICategoriesCache = &categoriesСache{}

type categoriesСache struct {
	*RedisCache
	cacheValid bool
	logger     *zap.Logger
}

type categoriesData struct {
	Categories []models.Category `json:"categories"`
}

func NewCategoriesСache(cache *RedisCache, logger *zap.Logger) *categoriesСache {
	logger.Debug("Enter in cache NewCategoriesСache()")
	return &categoriesСache{cache, false, logger}
}

// Checkcache checks for data in the cache
func (cache *categoriesСache) CheckCache(ctx context.Context, key string) bool {
	cache.logger.Sugar().Debugf("Enter in cache CheckCache() with args: ctx, key: %s", key)
	check := cache.Exists(ctx, key)
	result, err := check.Result()
	if err != nil {
		cache.logger.Error(fmt.Errorf("error on check cache: %w", err).Error())
		return false
	}
	cache.logger.Debug(fmt.Sprintf("Check cache with key: %s is %v", key, result))
	if result == 0 {
		cache.logger.Debug(fmt.Sprintf("Redis: key %s not exist", key))
		return false
	} else {
		cache.logger.Debug(fmt.Sprintf("Key %s in cache found success", key))
		return true
	}
}

// CreateCategoriesListСache creates cache of categories list
func (cache *categoriesСache) CreateCategoriesListСache(ctx context.Context, categories []models.Category, key string) error {
	cache.logger.Sugar().Debugf("Enter in CategoriesСache CreateСategoriesListCache() with args: ctx, categories []models.Category, key: %s", key)
	in := categoriesData{
		Categories: categories,
	}
	bytesData, err := json.Marshal(in)
	if err != nil {
		cache.logger.Sugar().Warnf("Error on json marshal data: %v", in)
		return fmt.Errorf("marshal unknown category: %w", err)
	}

	err = cache.Set(ctx, key, bytesData, cache.TTL).Err()
	if err != nil {
		cache.logger.Sugar().Warnf("Error on set cache with key: %s, error: %v", key, err)
		cache.cacheValid = false
		return fmt.Errorf("error on set cache with key: %v, error: %w", key, err)
	}
	cache.logger.Debug(fmt.Sprintf("cache with key %s write in redis success", key))
	return nil
}

// GetCategoriesListCache retrieves data from the cache
func (cache *categoriesСache) GetCategoriesListCache(ctx context.Context, key string) ([]models.Category, error) {
	cache.logger.Sugar().Debugf("Enter in cache GetCategoriesListCache() with args: ctx, key: %s", key)
	categories := categoriesData{}
	bytesData, err := cache.Get(ctx, key).Bytes()
	if err == redis.Nil {
		// we got empty result, it's not an error
		cache.logger.Debug("Success get nil result")
		return nil, nil
	} else if err != nil {
		cache.logger.Sugar().Errorf("Error on get cache: %v", err)
		return nil, err
	}
	err = json.Unmarshal(bytesData, &categories)
	if err != nil {
		cache.logger.Sugar().Warnf("Can't json unmarshal data: %v", bytesData)
		return nil, err
	}
	cache.logger.Debug("Get cache success")
	return categories.Categories, nil
}

// DeleteCache deleted cache by key
func (cache *categoriesСache) deleteCache(ctx context.Context, key string) error {
	cache.logger.Debug(fmt.Sprintf("Enter in cache DeleteCache with args: ctx, key: %s", key))
	err := cache.Del(ctx, key).Err()
	if err != nil {
		cache.cacheValid = false
		cache.logger.Sugar().Warnf("Error on delete cache with key: %s", key)
		return err
	}
	cache.logger.Sugar().Infof("Delete cache with key: %s success", key)
	return nil
}

// UpdateCache updating cache when creating or updating category
func (cache *categoriesСache) UpdateCategoryCache(ctx context.Context, newCategory *models.Category, op string) error {
	cache.logger.Sugar().Debugf("Enter in categoryCache UpdateCache() with args: ctx, newCategory: %v, op: %s", newCategory, op)
	// If the cache with such a key does not exist, we return the error, there is nothing to update
	if !cache.CheckCache(ctx, categoriesListKey) {
		return fmt.Errorf("cache is not exists")
	}

	// Get a list of categories from cache
	categories, err := cache.GetCategoriesListCache(ctx, categoriesListKey)
	if err != nil {
		return fmt.Errorf("error on get cache: %w", err)
	}
	// Change list of categories for update the cache
	if op == updateOp {
		for i, category := range categories {
			if category.Id == newCategory.Id {
				categories[i] = *newCategory
				break
			}
		}
	}
	if op == createOp {
		categories = append(categories, *newCategory)
	}
	if op == deleteOp {
		for i, category := range categories {
			if category.Id == newCategory.Id {
				categories = append(categories[:i], categories[i+1:]...)
				break
			}
		}
	}
	// Sort list of categories by name in alphabetical order
	sort.Slice(categories, func(i, j int) bool { return categories[i].Name < categories[j].Name })
	// Create new cache with list of categories
	err = cache.CreateCategoriesListСache(ctx, categories, categoriesListKey)
	if err != nil {
		return err
	}
	cache.logger.Info("Category cache update success")
	return nil
}

// DeleteCategoryCache deleted cache after deleting categories
func (cache *categoriesСache) DeleteCategoryCache(ctx context.Context, name string) error {
	cache.logger.Debug(fmt.Sprintf("Enter in usecase DeleteCategoryCache() with args: ctx, name: %s", name))
	// keys is a list of cache keys with items in deleted category sorting by name and price
	keys := []string{name + nameAsc, name + nameDesc, name + priceAsc, name + priceDesc}
	for _, key := range keys {
		// For each key from list delete cache
		err := cache.deleteCache(ctx, key)
		if err != nil {
			cache.logger.Error(fmt.Sprintf("error on delete cache with key: %s, error is %v", key, err))
			return err
		}
	}
	// Delete cache with quantity of items in deleted category
	err := cache.deleteCache(ctx, name+quantity)
	if err != nil {
		cache.logger.Error(fmt.Sprintf("error on delete cache with key: %s, error is %v", name, err))
		return err
	}
	cache.logger.Info("Category cache deleted success")
	return nil
}

func (cache *categoriesСache) Status(ctx context.Context) bool {
	cache.logger.Debug("Enter in cache Status()")
	return cache.cacheValid
}
