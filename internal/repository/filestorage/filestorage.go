package filestorage

import (
	"OnlineShopBackend/internal/models"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

type filestorage struct {
	serverURL string
	path      string
	logger    *zap.Logger
}

func New(url string, path string, logger *zap.Logger) *filestorage {
	logger.Sugar().Debugf("Enter in Newfilestorage() with args: url: %s, path: %s, logger", url, path)
	return &filestorage{serverURL: url, path: path, logger: logger}
}

func (f *filestorage) PutItemImage(id string, filename string, file []byte) (filePath string, err error) {
	f.logger.Sugar().Debugf("Enter in filestorage PutItemImage() with args: id: %s, filename: %s, file", id, filename)
	_, err = os.Stat(f.path + "items/" + id)
	if os.IsNotExist(err) {
		err = os.Mkdir(f.path+"items/"+id, 0700)
		if err != nil {
			f.logger.Debug(fmt.Sprintf("error on create dir for save image %v", err))
			return "", fmt.Errorf("error on create dir for save image: %w", err)
		}
	}
	filePath = f.path + "items/" + id + "/" + filename
	if err := os.WriteFile(filePath, file, os.ModePerm); err != nil {
		f.logger.Debug(fmt.Sprintf("error on filestorage put file: %v", err))
		return "", fmt.Errorf("error on filestorage put file: %w", err)
	}
	urlPath := f.serverURL + "/files/items/" + id + "/" + filename
	f.logger.Sugar().Debugf("Put item image success, urlPath: %s", urlPath)
	return urlPath, nil
}

func (f *filestorage) PutCategoryImage(id string, filename string, file []byte) (filePath string, err error) {
	f.logger.Sugar().Debugf("Enter in filestorage PutCategoryImage() with args: id: %s, filename: %s, file", id, filename)
	_, err = os.Stat(f.path + "categories/" + id)
	if os.IsNotExist(err) {
		err = os.Mkdir(f.path+"categories/"+id, 0700)
		if err != nil {
			f.logger.Debug(fmt.Sprintf("error on create dir for save image %v", err))
			return "", fmt.Errorf("error on create dir for save image: %w", err)
		}
	}
	filePath = f.path + "categories/" + id + "/" + filename
	if err := os.WriteFile(filePath, file, os.ModePerm); err != nil {
		f.logger.Debug(fmt.Sprintf("error on filestorage put file: %v", err))
		return "", fmt.Errorf("error on filestorage put file: %w", err)
	}
	urlPath := f.serverURL + "/files/categories/" + id + "/" + filename
	f.logger.Sugar().Debugf("Put category image success, urlPath: %s", urlPath)
	return urlPath, nil
}

func (f *filestorage) DeleteItemImage(id string, filename string) error {
	f.logger.Sugar().Debugf("Enter in filestorage DeleteItemImage() with args: id: %s, filename: %s", id, filename)
	f.logger.Debug(fmt.Sprintf("name of deleting image: %s", filename))
	err := os.Remove(f.path + "items/" + id + "/" + filename)
	if err != nil {
		f.logger.Debug(fmt.Sprintf("error on delete file: %v", err))
		return fmt.Errorf("error on delete file: %w", err)
	}
	f.logger.Info("Item image delete success")
	return nil
}

func (f *filestorage) DeleteCategoryImage(id string, filename string) error {
	f.logger.Sugar().Debugf("Enter in filestorage DeleteCategoryImage() with args: id: %s, filename: %s", id, filename)
	f.logger.Debug(fmt.Sprintf("name of deleting image: %s", filename))
	err := os.Remove(f.path + "categories/" + id + "/" + filename)
	if err != nil {
		f.logger.Debug(fmt.Sprintf("error on delete file: %v", err))
		return fmt.Errorf("error on delete file: %w", err)
	}
	f.logger.Info("Category image delete success")
	return nil
}

func (f *filestorage) DeleteCategoryImageById(id string) error {
	f.logger.Sugar().Debugf("Enter in filestorage DeleteCategoryImageById() with args: id: %s", id)
	f.logger.Debug(fmt.Sprintf("name of deleting folder: %s", id))
	err := os.RemoveAll(f.path + "categories/" + id)
	if err != nil {
		f.logger.Debug(fmt.Sprintf("error on delete folder: %v", err))
		return fmt.Errorf("error on delete folder: %w", err)
	}
	f.logger.Info("Category image folder delete success")
	return nil
}

func (f *filestorage) DeleteItemImagesFolderById(id string) error {
	f.logger.Sugar().Debugf("Enter in filestorage DeleteItemImageById() with args: id: %s", id)
	f.logger.Debug(fmt.Sprintf("name of deleting folder: %s", id))
	err := os.RemoveAll(f.path + "items/" + id)
	if err != nil {
		f.logger.Debug(fmt.Sprintf("error on delete folder: %v", err))
		return fmt.Errorf("error on delete folder: %w", err)
	}
	f.logger.Info("Item images folder delete success")
	return nil
}

func (f *filestorage) GetItemsImagesList() ([]*models.FileInfo, error) {
	f.logger.Debug("Enter in filestorage GetItemsImagesList()")

	result := make([]*models.FileInfo, 0)
	err := filepath.Walk(f.path+"/items", func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			result = append(result, &models.FileInfo{
				Name:       info.Name(),
				Path:       path,
				CreateDate: info.ModTime().String(),
				ModifyDate: info.ModTime().String(),
			})
		}
		return nil
	})
	if err != nil {
		f.logger.Error(fmt.Sprintf("error on get items images: %v", err))
		return nil, err
	}
	f.logger.Info("Get items images list success")
	return result, nil
}

func (f *filestorage) GetCategoriesImagesList() ([]*models.FileInfo, error) {
	f.logger.Debug("Enter in filestorage GetCategoriesImagesList()")

	result := make([]*models.FileInfo, 0)
	err := filepath.Walk(f.path+"/categoires", func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			result = append(result, &models.FileInfo{
				Name:       info.Name(),
				Path:       path,
				CreateDate: info.ModTime().String(),
				ModifyDate: info.ModTime().String(),
			})
		}
		return nil
	})
	if err != nil {
		f.logger.Error(fmt.Sprintf("error on get categories images list: %v", err))
		return nil, err
	}
	f.logger.Info("Get categories images list success")
	return result, nil
}
