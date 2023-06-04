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

type itemsCache struct {
	*RedisCache
	valid  valid
	logger *zap.SugaredLogger
}

type valid struct {
	list       bool
	inCategory bool
	search     bool
	favourite  bool
}

type results struct {
	Responses []models.Item
}

type options struct {
	op      string
	kind    string
	newItem *models.Item
	userId  uuid.UUID
}

func NewItemsCache(cache *RedisCache, logger *zap.SugaredLogger) *itemsCache {
	logger.Debug("Enter in cache NewItemsCache")
	valid := valid{}
	return &itemsCache{cache, valid, logger}
}

func (cache *itemsCache) ItemsToCache(ctx context.Context, items []models.Item, kind, param string) error {
	cache.logger.Debugf("Enter in cache CreateItemsCache() with args: ctx, items, kind: %s, param: %s", kind, param)

	if err := cache.itemsToItemsListCache(ctx, param, kind, items); err != nil {
		return fmt.Errorf("error on items to cache: %w", err)
	}
	return nil
}

func (cache *itemsCache) itemsToItemsListCache(ctx context.Context, param, kind string, items []models.Item) error {
	cache.logger.Debug("Enter in cache itemsToItemsListCache() with args: ctx, param: %s, kind: %s, items", param, kind)

	err := cache.createItemsQuantityCache(ctx, len(items), param+models.Quantity)
	if err != nil {
		cache.updateStatus(ctx, kind, false)
		return err
	}

	keys := []string{
		param + models.NameASC,
		param + models.NameDESC,
		param + models.PriceASC,
		param + models.PriceDESC,
	}

	for _, key := range keys {
		sortType, sortOrder := helpers.SortOptionsFromKey(key)
		helpers.SortItems(items, sortType, sortOrder)
		err := cache.createItemsCache(ctx, items, key)
		if err != nil {
			cache.updateStatus(ctx, kind, false)
			return err
		}
	}
	cache.updateStatus(ctx, kind, true)
	return nil
}

