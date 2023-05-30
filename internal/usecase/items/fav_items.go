package usecase

import (
	"OnlineShopBackend/internal/models"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	favourite = "favourite"
	fav       = "Fav"
	addOp     = "add"
)

// AddFavouriteItem added item in list of favourites items
func (usecase *itemUsecase) AddFavouriteItem(ctx context.Context, userId uuid.UUID, itemId uuid.UUID) error {
	usecase.logger.Sugar().Debugf("Enter in usecase AddFavouriteItem() with args: ctx, userId: %v, itemId: %v", userId, itemId)
	err := usecase.itemStore.AddFavouriteItem(ctx, userId, itemId)
	if err != nil {
		return err
	}
	item, err := usecase.GetItem(ctx, itemId)
	if err != nil {
		usecase.logger.Sugar().Errorf("error on get item: %v", err)
		return err
	}
	err = usecase.itemCache.UpdateCache(ctx, &models.ItemsCacheOptions{
		Op:      createOp,
		Kind:    []string{favourite},
		NewItem: item,
		UserId:  userId,
	})
	if err != nil {
		usecase.logger.Warn("error on update cache")
	}
	err = usecase.itemCache.UpdateFavIdsCache(ctx, userId, item, addOp)
	if err != nil {
		usecase.logger.Warn("error on update cache")
	}
	return nil
}

// DeleteFavouriteItem deleted item from list of favourites items
func (usecase *itemUsecase) DeleteFavouriteItem(ctx context.Context, userId uuid.UUID, itemId uuid.UUID) error {
	usecase.logger.Sugar().Debugf("Enter in usecase DeleteFavouriteItem() with args: ctx, userId: %v, itemId: %v", userId, itemId)
	err := usecase.itemStore.DeleteFavouriteItem(ctx, userId, itemId)
	if err != nil {
		return err
	}
	usecase.itemCache.UpdateCache(ctx, &models.ItemsCacheOptions{
		Op:      deleteOp,
		Kind:    []string{favourite},
		NewItem: &models.Item{Id: itemId},
		UserId:  userId,
	})
	err = usecase.itemCache.UpdateFavIdsCache(ctx, userId, &models.Item{Id: itemId}, deleteOp)
	if err != nil {
		usecase.logger.Warn("error on update cache")
	}
	return nil
}

// ItemsQuantityInFavourite check cache and if cache not exists call database
// method and write in cache and returns quantity of items in favourite
func (usecase *itemUsecase) ItemsQuantityInFavourite(ctx context.Context, userId uuid.UUID) (int, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase GetFavouriteQuantity() with args: ctx, userId: %v", userId)

	quantity, err := usecase.getItemsQuantity(ctx, userId.String(), favourite, usecase.itemStore.ItemsInFavouriteQuantity)
	if err != nil {
		return -1, err
	}
	return quantity, nil
}

// GetFavouriteItems call database method and returns chan with models.Item from list of favourites item or error
func (usecase *itemUsecase) GetFavouriteItems(ctx context.Context, userId uuid.UUID, opts *models.ItemsListOptions) ([]models.Item, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase GetFavouriteItems() with args: ctx, userId: %v, opts: %v", userId, opts)

	items, err := usecase.getItemsList(ctx, userId.String(), favourite, opts, usecase.itemStore.GetFavouriteItems)
	if err != nil {
		return nil, err
	}
	return items, nil
}

// GetFavouriteItemsId calls database method and returns map with identificators of favourite items of user or error
func (usecase *itemUsecase) GetFavouriteItemsId(ctx context.Context, userId uuid.UUID) (*map[uuid.UUID]uuid.UUID, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase GetFavouriteItemsId() with args: ctx, userId: %v", userId)

	favIds, err := usecase.getFavouriteItemsIds(ctx, favourite, userId, usecase.itemStore.GetFavouriteItemsId)
	if err != nil {
		return nil, err
	}
	return favIds, nil
}

func (usecase *itemUsecase) getFavouriteItemsIds(
	ctx context.Context,
	kind string,
	userId uuid.UUID,
	f func(ctx context.Context, userId uuid.UUID) (*map[uuid.UUID]uuid.UUID, error),
) (*map[uuid.UUID]uuid.UUID, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase getFavouribeItemsIds() with args: ctx, userId: %v, kind: %s, f()", userId, kind)

	// Context with timeout so as not to wait for an answer from the cache for too long
	ctxT, cancel := context.WithTimeout(ctx, timeout*time.Millisecond)
	defer cancel()
	favIds, err := usecase.itemCache.FavouriteItemsIdsFromCache(ctxT, userId.String()+fav, kind)
	if err == nil {
		return favIds, nil
	}
	if err != nil {
		usecase.logger.Sugar().Debugf("error on get favourite items ids from cache: %v", err)
	}
	favIds, err = f(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("error on get favourite items ids from store: %w", err)
	}
	err = usecase.itemCache.FavouriteItemsIdsToCache(ctxT, favIds, userId.String()+fav, kind)
	if err != nil {
		usecase.logger.Sugar().Errorf("error on recording favouribe items ids to cache: %v", err)
	}
	return favIds, nil
}
