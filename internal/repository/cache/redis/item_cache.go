package redis

import (
	"OnlineShopBackend/internal/models"
	"OnlineShopBackend/internal/usecase"
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var _ usecase.IItemsCache = &ItemsCache{}

type ItemsCache struct {
	*RedisCache
	logger *zap.Logger
}

type results struct {
	Responses []models.Item
}

func NewItemsCache(cache *RedisCache, logger *zap.Logger) usecase.IItemsCache {
	logger.Debug("Enter in cache NewItemsCache")
	return &ItemsCache{cache, logger}
}

// CheckCache checks for data in the cache
func (cache *ItemsCache) CheckCache(ctx context.Context, key string) bool {
	cache.logger.Sugar().Debugf("Enter in cache CheckCache() with args: ctx, key: %s", key)
	check := cache.Exists(ctx, key)
	result, err := check.Result()
	if err != nil {
		cache.logger.Error(fmt.Errorf("error on check cache: %w", err).Error())
		return false
	}
	cache.logger.Debug(fmt.Sprintf("Check Cache with key: %s is %v", key, result))
	if result == 0 {
		cache.logger.Debug(fmt.Sprintf("Redis: get record %s not exist", key))
		return false
	} else {
		cache.logger.Debug(fmt.Sprintf("Key %s in cache found success", key))
		return true
	}
}

// CreateCache add data in the cache
func (cache *ItemsCache) CreateItemsCache(ctx context.Context, res []models.Item, key string) error {
	cache.logger.Sugar().Debugf("Enter in cache CreateItemsCache() with args: ctx, res, key: %s", key)
	in := results{
		Responses: res,
	}
	data, err := json.Marshal(in)
	if err != nil {
		return fmt.Errorf("error on marshal items cache: %w", err)
	}

	err = cache.Set(ctx, key, data, cache.TTL).Err()
	if err != nil {
		return fmt.Errorf("redis: error on set key %s: %w", key, err)
	}
	cache.logger.Info(fmt.Sprintf("Cache with key: %s create success", key))
	return nil
}

// CreateFavouriteItemsIdCache add favourite items id in cache
func (cache *ItemsCache) CreateFavouriteItemsIdCache(ctx context.Context, res map[uuid.UUID]uuid.UUID, key string) error {
	cache.logger.Sugar().Debugf("Enter in cache CreateFavouriteItemsIdCache() with args: ctx, res, key: %s", key)
	data, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("error on marshal favourite items id cache: %w", err)
	}
	err = cache.Set(ctx, key, data, cache.TTL).Err()
	if err != nil {
		return fmt.Errorf("redis: error on set key: %s: %w", key, err)
	}
	return nil
}

// CreateItemsQuantityCache create cache for items quantity
func (cache *ItemsCache) CreateItemsQuantityCache(ctx context.Context, value int, key string) error {
	cache.logger.Sugar().Debugf("Enter in cache CreateItemsQuantityCache() with args: ctx, value: %d, key: %s", value, key)
	err := cache.Set(ctx, key, value, cache.TTL).Err()
	if err != nil {
		return fmt.Errorf("redis: error on set key %q: %w", key, err)
	}
	cache.logger.Info(fmt.Sprintf("Cache with key: %s create success", key))
	return nil
}

// GetItemsCache retrieves data from the cache
func (cache *ItemsCache) GetItemsCache(ctx context.Context, key string) ([]models.Item, error) {
	cache.logger.Sugar().Debugf("Enter in cache GetItemsCache() with args: ctx, key: %s", key)
	res := results{}
	data, err := cache.Get(ctx, key).Bytes()
	if err == redis.Nil {
		// we got empty result, it's not an error
		cache.logger.Debug("Success get nil result")
		return nil, nil
	} else if err != nil {
		cache.logger.Sugar().Errorf("Error on get cache: %v", err)
		return nil, err
	}
	err = json.Unmarshal(data, &res)
	if err != nil {
		cache.logger.Sugar().Warnf("Can't json unmarshal data: %v", data)
		return nil, err
	}
	cache.logger.Debug("Get cache success")
	return res.Responses, nil
}

// GetItemsQuantityCache retrieves data from the cache
func (cache *ItemsCache) GetItemsQuantityCache(ctx context.Context, key string) (int, error) {
	cache.logger.Sugar().Debugf("Enter in cache GetItemsQuantityCache() with args: ctx, key: %s", key)
	data, err := cache.Get(ctx, key).Int()
	if err != nil {
		cache.logger.Sugar().Errorf("Error on get cache: %v", err)
		return data, err
	}
	cache.logger.Debug("Get cache success")
	return data, nil
}

// GetItemsQuantityCache retrieves data from the cache
func (cache *ItemsCache) GetFavouriteItemsIdCache(ctx context.Context, key string) (*map[uuid.UUID]uuid.UUID, error) {
	cache.logger.Sugar().Debugf("Enter in cache GetFavouriteItemsIdCache() with args: ctx, key: %s", key)
	res := make(map[uuid.UUID]uuid.UUID)
	data, err := cache.Get(ctx, key).Bytes()
	if err == redis.Nil {
		// we got empty result, it's not an error
		cache.logger.Debug("Success get nil result")
		return nil, nil
	} else if err != nil {
		cache.logger.Sugar().Errorf("Error on get cache: %v", err)
		return nil, err
	}
	err = json.Unmarshal(data, &res)
	if err != nil {
		cache.logger.Sugar().Warnf("Can't json unmarshal data: %v", data)
		return nil, err
	}
	cache.logger.Debug("Get cache success")
	return &res, nil
}
