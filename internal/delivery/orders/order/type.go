package order

import (
	"OnlineShopBackend/internal/delivery/carts/cart"
	"sort"
	"time"
)

type Order struct {
	Id           string          `json:"id" binding:"required,uuid"  example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
	Items        []cart.CartItem `json:"items,omitempty" binding:"min=0" minimum:"0"`
	UserId       string          `json:"user_id,omitempty"  example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
	CreatedAt    time.Time       `json:"created_at" binding:"required" time_format:"2006-01-02"`
	ShipmentTime time.Time       `json:"shipment_time" binding:"required" time_format:"2006-01-02"`
	Address      OrderAddress    `json:"address" binding:"required"`
	Status       string          `json:"status,omitempty"`
}

func (order *Order) SortOrderItems() {
	sort.Slice(order.Items, func(i, j int) bool { return order.Items[i].Item.Title < order.Items[j].Item.Title })
}

type OrderAddress struct {
	Zipcode string `json:"zipcode" binding:"required"`
	Country string `json:"country,omitempty"`
	City    string `json:"city" binding:"required"`
	Street  string `json:"street"  binding:"required"`
}

type UserForCart struct {
	Id    string `json:"id" binding:"required,uuid"  example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role,omitempty"`
}

type CartAdressUser struct {
	Cart    cart.Cart    `json:"cart"`
	User    UserForCart  `json:"user"`
	Address OrderAddress `json:"address"`
}

type OrderId struct {
	Id        string `json:"id" binding:"required,uuid"  example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
	NewCartId string `json:"newCartId" binding:"required,uuid"  example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
}

type AddressWithUserAndId struct {
	User    UserForCart  `json:"user"`
	Address OrderAddress `json:"address"`
	OrderId string       `json:"order_id" binding:"required,uuid" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
}

type StatusWithUserAndId struct {
	User    UserForCart `json:"user"`
	Status  string      `json:"status"`
	OrderId string      `json:"order_id" binding:"required,uuid" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
}
