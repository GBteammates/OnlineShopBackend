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
	"OnlineShopBackend/internal/delivery/cart"
	"OnlineShopBackend/internal/models"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetCart - get a specific cart by id
//
//	@Summary		Get cart by id
//	@Description	The method allows you to get the cart by id.
//	@Tags			carts
//	@Accept			json
//	@Produce		json
//	@Param			cartID	path		string		true	"Id of cart"
//	@Success		200		{object}	cart.Cart	"Cart structure"
//	@Failure		400		{object}	ErrorResponse
//	@Failure		403		"Forbidden"
//	@Failure		404		{object}	ErrorResponse	"404 Not Found"
//	@Failure		500		{object}	ErrorResponse
//	@Router			/cart/{cartID} [get]
func (delivery *Delivery) GetCart(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery GetCart()")
	ctx := c.Request.Context()

	cartId, err := uuid.Parse(c.Param("cartID"))
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}

	modelCart, err := delivery.cartUsecase.GetCart(ctx, cartId)
	if err != nil && errors.Is(err, models.ErrorNotFound{}) {
		err := fmt.Errorf("cart with id: %v not found", cartId)
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusNotFound, err)
		return
	}
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}

	cartItems := make([]cart.CartItem, len(modelCart.Items))
	for idx, item := range modelCart.Items {
		cartItems[idx].Item.Id = item.Id.String()
		cartItems[idx].Item.Title = item.Title
		cartItems[idx].Item.Description = item.Description
		cartItems[idx].Item.Category.Id = item.Category.Id.String()
		cartItems[idx].Item.Category.Name = item.Category.Name
		cartItems[idx].Item.Category.Description = item.Category.Description
		cartItems[idx].Item.Category.Image = item.Category.Image
		cartItems[idx].Item.Price = item.Price
		cartItems[idx].Item.Vendor = item.Vendor
		cartItems[idx].Item.Images = item.Images
		cartItems[idx].Quantity.Quantity = item.Quantity
	}

	cart := cart.Cart{
		Id:     modelCart.Id.String(),
		UserId: modelCart.UserId.String(),
		Items:  cartItems,
	}
	cart.SortCartItems()

	c.JSON(http.StatusOK, cart)
}

// GetCartByUserId - get a specific cart by user id
//
//	@Summary		Get cart by user id
//	@Description	The method allows you to get the cart by user id.
//	@Tags			carts
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		string		true	"Id of user"
//	@Success		200		{object}	cart.Cart	"Cart structure"
//	@Failure		400		{object}	ErrorResponse
//	@Failure		403		"Forbidden"
//	@Failure		404		{object}	ErrorResponse	"404 Not Found"
//	@Failure		500		{object}	ErrorResponse
//	@Router			/cart/byUser/{userID} [get]
func (delivery *Delivery) GetCartByUserId(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery GetCartByUserId()")
	ctx := c.Request.Context()

	userId, err := uuid.Parse(c.Param("userID"))
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}

	modelCart, err := delivery.cartUsecase.GetCartByUserId(ctx, userId)
	if err != nil && errors.Is(err, models.ErrorNotFound{}) {
		err := fmt.Errorf("cart with user id: %v not found", userId)
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusNotFound, err)
		return
	}
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}

	cartItems := make([]cart.CartItem, len(modelCart.Items))
	for idx, item := range modelCart.Items {
		cartItems[idx].Item.Id = item.Id.String()
		cartItems[idx].Item.Title = item.Title
		cartItems[idx].Item.Description = item.Description
		cartItems[idx].Item.Category.Id = item.Category.Id.String()
		cartItems[idx].Item.Category.Name = item.Category.Name
		cartItems[idx].Item.Category.Description = item.Category.Description
		cartItems[idx].Item.Category.Image = item.Category.Image
		cartItems[idx].Item.Price = item.Price
		cartItems[idx].Item.Vendor = item.Vendor
		cartItems[idx].Item.Images = item.Images
		cartItems[idx].Quantity.Quantity = item.Quantity
	}

	cart := cart.Cart{
		Id:     modelCart.Id.String(),
		UserId: modelCart.UserId.String(),
		Items:  cartItems,
	}
	cart.SortCartItems()

	c.JSON(http.StatusOK, cart)
}

