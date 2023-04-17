package filestorage

import (
	"OnlineShopBackend/internal/models"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

type FileStorage struct {
	serverURL string
	path      string
	logger    *zap.Logger
}

func NewFileStorage(url string, path string, logger *zap.Logger) *FileStorage {
	logger.Sugar().Debugf("Enter in NewFileStorage() with args: url: %s, path: %s, logger", url, path)
	d := FileStorage{serverURL: url, path: path, logger: logger}
	return &d
}

func (filestorage *FileStorage) PutItemImage(id string, filename string, file []byte) (filePath string, err error) {
	filestorage.logger.Sugar().Debugf("Enter in filestorage PutItemImage() with args: id: %s, filename: %s, file", id, filename)
	_, err = os.Stat(filestorage.path + "items/" + id)
	if os.IsNotExist(err) {
		err = os.Mkdir(filestorage.path+"items/"+id, 0700)
		if err != nil {
			filestorage.logger.Debug(fmt.Sprintf("error on create dir for save image %v", err))
			return "", fmt.Errorf("error on create dir for save image: %w", err)
		}
	}
	filePath = filestorage.path + "items/" + id + "/" + filename
	if err := os.WriteFile(filePath, file, os.ModePerm); err != nil {
		filestorage.logger.Debug(fmt.Sprintf("error on filestorage put file: %v", err))
		return "", fmt.Errorf("error on filestorage put file: %w", err)
	}
	urlPath := filestorage.serverURL + "/files/items/" + id + "/" + filename
	filestorage.logger.Sugar().Debugf("Put item image success, urlPath: %s", urlPath)
	return urlPath, nil
}

func (filestorage *FileStorage) PutCategoryImage(id string, filename string, file []byte) (filePath string, err error) {
	filestorage.logger.Sugar().Debugf("Enter in filestorage PutCategoryImage() with args: id: %s, filename: %s, file", id, filename)
	_, err = os.Stat(filestorage.path + "categories/" + id)
	if os.IsNotExist(err) {
		err = os.Mkdir(filestorage.path+"categories/"+id, 0700)
		if err != nil {
			filestorage.logger.Debug(fmt.Sprintf("error on create dir for save image %v", err))
			return "", fmt.Errorf("error on create dir for save image: %w", err)
		}
	}
	filePath = filestorage.path + "categories/" + id + "/" + filename
	if err := os.WriteFile(filePath, file, os.ModePerm); err != nil {
		filestorage.logger.Debug(fmt.Sprintf("error on filestorage put file: %v", err))
		return "", fmt.Errorf("error on filestorage put file: %w", err)
	}
	urlPath := filestorage.serverURL + "/files/categories/" + id + "/" + filename
	filestorage.logger.Sugar().Debugf("Put category image success, urlPath: %s", urlPath)
	return urlPath, nil
}

func (filestorage *FileStorage) DeleteItemImage(id string, filename string) error {
	filestorage.logger.Sugar().Debugf("Enter in filestorage DeleteItemImage() with args: id: %s, filename: %s", id, filename)
	filestorage.logger.Debug(fmt.Sprintf("name of deleting image: %s", filename))
	err := os.Remove(filestorage.path + "items/" + id + "/" + filename)
	if err != nil {
		filestorage.logger.Debug(fmt.Sprintf("error on delete file: %v", err))
		return fmt.Errorf("error on delete file: %w", err)
	}
	filestorage.logger.Info("Item image delete success")
	return nil
}

func (filestorage *FileStorage) DeleteCategoryImage(id string, filename string) error {
	filestorage.logger.Sugar().Debugf("Enter in filestorage DeleteCategoryImage() with args: id: %s, filename: %s", id, filename)
	filestorage.logger.Debug(fmt.Sprintf("name of deleting image: %s", filename))
	err := os.Remove(filestorage.path + "categories/" + id + "/" + filename)
	if err != nil {
		filestorage.logger.Debug(fmt.Sprintf("error on delete file: %v", err))
		return fmt.Errorf("error on delete file: %w", err)
	}
	filestorage.logger.Info("Category image delete success")
	return nil
}

func (filestorage *FileStorage) DeleteCategoryImageById(id string) error {
	filestorage.logger.Sugar().Debugf("Enter in filestorage DeleteCategoryImageById() with args: id: %s", id)
	filestorage.logger.Debug(fmt.Sprintf("name of deleting folder: %s", id))
	err := os.RemoveAll(filestorage.path + "categories/" + id)
	if err != nil {
		filestorage.logger.Debug(fmt.Sprintf("error on delete folder: %v", err))
		return fmt.Errorf("error on delete folder: %w", err)
	}
	filestorage.logger.Info("Category image folder delete success")
	return nil
}

func (filestorage *FileStorage) DeleteItemImagesFolderById(id string) error {
	filestorage.logger.Sugar().Debugf("Enter in filestorage DeleteItemImageById() with args: id: %s", id)
	filestorage.logger.Debug(fmt.Sprintf("name of deleting folder: %s", id))
	err := os.RemoveAll(filestorage.path + "items/" + id)
	if err != nil {
		filestorage.logger.Debug(fmt.Sprintf("error on delete folder: %v", err))
		return fmt.Errorf("error on delete folder: %w", err)
	}
	filestorage.logger.Info("Item images folder delete success")
	return nil
}

func (filestorage *FileStorage) GetItemsImagesList() ([]models.FileInfo, error) {
	filestorage.logger.Debug("Enter in filestorage GetItemsImagesList()")

	result := make([]models.FileInfo, 0)
	err := filepath.Walk(filestorage.path+"/items", func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			result = append(result, models.FileInfo{
				Name:       info.Name(),
				Path:       path,
				CreateDate: info.ModTime().String(),
				ModifyDate: info.ModTime().String(),
			})
		}
		return nil
	})
	if err != nil {
		filestorage.logger.Error(fmt.Sprintf("error on get items images: %v", err))
		return nil, fmt.Errorf("error on get items images list: %w", err)
	}
	return result, nil
}

func (filestorage *FileStorage) GetCategoriesImagesList() ([]models.FileInfo, error) {
	filestorage.logger.Debug("Enter in filestorage GetCategoriesImagesList()")

	result := make([]models.FileInfo, 0)
	err := filepath.Walk(filestorage.path+"/categoires", func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			result = append(result, models.FileInfo{
				Name:       info.Name(),
				Path:       path,
				CreateDate: info.ModTime().String(),
				ModifyDate: info.ModTime().String(),
			})
		}
		return nil
	})
	if err != nil {
		filestorage.logger.Error(fmt.Sprintf("error on get categories images list: %v", err))
		return nil, fmt.Errorf("error on get categories images list: %w", err)
	}
	return result, nil
}
