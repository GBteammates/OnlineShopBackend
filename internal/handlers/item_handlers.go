package handlers

import (
	"OnlineShopBackend/internal/models"
	"context"
	"fmt"

	"github.com/google/uuid"
)

// Item is struct for DTO
type Item struct {
	Id          string `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Price       int32  `json:"price,omitempty"`
	Category    string `json:"category,omitempty"`
	Vendor      string `json:"vendor,omitempty"`
	Image       string `json:"image,omitempty"`
}

// CreateItem transform Item to models.Item and call usecase CreateItem
func (handlers *Handlers) CreateItem(ctx context.Context, item Item) (uuid.UUID, error) {
	handlers.logger.Debug("Enter in handlers CreateItem()")
	categoryId, err := uuid.Parse(item.Category)
	if err != nil {
		return uuid.UUID{}, err
	}
	newItem := &models.Item{
		Title:       item.Title,
		Description: item.Description,
		Price:       item.Price,
		Category:    categoryId,
		Vendor:      item.Vendor,
		Image:       item.Image,
	}
	id, err := handlers.repo.CreateItem(ctx, newItem)
	if err != nil {
		return id, err
	}
	return id, nil
}

// UpdateItem transform Item to models.Item and call usecase UpdateItem
func (handlers *Handlers) UpdateItem(ctx context.Context, item Item) error {
	handlers.logger.Debug("Enter in handlers UpdateItem()")
	id, err := uuid.Parse(item.Id)
	if err != nil {
		return fmt.Errorf("invalid item uuid: %w", err)
	}
	categoryId, err := uuid.Parse(item.Category)
	if err != nil {
		return fmt.Errorf("invalid category uuid: %w", err)
	}
	updateItem := &models.Item{
		Id:          id,
		Title:       item.Title,
		Description: item.Description,
		Price:       item.Price,
		Category:    categoryId,
		Vendor:      item.Vendor,
		Image:       item.Image,
	}
	return handlers.repo.UpdateItem(ctx, updateItem)
}

// GetItem returns Item on id
func (handlers *Handlers) GetItem(ctx context.Context, id string) (Item, error) {
	handlers.logger.Debug("Enter in handlers GetItem()")
	uid, err := uuid.Parse(id)
	if err != nil {
		return Item{}, fmt.Errorf("invalid item uuid: %w", err)
	}
	item, err := handlers.repo.GetItem(ctx, uid)
	if err != nil {
		return Item{}, err
	}
	return Item{
		Id:          item.Id.String(),
		Title:       item.Title,
		Description: item.Description,
		Price:       item.Price,
		Category:    item.Category.String(),
		Vendor:      item.Vendor,
		Image:       item.Image,
	}, nil
}

// ItemsList returns list of all Items
func (handlers *Handlers) ItemsList(ctx context.Context) ([]Item, error) {
	handlers.logger.Debug("Enter in handlers ItemsList()")
	res := make([]Item, 0, 100)
	items, err := handlers.repo.ItemsList(ctx)
	if err != nil {
		return res, err
	}
	for _, item := range items {
			res = append(res, Item{
				Id:          item.Id.String(),
				Title:       item.Title,
				Description: item.Description,
				Price:       item.Price,
				Category:    item.Category.String(),
				Vendor:      item.Vendor,
			})
		}
		return res, nil
	}

// SearchLine returns list of Items with parameters
func (handlers *Handlers) SearchLine(ctx context.Context, param string) ([]Item, error) {
	handlers.logger.Debug("Enter in handlers SearchLine()")
	res := make([]Item, 0, 100)
	items, err := handlers.repo.SearchLine(ctx, param)
	if err != nil {
		return res, err
	}
	for {
		select {
		case <-ctx.Done():
			handlers.logger.Debug("handlers SearchLine() ctx.Done recieved")
			return res, ctx.Err()
		case item, ok := <-items:
			if !ok {
				return res, nil
			}
			res = append(res, Item{
				Id:          item.Id.String(),
				Title:       item.Title,
				Description: item.Description,
				Price:       item.Price,
				Category:    item.Category.String(),
				Vendor:      item.Vendor,
			})
		}
	}
}
