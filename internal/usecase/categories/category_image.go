package usecase

import (
	"OnlineShopBackend/internal/models"
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func (usecase *categoryUsecase) UploadCategoryImage(ctx context.Context, id uuid.UUID, name string, file []byte) error {
	usecase.logger.Sugar().Debugf("Enter in usecase UploadCategoryImage() with args: ctx, id: %v, name: %s, file", id, name)

	category, err := usecase.categoryStore.GetCategory(ctx, id)
	if err != nil {
		return fmt.Errorf("error on get category: %w", err)
	}

	path, err := usecase.filestorage.PutCategoryImage(id.String(), name, file)
	if err != nil {
		return fmt.Errorf("error on put image to filestorage: %w", err)
	}

	category.Image = path

	err = usecase.categoryStore.UpdateCategory(ctx, category)
	if err != nil {
		return fmt.Errorf("error on update category: %w", err)
	}

	err = usecase.categoriesCache.UpdateCategoryCache(ctx, category, updateOp)
	if err != nil {
		usecase.logger.Error(fmt.Sprintf("error on update cache: %v", err))
	} else {
		usecase.logger.Info("Update cache success")
	}
	return nil
}

func (usecase *categoryUsecase) DeleteCategoryImage(ctx context.Context, id uuid.UUID, name string) error {
	usecase.logger.Sugar().Debugf("Enter in usecase DeleteCategoryImage() with args: ctx, id: %v, name: %s", id, name)

	category, err := usecase.categoryStore.GetCategory(ctx, id)
	if err != nil {
		return fmt.Errorf("error on get category: %w", err)
	}
	if strings.Contains(category.Image, name) {
		category.Image = ""
	}

	err = usecase.filestorage.DeleteCategoryImage(id.String(), name)
	if err != nil {
		return fmt.Errorf("error on delete category image: %w", err)
	}

	err = usecase.categoryStore.UpdateCategory(ctx, category)
	if err != nil {
		return fmt.Errorf("error on update category: %w", err)
	}
	err = usecase.categoriesCache.UpdateCategoryCache(ctx, category, updateOp)
	if err != nil {
		usecase.logger.Error(fmt.Sprintf("error on update cache: %v", err))
	} else {
		usecase.logger.Info("Update cache success")
	}
	return nil
}

func (usecase *categoryUsecase) GetCategoriesImagesList(ctx context.Context) ([]*models.FileInfo, error) {
	usecase.logger.Debug("Enter in usecase GetCategoriesImagesList()")

	result, err := usecase.filestorage.GetCategoriesImagesList()
	if err != nil {
		return nil, fmt.Errorf("error on get categories images list: %w", err)
	}
	return result, nil
}

func (usecase *categoryUsecase) DeleteCategoryImageById(ctx context.Context, id uuid.UUID) error {
	usecase.logger.Sugar().Debugf("Enter in usecase DeleteCategoryImageById() with arts: ctx, id: %v", id)

	err := usecase.filestorage.DeleteCategoryImageById(id.String())
	if err != nil {
		return fmt.Errorf("error on delete category image by id: %w", err)
	}
	return nil
}
