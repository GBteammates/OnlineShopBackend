package usecase

import (
	"OnlineShopBackend/internal/models"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// AddFavouriteItem added item in list of favourites items
func (usecase *itemUsecase) AddFavouriteItem(ctx context.Context, userId uuid.UUID, itemId uuid.UUID) error {
	usecase.logger.Sugar().Debugf("Enter in usecase AddFavouriteItem() with args: ctx, userId: %v, itemId: %v", userId, itemId)
	err := usecase.itemStore.AddFavouriteItem(ctx, userId, itemId)
	if err != nil {
		return err
	}
	usecase.UpdateFavouriteItemsCache(ctx, userId, itemId, "add")
	usecase.UpdateFavIdsCache(ctx, userId, itemId, "add")
	return nil
}

// DeleteFavouriteItem deleted item from list of favourites items
func (usecase *itemUsecase) DeleteFavouriteItem(ctx context.Context, userId uuid.UUID, itemId uuid.UUID) error {
	usecase.logger.Sugar().Debugf("Enter in usecase DeleteFavouriteItem() with args: ctx, userId: %v, itemId: %v", userId, itemId)
	err := usecase.itemStore.DeleteFavouriteItem(ctx, userId, itemId)
	if err != nil {
		return err
	}
	usecase.UpdateFavouriteItemsCache(ctx, userId, itemId, "delete")
	usecase.UpdateFavIdsCache(ctx, userId, itemId, "delete")
	return nil
}

// ItemsQuantityInFavourite check cache and if cache not exists call database
// method and write in cache and returns quantity of items in favourite
func (usecase *itemUsecase) ItemsQuantityInFavourite(ctx context.Context, userId uuid.UUID) (int, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase GetFavouriteQuantity() with args: ctx, userId: %v", userId)
	key := userId.String() + "Quantity"
	// Context with timeout so as not to wait for an answer from the cache for too long
	ctxT, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()
	// Сheck the existence of a cache with the quantity of items
	if ok := usecase.itemCache.CheckCache(ctxT, key); !ok {
		// If a cache with the quantity of items does not exist,
		// check whether there is a cache with a list of items in the basic sorting
		quantity, err := usecase.itemStore.ItemsInFavouriteQuantity(ctx, userId)
		if err != nil {
			usecase.logger.Sugar().Errorf("error on get items quantity in favourite with userId: %v from database: %v", userId, err)
			return -1, err
		}
		err = usecase.itemCache.CreateItemsQuantityCache(ctxT, quantity, key)
		if err != nil {
			usecase.logger.Sugar().Warnf("error on create items quantity in favourite cache with key: %s, err: %v", key, err)
		}
		usecase.logger.Info("Get items quantity in favourite success")
		return quantity, nil
	}
	quantity, err := usecase.itemCache.GetItemsQuantityCache(ctxT, key)
	if err != nil {
		usecase.logger.Sugar().Warnf("error on create items quantity in favourite cache with key: %s, err: %v", key, err)
		// If get cache impossible get items quantity from database
		quantity, err := usecase.itemStore.ItemsInFavouriteQuantity(ctx, userId)
		if err != nil {
			usecase.logger.Sugar().Errorf("error on get items quantity in favourite with userId: %v from database: %v", userId, err)
			return -1, err
		}
		usecase.logger.Info("Get items quantity in favourite success")
		return quantity, nil
	}
	usecase.logger.Info("Get items quantity in favourite success")
	return quantity, nil
}

// GetFavouriteItems call database method and returns chan with models.Item from list of favourites item or error
func (usecase *itemUsecase) GetFavouriteItems(ctx context.Context, userId uuid.UUID, limitOptions map[string]int, sortOptions map[string]string) ([]models.Item, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase GetFavouriteItems() with args: ctx, userId: %v", userId)

	// Context with timeout so as not to wait for an answer from the cache for too long
	ctxT, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	limit, offset := limitOptions["limit"], limitOptions["offset"]
	sortType, sortOrder := sortOptions["sortType"], sortOptions["sortOrder"]
	// Check whether there is a cache of items in favourites
	if ok := usecase.itemCache.CheckCache(ctxT, userId.String()+sortType+sortOrder); !ok {
		// If the cache does not exist, request a list of items in
		// favourites from the database
		itemIncomingChan, err := usecase.itemStore.GetFavouriteItems(ctx, userId)
		if err != nil {
			return nil, err
		}
		items := make([]models.Item, 0, 100)
		for item := range itemIncomingChan {
			items = append(items, item)
		}
		// Sort the list of items in favourites
		// based on the sorting parameters
		usecase.SortItems(items, sortType, sortOrder)
		// Create a cache with a sorted list of items in favourites
		err = usecase.itemCache.CreateItemsCache(ctxT, items, userId.String()+sortType+sortOrder)
		if err != nil {
			usecase.logger.Sugar().Warnf("error on create favourite items cache with key: %s, error: %v", userId.String()+sortType+sortOrder, err)
		} else {
			usecase.logger.Sugar().Infof("Create favourite items cache with key: %s success", userId.String()+sortType+sortOrder)
		}
		// Create a cache with a quantity of items in favourites
		err = usecase.itemCache.CreateItemsQuantityCache(ctxT, len(items), userId.String()+"Quantity")
		if err != nil {
			usecase.logger.Sugar().Warnf("error on create items in favourites quantity cache with key: %s, error: %v", userId.String()+"Quantity", err)
		} else {
			usecase.logger.Sugar().Infof("Create items in favourites quantity cache with key: %s success", userId.String()+"Quantity")
		}
	}
	// Get items list from cache
	items, err := usecase.itemCache.GetItemsCache(ctxT, userId.String()+sortType+sortOrder)
	if err != nil {
		usecase.logger.Sugar().Warnf("error on get items in favourite cache with key: %s, error: %v", userId.String()+sortType+sortOrder, err)
		// If error on get cache, request a list of items in favourite from the database
		itemIncomingChan, err := usecase.itemStore.GetFavouriteItems(ctx, userId)
		if err != nil {
			return nil, err
		}
		dbItems := make([]models.Item, 0, 100)
		for item := range itemIncomingChan {
			dbItems = append(dbItems, item)
		}
		// Sort the list of items based on the sorting parameters
		usecase.SortItems(dbItems, sortType, sortOrder)
		items = dbItems
	}
	if offset > len(items) {
		return nil, fmt.Errorf("error: offset bigger than lenght of items, offset: %d, lenght of items: %d", offset, len(items))
	}
	itemsWithLimit := make([]models.Item, 0, limit)
	var counter = 0
	for i := offset; i < len(items); i++ {
		if counter == limit {
			break
		}
		// Add items to the resulting list of items until the counter is equal to the limit
		itemsWithLimit = append(itemsWithLimit, items[i])
		counter++
	}
	return itemsWithLimit, nil
}

// GetFavouriteItemsId calls database method and returns map with identificators of favourite items of user or error
func (usecase *itemUsecase) GetFavouriteItemsId(ctx context.Context, userId uuid.UUID) (*map[uuid.UUID]uuid.UUID, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase GetFavouriteItemsId() with args: ctx, userId: %v", userId)

	// Context with timeout so as not to wait for an answer from the cache for too long
	ctxT, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	// Check whether there is a cache of identificators of favourite items
	if !usecase.itemCache.CheckCache(ctxT, userId.String()+"Fav") {
		// If the cache does not exist, request a quantity of
		// favourite items
		quantity, err := usecase.ItemsQuantityInFavourite(ctx, userId)
		if err != nil && quantity == -1 {
			usecase.logger.Warn(err.Error())
			return nil, err
		}
		if quantity == 0 {
			return nil, models.ErrorNotFound{}
		}
		// If quantity > 0 request a map with identificators of
		// favourite items from database
		favUids, err := usecase.itemStore.GetFavouriteItemsId(ctx, userId)
		if err != nil && errors.Is(err, models.ErrorNotFound{}) {
			return nil, models.ErrorNotFound{}
		}
		if err != nil {
			return nil, err
		}
		// Create cache with favourite items identificators
		err = usecase.itemCache.CreateFavouriteItemsIdCache(ctxT, *favUids, userId.String()+"Fav")
		if err != nil {
			usecase.logger.Sugar().Warnf("error on create favourite items id cache with key: %s, error: %v", userId.String()+"Fav", err)
		} else {
			usecase.logger.Sugar().Infof("Create favourite items id cache with key: %s success", userId.String()+"Fav")
		}
	}
	// Get favourite items identificators from cache
	favUids, err := usecase.itemCache.GetFavouriteItemsIdCache(ctxT, userId.String()+"Fav")
	if err != nil {
		usecase.logger.Sugar().Errorf("error on get favourite items id cache: %v", err)
		// If error on get cache, request a map of favourite items identificators
		// from the database
		dbFavUids, err := usecase.itemStore.GetFavouriteItemsId(ctx, userId)
		if err != nil && errors.Is(err, models.ErrorNotFound{}) {
			return nil, models.ErrorNotFound{}
		}
		if err != nil {
			return nil, err
		}
		favUids = dbFavUids
	}
	return favUids, nil
}

func (usecase *itemUsecase) UpdateFavouriteItemsCache(ctx context.Context, userId uuid.UUID, itemId uuid.UUID, op string) {
	usecase.logger.Sugar().Debugf("Enter in usecase UpdateFavouriteItemsCache() with args: ctx, userId: %v, itemId: %v, op: %s", userId, itemId, op)
	favouriteItemsKeyNameAsc := userId.String() + "nameasc"
	favouriteItemsKeyNameDesc := userId.String() + "namedesc"
	favouriteItemsKeyPriceAsc := userId.String() + "priceasc"
	favouriteItemsKeyPriceDesc := userId.String() + "pricedesc"
	favouriteItemsQuantityKey := userId.String() + "Quantity"

	keys := []string{favouriteItemsKeyNameAsc, favouriteItemsKeyNameDesc, favouriteItemsKeyPriceAsc, favouriteItemsKeyPriceDesc}
	// Check the presence of a cache with all possible keys
	if !usecase.itemCache.CheckCache(ctx, favouriteItemsKeyNameAsc) &&
		!usecase.itemCache.CheckCache(ctx, favouriteItemsKeyNameDesc) &&
		!usecase.itemCache.CheckCache(ctx, favouriteItemsKeyPriceAsc) &&
		!usecase.itemCache.CheckCache(ctx, favouriteItemsKeyPriceDesc) {
		// If the cache with any of the keys does not return the error
		usecase.logger.Error("cache is not exist")
		return
	}
	// Sort through all possible keys
	for _, key := range keys {
		// For each key get a cache
		items, err := usecase.itemCache.GetItemsCache(ctx, key)
		if err != nil {
			usecase.logger.Sugar().Errorf("error on get cache: %v", err)
			return
		}
		// Сhange the list of items in accordance with the operation
		if op == "add" {
			item, err := usecase.itemStore.GetItem(ctx, itemId)
			if err != nil {
				usecase.logger.Sugar().Errorf("error on get item: %v", err)
				return
			}
			items = append(items, *item)
			err = usecase.itemCache.CreateItemsQuantityCache(ctx, len(items), favouriteItemsQuantityKey)
			if err != nil {
				usecase.logger.Sugar().Errorf("error on create items quantity cache: %w", err)
				return
			}
		}
		if op == "delete" {
			for i, item := range items {
				if item.Id == itemId {
					items = append(items[:i], items[i+1:]...)
					err := usecase.itemCache.CreateItemsQuantityCache(ctx, len(items), favouriteItemsQuantityKey)
					if err != nil {
						usecase.logger.Sugar().Errorf("error on create items quantity cache: %w", err)
						return
					}
					break
				}
			}
		}
		// Sort the list of items
		switch {
		case key == favouriteItemsKeyNameAsc:
			usecase.SortItems(items, "name", "asc")
		case key == favouriteItemsKeyNameDesc:
			usecase.SortItems(items, "name", "desc")
		case key == favouriteItemsKeyPriceAsc:
			usecase.SortItems(items, "price", "asc")
		case key == favouriteItemsKeyPriceDesc:
			usecase.SortItems(items, "price", "desc")
		}
		// Record the updated cache
		err = usecase.itemCache.CreateItemsCache(ctx, items, key)
		if err != nil {
			usecase.logger.Sugar().Errorf("error on create favourite items cache: %v", err)
			return
		}
	}
	usecase.logger.Info("Update favourite items list cache success")
}

// UpdateFavIdsCache updates cache with favourite items identificators
func (usecase *itemUsecase) UpdateFavIdsCache(ctx context.Context, userId, itemId uuid.UUID, op string) {
	usecase.logger.Sugar().Debugf("Enter in usecase UpdateFavIdsCache() with args userId: %v, itemId: %v", userId, itemId)
	// Check the presence of a cache with key
	if !usecase.itemCache.CheckCache(ctx, userId.String()+"Fav") {
		// If cache doesn't exists create it
		favMap := make(map[uuid.UUID]uuid.UUID)
		// Add itemId in map of favourite
		// item's identificators
		favMap[itemId] = userId

		// Record the cache with favourite items identificators
		err := usecase.itemCache.CreateFavouriteItemsIdCache(ctx, favMap, userId.String()+"Fav")
		if err != nil {
			usecase.logger.Sugar().Warnf("error on create favourite items id cache: %v", err)
			return
		}
		usecase.logger.Info("create favourite items id cache success")
		return
	}
	// If cache exists get it
	favMapLink, err := usecase.itemCache.GetFavouriteItemsIdCache(ctx, userId.String()+"Fav")
	if err != nil {
		usecase.logger.Sugar().Warn("error on get favourite items id cache with key: %v, err: %v", userId.String()+"Fav", err)
		return
	}
	// Сhange the map of favourite items identificators
	// in accordance with the operation
	favMap := *favMapLink
	if op == "add" {
		favMap[itemId] = userId
	}
	if op == "delete" {
		delete(favMap, itemId)
	}
	// Record the updated cache
	err = usecase.itemCache.CreateFavouriteItemsIdCache(ctx, favMap, userId.String()+"Fav")
	if err != nil {
		usecase.logger.Sugar().Warn("error on create favourite items id cache: %v", err)
		return
	}
	usecase.logger.Info("Create favourite items id cache success")
}
