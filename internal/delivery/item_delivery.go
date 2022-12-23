/*
 * Backend for Online Shop
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package delivery

import (
	"OnlineShopBackend/internal/delivery/category"
	"OnlineShopBackend/internal/delivery/item"
	"OnlineShopBackend/internal/handlers"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon/v2"
	"go.uber.org/zap"
)

type Options struct {
	Offset int `form:"offset"`
	Limit  int `form:"limit"`
}

type SearchOptions struct {
	Param  string `form:"param"`
	Offset int    `form:"offset"`
	Limit  int    `form:"limit"`
}

type ImageOptions struct {
	Id   string `form:"id"`
	Name string `form:"name"`
}

type ItemsQuantity struct {
	Quantity int `json:"quantity" example:"10" default:"0" binding:"min=0" minimum:"0"`
}

type ShortDeliveryItem struct {
	Title       string `json:"title" binding:"required" example:"Пылесос"`
	Description string `json:"description" binding:"required" example:"Мощность всасывания 1.5 кВт"`
	Category    string `json:"category" binding:"required,uuid" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
	Price       int32  `json:"price" example:"1990" default:"0" binding:"min=0" minimum:"0"`
	Vendor      string `json:"vendor" binding:"required" example:"Витязь"`
}

type ItemResponse struct {
	Id string `json:"id" binding:"required,uuid" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
}

type DeliveryItem struct {
	Id          string   `json:"id,omitempty"`
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Price       int32    `json:"price,omitempty"`
	Category    string   `json:"category,omitempty"`
	Vendor      string   `json:"vendor,omitempty"`
	Images      []string `json:"image,omitempty"`
}

// CreateItem
// @Summary Method provides to create store item
// @Description Method provides to create store item
// @Tags 	items
// @Accept  json
// @Produce json
// @Param   item 		body 		item.ShortItem	true	"Data for creating item"
// @Success 200			{object}  	item.ItemId
// @Failure 400 		{object}    ErrorResponse
// @Failure 403	 		"Forbidden"
// @Failure 404 	    {object} 	ErrorResponse				"404 Not Found"
// @Failure 500			{object}    ErrorResponse
// @Router /items/create [post]
func (delivery *Delivery) CreateItem(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery CreateItem()")
	ctx := context.Background()
	fmt.Println(fmt.Sprintln(c))
	var deliveryItem item.ShortItem
	if err := c.ShouldBindJSON(&deliveryItem); err != nil {
		delivery.logger.Error(fmt.Sprintf("error on bind json from request: %v", err))
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	if deliveryItem.Title == "" && deliveryItem.Description == "" && deliveryItem.Category == "" && deliveryItem.Price == 0 && deliveryItem.Vendor == "" {
		delivery.logger.Error(fmt.Errorf("empty item in request").Error())
		delivery.SetError(c, http.StatusBadRequest, fmt.Errorf("empty item in request"))
		return
	}
	handlersItem := handlers.Item{
		Title:       deliveryItem.Title,
		Description: deliveryItem.Description,
		Price:       deliveryItem.Price,
		Category: handlers.Category{
			Id: deliveryItem.Category,
		},
		Vendor: deliveryItem.Vendor,
	}

	id, err := delivery.itemHandlers.CreateItem(ctx, handlersItem)
	if err != nil {
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, item.ItemId{Value: id.String()})
}

// GetItem - returns item by id
// @Summary Get item by id
// @Description The method allows you to get the product by id.
// @Tags items
// @Accept  json
// @Produce json
// @Param   id 			path 		string			true	"id of item"
// @Success 200			{object}  	item.Item				"Item structure"
// @Failure 400 		{object}    ErrorResponse
// @Failure 403	 		"Forbidden"
// @Failure 404 	    {object} 	ErrorResponse			"404 Not Found"
// @Failure 500			{object}    ErrorResponse
// @Router /items/{itemID} [get]
func (delivery *Delivery) GetItem(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery GetItem()")
	id := c.Param("itemID")
	delivery.logger.Debug(id)
	if id == "" {
		delivery.logger.Error(fmt.Errorf("empty id in request").Error())
		delivery.SetError(c, http.StatusBadRequest, fmt.Errorf("empty item in request"))
	}
	ctx := c.Request.Context()
	handlersItem, err := delivery.itemHandlers.GetItem(ctx, id)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, item.Item{
		Id:          handlersItem.Id,
		Title:       handlersItem.Title,
		Description: handlersItem.Description,
		Category: category.Category{
			Id:          handlersItem.Category.Id,
			Name:        handlersItem.Category.Name,
			Description: handlersItem.Category.Description,
			Image:       handlersItem.Category.Image,
		},
		Vendor: handlersItem.Vendor,
		Images: handlersItem.Images,
	})
}

// UpdateItem - update an item
// @Summary Method provides to update store item
// @Description Method provides to update store item
// @Tags 	items
// @Accept  json
// @Produce json
// @Param   item 		body 		item.Item	true	"Data for updateing item"
// @Success 200
// @Failure 400 		{object}    ErrorResponse
// @Failure 403	 		"Forbidden"
// @Failure 404 	    {object} 	ErrorResponse				"404 Not Found"
// @Failure 500			{object}    ErrorResponse
// @Router /items/update [put]
func (delivery *Delivery) UpdateItem(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery UpdateItem()")
	ctx := c.Request.Context()
	var deliveryItem item.Item
	if err := c.ShouldBindJSON(&deliveryItem); err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	err := delivery.itemHandlers.UpdateItem(ctx, handlers.Item{
		Id:          deliveryItem.Id,
		Title:       deliveryItem.Title,
		Description: deliveryItem.Description,
		Category:    handlers.Category(deliveryItem.Category),
		Price:       deliveryItem.Price,
		Vendor:      deliveryItem.Vendor,
		Images:      deliveryItem.Images,
	})
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// ItemsList - returns list of all items
// @Summary Get list of items
// @Description Method provides to get list of items
// @Tags items
// @Accept  json
// @Produce json
// @Param 	limit 		query 		int 					false "Quantity of recordings" default(10) minimum(0)
// @Param 	offset 		query 		int 					false "Offset when receiving records" default(0) mininum(0)
// @Success 200			{object}  	item.ItemsList		  	"List of items"
// @Failure 400 		{object}    ErrorResponse
// @Failure 403	 		"Forbidden"
// @Failure 404 	    {object} 	ErrorResponse			"404 Not Found"
// @Failure 500			{object}    ErrorResponse
// @Router /items/list [get]
func (delivery *Delivery) ItemsList(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery ItemsList()")
	var options Options
	err := c.Bind(&options)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	delivery.logger.Debug(fmt.Sprintf("options is %v", options))

	if options.Limit == 0 {
		quantity, err := delivery.itemHandlers.ItemsQuantity(c.Request.Context())
		if err != nil {
			delivery.logger.Error(err.Error())
		}
		if quantity == 0 {
			delivery.logger.Debug("quantity of items is 0")
			c.JSON(http.StatusOK, item.ItemsList{})
			return
		}
		if quantity <= 30 && quantity > 0 {
			options.Limit = quantity
		} else {
			options.Limit = 10
		}
	}
	delivery.logger.Sugar().Debugf("options limit is set in default value: %d", options.Limit)

	list, err := delivery.itemHandlers.ItemsList(c.Request.Context(), options.Offset, options.Limit)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	var itemsList item.ItemsList
	itemsList.List = make([]item.Item, len(list))
	for i, it := range list {
		itemsList.List[i] = item.Item{
			Id:          it.Id,
			Title:       it.Description,
			Description: it.Description,
			Category:    category.Category(it.Category),
			Price:       it.Price,
			Vendor:      it.Vendor,
			Images:      it.Images,
		}
	}
	c.JSON(http.StatusOK, itemsList)
}

// ItemsQuantity returns quantity of all items
// @Summary Get quantity of items
// @Description Method provides to get quantity of items
// @Tags items
// @Accept  json
// @Produce json
// @Success 200			{object}  	ItemsQuantity		  	"Quantity of items"
// @Failure 403	 		"Forbidden"
// @Failure 404 	    {object} 	ErrorResponse			"404 Not Found"
// @Failure 500			{object}    ErrorResponse
// @Router /items/quantity [get]
func (delivery *Delivery) ItemsQuantity(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery ItemsQuantity()")
	ctx := c.Request.Context()
	quantity, err := delivery.itemHandlers.ItemsQuantity(ctx)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	itemsQuantity := ItemsQuantity{Quantity: quantity}
	c.JSON(http.StatusOK, itemsQuantity)
}

// SearchLine - returns list of items with parameters
// @Summary Get list of items by search parameters
// @Description Method provides to get list of items by search parameters
// @Tags items
// @Accept  json
// @Produce json
// @Param	param		query 		string			false	"Search param"
// @Param 	limit 		query 		int 			false	"Quantity of recordings" default(10) minimum(0)
// @Param 	offset 		query 		int 			false 	"Offset when receiving records" default(0) mininum(0)
// @Success 200			{object}  	item.ItemsList		  	"List of items"
// @Failure 400 		{object}    ErrorResponse
// @Failure 403	 		"Forbidden"
// @Failure 404 	    {object} 	ErrorResponse			"404 Not Found"
// @Failure 500			{object}    ErrorResponse
// @Router /items/search [get]
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

	if options.Limit == 0 {
		options.Limit = 10
	}

	delivery.logger.Sugar().Debugf("options limit is set in default value: %d", options.Limit)

	list, err := delivery.itemHandlers.SearchLine(c.Request.Context(), options.Param, options.Offset, options.Limit)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	var itemsList item.ItemsList
	itemsList.List = make([]item.Item, len(list))
	for i, it := range list {
		itemsList.List[i] = item.Item{
			Id:          it.Id,
			Title:       it.Description,
			Description: it.Description,
			Category:    category.Category(it.Category),
			Price:       it.Price,
			Vendor:      it.Vendor,
			Images:      it.Images,
		}
	}
	c.JSON(http.StatusOK, itemsList)
}

// GetItemsByCategory returns list of items in category
// @Summary Get list of items by category name
// @Description Method provides to get list of items by category name
// @Tags items
// @Accept  json
// @Produce json
// @Param	param		query 		string			false	"Category name"
// @Param 	limit 		query 		int 			false	"Quantity of recordings" default(10) minimum(0)
// @Param 	offset 		query 		int 			false 	"Offset when receiving records" default(0) mininum(0)
// @Success 200			{object}  	item.ItemsList		  	"List of items"
// @Failure 400 		{object}    ErrorResponse
// @Failure 403	 		"Forbidden"
// @Failure 404 	    {object} 	ErrorResponse			"404 Not Found"
// @Failure 500			{object}    ErrorResponse
// @Router /items [get]
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
	if options.Limit == 0 {
		options.Limit = 10
	}
	delivery.logger.Sugar().Debugf("options limit is set in default value: %d", options.Limit)

	list, err := delivery.itemHandlers.GetItemsByCategory(c.Request.Context(), options.Param, options.Offset, options.Limit)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	var itemsList item.ItemsList
	itemsList.List = make([]item.Item, len(list))
	for i, it := range list {
		itemsList.List[i] = item.Item{
			Id:          it.Id,
			Title:       it.Description,
			Description: it.Description,
			Category:    category.Category(it.Category),
			Price:       it.Price,
			Vendor:      it.Vendor,
			Images:      it.Images,
		}
	}
	c.JSON(http.StatusOK, itemsList)
}

// UploadItemImage - upload an image
func (delivery *Delivery) UploadItemImage(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery UploadItemImage()")
	ctx := c.Request.Context()
	id := c.Param("itemID")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty item id"})
		return
	}
	var name string
	contentType := c.ContentType()

	if contentType == "image/jpeg" {
		name = carbon.Now().ToShortDateTimeString() + ".jpeg"
	} else if contentType == "image/png" {
		name = carbon.Now().ToShortDateTimeString() + ".png"
	} else {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{})
		return
	}

	file, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{})
		return
	}

	delivery.logger.Info("Read id", zap.String("id", id))
	delivery.logger.Info("File len=", zap.Int32("len", int32(len(file))))
	path, err := delivery.filestorage.PutItemImage(id, name, file)
	if err != nil {
		c.JSON(http.StatusInsufficientStorage, gin.H{})
		return
	}

	item, err := delivery.itemHandlers.GetItem(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	item.Images = append(item.Images, path)

	err = delivery.itemHandlers.UpdateItem(ctx, item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "upload image success"})
}

// DeleteItemImage delete an item image
func (delivery *Delivery) DeleteItemImage(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery DeleteItemImage()")
	var imageOptions ImageOptions
	err := c.Bind(&imageOptions)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	delivery.logger.Debug(fmt.Sprintf("image options is %v", imageOptions))

	if imageOptions.Id == "" || imageOptions.Name == "" {
		fmt.Println("empty item")
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("empty item id or file name")})
		return
	}
	err = delivery.filestorage.DeleteItemImage(imageOptions.Id, imageOptions.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}
	ctx := c.Request.Context()
	item, err := delivery.itemHandlers.GetItem(ctx, imageOptions.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for idx, imagePath := range item.Images {
		if strings.Contains(imagePath, imageOptions.Name) {
			item.Images = append(item.Images[:idx], item.Images[idx+1:]...)
			break
		}
	}
	err = delivery.itemHandlers.UpdateItem(ctx, item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "delete image success"})
}
