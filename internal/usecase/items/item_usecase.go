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

type itemUsecase struct {
	timeout     time.Duration
	store       usecase.ItemStore
	cache       usecase.IItemsCache
	filestorage usecase.Filestorage
	logger      *zap.SugaredLogger
}

func NewItemUsecase(timeout time.Duration, store usecase.ItemStore, cache usecase.IItemsCache, filestorage usecase.Filestorage, logger *zap.SugaredLogger) *itemUsecase {
	logger.Debug("Enter in usecase NewItemUsecase()")
	return &itemUsecase{
		timeout:     timeout,
		store:       store,
		cache:       cache,
		filestorage: filestorage,
		logger:      logger}
}

// CreateItem call database method and returns id of created item or error
func (u *itemUsecase) Create(ctx context.Context, item *models.Item) (uuid.UUID, error) {
	u.logger.Debugf("Enter in usecase Create() with args: ctx, item: %v", item)
	id, err := u.store.Create(ctx, item)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error on create item: %w", err)
	}

	item.Id = id

	err = u.cache.UpdateCache(ctx, &models.CacheOptions{
		Op:      models.CreateOp,
		Kind:    []string{models.List, models.InCategory},
		NewItem: item,
	})
	if err != nil {
		u.logger.Errorf("error on update cache: %v", err)
	}

	return id, nil
}

// UpdateItem call database method to update item and returns error or nil
func (u *itemUsecase) Update(ctx context.Context, item *models.Item) error {
	u.logger.Debugf("Enter in usecase Update() with args: ctx, item: %v", item)
	err := u.store.Update(ctx, item)
	if err != nil {
		return fmt.Errorf("error on update item: %w", err)
	}
	err = u.cache.UpdateCache(ctx, &models.CacheOptions{
		Op:      models.UpdateOp,
		Kind:    []string{models.List, models.InCategory},
		NewItem: item,
	})
	if err != nil {
		u.logger.Errorf("error on update cache: %v", err)
	}

	return nil
}

// GetItem call database and returns *models.Item with given id or returns error
func (u *itemUsecase) Get(ctx context.Context, id uuid.UUID) (*models.Item, error) {
	u.logger.Debugf("Enter in usecase Get() with args: ctx, id: %v", id)
	item, err := u.store.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error on get item: %w", err)
	}
	return item, nil
}

// DeleteItem call database method for deleting item
func (u *itemUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	u.logger.Debugf("Enter in usecase DeleteItem() with args: ctx, id: %v", id)
	err := u.store.Delete(ctx, id)
	if err != nil {
		return err
	}
	err = u.cache.UpdateCache(ctx, &models.CacheOptions{
		Op:      models.DeleteOp,
		Kind:    []string{models.List, models.InCategory},
		NewItem: &models.Item{Id: id},
	})
	if err != nil {
		u.logger.Error(fmt.Sprintf("error on update cache: %v", err))
	}
	return nil
}

// ItemsQuantity check cache and if cache not exists call database
// method and write in cache and returns quantity of all items
func (u *itemUsecase) Quantity(ctx context.Context, opts *models.QuantityOptions) (int, error) {
	u.logger.Debug("Enter in usecase ItemsQuantity() with args: ctx")

	switch opts.Kind {
	case models.List:
		opts.Handler = u.store.ListQuantity
	case models.InCategory:
		opts.Handler = u.store.ListByCategoryQuantity
	case models.Search:
		opts.Handler = u.store.InSearchQuantity
	case models.Favourites:
		opts.Handler = u.store.InFavouriteQuantity
	}

	quantity, err := u.getQuantity(ctx, opts)
	if err != nil {
		return -1, err
	}
	return quantity, nil
}

func (u *itemUsecase) getQuantity(ctx context.Context, opts *models.QuantityOptions) (int, error) {
	u.logger.Debugf("Enter in usecase getQuantity() with args: ctx, opts: %v", opts)

	// Context with timeout so as not to wait for an answer from the cache for too long
	ctxT, cancel := context.WithTimeout(ctx, u.timeout*time.Millisecond)
	defer cancel()

	quantity, err := u.cache.ItemsQuantityFromCache(ctxT, opts.Param+models.Quantity, opts.Kind)
	if err == nil {
		return quantity, nil
	}
	if err != nil {
		u.logger.Debugf("error on get items quantity from cache: %v", err)
	}
	quantity, err = opts.Handler(ctx, opts.Param)
	if err != nil {
		return -1, fmt.Errorf("error on get items quantity from store: %w", err)
	}
	err = u.cache.ItemsQuantityToCache(ctx, quantity, opts.Param+models.Quantity, opts.Kind)
	if err != nil {
		u.logger.Debugf("error on recording items quanitity to cache: %v", err)
	}
	return quantity, nil
}

// ItemsList call database method and returns slice with all models.Item or error
func (u *itemUsecase) List(ctx context.Context, opts *models.ListOptions) ([]models.Item, error) {
	u.logger.Debugf("Enter in usecase ItemsList() with args: ctx, opts: %v", opts)

	switch opts.Kind {
	case models.List:
		opts.Handler = u.store.List
	case models.InCategory:
		opts.Handler = u.store.ListByCategory
	case models.Search:
		opts.Handler = u.store.SearchLine
	case models.Favourites:
		opts.Handler = u.store.ListFavouriteItems
	}

	items, err := u.getList(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("error on get list of items: %w", err)
	}
	return items, nil
}

func (u *itemUsecase) getList(ctx context.Context, opts *models.ListOptions) ([]models.Item, error) {
	u.logger.Debugf("Enter in usecae getItemsList() with args: ctx, opts: %v", opts)

	// Context with timeout so as not to wait for an answer from the cache for too long
	ctxT, cancel := context.WithTimeout(ctx, u.timeout*time.Millisecond)
	defer cancel()

	// Get items list from cache
	items, err := u.cache.ItemsFromCache(ctxT, opts.Param+opts.SortType+opts.SortOrder, opts.Kind)
	if err == nil {
		return items, nil
	}
	if err != nil {
		u.logger.Debugf("error on get items from cache: %v", err)
	}

	itemsChan, err := opts.Handler(ctx, opts.Param)
	if err != nil {
		return nil, fmt.Errorf("error on get items from store: %w", err)
	}
	items = make([]models.Item, 0, 100)
	for item := range itemsChan {
		items = append(items, item)
	}

	err = u.cache.ItemsToCache(ctxT, items, opts.Kind, opts.Param)
	if err != nil {
		u.logger.Warnf("error on create cache: %v", err)
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
