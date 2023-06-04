package usecase

import (
	"OnlineShopBackend/internal/models"
	"context"

	"github.com/google/uuid"
)

type IItemUsecase interface {
	Create(ctx context.Context, item *models.Item) (uuid.UUID, error)
	Update(ctx context.Context, item *models.Item) error
	Get(ctx context.Context, id uuid.UUID) (*models.Item, error)
	List(ctx context.Context, opts *models.ListOptions) ([]models.Item, error)
	Quantity(ctx context.Context, opts *models.QuantityOptions) (int, error)
	Delete(ctx context.Context, id uuid.UUID) error
	AddFavouriteItem(ctx context.Context, userId uuid.UUID, itemId uuid.UUID) error
	DeleteFavouriteItem(ctx context.Context, userId uuid.UUID, itemId uuid.UUID) error
	ListFavouriteItemsId(ctx context.Context, userId uuid.UUID) (*map[uuid.UUID]uuid.UUID, error)
	UploadImage(ctx context.Context, id uuid.UUID, name string, file []byte) error
	DeleteImage(ctx context.Context, id uuid.UUID, name string) error
	ListImages(ctx context.Context) ([]*models.FileInfo, error)
	DeleteImagesFolder(ctx context.Context, id uuid.UUID) error
}

type ICategoryUsecase interface {
	Create(ctx context.Context, category *models.Category) (uuid.UUID, error)
	Update(ctx context.Context, category *models.Category) error
	Get(ctx context.Context, param string) (*models.Category, error)
	List(ctx context.Context) ([]models.Category, error)
	Delete(ctx context.Context, id uuid.UUID) error
	UploadImage(ctx context.Context, id uuid.UUID, name string, file []byte) error
	DeleteImage(ctx context.Context, id uuid.UUID, name string) error
	ListImages(ctx context.Context) ([]*models.FileInfo, error)
	DeleteImageById(ctx context.Context, id uuid.UUID) error
}

type IOrderUsecase interface {
	Create(ctx context.Context, cart *models.Cart, user models.User, address models.UserAddress) (*models.Order, error)
	ChangeStatus(ctx context.Context, order *models.Order) error
	List(ctx context.Context, user *models.User) ([]models.Order, error)
	Delete(ctx context.Context, order *models.Order) error
	ChangeAddress(ctx context.Context, order *models.Order) error
	Get(ctx context.Context, id uuid.UUID) (*models.Order, error)
}
type ICartUsecase interface {
	Get(ctx context.Context, cartId uuid.UUID) (*models.Cart, error)
	DeleteItem(ctx context.Context, cartId uuid.UUID, itemId uuid.UUID) error
	Create(ctx context.Context, userId uuid.UUID) (uuid.UUID, error)
	AddItem(ctx context.Context, cartId uuid.UUID, itemId uuid.UUID) error
	Delete(ctx context.Context, cartId uuid.UUID) error
	CartByUserId(ctx context.Context, userId uuid.UUID) (*models.Cart, error)
}

type IUserUsecase interface {
	Create(ctx context.Context, user *models.User) (uuid.UUID, error)
	UserByEmail(ctx context.Context, email string) (*models.User, error)
	GetRightsId(ctx context.Context, name string) (*models.Rights, error)
	UpdateUserData(ctx context.Context, user *models.User) (*models.User, error)
	UpdateUserRole(ctx context.Context, roleId uuid.UUID, email string) error
	ListRights(ctx context.Context) ([]models.Rights, error)
	CreateRights(ctx context.Context, rights *models.Rights) (uuid.UUID, error)
}
