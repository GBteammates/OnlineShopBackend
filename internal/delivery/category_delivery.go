package delivery

import (
	"OnlineShopBackend/internal/delivery/category"
	"OnlineShopBackend/internal/models"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CreateCategory - create a new category
//
//	@Summary		Method provides to create category of items
//	@Description	Method provides to create category of items.
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Param			category	body		category.ShortCategory	true	"Data for creating category"
//	@Success		201			{object}	category.CategoryId
//	@Failure		400			{object}	ErrorResponse
//	@Failure		403			"Forbidden"
//	@Failure		404			{object}	ErrorResponse	"404 Not Found"
//	@Failure		500			{object}	ErrorResponse
//	@Router			/categories/create [post]
func (delivery *Delivery) CreateCategory(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery CreateCategory()")
	ctx := c.Request.Context()
	var deliveryCategory category.ShortCategory
	if err := c.ShouldBindJSON(&deliveryCategory); err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	if deliveryCategory.Name == "NoCategory" {
		err := fmt.Errorf("can't create category with this name, name reserved by system")
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	delivery.logger.Sugar().Debugf("Binded struct: %v", deliveryCategory)
	if deliveryCategory.Name == "" && deliveryCategory.Description == "" {
		err := fmt.Errorf("empty category in request")
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	modelsCategory := models.Category{
		Name:        deliveryCategory.Name,
		Description: deliveryCategory.Description,
	}
	id, err := delivery.categoryUsecase.CreateCategory(ctx, &modelsCategory)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusCreated, category.CategoryId{Value: id.String()})
}

// UpdateCategory updating category
//
//	@Summary		Method provides to update category
//	@Description	Method provides to update category.
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Param			id			path	string					true	"id of category"
//	@Param			category	body	category.ShortCategory	true	"Data for updating category"
//	@Success		200
//	@Failure		400	{object}	ErrorResponse
//	@Failure		403	"Forbidden"
//	@Failure		404	{object}	ErrorResponse	"404 Not Found"
//	@Failure		500	{object}	ErrorResponse
//	@Router			/categories/update [put]
func (delivery *Delivery) UpdateCategory(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery UpdateCategory()")
	id := c.Param("categoryID")
	delivery.logger.Debug(fmt.Sprintf("Category id from request is %v", id))
	if id == "" {
		err := fmt.Errorf("empty id in request")
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	uid, err := uuid.Parse(id)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	var deliveryCategory category.ShortCategory
	if err := c.ShouldBindJSON(&deliveryCategory); err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	if deliveryCategory.Name == "NoCategory" {
		err := fmt.Errorf("this category is protected by changes")
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	modelsCategory := models.Category{
		Id:          uid,
		Name:        deliveryCategory.Name,
		Description: deliveryCategory.Description,
	}
	ctx := c.Request.Context()
	err = delivery.categoryUsecase.UpdateCategory(ctx, &modelsCategory)
	if err != nil && errors.Is(err, models.ErrorNotFound{}) {
		err = fmt.Errorf("category with id: %s not found", id)
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusNotFound, err)
		return
	}
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// UploadCategoryImage - upload an image
//
//	@Summary		Upload an image of category
//	@Description	Method provides to upload an image of category.
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"id of category"
//	@Param			image	formData	file	true	"picture of category"
//	@Success		201
//	@Failure		400	{object}	ErrorResponse
//	@Failure		403	"Forbidden"
//	@Failure		404	{object}	ErrorResponse	"404 Not Found"
//	@Failure		415	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Failure		507	{object}	ErrorResponse
//	@Router			/categories/image/upload/:categoryID [post]
func (delivery *Delivery) UploadCategoryImage(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery UploadCategoryImage()")
	id := c.Param("categoryID")
	if id == "" {
		err := fmt.Errorf("empty id in request")
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	var name string
	contentType := c.ContentType()

	if contentType == "image/jpeg" {
		name = carbon.Now().ToShortDateTimeString() + ".jpeg"
	} else if contentType == "image/png" {
		name = carbon.Now().ToShortDateTimeString() + ".png"
	} else {
		err := fmt.Errorf("unsupported media type: %s", contentType)
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusUnsupportedMediaType, err)
		return
	}

	file, err := io.ReadAll(c.Request.Body)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusUnsupportedMediaType, err)
		return
	}

	delivery.logger.Info("Read id", zap.String("id", id))
	delivery.logger.Info("File len=", zap.Int32("len", int32(len(file))))

	ctx := c.Request.Context()

	err = delivery.categoryUsecase.UploadCategoryImage(ctx, uid, name, file)
	if err != nil && errors.Is(err, models.ErrorNotFound{}) {
		err = fmt.Errorf("category with id: %s not found", uid)
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusNotFound, err)
		return
	}
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
}

// DeleteCategoryImage delete category image
//
//	@Summary		Delete a category image by category id
//	@Description	The method allows you to delete a category image by category id.
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Param			id		query	string	true	"Category id"
//	@Param			name	query	string	true	"Image name"
//	@Success		200
//	@Failure		400	{object}	ErrorResponse
//	@Failure		403	"Forbidden"
//	@Failure		404	{object}	ErrorResponse	"404 Not Found"
//	@Failure		500	{object}	ErrorResponse
//	@Router			/categories/image/delete [delete]
func (delivery *Delivery) DeleteCategoryImage(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery DeleteCategoryImage()")
	var imageOptions ImageOptions
	err := c.Bind(&imageOptions)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	if imageOptions.Id == "" || imageOptions.Name == "" {
		err = fmt.Errorf("empty id or image name")
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	delivery.logger.Debug(fmt.Sprintf("image options is %v", imageOptions))

	err = delivery.filestorage.DeleteCategoryImage(imageOptions.Id, imageOptions.Name)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return

	}
	ctx := c.Request.Context()

	uid, err := uuid.Parse(imageOptions.Id)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}

	err = delivery.categoryUsecase.DeleteCategoryImage(ctx, uid, imageOptions.Name)
	if err != nil && errors.Is(err, models.ErrorNotFound{}) {
		err = fmt.Errorf("category with id: %s not found", uid)
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusNotFound, err)
		return
	}
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// GetCategory - get a specific category by id
//
//	@Summary		Get category by id
//	@Description	The method allows you to get the category by id.
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Param			categoryID	path		string				true	"Id of category"
//	@Success		200			{object}	category.Category	"Category structure"
//	@Failure		400			{object}	ErrorResponse
//	@Failure		403			"Forbidden"
//	@Failure		404			{object}	ErrorResponse	"404 Not Found"
//	@Failure		500			{object}	ErrorResponse
//	@Router			/categories/{categoryID} [get]
func (delivery *Delivery) GetCategory(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery GetCategory()")
	id := c.Param("categoryID")
	if id == "" {
		err := fmt.Errorf("empty id from request")
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	uid, err := uuid.Parse(id)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	delivery.logger.Debug(fmt.Sprintf("Category id from request is %v", id))
	ctx := c.Request.Context()
	modelsCategory, err := delivery.categoryUsecase.GetCategory(ctx, uid)
	if err != nil && errors.Is(err, models.ErrorNotFound{}) {
		err = fmt.Errorf("category with id: %s not found", uid)
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusNotFound, err)
		return
	}
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, category.Category{
		Id:          modelsCategory.Id.String(),
		Name:        modelsCategory.Name,
		Description: modelsCategory.Description,
		Image:       modelsCategory.Image,
	})
}

