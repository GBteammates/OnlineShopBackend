package delivery

import (
	"OnlineShopBackend/internal/delivery/category"
	"OnlineShopBackend/internal/delivery/file"
	"OnlineShopBackend/internal/delivery/item"
	"OnlineShopBackend/internal/delivery/user/jwtauth"
	"OnlineShopBackend/internal/metrics"
	"OnlineShopBackend/internal/models"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/golang-module/carbon/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Options is the structure for parsing offset and sort parameters
type Options struct {
	Offset    int    `form:"offset"`
	Limit     int    `form:"limit"`
	SortType  string `form:"sortType"`
	SortOrder string `form:"sortOrder"`
}

// SearchOptions is the structure for search
// items and get items by category
type SearchOptions struct {
	Param string `form:"param"`
	Options
}

// ImageOptions is the structure for deleting item image
type ImageOptions struct {
	Id   string `form:"id"`
	Name string `form:"name"`
}

// CreateItem
//
//	@Summary		Method provides to create store item
//	@Description	Method provides to create store item
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Param			item	body		item.ShortItem	true	"Data for creating item"
//	@Success		201		{object}	item.ItemId
//	@Failure		400		{object}	ErrorResponse
//	@Failure		403		"Forbidden"
//	@Failure		404		{object}	ErrorResponse	"404 Not Found"
//	@Failure		500		{object}	ErrorResponse
//	@Router			/items/create/ [post]
func (delivery *Delivery) CreateItem(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery CreateItem()")
	ctx := context.Background()
	var deliveryItem item.ShortItem
	if err := c.ShouldBindJSON(&deliveryItem); err != nil {
		delivery.logger.Error(fmt.Sprintf("error on bind json from request: %v", err))
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	if deliveryItem.Title == "" || deliveryItem.Description == "" || deliveryItem.Price == 0 {
		delivery.logger.Error(fmt.Errorf("empty item fields in request").Error())
		delivery.SetError(c, http.StatusBadRequest, fmt.Errorf("empty item fields in request"))
		return
	}

	// If item is created without specifying the category, it falls into a special category for items without category
	if deliveryItem.Category == "" {
		// We check the existence of a category for items without a category in the database
		noCategory, err := delivery.categoryUsecase.GetCategoryByName(ctx, "NoCategory")
		if err != nil && !errors.Is(err, models.ErrorNotFound{}) {
			delivery.logger.Error(err.Error())
			delivery.SetError(c, http.StatusInternalServerError, err)
			return
		}
		//If NoСategory does not yet exist in the database, we create this category
		if err != nil && errors.Is(err, models.ErrorNotFound{}) {
			delivery.logger.Sugar().Errorf("NoCategory is not exists: %v", err)
			noCategory := models.Category{
				Name:        "NoCategory",
				Description: "Category for items without categories",
			}
			noCategoryId, err := delivery.categoryUsecase.CreateCategory(ctx, &noCategory)
			if err != nil {
				delivery.logger.Error(fmt.Sprintf("error on create no category: %v", err))
				delivery.SetError(c, http.StatusInternalServerError, err)
				return
			}
			// Record the Id of the created category in the new item
			deliveryItem.Category = noCategoryId.String()
		} else if err == nil {
			// If NoCategory already exists, we write it Id in a new item
			deliveryItem.Category = noCategory.Id.String()
		}
	}

	categoryId, err := uuid.Parse(deliveryItem.Category)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	modelsItem := models.Item{
		Title:       deliveryItem.Title,
		Description: deliveryItem.Description,
		Price:       deliveryItem.Price,
		Category: models.Category{
			Id: categoryId,
		},
		Vendor: deliveryItem.Vendor,
		Images: deliveryItem.Images,
	}

	id, err := delivery.itemUsecase.CreateItem(ctx, &modelsItem)
	if err != nil {
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, item.ItemId{Value: id.String()})

	metrics.ItemsMetrics.ItemsAddedTotal.Inc()
}

// GetItem - returns item by id
//
//	@Summary		Get item by id
//	@Description	The method allows you to get the product by id.
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Param			itemID	path		string			true	"id of item"
//	@Success		200		{object}	item.OutItem	"Item structure"
//	@Failure		400		{object}	ErrorResponse
//	@Failure		403		"Forbidden"
//	@Failure		404		{object}	ErrorResponse	"404 Not Found"
//	@Failure		500		{object}	ErrorResponse
//	@Router			/items/{itemID} [get]
func (delivery *Delivery) GetItem(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery GetItem()")
	id := c.Param("itemID")
	if id == "" {
		err := fmt.Errorf("empty item in request")
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	delivery.logger.Debug(id)

	uid, err := uuid.Parse(id)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	ctx := c.Request.Context()
	modelsItem, err := delivery.itemUsecase.GetItem(ctx, uid)
	if err != nil && errors.Is(err, models.ErrorNotFound{}) {
		delivery.logger.Sugar().Errorf("item with id: %v not found", uid)
		err = fmt.Errorf("item with id: %v not found", uid)
		delivery.SetError(c, http.StatusNotFound, err)
		return
	}
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}

	result := item.OutItem{
		Id:          modelsItem.Id.String(),
		Title:       modelsItem.Title,
		Description: modelsItem.Description,
		Category: category.Category{
			Id:          modelsItem.Category.Id.String(),
			Name:        modelsItem.Category.Name,
			Description: modelsItem.Category.Description,
			Image:       modelsItem.Category.Image,
		},
		Price:  modelsItem.Price,
		Vendor: modelsItem.Vendor,
		Images: modelsItem.Images,
		// If the item in the favourites, put true, if not, put false
		IsFavourite: delivery.IsFavourite(c, modelsItem.Id),
	}
	c.JSON(http.StatusOK, result)
}

// UpdateItem - update an item
//
//	@Summary		Method provides to update store item
//	@Description	Method provides to update store item
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Param			item	body	item.InItem	true	"Data for updating item"
//	@Success		200
//	@Failure		400	{object}	ErrorResponse
//	@Failure		403	"Forbidden"
//	@Failure		404	{object}	ErrorResponse	"404 Not Found"
//	@Failure		500	{object}	ErrorResponse
//	@Router			/items/update [put]
func (delivery *Delivery) UpdateItem(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery UpdateItem()")
	ctx := c.Request.Context()
	var deliveryItem item.InItem
	if err := c.ShouldBindJSON(&deliveryItem); err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	uid, err := uuid.Parse(deliveryItem.Id)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	categoryUid, err := uuid.Parse(deliveryItem.Category)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	// If the item list is empty, add an empty line to it so as not to cause a mistake on the frontend
	if len(deliveryItem.Images) == 0 {
		deliveryItem.Images = append(deliveryItem.Images, "")
	}
	if len(deliveryItem.Images) > 1 {
		for i, v := range deliveryItem.Images {
			// If there is an empty line in the image list, we delete it to correctly display the images of item on the frontend
			if v == "" {
				deliveryItem.Images = append(deliveryItem.Images[:i], deliveryItem.Images[i+1:]...)
			}
		}
	}
	// Get the condition of item before the update
	itemBeforUpdate, err := delivery.itemUsecase.GetItem(ctx, uid)
	if err != nil && errors.Is(err, models.ErrorNotFound{}) {
		delivery.logger.Sugar().Errorf("item with id: %v not found", uid)
		err = fmt.Errorf("item with id: %v not found", uid)
		delivery.SetError(c, http.StatusNotFound, err)
		return
	}
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}

	updatingItem := &models.Item{
		Id:          uid,
		Title:       deliveryItem.Title,
		Description: deliveryItem.Description,
		Category: models.Category{
			Id: categoryUid,
		},
		Price:  deliveryItem.Price,
		Vendor: deliveryItem.Vendor,
		Images: deliveryItem.Images,
	}

	if itemBeforUpdate.Category.Id != categoryUid {
		// If the updated item has changed the category, request this category, at the same time check its existence
		updCategory, err := delivery.categoryUsecase.GetCategory(ctx, categoryUid)
		if err != nil && errors.Is(err, models.ErrorNotFound{}) {
			delivery.logger.Sugar().Errorf("category with id: %v not found", categoryUid)
			err = fmt.Errorf("category with id: %v not found", categoryUid)
			delivery.SetError(c, http.StatusNotFound, err)
			return
		}
		if err != nil {
			delivery.logger.Error(err.Error())
			delivery.SetError(c, http.StatusInternalServerError, err)
			return
		}
		updatingItem.Category = *updCategory
	}

	err = delivery.itemUsecase.UpdateItem(ctx, updatingItem)
	if err != nil && errors.Is(err, models.ErrorNotFound{}) {
		delivery.logger.Sugar().Errorf("item with id: %v not found", updatingItem.Id)
		err = fmt.Errorf("item with id: %v not found", updatingItem.Id)
		delivery.SetError(c, http.StatusNotFound, err)
		return
	}
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	// Update the cache of the item list in the category
	err = delivery.itemUsecase.UpdateItemsInCategoryCache(ctx, updatingItem, "create")
	if err != nil {
		delivery.logger.Warn(err.Error())
	} else {
		delivery.logger.Info("Category cache of updating item updated success")
	}
	if itemBeforUpdate.Category.Id != categoryUid {
		err = delivery.itemUsecase.UpdateItemsInCategoryCache(ctx, itemBeforUpdate, "delete")
		if err != nil {
			delivery.logger.Warn(err.Error())
		} else {
			delivery.logger.Info("Category cache of old item updated success")
		}
	}
	c.JSON(http.StatusOK, gin.H{})
}

// ItemsList - returns list of all items
//
//	@Summary		Get list of items
//	@Description	Method provides to get list of items
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Param			offset		query		int				false	"Offset when receiving records"	default(0)	mininum(0)
//	@Param			limit		query		int				false	"Quantity of recordings"		default(10)	minimum(0)
//	@Param			sortType	query		string			false	"Sort type (name or price)"		default("name")
//	@Param			sortOrder	query		string			false	"Sort order (asc or desc)"		default("asc")
//	@Success		200			{object}	item.ItemsList	"List of items"
//	@Failure		400			{object}	ErrorResponse
//	@Failure		403			"Forbidden"
//	@Failure		404			{object}	ErrorResponse	"404 Not Found"
//	@Failure		500			{object}	ErrorResponse
//	@Router			/items/list [get]
func (delivery *Delivery) ItemsList(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery ItemsList()")
	ctx := c.Request.Context()
	var options Options
	err := c.Bind(&options)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	delivery.logger.Debug(fmt.Sprintf("options is %v", options))
	if options.Limit == 0 {
		// If the limit is not indicated, request the quantity of items
		quantity, err := delivery.itemUsecase.ItemsQuantity(ctx)
		if err != nil {
			delivery.logger.Error(err.Error())
		}
		if quantity == 0 {
			// If there are no items, send an empty list and complete the function
			delivery.logger.Debug("quantity of items is 0")
			c.JSON(http.StatusOK, item.ItemsList{})
			return
		}
		// If the quantity of items in the range is from 1 to 30, set the limit equal to the quantity of items
		if quantity <= 30 && quantity > 0 {
			options.Limit = quantity
		} else {
			// Otherwise, set the value of items equal to 10
			options.Limit = 10
			delivery.logger.Sugar().Debugf("options limit is set in default value: %d", options.Limit)
		}
	}

	// If the sorting parameters are not set, we set the sorting by name in alphabetical order
	if options.SortType == "" {
		options.SortType = "name"
		options.SortOrder = "asc"
	}

	limitOptions := map[string]int{"offset": options.Offset, "limit": options.Limit}
	sortOptions := map[string]string{"sortType": options.SortType, "sortOrder": options.SortOrder}
	list, err := delivery.itemUsecase.ItemsList(ctx, limitOptions, sortOptions)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}

	quantity, err := delivery.itemUsecase.ItemsQuantity(ctx)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}

	items := make([]item.OutItem, len(list))
	for idx, modelsItem := range list {
		items[idx] = item.OutItem{
			Id:          modelsItem.Id.String(),
			Title:       modelsItem.Title,
			Description: modelsItem.Description,
			Category: category.Category{
				Id:          modelsItem.Category.Id.String(),
				Name:        modelsItem.Category.Name,
				Description: modelsItem.Category.Description,
				Image:       modelsItem.Category.Image,
			},
			Price:  modelsItem.Price,
			Vendor: modelsItem.Vendor,
			Images: modelsItem.Images,
			// If the item in the favourites, put true, if not, put false
			IsFavourite: delivery.IsFavourite(c, modelsItem.Id),
		}
	}
	c.JSON(http.StatusOK, item.ItemsList{
		List:     items,
		Quantity: quantity,
	})
}

