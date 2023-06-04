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
func (u *itemUsecase) AddFavouriteItem(ctx context.Context, userId uuid.UUID, itemId uuid.UUID) error {
	u.logger.Debugf("Enter in usecase AddFavouriteItem() with args: ctx, userId: %v, itemId: %v", userId, itemId)
	err := u.store.AddFavouriteItem(ctx, userId, itemId)
	if err != nil {
		return err
	}
	item, err := u.Get(ctx, itemId)
	if err != nil {
		u.logger.Errorf("error on get item: %v", err)
		return err
	}
	err = u.cache.UpdateCache(ctx, &models.CacheOptions{
		Op:      models.CreateOp,
		Kind:    []string{models.Favourites},
		NewItem: item,
		UserId:  userId,
	})
	if err != nil {
		u.logger.Warn("error on update cache")
	}
	err = u.cache.UpdateFavIdsCache(ctx, userId, item, addOp)
	if err != nil {
		u.logger.Warn("error on update cache")
	}
	return nil
}

// DeleteFavouriteItem deleted item from list of favourites items
func (u *itemUsecase) DeleteFavouriteItem(ctx context.Context, userId uuid.UUID, itemId uuid.UUID) error {
	u.logger.Debugf("Enter in usecase DeleteFavouriteItem() with args: ctx, userId: %v, itemId: %v", userId, itemId)
	err := u.store.DeleteFavouriteItem(ctx, userId, itemId)
	if err != nil {
		return err
	}
	u.cache.UpdateCache(ctx, &models.CacheOptions{
		Op:      models.DeleteOp,
		Kind:    []string{models.Favourites},
		NewItem: &models.Item{Id: itemId},
		UserId:  userId,
	})
	err = u.cache.UpdateFavIdsCache(ctx, userId, &models.Item{Id: itemId}, models.DeleteOp)
	if err != nil {
		u.logger.Warn("error on update cache")
	}
	return nil
}

// GetFavouriteItemsId calls database method and returns map with identificators of favourite items of user or error
func (u *itemUsecase) ListFavouriteItemsId(ctx context.Context, userId uuid.UUID) (*map[uuid.UUID]uuid.UUID, error) {
	u.logger.Debugf("Enter in usecase GetFavouriteItemsId() with args: ctx, userId: %v", userId)

	favIds, err := u.getFavouriteItemsIds(ctx, models.Favourites, userId, u.store.FavouriteItemsId)
	if err != nil {
		return nil, err
	}
	return favIds, nil
}

func (u *itemUsecase) getFavouriteItemsIds(
	ctx context.Context,
	kind string,
	userId uuid.UUID,
	f func(ctx context.Context, userId uuid.UUID) (*map[uuid.UUID]uuid.UUID, error),
) (*map[uuid.UUID]uuid.UUID, error) {
	u.logger.Debugf("Enter in usecase getFavouribeItemsIds() with args: ctx, userId: %v, kind: %s, f()", userId, kind)

	// Context with timeout so as not to wait for an answer from the cache for too long
	ctxT, cancel := context.WithTimeout(ctx, u.timeout*time.Millisecond)
	defer cancel()
	favIds, err := u.cache.FavouriteItemsIdsFromCache(ctxT, userId.String()+fav, kind)
	if err == nil {
		return favIds, nil
	}
	if err != nil {
		u.logger.Debugf("error on get favourite items ids from cache: %v", err)
	}
	favIds, err = f(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("error on get favourite items ids from store: %w", err)
	}
	err = u.cache.FavouriteItemsIdsToCache(ctxT, favIds, userId.String()+fav, kind)
	if err != nil {
		u.logger.Errorf("error on recording favouribe items ids to cache: %v", err)
	}
	return favIds, nil
}
