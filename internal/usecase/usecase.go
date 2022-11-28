package usecase

import (
	"OnlineShopBackend/internal/cash"
	"OnlineShopBackend/internal/repository"

	"go.uber.org/zap"
)

type Storage struct {
	itemStore     repository.ItemStore
	categoryStore repository.CategoryStore
	itemCash      cash.Cash
	logger        *zap.Logger
}

func NewStorage(itemStore repository.ItemStore, categoryStore repository.CategoryStore, itemCash cash.Cash, logger *zap.Logger) *Storage {
	return &Storage{itemStore: itemStore, categoryStore: categoryStore, itemCash: itemCash, logger: logger}
}
