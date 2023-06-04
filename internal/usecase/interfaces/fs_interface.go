package usecase

import "OnlineShopBackend/internal/models"

type Filestorage interface {
	PutItemImage(id string, filename string, file []byte) (string, error)
	PutCategoryImage(id string, filename string, file []byte) (string, error)
	DeleteItemImage(id string, filename string) error
	DeleteCategoryImage(id string, filename string) error
	DeleteCategoryImageById(id string) error
	DeleteItemImagesFolderById(id string) error
	GetItemsImagesList() ([]*models.FileInfo, error)
	GetCategoriesImagesList() ([]*models.FileInfo, error)
}
