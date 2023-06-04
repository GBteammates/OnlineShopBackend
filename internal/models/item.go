package models

import (
	"context"

	"github.com/google/uuid"
)

type Item struct {
	Id          uuid.UUID
	Title       string
	Description string
	Price       int32
	Category    Category
	Vendor      string
	Images      []string
}

type ItemWithQuantity struct {
	Item
	Quantity int
}

type ListOptions struct {
	Limit     int
	Offset    int
	SortType  string
	SortOrder string
	Param     string
	Kind      string
	Handler   func(ctx context.Context, param string) (chan Item, error)
}

type CacheOptions struct {
	Op      string
	Kind    []string
	NewItem *Item
	UserId  uuid.UUID
}

type QuantityOptions struct {
	Kind    string
	Param   string
	Handler func(ctx context.Context, param string) (int, error)
}

const (
	List         = "List"
	InCategory   = "InCategory"
	Search       = "Search"
	Favourites   = "Favourites"
	Quantity     = "Quantity"
	CreateOp     = "Create"
	UpdateOp     = "Update"
	DeleteOp     = "Delete"
	ListItemsKey = "ListItems"
	ASC          = "ASC"
	DESC         = "DESC"
	Name         = "Name"
	Price        = "Price"
	NameASC      = Name + ASC
	NameDESC     = Name + DESC
	PriceASC     = Price + ASC
	PriceDESC    = Price + DESC
	Limit        = "Limit"
	Offset       = "Offset"
	SortType     = "SortType"
	SortOrder    = "SortOrder"
	FavIDs       = "FavIDs"
)
