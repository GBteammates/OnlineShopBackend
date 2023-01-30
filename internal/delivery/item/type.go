package item

import "OnlineShopBackend/internal/delivery/category"

type ShortItem struct {
	Title       string   `json:"title" binding:"required" example:"Пылесос"`
	Description string   `json:"description" binding:"required" example:"Мощность всасывания 1.5 кВт"`
	Category    string   `json:"category" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
	Price       int32    `json:"price" example:"1990" default:"10" binding:"required" minimum:"0"`
	Vendor      string   `json:"vendor" example:"Витязь"`
	Images      []string `json:"image,omitempty"`
}

type AddFavItem struct {
	UserId string `json:"userId" binding:"required,uuid" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
	ItemId string `json:"itemId" binding:"required,uuid" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
}

type ItemId struct {
	Value string `json:"id" uri:"itemID" binding:"required,uuid" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
}

type OutItem struct {
	Id          string            `json:"id" binding:"required,uuid" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
	Title       string            `json:"title" binding:"required" example:"Пылесос"`
	Description string            `json:"description" binding:"required" example:"Мощность всасывания 1.5 кВт"`
	Category    category.Category `json:"category" binding:"required"`
	Price       int32             `json:"price" example:"1990" default:"10" binding:"required" minimum:"0"`
	Vendor      string            `json:"vendor" binding:"required" example:"Витязь"`
	Images      []string          `json:"image,omitempty"`
	IsFavourite bool              `json:"isFavourite" example:"false"`
}

type InItem struct {
	Id          string   `json:"id" binding:"required,uuid" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
	Title       string   `json:"title" binding:"required" example:"Пылесос"`
	Description string   `json:"description" binding:"required" example:"Мощность всасывания 1.5 кВт"`
	Category    string   `json:"category" binding:"required,uuid" example:"00000000-0000-0000-0000-000000000000" format:"uuid"`
	Price       int32    `json:"price" example:"1990" default:"10" binding:"required" minimum:"0"`
	Vendor      string   `json:"vendor" binding:"required" example:"Витязь"`
	Images      []string `json:"image,omitempty"`
}

type ItemsQuantity struct {
	Quantity int `json:"quantity" example:"10" default:"0" binding:"min=0" minimum:"0"`
}

type ItemsList struct {
	List []OutItem `json:"items" binding:"min=0" minimum:"0"`
}
