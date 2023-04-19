package delivery

import (
	"OnlineShopBackend/internal/delivery/carts/cart"
	"OnlineShopBackend/internal/delivery/categories/category"
	"OnlineShopBackend/internal/delivery/helper"
	"OnlineShopBackend/internal/delivery/items/item"
	"OnlineShopBackend/internal/delivery/orders/order"
	"OnlineShopBackend/internal/models"
	usecase "OnlineShopBackend/internal/usecase/interfaces"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type OrderDelivery struct {
	orderUsecase usecase.IOrderUsecase
	cartUsecase  usecase.ICartUsecase
	logger       *zap.SugaredLogger
}

func NewOrderDelivery(orderUsecase usecase.IOrderUsecase, cartUsecase usecase.ICartUsecase, logger *zap.SugaredLogger) *OrderDelivery {
	return &OrderDelivery{
		orderUsecase: orderUsecase,
		cartUsecase:  cartUsecase,
		logger:       logger,
	}
}

// Create order - create an order out of cart and user
//
//	@Summary		Create order
//	@Description	The method allows you to create an order out of cart and user info
//	@Tags			order
//	@Accept			json
//	@Produce		json
//	@Param			cartAddressUser	body		order.CartAdressUser	true	"Data for creating order"
//	@Success		201				{object}	order.OrderId			"Order id and new cart id"
//	@Failure		400				{object}	ErrorResponse
//	@Failure		403				"Forbidden"
//	@Failure		404				{object}	ErrorResponse	"404 Not Found"
//	@Failure		500				{object}	ErrorResponse
//	@Router			/order/create/ [post]
func (delivery *OrderDelivery) CreateOrder(c *gin.Context) {
	delivery.logger.Debug("Eneter in delivery CreateOrder")
	ctx := c.Request.Context()
	var cart order.CartAdressUser
	if err := c.ShouldBindJSON(&cart); err != nil {
		delivery.logger.Errorf("can't bind json from request: %s", err)
		helper.SetError(c, http.StatusBadRequest, err)
		return
	}
	id, err := uuid.Parse(cart.User.Id)
	if err != nil {
		delivery.logger.Errorf("can't parse user id: %s", err)
		helper.SetError(c, http.StatusInternalServerError, err)
		return
	}
	user := models.User{
		Id:    id,
		Email: cart.User.Email,
	}
	id, err = uuid.Parse(cart.Cart.Id)
	if err != nil {
		delivery.logger.Errorf("can't parse cart id: %s", err)
		helper.SetError(c, http.StatusInternalServerError, err)
		return
	}
	cartModel := models.Cart{
		Id:     id,
		UserId: user.Id,
		Items:  make([]models.ItemWithQuantity, 0, len(cart.Cart.Items)),
	}
	for _, oitem := range cart.Cart.Items {
		id, err = uuid.Parse(oitem.Item.Id)
		if err != nil {
			delivery.logger.Errorf("can't parse item id: %s", err)
			helper.SetError(c, http.StatusInternalServerError, err)
			return
		}
		itemM := models.ItemWithQuantity{
			Item: models.Item{
				Id:    id,
				Title: oitem.Item.Title,
				Price: oitem.Item.Price,
			},
			Quantity: oitem.Quantity.Quantity,
		}
		cartModel.Items = append(cartModel.Items, itemM)
	}

	addressMdl := models.UserAddress{
		Country: cart.Address.Country,
		City:    cart.Address.City,
		Zipcode: cart.Address.Zipcode,
		Street:  cart.Address.Street,
	}

	ordr, err := delivery.orderUsecase.PlaceOrder(ctx, &cartModel, user, addressMdl)
	if err != nil {
		delivery.logger.Errorf("can't create order: %s", err)
		helper.SetError(c, http.StatusInternalServerError, err)
		return
	}
	err = delivery.cartUsecase.DeleteCart(ctx, cartModel.Id)
	if err != nil {
		delivery.logger.Warnf("error when deleting cart with id: %v, err: %v", id, err)
	} else {
		delivery.logger.Infof("Cart with id: %v delete success", cartModel.Id)
	}
	newCartId, err := delivery.cartUsecase.CreateCart(ctx, user.Id)
	if err != nil {
		delivery.logger.Warnf("error when creating new cart for user with id: %v, err: %v", user.Id, err)
	} else {
		delivery.logger.Infof("New cart with id: %v for user with id: %v create success", newCartId, user.Id)
	}

	c.JSON(http.StatusCreated, order.OrderId{
		Id:        ordr.Id.String(),
		NewCartId: newCartId.String(),
	})
}

// GetOrder - get a specific order by id
//
//	@Summary		Get order by id
//	@Description	The method allows you to get the order by id.
//	@Tags			order
//	@Accept			json
//	@Produce		json
//	@Param			orderID	path		string		true	"Id of order"
//	@Success		200		{object}	order.Order	"Order structure"
//	@Failure		400		{object}	ErrorResponse
//	@Failure		403		"Forbidden"
//	@Failure		404		{object}	ErrorResponse	"404 Not Found"
//	@Failure		500		{object}	ErrorResponse
//	@Router			/order/{orderID} [get]
func (delivery *OrderDelivery) GetOrder(c *gin.Context) {
	delivery.logger.Debug("Enter the delivery GetOrder()")
	ctx := c.Request.Context()
	orderId, err := uuid.Parse(c.Param(("orderID")))
	if err != nil {
		delivery.logger.Errorf("can't parse order id: %s", err)
		helper.SetError(c, http.StatusBadRequest, err)
		return
	}
	modelOrder, err := delivery.orderUsecase.GetOrder(ctx, orderId)
	if err != nil {
		delivery.logger.Errorf("can't get order: %s", err)
		helper.SetError(c, http.StatusInternalServerError, err)
		return
	}
	order := order.Order{
		Id:           modelOrder.Id.String(),
		UserId:       modelOrder.User.Id.String(),
		CreatedAt:    modelOrder.CreatedAt,
		ShipmentTime: modelOrder.ShipmentTime,
		Address:      order.OrderAddress(modelOrder.Address),
		Status:       string(modelOrder.Status),
		Items:        make([]cart.CartItem, 0, len(modelOrder.Items)),
	}
	for _, oitem := range modelOrder.Items {
		cartItem := cart.CartItem{
			Item: item.OutItem{
				Id:          oitem.Id.String(),
				Title:       oitem.Title,
				Description: oitem.Description,
				Category: category.Category{
					Id:          oitem.Category.Id.String(),
					Name:        oitem.Category.Name,
					Description: oitem.Category.Description,
					Image:       oitem.Category.Image,
				},
				Price:  oitem.Price,
				Vendor: oitem.Vendor,
				Images: oitem.Images,
			},
		}
		cartItem.Quantity.Quantity = oitem.Quantity
		order.Items = append(order.Items, cartItem)
	}
	order.SortOrderItems()
	c.JSON(http.StatusOK, order)
}

// GetOrdersForUser - get a specific order by UserId
//
//	@Summary		Get all orders by UserId
//	@Description	The method allows you to get all orders by UserId.
//	@Tags			order
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		string		true	"Id of the user"
//	@Success		200		{array}		order.Order	"List of orders"
//	@Failure		400		{object}	ErrorResponse
//	@Failure		403		"Forbidden"
//	@Failure		404		{object}	ErrorResponse	"404 Not Found"
//	@Failure		500		{object}	ErrorResponse
//	@Router			/order/list/{userID} [get]
func (delivery *OrderDelivery) GetOrdersForUser(c *gin.Context) {
	delivery.logger.Debug("Enter the delivery GetOrdersForUser()")
	ctx := c.Request.Context()
	userId, err := uuid.Parse(c.Param(("userID")))
	if err != nil {
		delivery.logger.Errorf("can't parse user id: %s", err)
		helper.SetError(c, http.StatusBadRequest, err)
		return
	}
	modelOrders, err := delivery.orderUsecase.GetOrdersByUser(ctx, &models.User{Id: userId})
	if err != nil {
		delivery.logger.Errorf("can't get order: %s", err)
		helper.SetError(c, http.StatusInternalServerError, err)
		return
	}
	orders := make([]order.Order, 0, len(modelOrders))
	for _, modelOrder := range modelOrders {
		order := order.Order{
			Id:           modelOrder.Id.String(),
			UserId:       modelOrder.User.Id.String(),
			CreatedAt:    modelOrder.CreatedAt,
			ShipmentTime: modelOrder.ShipmentTime,
			Address:      order.OrderAddress(modelOrder.Address),
			Status:       string(modelOrder.Status),
			Items:        make([]cart.CartItem, 0, len(modelOrder.Items)),
		}
		for _, oitem := range modelOrder.Items {
			cartItem := cart.CartItem{
				Item: item.OutItem{
					Id:          oitem.Id.String(),
					Title:       oitem.Title,
					Description: oitem.Description,
					Category: category.Category{
						Id:          oitem.Category.Id.String(),
						Name:        oitem.Category.Name,
						Description: oitem.Category.Description,
						Image:       oitem.Category.Image,
					},
					Price:  oitem.Price,
					Vendor: oitem.Vendor,
					Images: oitem.Images,
				},
			}
			cartItem.Quantity.Quantity = oitem.Quantity
			order.Items = append(order.Items, cartItem)
		}
		order.SortOrderItems()
		orders = append(orders, order)
	}
	c.JSON(http.StatusOK, orders)
}

// DeleteOrder - delete a specific order by id
//
//	@Summary		Delete an order by id
//	@Description	The method allows you to delete an order by id.
//	@Tags			order
//	@Accept			json
//	@Produce		json
//	@Param			orderID	path	string	true	"Id of the order to delete"
//	@Success		200
//	@Failure		400	{object}	ErrorResponse
//	@Failure		403	"Forbidden"
//	@Failure		404	{object}	ErrorResponse	"404 Not Found"
//	@Failure		500	{object}	ErrorResponse
//	@Router			/order/delete/{orderID} [delete]
func (delivery *OrderDelivery) DeleteOrder(c *gin.Context) {
	delivery.logger.Debug("Enter in delivery DeleteOrder()")

	ctx := c.Request.Context()

	orderId, err := uuid.Parse(c.Param("orderID"))
	if err != nil {
		delivery.logger.Errorf("Can't parse orderID %s", err)
		helper.SetError(c, http.StatusBadRequest, err)
		return
	}

	err = delivery.orderUsecase.DeleteOrder(ctx, &models.Order{Id: orderId})
	if err != nil {
		delivery.logger.Errorf("Can't delete order with orderID %s", err)
		helper.SetError(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// ChangeAddress - change address of a specific order by Id
//
//	@Summary		Change address of a  specific order by Id
//	@Description	The method allows you to change address of an order by Id.
//	@Tags			order
//	@Accept			json
//	@Produce		json
//	@Param			AddressWithUserAndId	body	order.AddressWithUserAndId	true	"New address with orderID and user structure"
//	@Success		200
//	@Failure		400	{object}	ErrorResponse
//	@Failure		403	"Forbidden"
//	@Failure		404	{object}	ErrorResponse	"404 Not Found"
//	@Failure		500	{object}	ErrorResponse
//	@Router			/order/changeaddress/ [patch]
func (delivery *OrderDelivery) ChangeAddress(c *gin.Context) {
	delivery.logger.Debug("Enter the delivery ChangeAddress()")
	ctx := c.Request.Context()
	var address order.AddressWithUserAndId
	if err := c.ShouldBindJSON(&address); err != nil {
		delivery.logger.Errorf("can't bind json from request: %s", err)
		helper.SetError(c, http.StatusBadRequest, err)
		return
	}
	orderID, err := uuid.Parse(address.OrderId)
	if err != nil {
		delivery.logger.Errorf("can't parse order id: %s", err)
		helper.SetError(c, http.StatusBadRequest, err)
		return
	}
	if strings.ToLower(address.User.Role) == "user" {
		delivery.logger.Errorf("the action not allowed: %s", err)
		helper.SetError(c, http.StatusForbidden, err)
		return
	}
	userID, err := uuid.Parse(address.User.Id)
	if err != nil {
		delivery.logger.Errorf("can't parse order id: %s", err)
		helper.SetError(c, http.StatusBadRequest, err)
		return
	}
	err = delivery.orderUsecase.ChangeAddress(ctx,
		&models.Order{
			Id:   orderID,
			User: models.User{Id: userID},
		}, models.UserAddress(address.Address))
	if err != nil {
		delivery.logger.Errorf("can't change address for order with id: %s %s", orderID, err)
		helper.SetError(c, http.StatusInternalServerError, err)
		return
	}
}

// ChangeStatus - change status of a specific order by Id
//
//	@Summary		Change status of a specific order by Id
//	@Description	The method allows you to change status of an order by Id.
//	@Tags			order
//	@Accept			json
//	@Produce		json
//	@Param			statusWithUserAndId	body	order.StatusWithUserAndId	true	"New status with orderID and User structure"
//	@Success		200
//	@Failure		400	{object}	ErrorResponse
//	@Failure		403	"Forbidden"
//	@Failure		404	{object}	ErrorResponse	"404 Not Found"
//	@Failure		500	{object}	ErrorResponse
//	@Router			/order/changestatus/ [patch]
func (delivery *OrderDelivery) ChangeStatus(c *gin.Context) {
	delivery.logger.Debug("Enter the delivery ChangeStatus()")
	ctx := c.Request.Context()
	var status order.StatusWithUserAndId
	if err := c.ShouldBindJSON(&status); err != nil {
		delivery.logger.Errorf("can't bind json from request: %s", err)
		helper.SetError(c, http.StatusBadRequest, err)
		return
	}
	orderID, err := uuid.Parse(status.OrderId)
	if err != nil {
		delivery.logger.Errorf("can't parse order id: %s", err)
		helper.SetError(c, http.StatusBadRequest, err)
		return
	}
	if strings.ToLower(status.User.Role) == "customer" {
		delivery.logger.Errorf("the action not allowed: %s", err)
		helper.SetError(c, http.StatusForbidden, err)
		return
	}
	userID, err := uuid.Parse(status.User.Id)
	if err != nil {
		delivery.logger.Errorf("can't parse order id: %s", err)
		helper.SetError(c, http.StatusBadRequest, err)
		return
	}
	err = delivery.orderUsecase.ChangeStatus(ctx, &models.Order{
		Id: orderID,
		User: models.User{
			Id: userID,
		},
	}, models.Status(status.Status))
	if err != nil {
		delivery.logger.Errorf("can't change address for order with id: %s %s", orderID, err)
		helper.SetError(c, http.StatusInternalServerError, err)
		return
	}
}
