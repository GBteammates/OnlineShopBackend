package usecase

import (
	"OnlineShopBackend/internal/models"
	"context"

	"github.com/google/uuid"
)

type IItemsCache interface {
	CheckCache(ctx context.Context, key string) bool
	CreateItemsCache(ctx context.Context, res []models.Item, key string) error
	CreateItemsQuantityCache(ctx context.Context, value int, key string) error
	GetItemsCache(ctx context.Context, key string) ([]models.Item, error)
	GetItemsQuantityCache(ctx context.Context, key string) (int, error)
	CreateFavouriteItemsIdCache(ctx context.Context, res map[uuid.UUID]uuid.UUID, key string) error
	GetFavouriteItemsIdCache(ctx context.Context, key string) (*map[uuid.UUID]uuid.UUID, error)
	UpdateCache(ctx context.Context, newItem *models.Item, op string) error
	UpdateFavouritesCache(ctx context.Context, userId uuid.UUID, item *models.Item, op string) error
}

type ICategoriesCache interface {
	CheckCache(ctx context.Context, key string) bool
	CreateCategoriesListСache(ctx context.Context, categories []models.Category, key string) error
	GetCategoriesListCache(ctx context.Context, key string) ([]models.Category, error)
	UpdateCategoryCache(ctx context.Context, newCategory *models.Category, op string) error
	DeleteCategoryCache(ctx context.Context, name string) error
}
