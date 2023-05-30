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
	list                  = "list"
	inCategory            = "inCategory"
	search                = "search"
	favourite             = "favourite"
)

type itemsCache struct {
	*RedisCache
	cacheValid cacheValid
	logger     *zap.Logger
}

type cacheValid struct {
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

func NewItemsCache(cache *RedisCache, logger *zap.Logger) *itemsCache {
	logger.Debug("Enter in cache NewItemsCache")
	cacheValid := cacheValid{}
	return &itemsCache{cache, cacheValid, logger}
}

func (cache *itemsCache) ItemsToCache(ctx context.Context, items []models.Item, kind, param string) error {
	cache.logger.Sugar().Debugf("Enter in cache CreateItemsCache() with args: ctx, items, kind: %s, param: %s", kind, param)

	switch kind {
	case list:
		return cache.itemsToItemsListCache(ctx, items)
	case inCategory:
		return cache.itemsToItemsByCategoryListCache(ctx, items, param)
	case search:
		return cache.itemsToItemsInSearchCache(ctx, items, param)
	default:
		return nil
	}
}

func (cache *itemsCache) itemsToItemsListCache(ctx context.Context, items []models.Item) error {
	cache.logger.Debug("Enter in cache itemsToItemsListCache()")

	err := cache.createItemsQuantityCache(ctx, len(items), itemsQuantityKey)
	if err != nil {
		cache.updateStatus(ctx, list, false)
		return err
	}

	keys := []string{itemsListKeyNameAsc, itemsListKeyNameDesc, itemsListKeyPriceAsc, itemsListKeyPriceDesc}
	for _, key := range keys {
		sortType, sortOrder := helpers.SortOptionsFromKey(key)
		helpers.SortItems(items, sortType, sortOrder)
		err := cache.createItemsCache(ctx, items, key)
		if err != nil {
			cache.updateStatus(ctx, list, false)
			return err
		}
	}
	cache.updateStatus(ctx, list, true)
	return nil
}

func (cache *itemsCache) itemsToItemsByCategoryListCache(ctx context.Context, items []models.Item, categoryName string) error {
	cache.logger.Sugar().Debugf("Enter in cache itemsToItemsByCategoryListCache() with args: ctx, items, categoryName: %s", categoryName)

	err := cache.createItemsQuantityCache(ctx, len(items), categoryName+quantity)
	if err != nil {
		cache.updateStatus(ctx, inCategory, false)
		return err
	}

	categoryItemsKeyNameAsc := categoryName + nameAsc
	categoryItemsKeyNameDesc := categoryName + nameDesc
	categoryItemsKeyPriceAsc := categoryName + priceAsc
	categoryItemsKeyPriceDesc := categoryName + priceDesc
	keys := []string{categoryItemsKeyNameAsc, categoryItemsKeyNameDesc, categoryItemsKeyPriceAsc, categoryItemsKeyPriceDesc}
	for _, key := range keys {
		sortType, sortOrder := helpers.SortOptionsFromKey(key)
		helpers.SortItems(items, sortType, sortOrder)
		err := cache.createItemsCache(ctx, items, key)
		if err != nil {
			cache.updateStatus(ctx, inCategory, false)
			return err
		}
	}

	cache.updateStatus(ctx, inCategory, true)
	return nil
}

func (cache *itemsCache) itemsToItemsInSearchCache(ctx context.Context, items []models.Item, searchRequest string) error {
	cache.logger.Sugar().Debugf("Enter in cache itemsToItemsInSearchCache() with args: ctx, items, searchRequest: %s", searchRequest)

	err := cache.createItemsQuantityCache(ctx, len(items), searchRequest+quantity)
	if err != nil {
		cache.updateStatus(ctx, search, false)
		return err
	}

	searchKeyNameAsc := searchRequest + nameAsc
	searchKeyNameDesc := searchRequest + nameDesc
	searchKeyPriceAsc := searchRequest + priceAsc
	searchKeyPriceDesc := searchRequest + priceDesc
	keys := []string{searchKeyNameAsc, searchKeyNameDesc, searchKeyPriceAsc, searchKeyPriceDesc}
	for _, key := range keys {
		sortType, sortOrder := helpers.SortOptionsFromKey(key)
		helpers.SortItems(items, sortType, sortOrder)
		err := cache.createItemsCache(ctx, items, key)
		if err != nil {
			cache.updateStatus(ctx, search, false)
			return err
		}
	}

	cache.updateStatus(ctx, search, true)
	return nil
}

func (cache *itemsCache) ItemsFromCache(ctx context.Context, cacheKey, kind string) ([]models.Item, error) {
	cache.logger.Sugar().Debugf("Enter in usecase itemsListFromCache() with args: ctx, cacheKey: %s, kind: %s", cacheKey, kind)

	if !cache.status(ctx, kind) {
		return nil, fmt.Errorf("cache with kind: %s is not valid", kind)
	}
	if !cache.checkCache(ctx, cacheKey) {
		return nil, fmt.Errorf("cache with key: %s is not exist", cacheKey)
	}
	items, err := cache.getItemsCache(ctx, cacheKey)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (cache *itemsCache) ItemsQuantityToCache(ctx context.Context, value int, key, kind string) error {
	cache.logger.Sugar().Debugf("Enter in cache ItemsQuantityToCache() with args: ctx, value: %d, key: %s, kind: %s", value, key, kind)
	err := cache.createItemsQuantityCache(ctx, value, key)
	if err != nil {
		cache.updateStatus(ctx, kind, false)
		return err
	}
	return nil
}

func (cache *itemsCache) ItemsQuantityFromCache(ctx context.Context, key string, kind string) (int, error) {
	cache.logger.Sugar().Debugf("Enter in cache ItemsQuantityFromCache() with args: ctx, key: %s, kind: %s", key, kind)
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
	cache.logger.Sugar().Debugf("Enter in cache FavouriteItemsIdsToCache() with args: ctx, favIds, key: %s, kind: %s", key, kind)

	err := cache.createFavouriteItemsIdCache(ctx, *favIds, key)
	if err != nil {
		cache.updateStatus(ctx, kind, false)
		return fmt.Errorf("error on create favourite items id cache: %w", err)
	}
	return nil
}

func (cache *itemsCache) FavouriteItemsIdsFromCache(ctx context.Context, key, kind string) (*map[uuid.UUID]uuid.UUID, error) {
	cache.logger.Sugar().Debugf("Enter in cache FavouriteItemsIdsFromCache() with args: ctx, key: %s, kind: %s", key, kind)

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
func (cache *itemsCache) UpdateCache(ctx context.Context, opts *models.ItemsCacheOptions) error {
	cache.logger.Sugar().Debugf("Enter in itemCache UpdateCache() with args: ctx, opts: %v", opts)

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
	cache.logger.Sugar().Debugf("Enter in cache updateFavouritesCache() with args: ctx, userId: %v, item: %v, op: %s", userId, item, op)

	err := cache.updateFavIdsCache(ctx, userId, item.Id, op)
	if err != nil {
		return err
	}

	return nil
}

// CheckCache checks for data in the cache
func (cache *itemsCache) checkCache(ctx context.Context, key string) bool {
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
func (cache *itemsCache) createItemsCache(ctx context.Context, res []models.Item, key string) error {
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
func (cache *itemsCache) createFavouriteItemsIdCache(ctx context.Context, res map[uuid.UUID]uuid.UUID, key string) error {
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
func (cache *itemsCache) createItemsQuantityCache(ctx context.Context, value int, key string) error {
	cache.logger.Sugar().Debugf("Enter in cache CreateItemsQuantityCache() with args: ctx, value: %d, key: %s", value, key)
	err := cache.Set(ctx, key, value, cache.TTL).Err()
	if err != nil {
		return fmt.Errorf("redis: error on set key %q: %w", key, err)
	}
	cache.logger.Info(fmt.Sprintf("Cache with key: %s create success", key))
	return nil
}

// GetItemsCache retrieves data from the cache
func (cache *itemsCache) getItemsCache(ctx context.Context, key string) ([]models.Item, error) {
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
func (cache *itemsCache) getItemsQuantityCache(ctx context.Context, key string) (int, error) {
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
func (cache *itemsCache) getFavouriteItemsIdCache(ctx context.Context, key string) (*map[uuid.UUID]uuid.UUID, error) {
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

func (cache *itemsCache) updateItemsCache(ctx context.Context, opts options) error {
	cache.logger.Sugar().Debugf("Enter in itemCache updateItemsListCache() with args: ctx, opts: %v", opts)

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
		if opts.op == updateOp {
			for i, item := range items {
				if item.Id == opts.newItem.Id {
					items[i] = *opts.newItem
					break
				}
			}
		}
		if opts.op == createOp {
			items = append(items, *opts.newItem)
			err := cache.createItemsQuantityCache(ctx, len(items), quantityKey)
			if err != nil {
				cache.updateStatus(ctx, opts.kind, false)
				return fmt.Errorf("error on create items quantity cache: %w", err)
			}
		}
		if opts.op == deleteOp {
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
		cache.logger.Sugar().Infof("Cache of items list with key: %s update success", key)
	}
	return nil
}

func (cache *itemsCache) getKeysByKind(ctx context.Context, opts options) ([]string, string) {
	cache.logger.Sugar().Debugf("Enter in cache getKeysByKind() with args: ctx, opts: %v", opts)

	switch opts.kind {
	case list:
		return []string{itemsListKeyNameAsc, itemsListKeyNameDesc, itemsListKeyPriceAsc, itemsListKeyPriceDesc}, itemsQuantityKey
	case inCategory:
		return []string{opts.newItem.Category.Name + nameAsc, opts.newItem.Category.Name + nameDesc, opts.newItem.Category.Name + priceAsc, opts.newItem.Category.Name + priceDesc}, opts.newItem.Category.Name + quantity
	case favourite:
		return []string{opts.userId.String() + nameAsc, opts.userId.String() + nameDesc, opts.userId.String() + priceAsc, opts.userId.String() + priceDesc}, opts.userId.String() + quantity
	default:
		cache.logger.Warn("Unexpected kind in getKeysByKind()")
		return []string{}, ""
	}
}

// UpdateFavIdsCache updates cache with favourite items identificators
func (cache *itemsCache) updateFavIdsCache(ctx context.Context, userId, itemId uuid.UUID, op string) error {
	cache.logger.Sugar().Debugf("Enter in usecase UpdateFavIdsCache() with args userId: %v, itemId: %v", userId, itemId)
	// Check the presence of a cache with key
	if !cache.checkCache(ctx, userId.String()+fav) {
		// If cache doesn't exists create it
		favMap := make(map[uuid.UUID]uuid.UUID)
		// Add itemId in map of favourite
		// item's identificators
		favMap[itemId] = userId

		// Record the cache with favourite items identificators
		err := cache.createFavouriteItemsIdCache(ctx, favMap, userId.String()+fav)
		if err != nil {
			return err
		}
		cache.logger.Info("create favourite items id cache success")
		return nil
	}
	// If cache exists get it
	favMapLink, err := cache.getFavouriteItemsIdCache(ctx, userId.String()+fav)
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
	err = cache.createFavouriteItemsIdCache(ctx, favMap, userId.String()+fav)
	if err != nil {
		return err
	}
	cache.logger.Info("Create favourite items id cache success")
	return nil
}

func (cache *itemsCache) status(ctx context.Context, kind string) bool {
	cache.logger.Sugar().Debugf("Enter in cache Status() with args: ctx, kind: %s", kind)
	switch kind {
	case list:
		return cache.cacheValid.list
	case inCategory:
		return cache.cacheValid.inCategory
	case search:
		return cache.cacheValid.search
	case favourite:
		return cache.cacheValid.favourite
	default:
		return false
	}
}

func (cache *itemsCache) updateStatus(ctx context.Context, kind string, status bool) {
	cache.logger.Sugar().Debugf("Enter in cache ChangeStatus() with args: ctx, kind: %s, status: %t", kind, status)
	switch kind {
	case list:
		cache.cacheValid.list = status
	case inCategory:
		cache.cacheValid.inCategory = status
	case search:
		cache.cacheValid.search = status
	case favourite:
		cache.cacheValid.favourite = status
	}
}
