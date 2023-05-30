package helpers

import (
	"OnlineShopBackend/internal/models"
	"sort"
	"strings"
)

const (
	name  = "name"
	price = "price"
	asc   = "asc"
	desc  = "desc"
)

// SortItems sorts list of items by sort parameters
func SortItems(items []models.Item, sortType string, sortOrder string) {
	sortType = strings.ToLower(sortType)
	sortOrder = strings.ToLower(sortOrder)
	switch {
	case sortType == "name" && sortOrder == "asc":
		sort.Slice(items, func(i, j int) bool { return items[i].Title < items[j].Title })
		return
	case sortType == "name" && sortOrder == "desc":
		sort.Slice(items, func(i, j int) bool { return items[i].Title > items[j].Title })
		return
	case sortType == "price" && sortOrder == "asc":
		sort.Slice(items, func(i, j int) bool { return items[i].Price < items[j].Price })
		return
	case sortType == "price" && sortOrder == "desc":
		sort.Slice(items, func(i, j int) bool { return items[i].Price > items[j].Price })
		return
	default:
	}
}

func SortOptionsFromKey(key string) (sortType, sortOrder string) {
	switch {
	case strings.Contains(key, name) && strings.Contains(key, asc):
		return name, asc
	case strings.Contains(key, name) && strings.Contains(key, desc):
		return name, desc
	case strings.Contains(key, price) && strings.Contains(key, asc):
		return price, asc
	case strings.Contains(key, price) && strings.Contains(key, desc):
		return price, desc
	default:
		return "", ""
	}
}
