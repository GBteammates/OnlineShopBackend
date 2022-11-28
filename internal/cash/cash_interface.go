package cash

import (
	"OnlineShopBackend/internal/models"
	"context"
)

type Cash interface {
	CheckCash(key string) bool
	CreateCash(ctx context.Context, res chan models.Item, key string) error
	GetCash(key string) ([]models.Item, error)
}
