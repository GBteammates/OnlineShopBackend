package cash

import (
	"OnlineShopBackend/internal/models"
	"context"
)

type Cash interface {
	CheckCash(key string) bool
	CreateCash(ctx context.Context, res []models.Item, key string) error
	GetCash(key string) ([]models.Item, error)
}
