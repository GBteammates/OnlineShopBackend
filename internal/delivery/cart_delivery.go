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

	cartId, err := uuid.FromBytes([]byte(c.Param("cartID")))
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}

	modelCart, err := delivery.cartUsecase.GetCart(ctx, cartId)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}

	cartItems := make([]cart.CartItem, len(modelCart.Items))
	for idx, item := range modelCart.Items {
		cartItems[idx].Id = item.Id.String()
		cartItems[idx].Title = item.Title
		cartItems[idx].Price = item.Price
		if len(item.Images) > 0 {
			cartItems[idx].Image = item.Images[0]
		}
	}

	cart := cart.Cart{
		Id:     modelCart.Id.String(),
		UserId: modelCart.UserId.String(),
		Items:  cartItems,
	}

	c.JSON(http.StatusOK, cart)
}

// CreateCart - create a new cart
//
//	@Summary		Method provides to create cart with items
//	@Description	Method provides to create cart with items.
//	@Tags			carts
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		string	false	"Id of user (if user autorized)"
//	@Success		201		{object}	cart.CartId
//	@Failure		400		{object}	ErrorResponse
//	@Failure		403		"Forbidden"
//	@Failure		404		{object}	ErrorResponse	"404 Not Found"
//	@Failure		500		{object}	ErrorResponse
//	@Router			/cart/create/:userID [get]
func (delivery *Delivery) CreateCart(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery CreateCart()")
	ctx := c.Request.Context()

	userId := c.Param("userID")
	if userId != "" {
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
	cartId, err := delivery.cartUsecase.Create(ctx, uuid.Nil)
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
	itemId, err := uuid.Parse(deliveryCart.CartId)
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

	cartId, err := uuid.FromBytes([]byte(c.Param("cartID")))
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
//	@Param			cart	body	cart.ShortCart	true	"Data for delete item from cart"
//	@Success		200
//	@Failure		400	{object}	ErrorResponse
//	@Failure		403	"Forbidden"
//	@Failure		404	{object}	ErrorResponse	"404 Not Found"
//	@Failure		500	{object}	ErrorResponse
//	@Router			/cart/deleteItem [put]
func (delivery *Delivery) DeleteItemFromCart(c *gin.Context) {
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
	itemId, err := uuid.Parse(deliveryCart.CartId)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusBadRequest, err)
		return
	}

	err = delivery.cartUsecase.DeleteItemFromCart(ctx, cartId, itemId)
	if err != nil {
		delivery.logger.Error(err.Error())
		delivery.SetError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
