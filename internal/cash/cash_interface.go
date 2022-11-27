package cash

import (
	"OnlineShopBackend/internal/handlers"
	"OnlineShopBackend/internal/models"
	"context"
)

type Cash interface {
	CheckCash(key string) bool
	CreateCash(ctx context.Context, res chan models.Item, key string) error
	GetCash(key string) ([]handlers.Item, error)
}
