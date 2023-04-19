package cart

import (
	"OnlineShopBackend/internal/delivery/items/item"
	"sort"
)

type Cart struct {
	Id     string     `json:"id" binding:"required,uuid" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
	UserId string     `json:"userId,omitempty" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
	Items  []CartItem `json:"items" binding:"min=0" minimum:"0"`
}

func (cart *Cart) SortCartItems() {
	sort.Slice(cart.Items, func(i, j int) bool { return cart.Items[i].Item.Title < cart.Items[j].Item.Title })
}

type ShortCart struct {
	CartId string `json:"cartId" binding:"required,uuid" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
	ItemId string `json:"itemId" binding:"required,uuid" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
}

type CartId struct {
	Value string `json:"id" uri:"cartID" binding:"required,uuid" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
}

type CartItem struct {
	Item item.OutItem `json:"item"`
	Quantity
}

type Quantity struct {
	Quantity int `json:"quantity" example:"3" default:"1" binding:"required" minimum:"1"`
}
