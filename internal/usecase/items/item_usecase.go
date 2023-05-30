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

const (
	timeout               = 100
	itemsListKey          = "ItemsList"
	itemsListKeyNameAsc   = "ItemsListnameasc"
	itemsListKeyNameDesc  = "ItemsListnamedesc"
	itemsListKeyPriceAsc  = "ItemsListpriceasc"
	itemsListKeyPriceDesc = "ItemsListpricedesc"
	itemsQuantityKey      = "ItemsQuantity"
	createOp              = "create"
	updateOp              = "update"
	deleteOp              = "delete"
	quantityKey           = "Quantity"
	list                  = "list"
	inCategory            = "inCategory"
	search                = "search"
	limit                 = "limit"
	offset                = "offset"
	sortType              = "sortType"
	sortOrder             = "sortOrder"
)

type itemUsecase struct {
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

	err = usecase.itemCache.UpdateCache(ctx, &models.ItemsCacheOptions{
		Op:      createOp,
		Kind:    []string{list, inCategory},
		NewItem: item,
	})
	if err != nil {
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
	err = usecase.itemCache.UpdateCache(ctx, &models.ItemsCacheOptions{
		Op:      updateOp,
		Kind:    []string{list, inCategory},
		NewItem: item,
	})
	if err != nil {
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
	err = usecase.itemCache.UpdateCache(ctx, &models.ItemsCacheOptions{
		Op:      deleteOp,
		Kind:    []string{list, inCategory},
		NewItem: &models.Item{Id: id},
	})
	if err != nil {
		usecase.logger.Error(fmt.Sprintf("error on update cache: %v", err))
	}
	return nil
}

// ItemsQuantity check cache and if cache not exists call database
// method and write in cache and returns quantity of all items
func (usecase *itemUsecase) ItemsQuantity(ctx context.Context) (int, error) {
	usecase.logger.Debug("Enter in usecase ItemsQuantity() with args: ctx")

	quantity, err := usecase.getItemsQuantity(ctx, itemsQuantityKey, list, usecase.itemStore.ItemsListQuantity)
	if err != nil {
		return -1, err
	}
	return quantity, nil
}

// ItemsQuantityInCategory check cache and if cache not exists call database
// method and write in cache and returns quantity of items in category
func (usecase *itemUsecase) ItemsQuantityInCategory(ctx context.Context, categoryName string) (int, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase ItemsQuantityInCategory() with args: ctx, categoryName: %s", categoryName)

	quantity, err := usecase.getItemsQuantity(ctx, categoryName, inCategory, usecase.itemStore.ItemsByCategoryQuantity)
	if err != nil {
		return -1, err
	}
	return quantity, nil
}

// ItemsQuantityInSearch check cache and if cache not exists call database method and write
// in cache and returns quantity of items in search request
func (usecase *itemUsecase) ItemsQuantityInSearch(ctx context.Context, searchRequest string) (int, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase ItemsQuantityInSearch() with args: ctx, searchRequest: %s", searchRequest)

	quantity, err := usecase.getItemsQuantity(ctx, searchRequest, search, usecase.itemStore.ItemsInSearchQuantity)
	if err != nil {
		return -1, err
	}
	return quantity, nil
}

func (usecase *itemUsecase) getItemsQuantity(
	ctx context.Context,
	param, kind string,
	f func(ctx context.Context, param string) (int, error),
) (int, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase getItemsQuantity() with args: ctx, param: %s, kind: %s", param, kind)

	// Context with timeout so as not to wait for an answer from the cache for too long
	ctxT, cancel := context.WithTimeout(ctx, timeout*time.Millisecond)
	defer cancel()

	quantity, err := usecase.itemCache.ItemsQuantityFromCache(ctxT, param+quantityKey, kind)
	if err == nil {
		return quantity, nil
	}
	if err != nil {
		usecase.logger.Sugar().Debugf("error on get items quantity from cache: %v", err)
	}
	quantity, err = f(ctx, param)
	if err != nil {
		return -1, fmt.Errorf("error on get items quantity from store: %w", err)
	}
	err = usecase.itemCache.ItemsQuantityToCache(ctx, quantity, param+quantityKey, kind)
	if err != nil {
		usecase.logger.Sugar().Debugf("error on recording items quanitity to cache: %v", err)
	}
	return quantity, nil
}

// ItemsList call database method and returns slice with all models.Item or error
func (usecase *itemUsecase) ItemsList(ctx context.Context, opts *models.ItemsListOptions) ([]models.Item, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase ItemsList() with args: ctx, opts: %v", opts)

	items, err := usecase.getItemsList(ctx, itemsListKey, list, opts, usecase.itemStore.ItemsList)
	if err != nil {
		return nil, fmt.Errorf("error on get list of items: %w", err)
	}
	return items, nil
}

// GetItemsByCategory call database method and returns chan with all models.Item in category or error
func (usecase *itemUsecase) GetItemsByCategory(ctx context.Context, categoryName string, opts *models.ItemsListOptions) ([]models.Item, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase GetItemsByCategory() with args: ctx, opts: %v", categoryName, opts)

	items, err := usecase.getItemsList(ctx, categoryName, inCategory, opts, usecase.itemStore.GetItemsByCategory)
	if err != nil {
		return nil, fmt.Errorf("error on get items by category: %w", err)
	}
	return items, nil
}

// SearchLine call database method and returns chan with all models.Item with given params or error
func (usecase *itemUsecase) SearchLine(ctx context.Context, param string, opts *models.ItemsListOptions) ([]models.Item, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase SearchLine() with args: ctx, param: %s, opts: %v", param, opts)

	items, err := usecase.getItemsList(ctx, param, search, opts, usecase.itemStore.SearchLine)
	if err != nil {
		return nil, fmt.Errorf("error on search items: %w", err)
	}
	return items, nil
}

func (usecase *itemUsecase) getItemsList(
	ctx context.Context,
	param string,
	kind string,
	opts *models.ItemsListOptions,
	f func(ctx context.Context, param string) (chan models.Item, error),
) ([]models.Item, error) {
	usecase.logger.Sugar().Debugf("Enter in usecae getItemsList() with args: ctx, param: %s, kind: %s, opts: %v, f: %v", param, kind, opts, f)

	// Context with timeout so as not to wait for an answer from the cache for too long
	ctxT, cancel := context.WithTimeout(ctx, timeout*time.Millisecond)
	defer cancel()

	// Get items list from cache
	items, err := usecase.itemCache.ItemsFromCache(ctxT, param+opts.SortType+opts.SortOrder, inCategory)
	if err == nil {
		return items, nil
	}
	if err != nil {
		usecase.logger.Sugar().Debugf("error on get items from cache: %v", err)
	}

	itemsChan, err := f(ctx, param)
	if err != nil {
		return nil, fmt.Errorf("error on get items from store: %w", err)
	}
	items = make([]models.Item, 0, 100)
	for item := range itemsChan {
		items = append(items, item)
	}

	err = usecase.itemCache.ItemsToCache(ctxT, items, kind, param)
	if err != nil {
		usecase.logger.Sugar().Warnf("error on create cache: %v", err)
	}

	helpers.SortItems(items, opts.SortType, opts.SortOrder)

	if opts.Offset > len(items) {
		return nil, fmt.Errorf("error: offset bigger than lenght of items, offset: %d, lenght of items: %d", opts.Offset, len(items))
	}
	itemsWithLimit := make([]models.Item, 0, opts.Limit)
	var counter = 0
	for i := opts.Offset; i < len(items); i++ {
		if counter == opts.Limit {
			break
		}
		// Add items to the resulting list of items until the counter is equal to the limit
		itemsWithLimit = append(itemsWithLimit, items[i])
		counter++
	}
	return itemsWithLimit, nil
}
