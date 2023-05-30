package usecase

import (
	"OnlineShopBackend/internal/models"
	"context"

	"github.com/google/uuid"
)

type IItemUsecase interface {
	CreateItem(ctx context.Context, item *models.Item) (uuid.UUID, error)
	UpdateItem(ctx context.Context, item *models.Item) error
	GetItem(ctx context.Context, id uuid.UUID) (*models.Item, error)
	ItemsList(ctx context.Context, opts *models.ItemsListOptions) ([]models.Item, error)
	ItemsQuantity(ctx context.Context) (int, error)
	ItemsQuantityInCategory(ctx context.Context, categoryName string) (int, error)
	SearchLine(ctx context.Context, param string, opts *models.ItemsListOptions) ([]models.Item, error)
	GetItemsByCategory(ctx context.Context, categoryName string, opts *models.ItemsListOptions) ([]models.Item, error)
	DeleteItem(ctx context.Context, id uuid.UUID) error
	AddFavouriteItem(ctx context.Context, userId uuid.UUID, itemId uuid.UUID) error
	DeleteFavouriteItem(ctx context.Context, userId uuid.UUID, itemId uuid.UUID) error
	GetFavouriteItems(ctx context.Context, userId uuid.UUID, opts *models.ItemsListOptions) ([]models.Item, error)
	ItemsQuantityInFavourite(ctx context.Context, userId uuid.UUID) (int, error)
	ItemsQuantityInSearch(ctx context.Context, search string) (int, error)
	GetFavouriteItemsId(ctx context.Context, userId uuid.UUID) (*map[uuid.UUID]uuid.UUID, error)
	UploadItemImage(ctx context.Context, id uuid.UUID, name string, file []byte) error
	DeleteItemImage(ctx context.Context, id uuid.UUID, name string) error
	GetItemsImagesList(ctx context.Context) ([]*models.FileInfo, error)
	DeleteItemImagesFolderById(ctx context.Context, id uuid.UUID) error
}

type ICategoryUsecase interface {
	CreateCategory(ctx context.Context, category *models.Category) (uuid.UUID, error)
	UpdateCategory(ctx context.Context, category *models.Category) error
	GetCategory(ctx context.Context, id uuid.UUID) (*models.Category, error)
	GetCategoryList(ctx context.Context) ([]models.Category, error)
	DeleteCategory(ctx context.Context, id uuid.UUID) error
	GetCategoryByName(ctx context.Context, name string) (*models.Category, error)
	UploadCategoryImage(ctx context.Context, id uuid.UUID, name string, file []byte) error
	DeleteCategoryImage(ctx context.Context, id uuid.UUID, name string) error
	GetCategoriesImagesList(ctx context.Context) ([]*models.FileInfo, error)
	DeleteCategoryImageById(ctx context.Context, id uuid.UUID) error
}

type IOrderUsecase interface {
	CreateOrder(ctx context.Context, cart *models.Cart, user models.User, address models.UserAddress) (*models.Order, error)
	ChangeStatus(ctx context.Context, order *models.Order) error
	GetOrdersByUser(ctx context.Context, user *models.User) ([]models.Order, error)
	DeleteOrder(ctx context.Context, order *models.Order) error
	ChangeAddress(ctx context.Context, order *models.Order) error
	GetOrder(ctx context.Context, id uuid.UUID) (*models.Order, error)
}
type ICartUsecase interface {
	GetCart(ctx context.Context, cartId uuid.UUID) (*models.Cart, error)
	DeleteItemFromCart(ctx context.Context, cartId uuid.UUID, itemId uuid.UUID) error
	CreateCart(ctx context.Context, userId uuid.UUID) (uuid.UUID, error)
	AddItemToCart(ctx context.Context, cartId uuid.UUID, itemId uuid.UUID) error
	DeleteCart(ctx context.Context, cartId uuid.UUID) error
	GetCartByUserId(ctx context.Context, userId uuid.UUID) (*models.Cart, error)
}

type IUserUsecase interface {
	CreateUser(ctx context.Context, user *models.User) (uuid.UUID, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetRightsId(ctx context.Context, name string) (*models.Rights, error)
	UpdateUserData(ctx context.Context, user *models.User) (*models.User, error)
	UpdateUserRole(ctx context.Context, roleId uuid.UUID, email string) error
	GetRightsList(ctx context.Context) ([]models.Rights, error)
	CreateRights(ctx context.Context, rights *models.Rights) (uuid.UUID, error)
}
