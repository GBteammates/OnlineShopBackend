package helpers

import (
	"OnlineShopBackend/internal/models"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSortItems(t *testing.T) {
	testItems := []models.Item{
		{Title: "A"},
		{Title: "C"},
		{Title: "B"},
	}
	testItems2 := []models.Item{
		{Price: 10},
		{Price: 30},
		{Price: 20},
	}

	SortItems(testItems, "name", "asc")
	require.Equal(t, testItems, []models.Item{
		{Title: "A"},
		{Title: "B"},
		{Title: "C"},
	})
	SortItems(testItems, "name", "desc")
	require.Equal(t, testItems, []models.Item{
		{Title: "C"},
		{Title: "B"},
		{Title: "A"},
	})
	SortItems(testItems2, "price", "asc")
	require.Equal(t, testItems2, []models.Item{
		{Price: 10},
		{Price: 20},
		{Price: 30},
	})
	SortItems(testItems2, "price", "desc")
	require.Equal(t, testItems2, []models.Item{
		{Price: 30},
		{Price: 20},
		{Price: 10},
	})
	SortItems(testItems, "pricee", "desc")
}
