package redis

import (
	"OnlineShopBackend/internal/helpers"
	"OnlineShopBackend/internal/models"
	usecase "OnlineShopBackend/internal/usecase/interfaces"
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var _ usecase.IItemsCache = &itemsCache{}

// Keys for create and get cache
const (
	itemsListKey          = "ItemsList"
	itemsListKeyNameAsc   = "ItemsListnameasc"
	itemsListKeyNameDesc  = "ItemsListnamedesc"
	itemsListKeyPriceAsc  = "ItemsListpriceasc"
	itemsListKeyPriceDesc = "ItemsListpricedesc"
	itemsQuantityKey      = "ItemsQuantity"
	asc                   = "asc"
	desc                  = "desc"
	nameAsc               = "nameasc"
	nameDesc              = "namedesc"
	priceAsc              = "priceasc"
	priceDesc             = "pricedesc"
	quantity              = "Quantity"
	price                 = "price"
	name                  = "name"
	updateOp              = "update"
	createOp              = "create"
	deleteOp              = "delete"
	addOp                 = "add"
	fav                   = "Fav"
)

type itemsCache struct {
	*RedisCache
	cacheValid bool
	logger     *zap.Logger
}

type results struct {
	Responses []models.Item
}

func NewItemsCache(cache *RedisCache, logger *zap.Logger) *itemsCache {
	logger.Debug("Enter in cache NewItemsCache")
	return &itemsCache{cache, false, logger}
}

// CheckCache checks for data in the cache
func (cache *itemsCache) CheckCache(ctx context.Context, key string) bool {
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
func (cache *itemsCache) CreateItemsCache(ctx context.Context, res []models.Item, key string) error {
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
func (cache *itemsCache) CreateFavouriteItemsIdCache(ctx context.Context, res map[uuid.UUID]uuid.UUID, key string) error {
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
func (cache *itemsCache) CreateItemsQuantityCache(ctx context.Context, value int, key string) error {
	cache.logger.Sugar().Debugf("Enter in cache CreateItemsQuantityCache() with args: ctx, value: %d, key: %s", value, key)
	err := cache.Set(ctx, key, value, cache.TTL).Err()
	if err != nil {
		return fmt.Errorf("redis: error on set key %q: %w", key, err)
	}
	cache.logger.Info(fmt.Sprintf("Cache with key: %s create success", key))
	return nil
}

// GetItemsCache retrieves data from the cache
func (cache *itemsCache) GetItemsCache(ctx context.Context, key string) ([]models.Item, error) {
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
func (cache *itemsCache) GetItemsQuantityCache(ctx context.Context, key string) (int, error) {
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
func (cache *itemsCache) GetFavouriteItemsIdCache(ctx context.Context, key string) (*map[uuid.UUID]uuid.UUID, error) {
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

// UpdateCache updating cache when creating, updating or deleting item
func (cache *itemsCache) UpdateCache(ctx context.Context, newItem *models.Item, op string) error {
	cache.logger.Sugar().Debugf("Enter in itemCache UpdateCache() with args: ctx, newItem: %v, op: %s", newItem, op)
	err := cache.updateItemsListCache(ctx, newItem, op)
	if err != nil {
		return err
	}
	err = cache.updateItemsInCategoryCache(ctx, newItem, op)
	if err != nil {
		return err
	}

	return nil
}

func (cache *itemsCache) updateItemsListCache(ctx context.Context, newItem *models.Item, op string) error {
	cache.logger.Sugar().Debugf("Enter in itemCache updateItemsListCache() with args: ctx, item: %v, op: %s", newItem, op)

	// Check the presence of a cache with all possible keys
	if !cache.CheckCache(ctx, itemsListKeyNameAsc) ||
		!cache.CheckCache(ctx, itemsListKeyNameDesc) ||
		!cache.CheckCache(ctx, itemsListKeyPriceAsc) ||
		!cache.CheckCache(ctx, itemsListKeyPriceDesc) {
		// If the cache with any of the keys does not return the error
		return fmt.Errorf("cache is not exists")
	}
	cacheKeys := []string{itemsListKeyNameAsc, itemsListKeyNameDesc, itemsListKeyPriceAsc, itemsListKeyPriceDesc}
	// Sort through all possible keys
	for _, key := range cacheKeys {
		// For each key get a cache
		items, err := cache.GetItemsCache(ctx, key)
		if err != nil {
			return fmt.Errorf("error on get cache: %w", err)
		}

		// Сhange the list of items in accordance with the operation
		if op == updateOp {
			for i, item := range items {
				if item.Id == newItem.Id {
					items[i] = *newItem
					break
				}
			}
		}
		if op == createOp {
			items = append(items, *newItem)
			err := cache.CreateItemsQuantityCache(ctx, len(items), itemsQuantityKey)
			if err != nil {
				return fmt.Errorf("error on create items quantity cache: %w", err)
			}
		}
		if op == deleteOp {
			for i, item := range items {
				if item.Id == newItem.Id {
					items = append(items[:i], items[i+1:]...)
					err := cache.CreateItemsQuantityCache(ctx, len(items), itemsQuantityKey)
					if err != nil {
						return fmt.Errorf("error on create items quantity cache: %w", err)
					}
					break
				}
			}
		}

		// Sort the list of items
		switch {
		case key == itemsListKeyNameAsc:
			helpers.SortItems(items, name, asc)
		case key == itemsListKeyNameDesc:
			helpers.SortItems(items, name, desc)
		case key == itemsListKeyPriceAsc:
			helpers.SortItems(items, price, asc)
		case key == itemsListKeyPriceDesc:
			helpers.SortItems(items, price, desc)
		}
		// Record the updated cache
		err = cache.CreateItemsCache(ctx, items, key)
		if err != nil {
			return err
		}
		cache.logger.Sugar().Infof("Cache of items list with key: %s update success", key)
	}
	return nil
}

// updateItemsInCategoryCache update cache items from category
func (cache *itemsCache) updateItemsInCategoryCache(ctx context.Context, newItem *models.Item, op string) error {
	cache.logger.Debug(fmt.Sprintf("Enter in cache updateItemsInCategoryCache() with args: ctx, newItem: %v, op: %s", newItem, op))
	categoryItemsKeyNameAsc := newItem.Category.Name + nameAsc
	categoryItemsKeyNameDesc := newItem.Category.Name + nameDesc
	categoryItemsKeyPriceAsc := newItem.Category.Name + priceAsc
	categoryItemsKeyPriceDesc := newItem.Category.Name + priceDesc
	categoryItemsQuantityKey := newItem.Category.Name + quantity

	keys := []string{categoryItemsKeyNameAsc, categoryItemsKeyNameDesc, categoryItemsKeyPriceAsc, categoryItemsKeyPriceDesc}
	// Check the presence of a cache with all possible keys
	if !cache.CheckCache(ctx, categoryItemsKeyNameAsc) ||
		!cache.CheckCache(ctx, categoryItemsKeyNameDesc) ||
		!cache.CheckCache(ctx, categoryItemsKeyPriceAsc) ||
		!cache.CheckCache(ctx, categoryItemsKeyPriceDesc) {
		// If the cache with any of the keys does not return the error
		return fmt.Errorf("cache is not exist")
	}
	// Sort through all possible keys
	for _, key := range keys {
		// For each key get a cache
		items, err := cache.GetItemsCache(ctx, key)
		if err != nil {
			return fmt.Errorf("error on get cache: %w", err)
		}
		// Сhange the list of items in accordance with the operation
		if op == updateOp {
			for i, item := range items {
				if item.Id == newItem.Id {
					items[i] = *newItem
					break
				}
			}
		}
		if op == createOp {
			items = append(items, *newItem)
			err := cache.CreateItemsQuantityCache(ctx, len(items), categoryItemsQuantityKey)
			if err != nil {
				return fmt.Errorf("error on create items quantity cache: %w", err)
			}
		}
		if op == deleteOp {
			for i, item := range items {
				if item.Id == newItem.Id {
					items = append(items[:i], items[i+1:]...)
					err := cache.CreateItemsQuantityCache(ctx, len(items), categoryItemsQuantityKey)
					if err != nil {
						return fmt.Errorf("error on create items quantity cache: %w", err)
					}
					break
				}
			}
		}
		// Sort the list of items
		switch {
		case key == categoryItemsKeyNameAsc:
			helpers.SortItems(items, name, asc)
		case key == categoryItemsKeyNameDesc:
			helpers.SortItems(items, name, desc)
		case key == categoryItemsKeyPriceAsc:
			helpers.SortItems(items, price, asc)
		case key == categoryItemsKeyPriceDesc:
			helpers.SortItems(items, price, desc)
		}
		// Record the updated cache
		err = cache.CreateItemsCache(ctx, items, key)
		if err != nil {
			return err
		}
	}
	cache.logger.Info("Update category list cache success")
	return nil
}

func (cache *itemsCache) UpdateFavouritesCache(ctx context.Context, userId uuid.UUID, item *models.Item, op string) error {
	cache.logger.Sugar().Debugf("Enter in cache updateFavouritesCache() with args: ctx, userId: %v, item: %v, op: %s", userId, item, op)

	err := cache.updateFavouriteItemsCache(ctx, userId, item, op)
	if err != nil {
		return err
	}
	err = cache.updateFavIdsCache(ctx, userId, item.Id, op)
	if err != nil {
		return err
	}

	return nil
}

func (cache *itemsCache) updateFavouriteItemsCache(ctx context.Context, userId uuid.UUID, item *models.Item, op string) error {
	cache.logger.Sugar().Debugf("Enter in usecase UpdateFavouriteItemsCache() with args: ctx, userId: %v, item: %v, op: %s", userId, item, op)
	favouriteItemsKeyNameAsc := userId.String() + nameAsc
	favouriteItemsKeyNameDesc := userId.String() + nameDesc
	favouriteItemsKeyPriceAsc := userId.String() + priceAsc
	favouriteItemsKeyPriceDesc := userId.String() + priceDesc
	favouriteItemsQuantityKey := userId.String() + quantity

	keys := []string{favouriteItemsKeyNameAsc, favouriteItemsKeyNameDesc, favouriteItemsKeyPriceAsc, favouriteItemsKeyPriceDesc}
	// Check the presence of a cache with all possible keys
	if !cache.CheckCache(ctx, favouriteItemsKeyNameAsc) ||
		!cache.CheckCache(ctx, favouriteItemsKeyNameDesc) ||
		!cache.CheckCache(ctx, favouriteItemsKeyPriceAsc) ||
		!cache.CheckCache(ctx, favouriteItemsKeyPriceDesc) {
		// If the cache with any of the keys does not return the error
		return fmt.Errorf("cache is not exist")
	}
	// Sort through all possible keys
	for _, key := range keys {
		// For each key get a cache
		items, err := cache.GetItemsCache(ctx, key)
		if err != nil {
			cache.logger.Sugar().Errorf("error on get cache: %v", err)
			return err
		}
		// Сhange the list of items in accordance with the operation
		if op == addOp {
			items = append(items, *item)
			err = cache.CreateItemsQuantityCache(ctx, len(items), favouriteItemsQuantityKey)
			if err != nil {
				cache.logger.Sugar().Errorf("error on create items quantity cache: %w", err)
				return err
			}
		}
		if op == deleteOp {
			for i, itm := range items {
				if itm.Id == item.Id {
					items = append(items[:i], items[i+1:]...)
					err := cache.CreateItemsQuantityCache(ctx, len(items), favouriteItemsQuantityKey)
					if err != nil {
						return err
					}
					break
				}
			}
		}
		// Sort the list of items
		switch {
		case key == favouriteItemsKeyNameAsc:
			helpers.SortItems(items, name, asc)
		case key == favouriteItemsKeyNameDesc:
			helpers.SortItems(items, name, desc)
		case key == favouriteItemsKeyPriceAsc:
			helpers.SortItems(items, price, asc)
		case key == favouriteItemsKeyPriceDesc:
			helpers.SortItems(items, price, asc)
		}
		// Record the updated cache
		err = cache.CreateItemsCache(ctx, items, key)
		if err != nil {
			return err
		}
	}
	cache.logger.Info("Update favourite items list cache success")
	return nil
}

// UpdateFavIdsCache updates cache with favourite items identificators
func (cache *itemsCache) updateFavIdsCache(ctx context.Context, userId, itemId uuid.UUID, op string) error {
	cache.logger.Sugar().Debugf("Enter in usecase UpdateFavIdsCache() with args userId: %v, itemId: %v", userId, itemId)
	// Check the presence of a cache with key
	if !cache.CheckCache(ctx, userId.String()+fav) {
		// If cache doesn't exists create it
		favMap := make(map[uuid.UUID]uuid.UUID)
		// Add itemId in map of favourite
		// item's identificators
		favMap[itemId] = userId

		// Record the cache with favourite items identificators
		err := cache.CreateFavouriteItemsIdCache(ctx, favMap, userId.String()+fav)
		if err != nil {
			return err
		}
		cache.logger.Info("create favourite items id cache success")
		return nil
	}
	// If cache exists get it
	favMapLink, err := cache.GetFavouriteItemsIdCache(ctx, userId.String()+fav)
	if err != nil {
		cache.logger.Sugar().Warn("error on get favourite items id cache with key: %v, err: %v", userId.String()+fav, err)
		return err
	}
	// Сhange the map of favourite items identificators
	// in accordance with the operation
	favMap := *favMapLink
	if op == addOp {
		favMap[itemId] = userId
	}
	if op == deleteOp {
		delete(favMap, itemId)
	}
	// Record the updated cache
	err = cache.CreateFavouriteItemsIdCache(ctx, favMap, userId.String()+fav)
	if err != nil {
		return err
	}
	cache.logger.Info("Create favourite items id cache success")
	return nil
}
