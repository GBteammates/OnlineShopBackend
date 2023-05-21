package usecase

import (
	"OnlineShopBackend/internal/helpers"
	"OnlineShopBackend/internal/models"
	usecase "OnlineShopBackend/internal/usecase/interfaces"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var _ usecase.IItemUsecase = (*itemUsecase)(nil)

// Keys for create and get cache
const (
	itemsListKey          = "ItemsList"
	itemsListKeyNameAsc   = "ItemsListnameasc"
	itemsListKeyNameDesc  = "ItemsListnamedesc"
	itemsListKeyPriceAsc  = "ItemsListpriceasc"
	itemsListKeyPriceDesc = "ItemsListpricedesc"
	itemsQuantityKey      = "ItemsQuantity"
	timeout               = 100
	createOp              = "create"
	updateOp              = "update"
	deleteOp              = "delete"
	quantityKey           = "Quantity"
)

type itemUsecase struct {
	validCache  bool
	itemStore   usecase.ItemStore
	itemCache   usecase.IItemsCache
	filestorage usecase.FileStorager
	logger      *zap.Logger
}

func NewItemUsecase(itemStore usecase.ItemStore, itemCache usecase.IItemsCache, filestorage usecase.FileStorager, logger *zap.Logger) *itemUsecase {
	logger.Debug("Enter in usecase NewItemUsecase()")
	return &itemUsecase{
		itemStore:   itemStore,
		itemCache:   itemCache,
		filestorage: filestorage,
		logger:      logger}
}

// CreateItem call database method and returns id of created item or error
func (usecase *itemUsecase) CreateItem(ctx context.Context, item *models.Item) (uuid.UUID, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase CreateItem() with args: ctx, item: %v", item)
	id, err := usecase.itemStore.CreateItem(ctx, item)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error on create item: %w", err)
	}

	item.Id = id

	err = usecase.itemCache.UpdateCache(ctx, item, createOp)
	if err != nil {
		usecase.validCache = false
		usecase.logger.Sugar().Errorf("error on update cache: %v", err)
	}

	return id, nil
}

// UpdateItem call database method to update item and returns error or nil
func (usecase *itemUsecase) UpdateItem(ctx context.Context, item *models.Item) error {
	usecase.logger.Sugar().Debugf("Enter in usecase UpdateItem() with args: ctx, item: %v", item)
	err := usecase.itemStore.UpdateItem(ctx, item)
	if err != nil {
		return fmt.Errorf("error on update item: %w", err)
	}
	err = usecase.itemCache.UpdateCache(ctx, item, updateOp)
	if err != nil {
		usecase.validCache = false
		usecase.logger.Sugar().Errorf("error on update cache: %v", err)
	}

	return nil
}

// GetItem call database and returns *models.Item with given id or returns error
func (usecase *itemUsecase) GetItem(ctx context.Context, id uuid.UUID) (*models.Item, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase GetItem() with args: ctx, id: %v", id)
	item, err := usecase.itemStore.GetItem(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error on get item: %w", err)
	}
	return item, nil
}

// DeleteItem call database method for deleting item
func (usecase *itemUsecase) DeleteItem(ctx context.Context, id uuid.UUID) error {
	usecase.logger.Sugar().Debugf("Enter in usecase DeleteItem() with args: ctx, id: %v", id)
	err := usecase.itemStore.DeleteItem(ctx, id)
	if err != nil {
		return err
	}
	err = usecase.itemCache.UpdateCache(ctx, &models.Item{Id: id}, "delete")
	if err != nil {
		usecase.logger.Error(fmt.Sprintf("error on update cache: %v", err))
	}
	return nil
}

