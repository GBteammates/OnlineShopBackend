package usecase

import (
	"OnlineShopBackend/internal/models"
	mocks "OnlineShopBackend/internal/usecase/repo_mocks"
	"context"
	"fmt"
	"strings"
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

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepo := mocks.NewMockUserStore(ctrl)
	logger := zap.L()
	ctx := context.Background()
	id := uuid.New()

	testCases := []struct {
		name   string
		user   *models.User
		id     uuid.UUID
		expect func(*mocks.MockUserStore)
	}{
		{
			name: "error on create user",
			user: &models.User{},
			expect: func(mos *mocks.MockUserStore) {
				mos.EXPECT().CreateUser(ctx, &models.User{}).Return(uuid.Nil, fmt.Errorf("error"))
			},
		},
		{
			name: "success create user",
			user: &models.User{},
			expect: func(mos *mocks.MockUserStore) {
				mos.EXPECT().CreateUser(ctx, &models.User{}).Return(id, nil)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			usecase := NewUserUsecase(userRepo, logger)
			if tc.expect != nil {
				tc.expect(userRepo)
			}
			res, err := usecase.CreateUser(ctx, tc.user)
			if strings.Contains(tc.name, "error") {
				require.Error(t, err)
				require.Equal(t, uuid.Nil, res)
			} else {
				require.NoError(t, err)
				require.Equal(t, id, res)
			}
		})
	}
}

func TestUpdateUserData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepo := mocks.NewMockUserStore(ctrl)
	logger := zap.L()
	ctx := context.Background()

	testCases := []struct {
		name   string
		user   *models.User
		expect func(*mocks.MockUserStore)
	}{
		{
			name: "error on update user data",
			user: &models.User{},
			expect: func(mos *mocks.MockUserStore) {
				mos.EXPECT().UpdateUserData(ctx, &models.User{}).Return(nil, fmt.Errorf("error"))
			},
		},
		{
			name: "success update user data",
			user: &models.User{},
			expect: func(mos *mocks.MockUserStore) {
				mos.EXPECT().UpdateUserData(ctx, &models.User{}).Return(&models.User{}, nil)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			usecase := NewUserUsecase(userRepo, logger)
			if tc.expect != nil {
				tc.expect(userRepo)
			}
			res, err := usecase.UpdateUserData(ctx, tc.user)
			if strings.Contains(tc.name, "error") {
				require.Error(t, err)
				require.Nil(t, res)
			} else {
				require.NoError(t, err)
				require.NotNil(t, res)
			}
		})
	}
}

func TestUpdateUserRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepo := mocks.NewMockUserStore(ctrl)
	logger := zap.L()
	ctx := context.Background()

	roleId := uuid.New()
	email := "testEmail"

	testCases := []struct {
		name   string
		expect func(*mocks.MockUserStore)
	}{
		{
			name: "error on update user role",
			expect: func(mos *mocks.MockUserStore) {
				mos.EXPECT().UpdateUserRole(ctx, roleId, email).Return(fmt.Errorf("error"))
			},
		},
		{
			name: "success update user role",
			expect: func(mos *mocks.MockUserStore) {
				mos.EXPECT().UpdateUserRole(ctx, roleId, email).Return(nil)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			usecase := NewUserUsecase(userRepo, logger)
			if tc.expect != nil {
				tc.expect(userRepo)
			}
			err := usecase.UpdateUserRole(ctx, roleId, email)
			if strings.Contains(tc.name, "error") {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetRightsList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepo := mocks.NewMockUserStore(ctrl)
	logger := zap.L()
	ctx := context.Background()
	ch := make(chan models.Rights, 1)
	testRights := models.Rights{}
	ch <- testRights
	close(ch)
	expect := make([]models.Rights, 0, 100)
	expect = append(expect, testRights)

	testCases := []struct {
		name   string
		expect func(*mocks.MockUserStore)
	}{
		{
			name: "error on get rights list",
			expect: func(mos *mocks.MockUserStore) {
				mos.EXPECT().GetRightsList(ctx).Return(nil, fmt.Errorf("error"))
			},
		},
		{
			name: "success get rights list",
			expect: func(mos *mocks.MockUserStore) {
				mos.EXPECT().GetRightsList(ctx).Return(ch, nil)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			usecase := NewUserUsecase(userRepo, logger)
			if tc.expect != nil {
				tc.expect(userRepo)
			}
			res, err := usecase.GetRightsList(ctx)
			if strings.Contains(tc.name, "error") {
				require.Error(t, err)
				require.Nil(t, res)
			} else {
				require.NoError(t, err)
				require.Equal(t, expect, res)
			}
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepo := mocks.NewMockUserStore(ctrl)
	logger := zap.L()
	ctx := context.Background()

	email := "testEmail"
	user := &models.User{}

	testCases := []struct {
		name   string
		expect func(*mocks.MockUserStore)
	}{
		{
			name: "error on get user by email",
			expect: func(mos *mocks.MockUserStore) {
				mos.EXPECT().GetUserByEmail(ctx, email).Return(nil, fmt.Errorf("error"))
			},
		},
		{
			name: "success get user by email",
			expect: func(mos *mocks.MockUserStore) {
				mos.EXPECT().GetUserByEmail(ctx, email).Return(user, nil)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			usecase := NewUserUsecase(userRepo, logger)
			if tc.expect != nil {
				tc.expect(userRepo)
			}
			res, err := usecase.GetUserByEmail(ctx, email)
			if strings.Contains(tc.name, "error") {
				require.Error(t, err)
				require.Nil(t, res)
			} else {
				require.NoError(t, err)
				require.Equal(t, user, res)
			}
		})
	}
}

func TestGetRightsId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepo := mocks.NewMockUserStore(ctrl)
	logger := zap.L()
	ctx := context.Background()

	name := "testName"
	rights := models.Rights{}

	testCases := []struct {
		name   string
		expect func(*mocks.MockUserStore)
	}{
		{
			name: "error on get rights id",
			expect: func(mos *mocks.MockUserStore) {
				mos.EXPECT().GetRightsId(ctx, name).Return(rights, fmt.Errorf("error"))
			},
		},
		{
			name: "success get rights id",
			expect: func(mos *mocks.MockUserStore) {
				mos.EXPECT().GetRightsId(ctx, name).Return(rights, nil)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			usecase := NewUserUsecase(userRepo, logger)
			if tc.expect != nil {
				tc.expect(userRepo)
			}
			res, err := usecase.GetRightsId(ctx, name)
			if strings.Contains(tc.name, "error") {
				require.Error(t, err)
				require.Nil(t, res)
			} else {
				require.NoError(t, err)
				require.Equal(t, &rights, res)
			}
		})
	}
}
