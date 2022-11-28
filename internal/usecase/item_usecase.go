package usecase

import (
	"OnlineShopBackend/internal/models"
	"context"
	"fmt"

	"github.com/google/uuid"
)

const cashKey = "ItemsList"

// CreateItem call database method and returns id of created item or error
func (usecase *Usecase) CreateItem(ctx context.Context, item *models.Item) (uuid.UUID, error) {
	usecase.logger.Debug("Enter in usecase CreateItem()")
	id, err := usecase.itemStore.CreateItem(ctx, item)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error on create item: %w", err)
	}
	err = usecase.updateCash(ctx, id, "create")
	if err != nil {
		usecase.logger.Error(fmt.Sprintf("error on update cash: %v", err))
	} else {
		usecase.logger.Info("Update cash success")
	}
	return id, nil
}

// UpdateItem call database method to update item and returns error or nil
func (usecase *Usecase) UpdateItem(ctx context.Context, item *models.Item) error {
	usecase.logger.Debug("Enter in usecase UpdateItem()")
	err := usecase.itemStore.UpdateItem(ctx, item)
	if err != nil {
		return fmt.Errorf("error on update item: %w", err)
	}
	err = usecase.updateCash(ctx, item.Id, "update")
	if err != nil {
		usecase.logger.Error(fmt.Sprintf("error on update cash: %v", err))
	} else {
		usecase.logger.Info("Update cash success")
	}
	return nil
}

// GetItem call database and returns *models.Item with given id or returns error
func (usecase *Usecase) GetItem(ctx context.Context, id uuid.UUID) (*models.Item, error) {
	usecase.logger.Debug("Enter in usecase GetItem()")
	item, err := usecase.itemStore.GetItem(ctx, id)
	if err != nil {
		return &models.Item{}, fmt.Errorf("error on get item: %w", err)
	}
	return item, nil
}

// ItemsList call database method and returns chan with all models.Item or error
func (usecase *Usecase) ItemsList(ctx context.Context) ([]models.Item, error) {
	usecase.logger.Debug("Enter in usecase ItemsList()")
	if !usecase.itemCash.CheckCash(cashKey) {
		itemIncomingChan, err := usecase.itemStore.ItemsList(ctx)
		if err != nil {
			return nil, err
		}
		items := make([]models.Item, 0, 100)
		for item := range itemIncomingChan {
			items = append(items, item)
		}

		err = usecase.itemCash.CreateCash(ctx, items, cashKey)
		if err != nil {
			return nil, fmt.Errorf("error on create cash: %w", err)
		}
	}
	return usecase.itemCash.GetCash(cashKey)
}

// SearchLine call database method and returns chan with all models.Item with given params or error
func (usecase *Usecase) SearchLine(ctx context.Context, param string) (chan models.Item, error) {
	usecase.logger.Debug("Enter in usecase SearchLine()")
	itemIncomingChan, err := usecase.itemStore.SearchLine(ctx, param)
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

// updateCash updating cash when creating or updating item
func (usecase *Usecase) updateCash(ctx context.Context, id uuid.UUID, op string) error {
	usecase.logger.Debug("Enter in usecase UpdateCash()")
	if !usecase.itemCash.CheckCash(cashKey) {
		return fmt.Errorf("cash is not exists")
	}
	newItem, err := usecase.itemStore.GetItem(ctx, id)
	if err != nil {
		return fmt.Errorf("error on get item: %w", err)
	}
	items, err := usecase.itemCash.GetCash(cashKey)
	if err != nil {
		return fmt.Errorf("error on get cash: %w", err)
	}
	if op == "update" {
		for i, item := range items {
			if item.Id == id {
				items[i] = *newItem
				break
			}
		}
	}
	if op == "create" {
		items = append(items, *newItem)
	}

	return usecase.itemCash.CreateCash(ctx, items, cashKey)
}
