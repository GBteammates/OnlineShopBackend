package usecase

import (
	"OnlineShopBackend/internal/models"
	usecase "OnlineShopBackend/internal/usecase/interfaces"
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var _ usecase.IUserUsecase = (*userUsecase)(nil)

type userUsecase struct {
	userStore usecase.UserStore
	logger    *zap.Logger
}

func NewUserUsecase(userStore usecase.UserStore, logger *zap.Logger) *userUsecase {
	return &userUsecase{userStore: userStore, logger: logger}
}

func (usecase *userUsecase) CreateUser(ctx context.Context, user *models.User) (uuid.UUID, error) {
	usecase.logger.Debug("Enter in usecase CreateUser()")

	id, err := usecase.userStore.CreateUser(ctx, user)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error on create user: %w", err)
	}
	return id, nil
}

func (usecase *userUsecase) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase GetUserByEmail() with args: ctx, email: %s", email)

	user, err := usecase.userStore.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (usecase *userUsecase) GetRightsId(ctx context.Context, name string) (*models.Rights, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase GetRightsId() with args: ctx, name: %s", name)

	rights, err := usecase.userStore.GetRightsId(ctx, name)
	if err != nil {
		return nil, err
	}

	return &rights, nil
}

func (usecase *userUsecase) UpdateUserData(ctx context.Context, user *models.User) (*models.User, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase UpdateUserData() with args: ctx, user: %v", user)

	userUpdated, err := usecase.userStore.UpdateUserData(ctx, user)
	if err != nil {
		return nil, err
	}
	return userUpdated, nil
}

func (usecase *userUsecase) UpdateUserRole(ctx context.Context, roleId uuid.UUID, email string) error {
	err := usecase.userStore.UpdateUserRole(ctx, roleId, email)
	if err != nil {
		return err
	}
	return nil
}

func (usecase *userUsecase) GetRightsList(ctx context.Context) ([]models.Rights, error) {
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

func (usecase *userUsecase) CreateRights(ctx context.Context, rights *models.Rights) (uuid.UUID, error) {
	usecase.logger.Sugar().Debugf("Enter in usecase CreateRights() with args: ctx, rights: %v", rights)

	id, err := usecase.userStore.CreateRights(ctx, rights)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
