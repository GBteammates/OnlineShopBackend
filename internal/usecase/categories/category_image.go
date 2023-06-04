package usecase

import (
	"OnlineShopBackend/internal/models"
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func (u *categoryUsecase) UploadImage(ctx context.Context, id uuid.UUID, name string, file []byte) error {
	u.logger.Debugf("Enter in usecase UploadCategoryImage() with args: ctx, id: %v, name: %s, file", id, name)

	category, err := u.store.Get(ctx, id.String())
	if err != nil {
		return fmt.Errorf("error on get category: %w", err)
	}

	path, err := u.filestorage.PutCategoryImage(id.String(), name, file)
	if err != nil {
		return fmt.Errorf("error on put image to filestorage: %w", err)
	}

	category.Image = path

	err = u.store.Update(ctx, category)
	if err != nil {
		return fmt.Errorf("error on update category: %w", err)
	}

	err = u.cache.UpdateCache(ctx, category, models.UpdateOp)
	if err != nil {
		u.logger.Error(fmt.Sprintf("error on update cache: %v", err))
	} else {
		u.logger.Info("Update cache success")
	}
	return nil
}

func (u *categoryUsecase) DeleteImage(ctx context.Context, id uuid.UUID, name string) error {
	u.logger.Debugf("Enter in usecase DeleteCategoryImage() with args: ctx, id: %v, name: %s", id, name)

	category, err := u.store.Get(ctx, id.String())
	if err != nil {
		return fmt.Errorf("error on get category: %w", err)
	}
	if strings.Contains(category.Image, name) {
		category.Image = ""
	}

	err = u.filestorage.DeleteCategoryImage(id.String(), name)
	if err != nil {
		return fmt.Errorf("error on delete category image: %w", err)
	}

	err = u.store.Update(ctx, category)
	if err != nil {
		return fmt.Errorf("error on update category: %w", err)
	}
	err = u.cache.UpdateCache(ctx, category, models.UpdateOp)
	if err != nil {
		u.logger.Error(fmt.Sprintf("error on update cache: %v", err))
	} else {
		u.logger.Info("Update cache success")
	}
	return nil
}

func (u *categoryUsecase) ListImages(ctx context.Context) ([]*models.FileInfo, error) {
	u.logger.Debug("Enter in usecase GetCategoriesImagesList()")

	result, err := u.filestorage.GetCategoriesImagesList()
	if err != nil {
		return nil, fmt.Errorf("error on get categories images list: %w", err)
	}
	return result, nil
}

func (u *categoryUsecase) DeleteImageById(ctx context.Context, id uuid.UUID) error {
	u.logger.Debugf("Enter in usecase DeleteCategoryImageById() with arts: ctx, id: %v", id)

	err := u.filestorage.DeleteCategoryImageById(id.String())
	if err != nil {
		return fmt.Errorf("error on delete category image by id: %w", err)
	}
	return nil
}
