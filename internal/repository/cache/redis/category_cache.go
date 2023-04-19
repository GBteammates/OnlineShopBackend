package redis

import (
	"OnlineShopBackend/internal/models"
	"OnlineShopBackend/internal/usecase/interfaces"
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var _ usecase.ICategoriesCache = &CategoriesСache{}

type CategoriesСache struct {
	*RedisCache
	logger *zap.Logger
}

type categoriesData struct {
	Categories []models.Category `json:"categories"`
}

func NewCategoriesСache(cache *RedisCache, logger *zap.Logger) usecase.ICategoriesCache {
	logger.Debug("Enter in cache NewCategoriesСache()")
	return &CategoriesСache{cache, logger}
}

// Checkcache checks for data in the cache
func (cache *CategoriesСache) CheckCache(ctx context.Context, key string) bool {
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
func (cache *CategoriesСache) CreateCategoriesListСache(ctx context.Context, categories []models.Category, key string) error {
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
		return fmt.Errorf("error on set cache with key: %v, error: %w", key, err)
	}
	cache.logger.Debug(fmt.Sprintf("cache with key %s write in redis success", key))
	return nil
}

// GetCategoriesListCache retrieves data from the cache
func (cache *CategoriesСache) GetCategoriesListCache(ctx context.Context, key string) ([]models.Category, error) {
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
func (cache *CategoriesСache) DeleteCache(ctx context.Context, key string) error {
	cache.logger.Debug(fmt.Sprintf("Enter in cache DeleteCache with args: ctx, key: %s", key))
	err := cache.Del(ctx, key).Err()
	if err != nil {
		cache.logger.Sugar().Warnf("Error on delete cache with key: %s", key)
		return err
	}
	cache.logger.Sugar().Infof("Delete cache with key: %s success", key)
	return nil
}
