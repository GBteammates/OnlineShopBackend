package user_usecase

import (
	"OnlineShopBackend/internal/models"
	mocks "OnlineShopBackend/internal/usecase/repo_mocks"
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var (
	testRightsNoId = &models.Rights{
		Name: "Test",
	}
	testRightsId = uuid.New()
)

func TestCreateRights(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := zap.L()
	userRepo := mocks.NewMockUserStore(ctrl)
	usecase := NewUserUsecase(userRepo, logger)
	ctx := context.Background()

	userRepo.EXPECT().CreateRights(ctx, testRightsNoId).Return(uuid.Nil, fmt.Errorf("error"))
	res, err := usecase.CreateRights(ctx, testRightsNoId)
	require.Error(t, err)
	require.Equal(t, res, uuid.Nil)

	userRepo.EXPECT().CreateRights(ctx, testRightsNoId).Return(testRightsId, nil)
	res, err = usecase.CreateRights(ctx, testRightsNoId)
	require.NoError(t, err)
	require.Equal(t, res, testRightsId)
}
