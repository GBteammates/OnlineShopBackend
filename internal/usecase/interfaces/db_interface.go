package usecase

import (
	"OnlineShopBackend/internal/models"
	"context"

	"github.com/google/uuid"
)

type ItemStore interface {
	CreateItem(ctx context.Context, item *models.Item) (uuid.UUID, error)
	UpdateItem(ctx context.Context, item *models.Item) error
	GetItem(ctx context.Context, id uuid.UUID) (*models.Item, error)
	ItemsList(ctx context.Context) (chan models.Item, error)
	SearchLine(ctx context.Context, param string) (chan models.Item, error)
	GetItemsByCategory(ctx context.Context, categoryName string) (chan models.Item, error)
	DeleteItem(ctx context.Context, id uuid.UUID) error
	AddFavouriteItem(ctx context.Context, userId uuid.UUID, itemId uuid.UUID) error
	DeleteFavouriteItem(ctx context.Context, userId uuid.UUID, itemId uuid.UUID) error
	GetFavouriteItems(ctx context.Context, userId uuid.UUID) (chan models.Item, error)
	GetFavouriteItemsId(ctx context.Context, userId uuid.UUID) (*map[uuid.UUID]uuid.UUID, error)
	ItemsListQuantity(ctx context.Context) (int, error)
	ItemsByCategoryQuantity(ctx context.Context, categoryName string) (int, error)
	ItemsInSearchQuantity(ctx context.Context, searchRequest string) (int, error)
	ItemsInFavouriteQuantity(ctx context.Context, userId uuid.UUID) (int, error)
}

type CategoryStore interface {
	CreateCategory(ctx context.Context, category *models.Category) (uuid.UUID, error)
	UpdateCategory(ctx context.Context, category *models.Category) error
	GetCategory(ctx context.Context, id uuid.UUID) (*models.Category, error)
	GetCategoryList(ctx context.Context) (chan models.Category, error)
	DeleteCategory(ctx context.Context, id uuid.UUID) error
	GetCategoryByName(ctx context.Context, name string) (*models.Category, error)
}

type UserStore interface {
	CreateUser(ctx context.Context, user *models.User) (uuid.UUID, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetRightsId(ctx context.Context, name string) (models.Rights, error)
	UpdateUserData(ctx context.Context, user *models.User) (*models.User, error)
	UpdateUserRole(ctx context.Context, roleId uuid.UUID, email string) error
	GetRightsList(ctx context.Context) (chan models.Rights, error)
	CreateRights(ctx context.Context, rights *models.Rights) (uuid.UUID, error)
}

type CartStore interface {
	CreateCart(ctx context.Context, userId uuid.UUID) (uuid.UUID, error)
	AddItemToCart(ctx context.Context, cartId uuid.UUID, itemId uuid.UUID) error
	DeleteCart(ctx context.Context, cartId uuid.UUID) error
	DeleteItemFromCart(ctx context.Context, cartId uuid.UUID, itemId uuid.UUID) error
	GetCart(ctx context.Context, cartId uuid.UUID) (*models.Cart, error)
	GetCartByUserId(ctx context.Context, userId uuid.UUID) (*models.Cart, error)
}

type OrderStore interface {
	CreateOrder(ctx context.Context, order *models.Order) (*models.Order, error)
	DeleteOrder(ctx context.Context, order *models.Order) error
	ChangeAddress(ctx context.Context, order *models.Order) error
	ChangeStatus(ctx context.Context, order *models.Order) error
	GetOrderById(ctx context.Context, id uuid.UUID) (models.Order, error)
	GetOrdersByUser(ctx context.Context, user *models.User) (chan models.Order, error)
}