// ItemsQuantity returns quantity of all items
//
//	@Summary		Get quantity of items
//	@Description	Method provides to get quantity of items
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	item.ItemsQuantity	"Quantity of items"
//	@Failure		403	"Forbidden"
//	@Failure		404	{object}	ErrorResponse	"404 Not Found"
//	@Failure		500	{object}	ErrorResponse
//	@Router			/items/quantity [get]
func (delivery *Delivery) ItemsQuantity(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery ItemsListQuantity()")
	ctx := c.Request.Context()
	quantity, err := delivery.itemUsecase.ItemsQuantity(ctx)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	itemsQuantity := item.ItemsQuantity{Quantity: quantity}
	c.JSON(http.StatusOK, itemsQuantity)
}

// ItemsQuantityInCategory returns quantity of items in category
//
//	@Summary		Get quantity of items in category
//	@Description	Method provides to get quantity of items in category
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Param			categoryName	path		string				true	"Name of category"
//	@Success		200				{object}	item.ItemsQuantity	"Quantity of items"
//	@Failure		403				"Forbidden"
//	@Failure		404				{object}	ErrorResponse	"404 Not Found"
//	@Failure		500				{object}	ErrorResponse
//	@Router			/items/quantityCat/{categoryName} [get]
func (delivery *Delivery) ItemsQuantityInCategory(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery ItemsQuantityInCategory()")
	categoryName := c.Param("categoryName")
	if categoryName == "" {
		err := fmt.Errorf("empty  categoryName is not correct")
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	ctx := c.Request.Context()
	quantity, err := delivery.itemUsecase.ItemsQuantityInCategory(ctx, categoryName)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	itemsQuantity := item.ItemsQuantity{Quantity: quantity}
	c.JSON(http.StatusOK, itemsQuantity)
}

// ItemsQuantityInFavourite returns quantity of favourite items
//
//	@Summary		Get quantity of favourite items
//	@Description	Method provides to get quantity favourite items
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		string				true	"id of user"
//	@Success		200		{object}	item.ItemsQuantity	"Quantity of items"
//	@Failure		403		"Forbidden"
//	@Failure		404		{object}	ErrorResponse	"404 Not Found"
//	@Failure		500		{object}	ErrorResponse
//	@Router			/items/quantityFav/{userID} [get]
func (delivery *Delivery) ItemsQuantityInFavourite(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery ItemsQuantityInFavourite()")
	userId, err := uuid.Parse(c.Param("userID"))
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	ctx := c.Request.Context()
	quantity, err := delivery.itemUsecase.ItemsQuantityInFavourite(ctx, userId)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	itemsQuantity := item.ItemsQuantity{Quantity: quantity}
	c.JSON(http.StatusOK, itemsQuantity)
}

// ItemsQuantityInSearch returns quantity of items in search request
//
//	@Summary		Get quantity of items in search request
//	@Description	Method provides to get quantity of items in search request
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Param			searchRequest	path		string				true	"Search request"
//	@Success		200				{object}	item.ItemsQuantity	"Quantity of items"
//	@Failure		403				"Forbidden"
//	@Failure		404				{object}	ErrorResponse	"404 Not Found"
//	@Failure		500				{object}	ErrorResponse
//	@Router			/items/quantitySearch/{searchRequest} [get]
func (delivery *Delivery) ItemsQuantityInSearch(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery ItemsQuantityInSearch()")
	searchRequest := c.Param("searchRequest")
	if searchRequest == "" {
		err := fmt.Errorf("empty  searchRequest is not correct")
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	ctx := c.Request.Context()
	quantity, err := delivery.itemUsecase.ItemsQuantityInSearch(ctx, searchRequest)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	itemsQuantity := item.ItemsQuantity{Quantity: quantity}
	c.JSON(http.StatusOK, itemsQuantity)
}

// SearchLine - returns list of items with parameters
//
//	@Summary		Get list of items by search parameters
//	@Description	Method provides to get list of items by search parameters
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Param			param		query		string			false	"Search param"
//	@Param			offset		query		int				false	"Offset when receiving records"	default(0)	mininum(0)
//	@Param			limit		query		int				false	"Quantity of recordings"		default(10)	minimum(0)
//	@Param			sortType	query		string			false	"Sort type (name or price)"		default("name")
//	@Param			sortOrder	query		string			false	"Sort order (asc or desc)"		default("asc")
//	@Success		200			{object}	item.ItemsList	"List of items"
//	@Failure		400			{object}	ErrorResponse
//	@Failure		403			"Forbidden"
//	@Failure		404			{object}	ErrorResponse	"404 Not Found"
//	@Failure		500			{object}	ErrorResponse
//	@Router			/items/search [get]
func (delivery *Delivery) SearchLine(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery SearchLine()")
	var options SearchOptions
	err := c.Bind(&options)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	delivery.logger.Debug(fmt.Sprintf("options is %v", options))
	if options.Param == "" {
		err = fmt.Errorf("empty search request")
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	// If the limit is not set to set the value of 10
	if options.Limit == 0 {
		options.Limit = 10
		delivery.logger.Sugar().Debugf("options limit is set in default value: %d", options.Limit)
	}

	// If sorting parameters are not set, sorting by name in alphabetical order is set
	if options.SortType == "" {
		options.SortType = "name"
		options.SortOrder = "asc"
		delivery.logger.Sugar().Debugf("options sort params is set in default values: sortType: %s, sortOrder: %s", options.SortType, options.SortOrder)
	}

	ctx := c.Request.Context()

	limitOptions := map[string]int{"offset": options.Offset, "limit": options.Limit}
	sortOptions := map[string]string{"sortType": options.SortType, "sortOrder": options.SortOrder}
	list, err := delivery.itemUsecase.SearchLine(ctx, options.Param, limitOptions, sortOptions)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}

	quantity, err := delivery.itemUsecase.ItemsQuantityInSearch(ctx, options.Param)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}

	items := make([]item.OutItem, len(list))
	for idx, modelsItem := range list {
		items[idx] = item.OutItem{
			Id:          modelsItem.Id.String(),
			Title:       modelsItem.Title,
			Description: modelsItem.Description,
			Category: category.Category{
				Id:          modelsItem.Category.Id.String(),
				Name:        modelsItem.Category.Name,
				Description: modelsItem.Category.Description,
				Image:       modelsItem.Category.Image,
			},
			Price:  modelsItem.Price,
			Vendor: modelsItem.Vendor,
			Images: modelsItem.Images,
			// If the item in the favourites, put true, if not, put false
			IsFavourite: delivery.IsFavourite(c, modelsItem.Id),
		}
	}
	c.JSON(http.StatusOK, item.ItemsList{
		List:     items,
		Quantity: quantity,
	})
}

// GetItemsByCategory returns list of items in category
//
//	@Summary		Get list of items by category name
//	@Description	Method provides to get list of items by category name
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Param			param		query		string			false	"Category name"
//	@Param			offset		query		int				false	"Offset when receiving records"	default(0)	mininum(0)
//	@Param			limit		query		int				false	"Quantity of recordings"		default(10)	minimum(0)
//	@Param			sortType	query		string			false	"Sort type (name or price)"		default("name")
//	@Param			sortOrder	query		string			false	"Sort order (asc or desc)"		default("asc")
//	@Success		200			{object}	item.ItemsList	"List of items"
//	@Failure		400			{object}	ErrorResponse
//	@Failure		403			"Forbidden"
//	@Failure		404			{object}	ErrorResponse	"404 Not Found"
//	@Failure		500			{object}	ErrorResponse
//	@Router			/items [get]
func (delivery *Delivery) GetItemsByCategory(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery GetItemsByCategory()")
	var options SearchOptions
	err := c.Bind(&options)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	delivery.logger.Debug(fmt.Sprintf("options is %v", options))
	if options.Param == "" {
		err = fmt.Errorf("empty search request")
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	// If the limit is not set to set the value of 10
	if options.Limit == 0 {
		options.Limit = 10
		delivery.logger.Sugar().Debugf("options limit is set in default value: %d", options.Limit)
	}

	// If sorting parameters are not set, sorting by name in alphabetical order is set
	if options.SortType == "" {
		options.SortType = "name"
		options.SortOrder = "asc"
		delivery.logger.Sugar().Debugf("options sort params is set in default values: sortType: %s, sortOrder: %s", options.SortType, options.SortOrder)
	}

	ctx := c.Request.Context()
	limitOptions := map[string]int{"offset": options.Offset, "limit": options.Limit}
	sortOptions := map[string]string{"sortType": options.SortType, "sortOrder": options.SortOrder}
	list, err := delivery.itemUsecase.GetItemsByCategory(ctx, options.Param, limitOptions, sortOptions)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}

	quantity, err := delivery.itemUsecase.ItemsQuantityInCategory(ctx, options.Param)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	items := make([]item.OutItem, len(list))
	for idx, modelsItem := range list {
		items[idx] = item.OutItem{
			Id:          modelsItem.Id.String(),
			Title:       modelsItem.Title,
			Description: modelsItem.Description,
			Category: category.Category{
				Id:          modelsItem.Category.Id.String(),
				Name:        modelsItem.Category.Name,
				Description: modelsItem.Category.Description,
				Image:       modelsItem.Category.Image,
			},
			Price:  modelsItem.Price,
			Vendor: modelsItem.Vendor,
			Images: modelsItem.Images,
			// If the item in the favourites, put true, if not, put false
			IsFavourite: delivery.IsFavourite(c, modelsItem.Id),
		}
	}
	c.JSON(http.StatusOK, item.ItemsList{
		List:     items,
		Quantity: quantity,
	})
}

// UploadItemImage - upload an image
//
//	@Summary		Upload an image of item
//	@Description	Method provides to upload an image of item
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"id of item"
//	@Param			image	formData	file	true	"picture of item"
//	@Success		201
//	@Failure		400	{object}	ErrorResponse
//	@Failure		403	"Forbidden"
//	@Failure		404	{object}	ErrorResponse	"404 Not Found"
//	@Failure		415	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Failure		507	{object}	ErrorResponse
//	@Router			/items/image/upload/:itemID [post]
func (delivery *Delivery) UploadItemImage(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery UploadItemImage()")
	ctx := c.Request.Context()
	id := c.Param("itemID")
	if id == "" {
		err := fmt.Errorf("empty search request")
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

	// By value Content Type, determine the extension of the file
	if contentType == "image/jpeg" {
		// The file name is compiled from the date and time at the time of creation
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

	err = delivery.itemUsecase.UploadItemImage(ctx, uid, name, file)
	if err != nil && errors.Is(err, models.ErrorNotFound{}) {
		delivery.logger.Sugar().Errorf("item with id: %v not found", uid)
		err = fmt.Errorf("item with id: %v not found", uid)
		delivery.SetError(c, http.StatusNotFound, err)
		return
	}
	if err != nil {
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
}

// DeleteItemImage delete an item image
//
//	@Summary		Delete an item image by item id
//	@Description	The method allows you to delete an item image by item id.
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Param			id		query	string	true	"Item id"
//	@Param			name	query	string	true	"Image name"
//	@Success		200
//	@Failure		400	{object}	ErrorResponse
//	@Failure		403	"Forbidden"
//	@Failure		404	{object}	ErrorResponse	"404 Not Found"
//	@Failure		500	{object}	ErrorResponse
//	@Router			/items/image/delete [delete]
func (delivery *Delivery) DeleteItemImage(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery DeleteItemImage()")
	var imageOptions ImageOptions
	err := c.Bind(&imageOptions)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}

	delivery.logger.Debug(fmt.Sprintf("image options is %v", imageOptions))

	if imageOptions.Id == "" || imageOptions.Name == "" {
		err := fmt.Errorf("empty image options in request")
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	uid, err := uuid.Parse(imageOptions.Id)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
	}

	ctx := c.Request.Context()

	err = delivery.itemUsecase.DeleteItemImage(ctx, uid, imageOptions.Name)
	if err != nil && errors.Is(err, models.ErrorNotFound{}) {
		delivery.logger.Sugar().Errorf("item with id: %v not found", uid)
		err = fmt.Errorf("item with id: %v not found", uid)
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

// DeleteItem deleted item by id
//
//	@Summary		Method provides to delete item
//	@Description	Method provides to delete item.
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Param			itemID	path	string	true	"id of item"
//	@Success		200
//	@Failure		400	{object}	ErrorResponse
//	@Failure		403	"Forbidden"
//	@Failure		404	{object}	ErrorResponse	"404 Not Found"
//	@Failure		500	{object}	ErrorResponse
//	@Router			/items/delete/{itemID} [delete]
func (delivery *Delivery) DeleteItem(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery DeleteItem()")
	id := c.Param("itemID")
	if id == "" {
		err := fmt.Errorf("empty item id in request")
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

	// Get the removed item
	deletedItem, err := delivery.itemUsecase.GetItem(ctx, uid)
	if err != nil && errors.Is(err, models.ErrorNotFound{}) {
		delivery.logger.Sugar().Errorf("item with id: %v not found", uid)
		err = fmt.Errorf("item with id: %v not found", uid)
		delivery.SetError(c, http.StatusNotFound, err)
		return
	}
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	delivery.logger.Debug(fmt.Sprintf("deletedItem: %v", deletedItem))

	err = delivery.itemUsecase.DeleteItem(ctx, uid)
	if err != nil && errors.Is(err, models.ErrorNotFound{}) {
		delivery.logger.Sugar().Errorf("item with id: %v not found", uid)
		err = fmt.Errorf("item with id: %v not found", uid)
		delivery.SetError(c, http.StatusNotFound, err)
		return
	}
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}

	// Update the item list cache in the deleting item's category
	err = delivery.itemUsecase.UpdateItemsInCategoryCache(ctx, deletedItem, "delete")
	if err != nil {
		delivery.logger.Sugar().Errorf("error on update cache in category items list: %v", err)
	}

	// If item has pictures, we remove them from the storage of pictures
	if len(deletedItem.Images) > 0 {
		err = delivery.itemUsecase.DeleteItemImagesFolderById(ctx, uid)
		if err != nil {
			delivery.logger.Error(err.Error())
		}
	}
	delivery.logger.Sugar().Infof("Item with id: %s deleted success", id)
	c.JSON(http.StatusOK, gin.H{})

	metrics.ItemsMetrics.ItemsDeleted.Inc()
}

// AddFavouriteItem add item in fauvorites
//
//	@Summary		Method provides add item in favourites
//	@Description	Method provides add item in favourites.
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Param			item	body	item.AddFavItem	true	"Data for add item to favourite"
//	@Success		200
//	@Failure		400	{object}	ErrorResponse
//	@Failure		403	"Forbidden"
//	@Failure		404	{object}	ErrorResponse	"404 Not Found"
//	@Failure		500	{object}	ErrorResponse
//	@Router			/items/addFavItem [post]
func (delivery *Delivery) AddFavouriteItem(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery AddFavouriteItem()")
	var addFavItem item.AddFavItem
	if err := c.ShouldBindJSON(&addFavItem); err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	userId, err := uuid.Parse(addFavItem.UserId)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	itemId, err := uuid.Parse(addFavItem.ItemId)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	ctx := c.Request.Context()
	err = delivery.itemUsecase.AddFavouriteItem(ctx, userId, itemId)
	if err != nil && errors.Is(err, models.ErrorNotFound{}) {
		delivery.logger.Sugar().Errorf("user with id: %v or item with id: %v not found", userId, itemId)
		err = fmt.Errorf("user with id: %v or item with id: %v not found", userId, itemId)
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

// DelteFavouriteItem delete item from fauvorites
//
//	@Summary		Method provides delete item from favourites
//	@Description	Method provides delete item from favourites.
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Param			userID	path	string	true	"id of user"
//	@Param			itemID	path	string	true	"id of item"
//	@Success		200
//	@Failure		400	{object}	ErrorResponse
//	@Failure		403	"Forbidden"
//	@Failure		404	{object}	ErrorResponse	"404 Not Found"
//	@Failure		500	{object}	ErrorResponse
//	@Router			/items/deleteFav/{userID}/{itemID} [delete]
func (delivery *Delivery) DeleteFavouriteItem(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery DeleteFavouriteItem()")
	userId, err := uuid.Parse(c.Param("userID"))
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	itemId, err := uuid.Parse(c.Param("itemID"))
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	ctx := c.Request.Context()
	err = delivery.itemUsecase.DeleteFavouriteItem(ctx, userId, itemId)
	if err != nil && errors.Is(err, models.ErrorNotFound{}) {
		delivery.logger.Sugar().Errorf("user with id: %v or item with id: %v not found", userId, itemId)
		err = fmt.Errorf("user with id: %v or item with id: %v not found", userId, itemId)
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

// GetFavouriteItems - returns list of all favourite items
//
//	@Summary		Get list of favourite items
//	@Description	Method provides to get list of favourite items
//	@Tags			items
//	@Accept			json
//	@Produce		json
//	@Param			param		query		string			false	"ID of user"
//	@Param			limit		query		int				false	"Quantity of recordings"		default(10)	minimum(0)
//	@Param			offset		query		int				false	"Offset when receiving records"	default(0)	mininum(0)
//	@Param			sortType	query		string			false	"Sort type (name or price)"
//	@Param			sortOrder	query		string			false	"Sort order (asc or desc)"
//	@Success		200			{object}	item.ItemsList	"List of items"
//	@Failure		400			{object}	ErrorResponse
//	@Failure		403			"Forbidden"
//	@Failure		404			{object}	ErrorResponse	"404 Not Found"
//	@Failure		500			{object}	ErrorResponse
//	@Router			/items/favList [get]
func (delivery *Delivery) GetFavouriteItems(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery GetFavouriteItems()")
	var options SearchOptions
	err := c.Bind(&options)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	delivery.logger.Debug(fmt.Sprintf("options is %v", options))
	if options.Param == "" {
		err = fmt.Errorf("empty search request")
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	// If the limit is not set to set the value of 10
	if options.Limit == 0 {
		options.Limit = 10
		delivery.logger.Sugar().Debugf("options limit is set in default value: %d", options.Limit)
	}

	// If sorting parameters are not set, sorting by name in alphabetical order is set
	if options.SortType == "" {
		options.SortType = "name"
		options.SortOrder = "asc"
		delivery.logger.Sugar().Debugf("options sort params is set in default values: sortType: %s, sortOrder: %s", options.SortType, options.SortOrder)
	}

	userId, err := uuid.Parse(options.Param)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}

	limitOptions := map[string]int{"offset": options.Offset, "limit": options.Limit}
	sortOptions := map[string]string{"sortType": options.SortType, "sortOrder": options.SortOrder}

	ctx := c.Request.Context()
	list, err := delivery.itemUsecase.GetFavouriteItems(ctx, userId, limitOptions, sortOptions)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	quantity, err := delivery.itemUsecase.ItemsQuantityInFavourite(ctx, userId)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	items := make([]item.OutItem, len(list))
	for idx, modelsItem := range list {
		items[idx] = item.OutItem{
			Id:          modelsItem.Id.String(),
			Title:       modelsItem.Title,
			Description: modelsItem.Description,
			Category: category.Category{
				Id:          modelsItem.Category.Id.String(),
				Name:        modelsItem.Category.Name,
				Description: modelsItem.Category.Description,
				Image:       modelsItem.Category.Image,
			},
			Price:       modelsItem.Price,
			Vendor:      modelsItem.Vendor,
			Images:      modelsItem.Images,
			IsFavourite: true,
		}
	}
	c.JSON(http.StatusOK, item.ItemsList{
		List:     items,
		Quantity: quantity,
	})
}

// IsFavourite checks whether item is the favourite
func (delivery *Delivery) IsFavourite(c *gin.Context, itemId uuid.UUID) bool {
	delivery.logger.Debug("Enter in delivery IsFavourite()")
	if !delivery.IsAuthorize(c) {
		return false
	}
	userId, err := delivery.GetUserId(c)
	if err != nil {
		return false
	}
	ctx := c.Request.Context()
	// Suspend the map containing the id's of the favourite items of the current user
	favIds, err := delivery.itemUsecase.GetFavouriteItemsId(ctx, userId)
	if err != nil && errors.Is(err, models.ErrorNotFound{}) {
		delivery.logger.Debug("User haven't favourite items")
		return false
	}
	if err != nil {
		delivery.logger.Error(err.Error())
		return false
	}
	favMap := *favIds
	// Check if there is an item id in the list of favourites
	_, ok := favMap[itemId]
	return ok
}

// IsAuthorize checks authorized whether the user who makes a request is
func (delivery *Delivery) IsAuthorize(c *gin.Context) bool {
	delivery.logger.Debug("Enter in delivery item IsAuthorize()")

	tokenString := c.GetHeader(authorizationHeader)

	if tokenString == "" {
		delivery.logger.Debug("Token string is empty, user not authorized")
		return false
	}

	headerSplit := strings.Split(tokenString, " ")
	if len(headerSplit) != 2 || headerSplit[0] != "Bearer" {
		delivery.logger.Debug("Header[0] is not Bearer")
		return false
	}
	if len(headerSplit[1]) == 0 {
		delivery.logger.Debug("Header[1] is empty")
		return false
	}
	return true
}

// GetUserId returns id of authorized user or error
func (delivery *Delivery) GetUserId(c *gin.Context) (uuid.UUID, error) {
	delivery.logger.Debug("Enter in delivery GetUserId()")

	tokenString := c.GetHeader(authorizationHeader)
	headerSplit := strings.Split(tokenString, " ")
	jwtKey, err := jwtauth.NewJWTKeyConfig()
	if err != nil {
		delivery.logger.Warn("Empty JWTKeyConfig")
		return uuid.Nil, err
	}

	claims := &jwtauth.Payload{}
	token, err := jwt.ParseWithClaims(headerSplit[1], claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey.Key), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			delivery.logger.Warn("Invalid signature of token")
			return uuid.Nil, err
		}
	}

	if !token.Valid {
		delivery.logger.Warn("Invalid token")
		return uuid.Nil, err
	}
	return claims.UserId, nil
}

// GetItemsImagesList returns list of items images
//
//	@Summary		Get list of items images
//	@Description	Method provides to get list of items images
//	@Tags			files
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	file.FileListResponse	"List of items images"
//	@Failure		400	{object}	ErrorResponse
//	@Failure		403	"Forbidden"
//	@Failure		404	{object}	ErrorResponse	"404 Not Found"
//	@Failure		500	{object}	ErrorResponse
//	@Router			/items/images/list [get]
func (delivery *Delivery) GetItemsImagesList(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery GetItemsImagesList()")

	ctx := c.Request.Context()
	fileInfos, err := delivery.itemUsecase.GetItemsImagesList(ctx)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	var files file.FileListResponse
	files.Files = make([]file.FilesInfo, len(fileInfos))
	for i, info := range fileInfos {
		files.Files[i] = file.FilesInfo{
			Name:       info.Name,
			Path:       info.Path,
			CreateDate: info.CreateDate,
			ModifyDate: info.ModifyDate,
		}
	}

	c.JSON(http.StatusOK, files)
}
