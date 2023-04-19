package user_usecase

import (
	"OnlineShopBackend/internal/delivery/user"
	"OnlineShopBackend/internal/models"
	"OnlineShopBackend/internal/usecase/interfaces"
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var _ usecase.IUserUsecase = &UserUsecase{}

type UserUsecase struct {
	userStore usecase.UserStore
	logger    *zap.Logger
}

func NewUserUsecase(userStore usecase.UserStore, logger *zap.Logger) usecase.IUserUsecase {
	return &UserUsecase{userStore: userStore, logger: logger}
}

type Credentials struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Profile struct {
	Email     string  `json:"email,omitempty"`
	FirstName string  `json:"firstname,omitempty"`
	LastName  string  `json:"lastname,omitempty"`
	Address   Address `json:"address,omitempty"`
	Rights    Rights  `json:"rights,omitempty"`
}

type Address struct {
	Zipcode string `json:"zipcode,omitempty"`
	Country string `json:"country,omitempty"`
	City    string `json:"city,omitempty"`
	Street  string `json:"street,omitempty"`
}

type Rights struct {
	ID    uuid.UUID `json:"id,omitempty"`
	Name  string    `json:"name,omitempty"`
	Rules []string  `json:"rules,omitempty"`
}

func (usecase *UserUsecase) CreateUser(ctx context.Context, user *user.CreateUserData) (*models.User, error) {
	usecase.logger.Debug("Enter in usecase CreateUser()")

	rights, err := usecase.userStore.GetRightsId(ctx, "Customer")
	if err != nil {
		return &models.User{}, err
	}

	usecaseUser := &models.User{
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		Password:  user.Password,
		Address: models.UserAddress{
			Zipcode: user.Address.Zipcode,
			Country: user.Address.Country,
			City:    user.Address.City,
			Street:  user.Address.Street,
		},
		Rights: models.Rights{
			ID:    rights.ID,
			Name:  rights.Name,
			Rules: rights.Rules,
		},
	}

	id, err := usecase.userStore.Create(ctx, usecaseUser)
	if err != nil {
		return &models.User{}, fmt.Errorf("error on create user: %w", err)
	}
	return id, nil
}

func (usecase *UserUsecase) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user *models.User

	user, err := usecase.userStore.GetUserByEmail(ctx, email)
	if err != nil {
		return &models.User{}, err
	}

	return user, nil
}

func (usecase *UserUsecase) GetRightsId(ctx context.Context, name string) (*models.Rights, error) {
	//var rights models.Rights

	rights, err := usecase.userStore.GetRightsId(ctx, name)
	if err != nil {
		return nil, err
	}

	return &rights, nil
}

func (usecase *UserUsecase) UpdateUserData(ctx context.Context, id uuid.UUID, user *user.CreateUserData) (*models.User, error) {
	userData := &models.User{
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Address: models.UserAddress{
			Zipcode: user.Address.Zipcode,
			Country: user.Address.Country,
			City:    user.Address.City,
			Street:  user.Address.Street,
		},
	}
	userUpdated, err := usecase.userStore.UpdateUserData(ctx, id, userData)
	if err != nil {
		return nil, err
	}
	return userUpdated, nil
}

func (usecase *UserUsecase) UpdateUserRole(ctx context.Context, roleId uuid.UUID, email string) error {
	err := usecase.userStore.UpdateUserRole(ctx, roleId, email)
	if err != nil {
		return err
	}
	return nil
}

func (usecase *UserUsecase) GetRightsList(ctx context.Context) ([]models.Rights, error) {
	usecase.logger.Debug("Enter in usecase GetRightsList() with args: ctx")
	rightsIncomingChan, err := usecase.userStore.GetRightsList(ctx)
	if err != nil {
		return nil, err
	}

	rights := make([]models.Rights, 0, 100)
	for rule := range rightsIncomingChan {
		rights = append(rights, rule)
	}

	return rights, nil

}

func (usecase *UserUsecase) CreateRights(ctx context.Context, rights *models.Rights) (uuid.UUID, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase CreateRights() with args: ctx, rights: %v", rights)

	id, err := usecase.userStore.CreateRights(ctx, rights)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

// GetCartByUserId creates request in db and returns cart or error
func (usecase *UserUsecase) GetCartByUserId(ctx context.Context, userId uuid.UUID) (*models.Cart, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase GetCart() with args: ctx, userId: %v", userId)
	cart, err := usecase.userStore.GetCartByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	return cart, nil
}

// Create create new cart
func (usecase *UserUsecase) CreateCart(ctx context.Context, userId uuid.UUID) (uuid.UUID, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase cart Create() with args: ctx, userId: %v", userId)
	cartId, err := usecase.userStore.CreateCart(ctx, userId)
	if err != nil {
		return cartId, err
	}
	return cartId, nil
}