func (cache *itemsCache) ItemsFromCache(ctx context.Context, key, kind string) ([]models.Item, error) {
	cache.logger.Debugf("Enter in cache itemsListFromCache() with args: ctx, key: %s, kind: %s", key, kind)

	if !cache.status(ctx, kind) {
		return nil, fmt.Errorf("cache with kind: %s is not valid", kind)
	}
	if !cache.checkCache(ctx, key) {
		return nil, fmt.Errorf("cache with key: %s is not exist", key)
	}
	items, err := cache.getItemsCache(ctx, key)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (cache *itemsCache) ItemsQuantityToCache(ctx context.Context, value int, key, kind string) error {
	cache.logger.Debugf("Enter in cache ItemsQuantityToCache() with args: ctx, value: %d, key: %s, kind: %s", value, key, kind)
	err := cache.createItemsQuantityCache(ctx, value, key)
	if err != nil {
		cache.updateStatus(ctx, kind, false)
		return err
	}
	return nil
}

func (cache *itemsCache) ItemsQuantityFromCache(ctx context.Context, key string, kind string) (int, error) {
	cache.logger.Debugf("Enter in cache ItemsQuantityFromCache() with args: ctx, key: %s, kind: %s", key, kind)
	if !cache.status(ctx, kind) {
		return -1, fmt.Errorf("cache with kind: %s is not valid", kind)
	}
	if !cache.checkCache(ctx, key) {
		return -1, fmt.Errorf("cache with key: %s is not exist", key)
	}
	res, err := cache.getItemsQuantityCache(ctx, key)
	if err != nil {
		return -1, fmt.Errorf("error on get items quantity cache: %w", err)
	}
	return res, nil
}

func (cache *itemsCache) FavouriteItemsIdsToCache(ctx context.Context, favIds *map[uuid.UUID]uuid.UUID, key, kind string) error {
	cache.logger.Debugf("Enter in cache FavouriteItemsIdsToCache() with args: ctx, favIds, key: %s, kind: %s", key, kind)

	err := cache.createFavouriteItemsIdCache(ctx, *favIds, key)
	if err != nil {
		cache.updateStatus(ctx, kind, false)
		return fmt.Errorf("error on create favourite items id cache: %w", err)
	}
	return nil
}

func (cache *itemsCache) FavouriteItemsIdsFromCache(ctx context.Context, key, kind string) (*map[uuid.UUID]uuid.UUID, error) {
	cache.logger.Debugf("Enter in cache FavouriteItemsIdsFromCache() with args: ctx, key: %s, kind: %s", key, kind)

	if !cache.status(ctx, kind) {
		return nil, fmt.Errorf("cache with kind: %s is not valid", kind)
	}
	if !cache.checkCache(ctx, key) {
		return nil, fmt.Errorf("cache with key: %s is not exist", key)
	}
	res, err := cache.getFavouriteItemsIdCache(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("error on get favourite items ids from cache: %w", err)
	}
	return res, nil
}

// UpdateCache updating cache when creating, updating or deleting item
func (cache *itemsCache) UpdateCache(ctx context.Context, opts *models.CacheOptions) error {
	cache.logger.Debugf("Enter in itemCache UpdateCache() with args: ctx, opts: %v", opts)

	for _, kind := range opts.Kind {
		if !cache.status(ctx, kind) {
			return fmt.Errorf("cache with kind: %s is not valid", kind)
		}
		options := options{
			op:      opts.Op,
			kind:    kind,
			newItem: opts.NewItem,
			userId:  opts.UserId,
		}
		err := cache.updateItemsCache(ctx, options)
		if err != nil {
			return fmt.Errorf("error on update cache: %w", err)
		}
	}

	return nil
}

func (cache *itemsCache) UpdateFavIdsCache(ctx context.Context, userId uuid.UUID, item *models.Item, op string) error {
	cache.logger.Debugf("Enter in cache updateFavouritesCache() with args: ctx, userId: %v, item: %v, op: %s", userId, item, op)

	err := cache.updateFavIdsCache(ctx, userId, item.Id, op)
	if err != nil {
		return err
	}

	return nil
}

// CheckCache checks for data in the cache
func (cache *itemsCache) checkCache(ctx context.Context, key string) bool {
	cache.logger.Debugf("Enter in cache CheckCache() with args: ctx, key: %s", key)
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
func (cache *itemsCache) createItemsCache(ctx context.Context, res []models.Item, key string) error {
	cache.logger.Debugf("Enter in cache CreateItemsCache() with args: ctx, res, key: %s", key)
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
func (cache *itemsCache) createFavouriteItemsIdCache(ctx context.Context, res map[uuid.UUID]uuid.UUID, key string) error {
	cache.logger.Debugf("Enter in cache CreateFavouriteItemsIdCache() with args: ctx, res, key: %s", key)
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
func (cache *itemsCache) createItemsQuantityCache(ctx context.Context, value int, key string) error {
	cache.logger.Debugf("Enter in cache CreateItemsQuantityCache() with args: ctx, value: %d, key: %s", value, key)
	err := cache.Set(ctx, key, value, cache.TTL).Err()
	if err != nil {
		return fmt.Errorf("redis: error on set key %q: %w", key, err)
	}
	cache.logger.Info(fmt.Sprintf("Cache with key: %s create success", key))
	return nil
}

// GetItemsCache retrieves data from the cache
func (cache *itemsCache) getItemsCache(ctx context.Context, key string) ([]models.Item, error) {
	cache.logger.Debugf("Enter in cache GetItemsCache() with args: ctx, key: %s", key)
	res := results{}
	data, err := cache.Get(ctx, key).Bytes()
	if err == redis.Nil {
		// we got empty result, it's not an error
		cache.logger.Debug("Success get nil result")
		return nil, nil
	} else if err != nil {
		cache.logger.Errorf("Error on get cache: %v", err)
		return nil, err
	}
	err = json.Unmarshal(data, &res)
	if err != nil {
		cache.logger.Warnf("Can't json unmarshal data: %v", data)
		return nil, err
	}
	cache.logger.Debug("Get cache success")
	return res.Responses, nil
}

// GetItemsQuantityCache retrieves data from the cache
func (cache *itemsCache) getItemsQuantityCache(ctx context.Context, key string) (int, error) {
	cache.logger.Debugf("Enter in cache GetItemsQuantityCache() with args: ctx, key: %s", key)
	data, err := cache.Get(ctx, key).Int()
	if err != nil {
		cache.logger.Errorf("Error on get cache: %v", err)
		return data, err
	}
	cache.logger.Debug("Get cache success")
	return data, nil
}

// GetItemsQuantityCache retrieves data from the cache
func (cache *itemsCache) getFavouriteItemsIdCache(ctx context.Context, key string) (*map[uuid.UUID]uuid.UUID, error) {
	cache.logger.Debugf("Enter in cache GetFavouriteItemsIdCache() with args: ctx, key: %s", key)
	res := make(map[uuid.UUID]uuid.UUID)
	data, err := cache.Get(ctx, key).Bytes()
	if err == redis.Nil {
		// we got empty result, it's not an error
		cache.logger.Debug("Success get nil result")
		return nil, nil
	} else if err != nil {
		cache.logger.Errorf("Error on get cache: %v", err)
		return nil, err
	}
	err = json.Unmarshal(data, &res)
	if err != nil {
		cache.logger.Warnf("Can't json unmarshal data: %v", data)
		return nil, err
	}
	cache.logger.Debug("Get cache success")
	return &res, nil
}

func (cache *itemsCache) updateItemsCache(ctx context.Context, opts options) error {
	cache.logger.Debugf("Enter in itemCache updateItemsListCache() with args: ctx, opts: %v", opts)

	if !cache.status(ctx, opts.kind) {
		return fmt.Errorf("cache with kind: %s is not valid", opts.kind)
	}

	keys, quantityKey := cache.getKeysByKind(ctx, opts)
	// If the cache with any of the keys does not exist return the error
	for _, key := range keys {
		if !cache.checkCache(ctx, key) {
			return fmt.Errorf("cache with key: %s is not exist", key)
		}
	}

	// Sort through all possible keys
	for _, key := range keys {
		// For each key get a cache
		items, err := cache.getItemsCache(ctx, key)
		if err != nil {
			return fmt.Errorf("error on get cache: %w", err)
		}

		// Сhange the list of items in accordance with the operation
		if opts.op == models.UpdateOp {
			for i, item := range items {
				if item.Id == opts.newItem.Id {
					items[i] = *opts.newItem
					break
				}
			}
		}
		if opts.op == models.CreateOp {
			items = append(items, *opts.newItem)
			err := cache.createItemsQuantityCache(ctx, len(items), quantityKey)
			if err != nil {
				cache.updateStatus(ctx, opts.kind, false)
				return fmt.Errorf("error on create items quantity cache: %w", err)
			}
		}
		if opts.op == models.DeleteOp {
			for i, item := range items {
				if item.Id == opts.newItem.Id {
					items = append(items[:i], items[i+1:]...)
					err := cache.createItemsQuantityCache(ctx, len(items), quantityKey)
					if err != nil {
						cache.updateStatus(ctx, opts.kind, false)
						return fmt.Errorf("error on create items quantity cache: %w", err)
					}
					break
				}
			}
		}

		// Sort the list of items
		sortType, sortOrder := helpers.SortOptionsFromKey(key)
		helpers.SortItems(items, sortType, sortOrder)
		// Record the updated cache
		err = cache.createItemsCache(ctx, items, key)
		if err != nil {
			cache.updateStatus(ctx, opts.kind, false)
			return err
		}
		cache.logger.Infof("Cache of items list with key: %s update success", key)
	}
	return nil
}

func (cache *itemsCache) getKeysByKind(ctx context.Context, opts options) ([]string, string) {
	cache.logger.Debugf("Enter in cache getKeysByKind() with args: ctx, opts: %v", opts)

	switch opts.kind {
	case models.List:
		return []string{
			models.ListItemsKey + models.NameASC,
			models.ListItemsKey + models.NameDESC,
			models.ListItemsKey + models.PriceASC,
			models.ListItemsKey + models.PriceDESC,
		}, models.ListItemsKey + models.Quantity
	case models.InCategory:
		return []string{
			opts.newItem.Category.Name + models.NameASC,
			opts.newItem.Category.Name + models.NameDESC,
			opts.newItem.Category.Name + models.PriceASC,
			opts.newItem.Category.Name + models.PriceDESC,
		}, opts.newItem.Category.Name + models.Quantity
	case models.Favourites:
		return []string{
			opts.userId.String() + models.NameASC,
			opts.userId.String() + models.NameDESC,
			opts.userId.String() + models.PriceASC,
			opts.userId.String() + models.PriceDESC,
		}, opts.userId.String() + models.Quantity
	default:
		cache.logger.Warn("Unexpected opts.kind: %s", opts.kind)
		return []string{""}, ""
	}
}

// UpdateFavIdsCache updates cache with favourite items identificators
func (cache *itemsCache) updateFavIdsCache(ctx context.Context, userId, itemId uuid.UUID, op string) error {
	cache.logger.Debugf("Enter in usecase UpdateFavIdsCache() with args userId: %v, itemId: %v", userId, itemId)
	// Check the presence of a cache with key
	if !cache.checkCache(ctx, userId.String()+models.FavIDs) {
		// If cache doesn't exists create it
		favMap := make(map[uuid.UUID]uuid.UUID)
		// Add itemId in map of favourite
		// item's identificators
		favMap[itemId] = userId

		// Record the cache with favourite items identificators
		err := cache.createFavouriteItemsIdCache(ctx, favMap, userId.String()+models.FavIDs)
		if err != nil {
			return err
		}
		cache.logger.Info("create favourite items id cache success")
		return nil
	}
	// If cache exists get it
	favMapLink, err := cache.getFavouriteItemsIdCache(ctx, userId.String()+models.FavIDs)
	if err != nil {
		cache.logger.Warn("error on get favourite items id cache with key: %v, err: %v", userId.String()+models.FavIDs, err)
		return err
	}
	// Сhange the map of favourite items identificators
	// in accordance with the operation
	favMap := *favMapLink
	if op == models.CreateOp {
		favMap[itemId] = userId
	}
	if op == models.DeleteOp {
		delete(favMap, itemId)
	}
	// Record the updated cache
	err = cache.createFavouriteItemsIdCache(ctx, favMap, userId.String()+models.FavIDs)
	if err != nil {
		return err
	}
	cache.logger.Info("Create favourite items id cache success")
	return nil
}

func (cache *itemsCache) status(ctx context.Context, kind string) bool {
	cache.logger.Debugf("Enter in cache Status() with args: ctx, kind: %s", kind)
	switch kind {
	case models.List:
		return cache.valid.list
	case models.InCategory:
		return cache.valid.inCategory
	case models.Search:
		return cache.valid.search
	case models.Favourites:
		return cache.valid.favourite
	default:
		return false
	}
}

func (cache *itemsCache) updateStatus(ctx context.Context, kind string, status bool) {
	cache.logger.Debugf("Enter in cache ChangeStatus() with args: ctx, kind: %s, status: %t", kind, status)
	switch kind {
	case models.List:
		cache.valid.list = status
	case models.InCategory:
		cache.valid.inCategory = status
	case models.Search:
		cache.valid.search = status
	case models.Favourites:
		cache.valid.favourite = status
	}
}
