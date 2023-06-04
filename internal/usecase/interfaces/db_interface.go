package usecase

import (
	"OnlineShopBackend/internal/models"
	"context"

	"github.com/google/uuid"
)

type ItemStore interface {
	Create(ctx context.Context, item *models.Item) (uuid.UUID, error)
	Update(ctx context.Context, item *models.Item) error
	Get(ctx context.Context, id uuid.UUID) (*models.Item, error)
	List(ctx context.Context, param string) (chan models.Item, error)
	SearchLine(ctx context.Context, param string) (chan models.Item, error)
	ListByCategory(ctx context.Context, param string) (chan models.Item, error)
	Delete(ctx context.Context, id uuid.UUID) error
	AddFavouriteItem(ctx context.Context, userId uuid.UUID, itemId uuid.UUID) error
	DeleteFavouriteItem(ctx context.Context, userId uuid.UUID, itemId uuid.UUID) error
	ListFavouriteItems(ctx context.Context, userId string) (chan models.Item, error)
	FavouriteItemsId(ctx context.Context, userId uuid.UUID) (*map[uuid.UUID]uuid.UUID, error)
	ListQuantity(ctx context.Context, param string) (int, error)
	ListByCategoryQuantity(ctx context.Context, param string) (int, error)
	InSearchQuantity(ctx context.Context, param string) (int, error)
	InFavouriteQuantity(ctx context.Context, userId string) (int, error)
}

type CategoryStore interface {
	Create(ctx context.Context, category *models.Category) (uuid.UUID, error)
	Update(ctx context.Context, category *models.Category) error
	Get(ctx context.Context, param string) (*models.Category, error)
	List(ctx context.Context) (chan models.Category, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type UserStore interface {
	Create(ctx context.Context, user *models.User) (uuid.UUID, error)
	UserByEmail(ctx context.Context, email string) (*models.User, error)
	GetRightsId(ctx context.Context, name string) (models.Rights, error)
	UpdateUserData(ctx context.Context, user *models.User) (*models.User, error)
	UpdateUserRole(ctx context.Context, roleId uuid.UUID, email string) error
	ListRights(ctx context.Context) (chan models.Rights, error)
	CreateRights(ctx context.Context, rights *models.Rights) (uuid.UUID, error)
}

type CartStore interface {
	Create(ctx context.Context, userId uuid.UUID) (uuid.UUID, error)
	AddItem(ctx context.Context, cartId uuid.UUID, itemId uuid.UUID) error
	Delete(ctx context.Context, cartId uuid.UUID) error
	DeleteItem(ctx context.Context, cartId uuid.UUID, itemId uuid.UUID) error
	Get(ctx context.Context, cartId uuid.UUID) (*models.Cart, error)
	CartByUserId(ctx context.Context, userId uuid.UUID) (*models.Cart, error)
}

type OrderStore interface {
	Create(ctx context.Context, order *models.Order) (*models.Order, error)
	Delete(ctx context.Context, order *models.Order) error
	ChangeAddress(ctx context.Context, order *models.Order) error
	ChangeStatus(ctx context.Context, order *models.Order) error
	Get(ctx context.Context, id uuid.UUID) (models.Order, error)
	List(ctx context.Context, user *models.User) (chan models.Order, error)
}