// ItemsQuantity check cache and if cache not exists call database
// method and write in cache and returns quantity of all items
func (usecase *itemUsecase) ItemsQuantity(ctx context.Context) (int, error) {
	usecase.logger.Debug("Enter in usecase ItemsQuantity() with args: ctx")
	// Context with timeout so as not to wait for an answer from the cache for too long
	ctxT, cancel := context.WithTimeout(ctx, timeout*time.Millisecond)
	defer cancel()
	// Сheck the existence of a cache with the quantity of items
	if ok := usecase.itemCache.CheckCache(ctxT, itemsQuantityKey); !ok {
		// If a cache with the quantity of items does not exist,
		// check whether there is a cache with a list of items in the basic sorting
		quantity, err := usecase.itemStore.ItemsListQuantity(ctx)
		if err != nil {
			usecase.logger.Sugar().Errorf("error on get items list quantity from database: %v", err)
			return -1, err
		}
		err = usecase.itemCache.CreateItemsQuantityCache(ctxT, quantity, itemsQuantityKey)
		if err != nil {
			usecase.logger.Sugar().Warnf("error on create items quantity cache with key: %s, err: %v", itemsQuantityKey, err)
		}
		usecase.logger.Info("Get items quantity success")
		return quantity, nil
	}
	quantity, err := usecase.itemCache.GetItemsQuantityCache(ctxT, itemsQuantityKey)
	if err != nil {
		usecase.logger.Sugar().Warnf("error on get items quantity cache with key: %s, error: %v", itemsQuantityKey, err)
		// If get cache impossible get items quantity from database
		quantity, err := usecase.itemStore.ItemsListQuantity(ctx)
		if err != nil {
			usecase.logger.Sugar().Errorf("error on get items list quantity from database: %v", err)
			return -1, err
		}
		usecase.logger.Info("Get items quantity success")
		return quantity, nil
	}
	usecase.logger.Info("Get items quantity success")
	return quantity, nil
}

// ItemsQuantityInCategory check cache and if cache not exists call database
// method and write in cache and returns quantity of items in category
func (usecase *itemUsecase) ItemsQuantityInCategory(ctx context.Context, categoryName string) (int, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase ItemsQuantityInCategory() with args: ctx, categoryName: %s", categoryName)
	key := categoryName + quantityKey
	// Context with timeout so as not to wait for an answer from the cache for too long
	ctxT, cancel := context.WithTimeout(ctx, timeout*time.Millisecond)
	defer cancel()
	// Сheck the existence of a cache with the quantity of items
	if ok := usecase.itemCache.CheckCache(ctxT, key); !ok {
		// If a cache with the quantity of items does not exist,
		// check whether there is a cache with a list of items in the basic sorting
		quantity, err := usecase.itemStore.ItemsByCategoryQuantity(ctx, categoryName)
		if err != nil {
			usecase.logger.Sugar().Errorf("error on get items quantity in category from database: %v", err)
			return -1, err
		}
		err = usecase.itemCache.CreateItemsQuantityCache(ctxT, quantity, key)
		if err != nil {
			usecase.logger.Sugar().Warnf("error on create items quantity in category  cache with key: %s, err: %v", key, err)
		}
		usecase.logger.Info("Get items quantity in category success")
		return quantity, nil
	}
	quantity, err := usecase.itemCache.GetItemsQuantityCache(ctxT, key)
	if err != nil {
		usecase.logger.Sugar().Warnf("error on get items quantity in category cache with key: %s, error: %v", key, err)
		// If get cache impossible get items quantity from database
		quantity, err := usecase.itemStore.ItemsByCategoryQuantity(ctx, categoryName)
		if err != nil {
			usecase.logger.Sugar().Errorf("error on get items quantity in category from database: %v", err)
			return -1, err
		}
		usecase.logger.Info("Get items quantity in category success")
		return quantity, nil
	}
	usecase.logger.Info("Get items quantity in category success")
	return quantity, nil
}

// ItemsQuantityInSearch check cache and if cache not exists call database method and write
// in cache and returns quantity of items in search request
func (usecase *itemUsecase) ItemsQuantityInSearch(ctx context.Context, searchRequest string) (int, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase ItemsQuantityInSearch() with args: ctx, searchRequest: %s", searchRequest)
	key := searchRequest + quantityKey
	// Context with timeout so as not to wait for an answer from the cache for too long
	ctxT, cancel := context.WithTimeout(ctx, timeout*time.Millisecond)
	defer cancel()
	// Сheck the existence of a cache with the quantity of items
	if ok := usecase.itemCache.CheckCache(ctxT, key); !ok {
		// If a cache with the quantity of items does not exist,
		// check whether there is a cache with a list of items in the basic sorting
		quantity, err := usecase.itemStore.ItemsInSearchQuantity(ctx, searchRequest)
		if err != nil {
			usecase.logger.Sugar().Errorf("error on get items quantity in search by search request: %s from database: %v", searchRequest, err)
			return -1, err
		}
		err = usecase.itemCache.CreateItemsQuantityCache(ctxT, quantity, key)
		if err != nil {
			usecase.logger.Sugar().Warnf("error on create items quantity in search cache with key: %s, err: %v", key, err)
		}
		usecase.logger.Info("Get items quantity in search success")
		return quantity, nil
	}
	quantity, err := usecase.itemCache.GetItemsQuantityCache(ctxT, key)
	if err != nil {
		usecase.logger.Sugar().Warnf("error on get items quantity in search cache with key: %s, error: %v", key, err)
		// If get cache impossible get items quantity from database
		quantity, err := usecase.itemStore.ItemsInSearchQuantity(ctx, searchRequest)
		if err != nil {
			usecase.logger.Sugar().Errorf("error on get items quantity in search by search request: %s from database: %v", searchRequest, err)
			return -1, err
		}
		usecase.logger.Info("Get items quantity in search success")
		return quantity, nil
	}
	usecase.logger.Info("Get items quantity in search success")
	return quantity, nil
}

