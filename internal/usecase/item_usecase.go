package usecase

import (
	"OnlineShopBackend/internal/models"
	"context"
	"fmt"

	"github.com/google/uuid"
)

// CreateItem call database method and returns id of created item or error
func (storage *Storage) CreateItem(ctx context.Context, item *models.Item) (uuid.UUID, error) {
	storage.logger.Debug("Enter in usecase CreateItem()")
	id, err := storage.itemStore.CreateItem(ctx, item)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error on create item: %w", err)
	}
	err = storage.updateCash(ctx, id, "create")
	if err != nil {
		storage.logger.Error(fmt.Sprintf("error on update cash: %v", err))
	} else {
		storage.logger.Info("Update cash success")
	}
	return id, nil
}

// UpdateItem call database method to update item and returns error or nil
func (storage *Storage) UpdateItem(ctx context.Context, item *models.Item) error {
	storage.logger.Debug("Enter in usecase UpdateItem()")
	err := storage.itemStore.UpdateItem(ctx, item)
	if err != nil {
		return fmt.Errorf("error on update item: %w", err)
	}
	err = storage.updateCash(ctx, item.Id, "update")
	if err != nil {
		storage.logger.Error(fmt.Sprintf("error on update cash: %v", err))
	} else {
		storage.logger.Info("Update cash success")
	}
	return nil
}

// GetItem call database and returns *models.Item with given id or returns error
func (storage *Storage) GetItem(ctx context.Context, id uuid.UUID) (*models.Item, error) {
	storage.logger.Debug("Enter in usecase GetItem()")
	item, err := storage.itemStore.GetItem(ctx, id)
	if err != nil {
		return &models.Item{}, fmt.Errorf("error on get item: %w", err)
	}
	return item, nil
}

// ItemsList call database method and returns chan with all models.Item or error
func (storage *Storage) ItemsList(ctx context.Context) ([]models.Item, error) {
	storage.logger.Debug("Enter in usecase ItemsList()")
	cashKey := "ItemsList"
	if !storage.itemCash.CheckCash(cashKey) {
		itemIncomingChan, err := storage.itemStore.ItemsList(ctx)
		if err != nil {
			return nil, err
		}
		err = storage.itemCash.CreateCash(ctx, itemIncomingChan, cashKey)
		if err != nil {
			return nil, fmt.Errorf("error on create cash: %w", err)
		}
	}
	return storage.itemCash.GetCash(cashKey)
}

// SearchLine call database method and returns chan with all models.Item with given params or error
func (storage *Storage) SearchLine(ctx context.Context, param string) (chan models.Item, error) {
	storage.logger.Debug("Enter in usecase SearchLine()")
	itemIncomingChan, err := storage.itemStore.SearchLine(ctx, param)
	if err != nil {
		return nil, err
	}
	itemOutChan := make(chan models.Item, 100)
	go func() {
		defer close(itemOutChan)
		for {
			select {
			case <-ctx.Done():
				return
			case item, ok := <-itemIncomingChan:
				if !ok {
					return
				}
				itemOutChan <- item
			}
		}
	}()
	return itemOutChan, nil
}

func (storage *Storage) updateCash(ctx context.Context, id uuid.UUID, op string) error {
	storage.logger.Debug("Enter in usecase UpdateCash()")
	key := "ItemsList"
	if !storage.itemCash.CheckCash(key) {
		return fmt.Errorf("cash is not exists")
	}
	item, err := storage.itemStore.GetItem(ctx, id)
	if err != nil {
		return fmt.Errorf("error on get item: %w", err)
	}
	items, err := storage.itemCash.GetCash(key)
	if err != nil {
		return fmt.Errorf("error on get cash: %w", err)
	}
	if op == "update" {
		for i, item := range items {
			if item.Id == id {
				items[i] = item
				break
			}
		}
	}
	if op == "create" {
		items = append(items, *item)
	}

	itemsChan := make(chan models.Item, len(items))
	go func() {
		defer close(itemsChan)
		for _, item := range items {
			itemsChan <- item
		}
	}()
	return storage.itemCash.CreateCash(ctx, itemsChan, key)
}