// GetCategoryList - get a list of categories
//
//	@Summary		Get list of categories
//	@Description	Method provides to get list of categories
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Success		200	array		category.Category	"List of categories"
//	@Failure		400	{object}	ErrorResponse
//	@Failure		403	"Forbidden"
//	@Failure		404	{object}	ErrorResponse	"404 Not Found"
//	@Failure		500	{object}	ErrorResponse
//	@Router			/categories/list [get]
func (delivery *Delivery) GetCategoryList(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery GetCategoryList()")
	list, err := delivery.categoryUsecase.GetCategoryList(c.Request.Context())
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	ctx := c.Request.Context()
	categories := make([]category.Category, 0, len(list))
	// NoCategory - a special category for items without a category
	for _, cat := range list {
		// If there is a NoCategory among the categories
		if cat.Name == "NoCategory" {
			// Requesting the number of items in this category
			quantity, err := delivery.itemUsecase.ItemsQuantityInCategory(ctx, cat.Name)
			if err != nil {
				delivery.logger.Error(err.Error())
				continue
			}
			// If there are no items, this category is not added to the final list
			if quantity == 0 {
				delivery.logger.Info("NoCategory is empty")
				continue
			}
		}
		categories = append(categories, category.Category{
			Id:          cat.Id.String(),
			Name:        cat.Name,
			Description: cat.Description,
			Image:       cat.Image,
		})
	}
	c.JSON(http.StatusOK, categories)
}