// ItemsList call database method and returns slice with all models.Item or error
func (usecase *itemUsecase) ItemsList(ctx context.Context, limitOptions map[string]int, sortOptions map[string]string) ([]models.Item, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase ItemsList() with args: ctx, limitOptions: %v, sortOptions: %v", limitOptions, sortOptions)
	// Context with timeout so as not to wait for an answer from the cache for too long
	ctxT, cancel := context.WithTimeout(ctx, timeout*time.Millisecond)
	defer cancel()
	limit, offset := limitOptions["limit"], limitOptions["offset"]
	sortType, sortOrder := sortOptions["sortType"], sortOptions["sortOrder"]
	// Check whether there is a cache with that name
	if ok := usecase.itemCache.CheckCache(ctxT, itemsListKey+sortType+sortOrder); !ok {
		// If the cache does not exist, request a list of items from the database
		itemIncomingChan, err := usecase.itemStore.ItemsList(ctx)
		if err != nil {
			return nil, err
		}
		items := make([]models.Item, 0, 100)
		for item := range itemIncomingChan {
			items = append(items, item)
		}
		// Sort the list of items based on the sorting parameters
		helpers.SortItems(items, sortType, sortOrder)
		// Create a cache with a sorted list of items
		err = usecase.itemCache.CreateItemsCache(ctxT, items, itemsListKey+sortType+sortOrder)
		if err != nil {
			usecase.logger.Sugar().Warnf("error on create items list cache with key: %s, error: %v", itemsListKey+sortType+sortOrder, err)
		} else {
			usecase.logger.Sugar().Infof("Create items list cache with key: %s success", itemsListKey+sortType+sortOrder)
		}
		// Create a cache with a quantity of items in list
		err = usecase.itemCache.CreateItemsQuantityCache(ctxT, len(items), itemsQuantityKey)
		if err != nil {
			usecase.logger.Sugar().Warnf("error on create items quantity cache with key: %s, error: %v", itemsQuantityKey, err)
		} else {
			usecase.logger.Sugar().Infof("Create items quantity cache with key: %s success", itemsQuantityKey)
		}
	}
	// Get items list from cache
	items, err := usecase.itemCache.GetItemsCache(ctxT, itemsListKey+sortType+sortOrder)
	if err != nil {
		usecase.logger.Sugar().Warnf("error on get cache with key: %s, err: %v", itemsListKey+sortType+sortOrder, err)
		// If error on get cache, request a list of items from the database
		itemIncomingChan, err := usecase.itemStore.ItemsList(ctx)
		if err != nil {
			usecase.logger.Sugar().Warnf("error on itemStore.ItemsList: %v", err)
			return nil, err
		}
		dbItems := make([]models.Item, 0, 100)
		for item := range itemIncomingChan {
			dbItems = append(dbItems, item)
		}
		// Sort the list of items based on the sorting parameters
		helpers.SortItems(dbItems, sortType, sortOrder)

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

// GetItemsByCategory call database method and returns chan with all models.Item in category or error
func (usecase *itemUsecase) GetItemsByCategory(ctx context.Context, categoryName string, limitOptions map[string]int, sortOptions map[string]string) ([]models.Item, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase GetItemsByCategory() with args: ctx, categoryName: %s, limitOptions: %v, sortOptions: %v", categoryName, limitOptions, sortOptions)

	// Context with timeout so as not to wait for an answer from the cache for too long
	ctxT, cancel := context.WithTimeout(ctx, timeout*time.Millisecond)
	defer cancel()

	limit, offset := limitOptions["limit"], limitOptions["offset"]
	sortType, sortOrder := sortOptions["sortType"], sortOptions["sortOrder"]

	// Check whether there is a cache of items in category
	if ok := usecase.itemCache.CheckCache(ctxT, categoryName+sortType+sortOrder); !ok {
		// If the cache does not exist, request a list of items in
		// category from the database
		itemIncomingChan, err := usecase.itemStore.GetItemsByCategory(ctx, categoryName)
		if err != nil {
			return nil, err
		}
		items := make([]models.Item, 0, 100)
		for item := range itemIncomingChan {
			items = append(items, item)
		}
		// Sort the list of items in category based on the sorting parameters
		helpers.SortItems(items, sortType, sortOrder)
		// Create a cache with a sorted list of items in category
		err = usecase.itemCache.CreateItemsCache(ctxT, items, categoryName+sortType+sortOrder)
		if err != nil {
			usecase.logger.Sugar().Warnf("error on create items cache with key: %s, error: %v", categoryName+sortType+sortOrder, err)
		} else {
			usecase.logger.Sugar().Infof("Create items cache with key: %s success", categoryName+sortType+sortOrder)
		}
		// Create a cache with a quantity of items in category
		err = usecase.itemCache.CreateItemsQuantityCache(ctxT, len(items), categoryName+"Quantity")
		if err != nil {
			usecase.logger.Sugar().Warnf("error on create items in category quantity cache with key: %s, error: %v", categoryName+"Quantity", err)
		} else {
			usecase.logger.Sugar().Infof("Create items in category quantity cache with key: %s success", categoryName+"Quantity")
		}
	}
	// Get items list from cache
	items, err := usecase.itemCache.GetItemsCache(ctxT, categoryName+sortType+sortOrder)
	if err != nil {
		usecase.logger.Sugar().Warnf("error on get cache with key: %s, error: %v", categoryName+sortType+sortOrder, err)
		// If error on get cache, request a list of items from the database
		itemIncomingChan, err := usecase.itemStore.GetItemsByCategory(ctx, categoryName)
		if err != nil {
			return nil, err
		}
		dbItems := make([]models.Item, 0, 100)
		for item := range itemIncomingChan {
			dbItems = append(dbItems, item)
		}
		// Sort the list of items based on the sorting parameters
		helpers.SortItems(dbItems, sortType, sortOrder)
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

// SearchLine call database method and returns chan with all models.Item with given params or error
func (usecase *itemUsecase) SearchLine(ctx context.Context, param string, limitOptions map[string]int, sortOptions map[string]string) ([]models.Item, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase SearchLine() with args: ctx, param: %s, limitOptions: %v, sortOptions: %v", param, limitOptions, sortOptions)

	// Context with timeout so as not to wait for an answer from the cache for too long
	ctxT, cancel := context.WithTimeout(ctx, timeout*time.Millisecond)
	defer cancel()

	limit, offset := limitOptions["limit"], limitOptions["offset"]
	sortType, sortOrder := sortOptions["sortType"], sortOptions["sortOrder"]

	// Check whether there is a cache of this search request
	if ok := usecase.itemCache.CheckCache(ctxT, param+sortType+sortOrder); !ok {
		// If the cache does not exist, request a list of items by
		// search request from the database
		itemIncomingChan, err := usecase.itemStore.SearchLine(ctx, param)
		if err != nil {
			return nil, err
		}
		items := make([]models.Item, 0, 100)
		for item := range itemIncomingChan {
			items = append(items, item)
		}
		// Create a cache with a quantity of items in list by search request
		err = usecase.itemCache.CreateItemsQuantityCache(ctxT, len(items), param+"Quantity")
		if err != nil {
			usecase.logger.Warn("can't create items quantity cache: %v", zap.Error(err))
		} else {
			usecase.logger.Info("Items quantity cache create success")
		}
		// Sort the list of items in search request based on the sorting parameters
		helpers.SortItems(items, sortType, sortOrder)
		// Create a cache with a sorted list of items in search request
		err = usecase.itemCache.CreateItemsCache(ctxT, items, param+sortType+sortOrder)
		if err != nil {
			usecase.logger.Sugar().Warnf("error on create cache of items in search with key: %s, error: %v", param+sortType+sortOrder, err)
		} else {
			usecase.logger.Sugar().Infof("Create cache of items in search with key: %s success", param+sortType+sortOrder)
		}
	}
	// Get items list from cache
	items, err := usecase.itemCache.GetItemsCache(ctxT, param+sortType+sortOrder)
	if err != nil {
		usecase.logger.Sugar().Warnf("error on get cache with key: %s, error: %v", param+sortType+sortOrder, err)
		// If error on get cache, request a list of items from the database
		itemIncomingChan, err := usecase.itemStore.SearchLine(ctx, param)
		if err != nil {
			return nil, err
		}
		dbItems := make([]models.Item, 0, 100)
		for item := range itemIncomingChan {
			dbItems = append(dbItems, item)
		}
		// Sort the list of items based on the sorting parameters
		helpers.SortItems(dbItems, sortType, sortOrder)
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