// CreateCart - create a new cart
//
//	@Summary		Method provides to create cart with items
//	@Description	Method provides to create cart with items.
//	@Tags			carts
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		string	true	"Id of user (if user autorized)"
//	@Success		201		{object}	cart.CartId
//	@Failure		400		{object}	ErrorResponse
//	@Failure		403		"Forbidden"
//	@Failure		404		{object}	ErrorResponse	"404 Not Found"
//	@Failure		500		{object}	ErrorResponse
//	@Router			/cart/create/{userID} [post]
func (delivery *Delivery) CreateCart(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery CreateCart()")
	ctx := c.Request.Context()

	userId := c.Param("userID")
	userUid, err := uuid.Parse(userId)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	cartId, err := delivery.cartUsecase.Create(ctx, userUid)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, cart.CartId{Value: cartId.String()})

}

// AddItemToCart - add new item to cart
//
//	@Summary		Method provides to add item to cart
//	@Description	Method provides to add item to cart.
//	@Tags			carts
//	@Accept			json
//	@Produce		json
//	@Param			cart	body	cart.ShortCart	true	"Data for add item to cart"
//	@Success		200
//	@Failure		400	{object}	ErrorResponse
//	@Failure		403	"Forbidden"
//	@Failure		404	{object}	ErrorResponse	"404 Not Found"
//	@Failure		500	{object}	ErrorResponse
//	@Router			/cart/addItem [put]
func (delivery *Delivery) AddItemToCart(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery AddItemToCart()")
	ctx := c.Request.Context()

	var deliveryCart cart.ShortCart
	if err := c.ShouldBindJSON(&deliveryCart); err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	if deliveryCart.CartId == "" || deliveryCart.ItemId == "" {
		err := fmt.Errorf("empty value of cart id or item id")
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	cartId, err := uuid.Parse(deliveryCart.CartId)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	itemId, err := uuid.Parse(deliveryCart.ItemId)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	err = delivery.cartUsecase.AddItemToCart(ctx, cartId, itemId)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// DeleteCart deleted cart by id
//
//	@Summary		Method provides to delete cart
//	@Description	Method provides to delete cart.
//	@Tags			carts
//	@Accept			json
//	@Produce		json
//	@Param			cartID	path	string	true	"id of cart"
//	@Success		200
//	@Failure		400	{object}	ErrorResponse
//	@Failure		403	"Forbidden"
//	@Failure		404	{object}	ErrorResponse	"404 Not Found"
//	@Failure		500	{object}	ErrorResponse
//	@Router			/cart/delete/{cartID} [delete]
func (delivery *Delivery) DeleteCart(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery DeleteCart()")

	ctx := c.Request.Context()

	cartId, err := uuid.Parse(c.Param("cartID"))
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}

	err = delivery.cartUsecase.DeleteCart(ctx, cartId)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// DeleteItemFromCart - delete item from cart
//
//	@Summary		Method provides to delete item from cart
//	@Description	Method provides to delete item from cart.
//	@Tags			carts
//	@Accept			json
//	@Produce		json
//	@Param			cartID	path	string	true	"id of cart"
//	@Param			itemID	path	string	true	"id of item"
//	@Success		200
//	@Failure		400	{object}	ErrorResponse
//	@Failure		403	"Forbidden"
//	@Failure		404	{object}	ErrorResponse	"404 Not Found"
//	@Failure		500	{object}	ErrorResponse
//	@Router			/cart/delete/{cartID}/{itemID} [delete]
func (delivery *Delivery) DeleteItemFromCart(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery DeleteItemFromCart()")

	ctx := c.Request.Context()

	cartId, err := uuid.Parse(c.Param("cartID"))
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	delivery.logger.Sugar().Debugf("cartId: %v", cartId)
	itemId, err := uuid.Parse(c.Param("itemID"))
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}
	delivery.logger.Sugar().Debugf("itemId: %v", itemId)

	err = delivery.cartUsecase.DeleteItemFromCart(ctx, cartId, itemId)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
