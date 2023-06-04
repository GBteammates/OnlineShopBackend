package usecase

import (
	"OnlineShopBackend/internal/models"
	"context"

	"github.com/google/uuid"
)

type IItemsCache interface {
	ItemsToCache(ctx context.Context, items []models.Item, kind, param string) error
	ItemsFromCache(ctx context.Context, cacheKey, kind string) ([]models.Item, error)
	ItemsQuantityToCache(ctx context.Context, value int, key, kind string) error
	ItemsQuantityFromCache(ctx context.Context, key string, kind string) (int, error)
	FavouriteItemsIdsToCache(ctx context.Context, favIds *map[uuid.UUID]uuid.UUID, key, kind string) error
	FavouriteItemsIdsFromCache(ctx context.Context, key, kind string) (*map[uuid.UUID]uuid.UUID, error)
	UpdateCache(ctx context.Context, opts *models.CacheOptions) error
	UpdateFavIdsCache(ctx context.Context, userId uuid.UUID, item *models.Item, op string) error
}

type ICategoriesCache interface {
	CategoriesToCache(ctx context.Context, categories []models.Category) error
	CategoriesFromCache(ctx context.Context, key string) ([]models.Category, error)
	UpdateCache(ctx context.Context, newCategory *models.Category, op string) error
	DeleteCache(ctx context.Context, name string) error
}
