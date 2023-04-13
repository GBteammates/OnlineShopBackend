package filestorage

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

type FileStorager interface {
	GetFileList() ([]FileInStorageInfo, error)
	PutItemImage(id string, filename string, file []byte) (string, error)
	PutCategoryImage(id string, filename string, file []byte) (string, error)
	DeleteItemImage(id string, filename string) error
	DeleteCategoryImage(id string, filename string) error
	DeleteCategoryImageById(id string) error
	DeleteItemImagesFolderById(id string) error
}

type FileInStorageInfo struct {
	Name       string `json:"Name"`
	Path       string `json:"Path"`
	CreateDate string `json:"CreateDate"`
	ModifyDate string `json:"ModifyDate"`
}

type OnDiskLocalStorage struct {
	serverURL string
	path      string
	logger    *zap.Logger
}

func NewOnDiskLocalStorage(url string, path string, logger *zap.Logger) *OnDiskLocalStorage {
	logger.Sugar().Debugf("Enter in NewOnDiskLocalStorage() with args: url: %s, path: %s, logger", url, path)
	d := OnDiskLocalStorage{serverURL: url, path: path, logger: logger}
	return &d
}

func (imagestorage *OnDiskLocalStorage) PutItemImage(id string, filename string, file []byte) (filePath string, err error) {
	imagestorage.logger.Sugar().Debugf("Enter in filestorage PutItemImage() with args: id: %s, filename: %s, file", id, filename)
	_, err = os.Stat(imagestorage.path + "items/" + id)
	if os.IsNotExist(err) {
		err = os.Mkdir(imagestorage.path+"items/"+id, 0700)
		if err != nil {
			imagestorage.logger.Debug(fmt.Sprintf("error on create dir for save image %v", err))
			return "", fmt.Errorf("error on create dir for save image: %w", err)
		}
	}
	filePath = imagestorage.path + "items/" + id + "/" + filename
	if err := os.WriteFile(filePath, file, os.ModePerm); err != nil {
		imagestorage.logger.Debug(fmt.Sprintf("error on filestorage put file: %v", err))
		return "", fmt.Errorf("error on filestorage put file: %w", err)
	}
	urlPath := imagestorage.serverURL + "/files/items/" + id + "/" + filename
	imagestorage.logger.Sugar().Debugf("Put item image success, urlPath: %s", urlPath)
	return urlPath, nil
}

func (imagestorage *OnDiskLocalStorage) PutCategoryImage(id string, filename string, file []byte) (filePath string, err error) {
	imagestorage.logger.Sugar().Debugf("Enter in filestorage PutCategoryImage() with args: id: %s, filename: %s, file", id, filename)
	_, err = os.Stat(imagestorage.path + "categories/" + id)
	if os.IsNotExist(err) {
		err = os.Mkdir(imagestorage.path+"categories/"+id, 0700)
		if err != nil {
			imagestorage.logger.Debug(fmt.Sprintf("error on create dir for save image %v", err))
			return "", fmt.Errorf("error on create dir for save image: %w", err)
		}
	}
	filePath = imagestorage.path + "categories/" + id + "/" + filename
	if err := os.WriteFile(filePath, file, os.ModePerm); err != nil {
		imagestorage.logger.Debug(fmt.Sprintf("error on filestorage put file: %v", err))
		return "", fmt.Errorf("error on filestorage put file: %w", err)
	}
	urlPath := imagestorage.serverURL + "/files/categories/" + id + "/" + filename
	imagestorage.logger.Sugar().Debugf("Put category image success, urlPath: %s", urlPath)
	return urlPath, nil
}

func (imagestorage *OnDiskLocalStorage) DeleteItemImage(id string, filename string) error {
	imagestorage.logger.Sugar().Debugf("Enter in filestorage DeleteItemImage() with args: id: %s, filename: %s", id, filename)
	imagestorage.logger.Debug(fmt.Sprintf("name of deleting image: %s", filename))
	err := os.Remove(imagestorage.path + "items/" + id + "/" + filename)
	if err != nil {
		imagestorage.logger.Debug(fmt.Sprintf("error on delete file: %v", err))
		return fmt.Errorf("error on delete file: %w", err)
	}
	imagestorage.logger.Info("Item image delete success")
	return nil
}

func (imagestorage *OnDiskLocalStorage) DeleteCategoryImage(id string, filename string) error {
	imagestorage.logger.Sugar().Debugf("Enter in filestorage DeleteCategoryImage() with args: id: %s, filename: %s", id, filename)
	imagestorage.logger.Debug(fmt.Sprintf("name of deleting image: %s", filename))
	err := os.Remove(imagestorage.path + "categories/" + id + "/" + filename)
	if err != nil {
		imagestorage.logger.Debug(fmt.Sprintf("error on delete file: %v", err))
		return fmt.Errorf("error on delete file: %w", err)
	}
	imagestorage.logger.Info("Category image delete success")
	return nil
}

func (imagestorage *OnDiskLocalStorage) DeleteCategoryImageById(id string) error {
	imagestorage.logger.Sugar().Debugf("Enter in filestorage DeleteCategoryImageById() with args: id: %s", id)
	imagestorage.logger.Debug(fmt.Sprintf("name of deleting folder: %s", id))
	err := os.RemoveAll(imagestorage.path + "categories/" + id)
	if err != nil {
		imagestorage.logger.Debug(fmt.Sprintf("error on delete folder: %v", err))
		return fmt.Errorf("error on delete folder: %w", err)
	}
	imagestorage.logger.Info("Category image folder delete success")
	return nil
}

func (imagestorage *OnDiskLocalStorage) DeleteItemImagesFolderById(id string) error {
	imagestorage.logger.Sugar().Debugf("Enter in filestorage DeleteItemImageById() with args: id: %s", id)
	imagestorage.logger.Debug(fmt.Sprintf("name of deleting folder: %s", id))
	err := os.RemoveAll(imagestorage.path + "items/" + id)
	if err != nil {
		imagestorage.logger.Debug(fmt.Sprintf("error on delete folder: %v", err))
		return fmt.Errorf("error on delete folder: %w", err)
	}
	imagestorage.logger.Info("Item images folder delete success")
	return nil
}

func (imagestorage *OnDiskLocalStorage) GetFileList() ([]FileInStorageInfo, error) {
	imagestorage.logger.Debug("Enter in filestorage GetFileList()")
	result := make([]FileInStorageInfo, 0)
	err := filepath.Walk(imagestorage.path, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			result = append(result, FileInStorageInfo{
				Name:       info.Name(),
				Path:       path,
				CreateDate: info.ModTime().String(),
				ModifyDate: info.ModTime().String(),
			})
		}
		return nil
	})
	if err != nil {
		imagestorage.logger.Error(fmt.Sprintf("error on filepath.Walk: %v", err))
		return nil, err
	}
	return result, nil
}
