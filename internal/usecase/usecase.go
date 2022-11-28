package usecase

import (
	"OnlineShopBackend/internal/cash"
	"OnlineShopBackend/internal/repository"

	"go.uber.org/zap"
)

type Usecase struct {
	itemStore     repository.ItemStore
	categoryStore repository.CategoryStore
	itemCash      cash.Cash
	logger        *zap.Logger
}

func NewUsecase(itemStore repository.ItemStore, categoryStore repository.CategoryStore, itemCash cash.Cash, logger *zap.Logger) *Usecase {
	return &Usecase{itemStore: itemStore, categoryStore: categoryStore, itemCash: itemCash, logger: logger}
}
