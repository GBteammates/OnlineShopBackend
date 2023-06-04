package usecase

import (
	"OnlineShopBackend/internal/models"
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func (u *itemUsecase) UploadImage(ctx context.Context, id uuid.UUID, name string, file []byte) error {
	u.logger.Debugf("Enter in usecase UploadImage() with args: ctx, id: %v, name: %s, file", id, name)
	// Request item for which the picture is installed
	item, err := u.store.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("error on get item: %w", err)
	}

	// Put the picture in the file storage and get it url
	path, err := u.filestorage.PutItemImage(id.String(), name, file)
	if err != nil {
		return fmt.Errorf("error on put image to filestorage: %w", err)
	}

	// Add url of picture to the item pictures list
	item.Images = append(item.Images, path)
	for i, v := range item.Images {
		// If the list of pictures has an empty line, remove it from the list
		if v == "" {
			item.Images = append(item.Images[:i], item.Images[i+1:]...)
		}
	}

	err = u.store.Update(ctx, item)
	if err != nil {
		return fmt.Errorf("error on update item: %w", err)
	}
	err = u.cache.UpdateCache(ctx, &models.CacheOptions{
		Op:      models.UpdateOp,
		Kind:    []string{models.List, models.InCategory},
		NewItem: item,
	})
	if err != nil {
		u.logger.Debugf("error on update cache: %v", err)
	}
	return nil
}

func (u *itemUsecase) DeleteImage(ctx context.Context, id uuid.UUID, name string) error {
	u.logger.Debugf("Enter in usecase DeleteItemImage() with args: ctx, id: %v, name: %s", id, name)

	// Get item from which the picture is deleted
	item, err := u.store.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("error on get item: %w", err)
	}

	err = u.filestorage.DeleteItemImage(id.String(), name)
	if err != nil {
		return fmt.Errorf("error on delete image from filestorage: %w", err)
	}

	// Delete the address of the picture from the list of pictures of item
	for idx, imagePath := range item.Images {
		if strings.Contains(imagePath, name) {
			item.Images = append(item.Images[:idx], item.Images[idx+1:]...)
			break
		}
	}
	// If, after deleting the picture from the list, the list is empty - add
	// an empty line there so that item is correctly displayed on the frontend
	if len(item.Images) == 0 {
		item.Images = append(item.Images, "")
	}
	err = u.store.Update(ctx, item)
	if err != nil {
		return fmt.Errorf("error on update item: %w", err)
	}
	err = u.cache.UpdateCache(ctx, &models.CacheOptions{
		Op:      models.UpdateOp,
		Kind:    []string{models.List, models.InCategory},
		NewItem: item,
	})
	if err != nil {
		u.logger.Debugf("error on update cache: %v", err)
	}
	return nil
}

func (u *itemUsecase) ListImages(ctx context.Context) ([]*models.FileInfo, error) {
	u.logger.Debug("Enter in usecase GetItemsImagesList()")

	result, err := u.filestorage.GetItemsImagesList()
	if err != nil {
		return nil, fmt.Errorf("error on get items images list: %w", err)
	}
	return result, nil
}

func (u *itemUsecase) DeleteImagesFolder(ctx context.Context, id uuid.UUID) error {
	u.logger.Debugf("Enter in usecase DeleteItemImagesFolderById() with args: ctx, id: %v", id)

	err := u.filestorage.DeleteItemImagesFolderById(id.String())
	if err != nil {
		return fmt.Errorf("error on delete item images folder: %w", err)
	}
	return nil
}