// DeleteCategory deleted category by id
//
//	@Summary		Method provides to delete category
//	@Description	Method provides to delete category.
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Param			categoryID	path	string	true	"id of category"
//	@Success		200
//	@Failure		400	{object}	ErrorResponse
//	@Failure		403	"Forbidden"
//	@Failure		404	{object}	ErrorResponse	"404 Not Found"
//	@Failure		500	{object}	ErrorResponse
//	@Router			/categories/delete/{categoryID} [delete]
func (delivery *Delivery) DeleteCategory(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery DeleteCategory()")
	id := c.Param("categoryID")
	if id == "" {
		err := fmt.Errorf("empty category id in request")
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	uid, err := uuid.Parse(id)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	ctx := c.Request.Context()
	// We get the deleted category, at the same time we check its existence
	deletedCategory, err := delivery.categoryUsecase.GetCategory(ctx, uid)
	if err != nil && errors.Is(err, models.ErrorNotFound{}) {
		err = fmt.Errorf("category with id: %s not found", uid)
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusNotFound, err)
		return
	}
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	delivery.logger.Debug(fmt.Sprintf("deletedCategory: %v", deletedCategory))

	// NoCategory is a special category for items without a category that cannot be deleted
	if deletedCategory.Name == "NoCategory" {
		err = fmt.Errorf("category NoCategory is a system category and it protected by deleting")
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}

	// Requesting the number of items in the category to be deleted
	quantity, err := delivery.itemUsecase.ItemsQuantityInCategory(ctx, deletedCategory.Name)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	var items []models.Item
	// If the quantity is greater than zero, we request a list of products from this category
	if quantity > 0 {
		limitOptions := map[string]int{"offset": 0, "limit": quantity}
		sortOptions := map[string]string{"sortType": "name", "sortOrder": "asc"}
		items, err = delivery.itemUsecase.GetItemsByCategory(ctx, deletedCategory.Name, limitOptions, sortOptions)
		if err != nil {
			delivery.logger.Error(err.Error())
			delivery.SetError(c, http.StatusInternalServerError, err)
			return
		}
	}

	// Deleting a category
	err = delivery.categoryUsecase.DeleteCategory(ctx, uid)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}

	// Deleting the cache of the list of products from this category
	err = delivery.categoryUsecase.DeleteCategoryCache(ctx, deletedCategory.Name)
	if err != nil {
		delivery.logger.Error(fmt.Sprintf("error on delete category cache: %v", err))
	}

	// If the category has a picture, delete this picture
	if deletedCategory.Image != "" {
		err = delivery.filestorage.DeleteCategoryImageById(id)
		if err != nil {
			delivery.logger.Error(err.Error())
		}
	}

	// If there were no items in the category, we terminate the function
	if quantity == 0 {
		delivery.logger.Sugar().Infof("Category with id: %s deleted success", id)
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	// We are clarifying whether the NoCategory category exists in the database
	noCategory, err := delivery.categoryUsecase.GetCategoryByName(ctx, "NoCategory")
	if err != nil && errors.Is(err, models.ErrorNotFound{}) {
		delivery.logger.Error("NoCategory is not exists")
		// If such a category does not exist, create it
		noCategory := models.Category{
			Name:        "NoCategory",
			Description: "Category for items from deleting categories",
		}
		noCategoryId, err := delivery.categoryUsecase.CreateCategory(ctx, &noCategory)
		if err != nil {
			delivery.logger.Error(err.Error())
			delivery.SetError(c, http.StatusInternalServerError, err)
			return
		}
		noCategory.Id = noCategoryId
		// We iterate through the items from the category being deleted in the cycle
		for _, item := range items {
			// In each item, we change the deleted category to NoCategory
			item.Category = noCategory
			// Updating the item in the database
			err := delivery.itemUsecase.UpdateItem(ctx, &item)
			if err != nil {
				delivery.logger.Error(fmt.Sprintf("error on update item: %v", err))
			}
			// Updating the cache of the list of items in the category
			err = delivery.itemUsecase.UpdateItemsInCategoryCache(ctx, &item, "create")
			if err != nil {
				delivery.logger.Error(fmt.Sprintf("error on update cache of no category: %v", err))
			}
		}
		delivery.logger.Sugar().Infof("Category with id: %s deleted success", id)
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	if err != nil {
		delivery.logger.Error(fmt.Sprintf("error on get category by name: %v", err))
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	// We perform the same operations with items if the NoCategory already exists in the database
	for _, item := range items {
		item.Category = *noCategory
		err := delivery.itemUsecase.UpdateItem(ctx, &item)
		if err != nil {
			delivery.logger.Error(fmt.Sprintf("error on update item: %v", err))
		}
		err = delivery.itemUsecase.UpdateItemsInCategoryCache(ctx, &item, "create")
		if err != nil {
			delivery.logger.Error(fmt.Sprintf("error on update cache of no category: %v", err))
		}
	}
	delivery.logger.Sugar().Infof("Category with id: %s deleted success", id)
	c.JSON(http.StatusOK, gin.H{})
}
