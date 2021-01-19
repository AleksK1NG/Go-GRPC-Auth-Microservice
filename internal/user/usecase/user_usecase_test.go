package usecase

import (
	"context"
	"database/sql"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/AleksK1NG/auth-microservice/internal/models"
	"github.com/AleksK1NG/auth-microservice/internal/user/mock"
	"github.com/AleksK1NG/auth-microservice/pkg/logger"
)

func TestUserUseCase_Register(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userPGRepository := mock.NewMockUserPGRepository(ctrl)
	userRedisRepository := mock.NewMockUserRedisRepository(ctrl)
	apiLogger := logger.NewAPILogger(nil)
	userUC := NewUserUseCase(apiLogger, userPGRepository, userRedisRepository)

	userID := uuid.New()
	mockUser := &models.User{
		Email:     "email@gmail.com",
		FirstName: "FirstName",
		LastName:  "LastName",
		Role:      "admin",
		Avatar:    nil,
		Password:  "123456",
	}

	ctx := context.Background()

	userPGRepository.EXPECT().FindByEmail(gomock.Any(), mockUser.Email).Return(nil, sql.ErrNoRows)

	userPGRepository.EXPECT().Create(gomock.Any(), mockUser).Return(&models.User{
		UserID:    userID,
		Email:     "email@gmail.com",
		FirstName: "FirstName",
		LastName:  "LastName",
		Role:      "admin",
		Avatar:    nil,
		Password:  "123456",
	}, nil)

	createdUser, err := userUC.Register(ctx, mockUser)
	require.NoError(t, err)
	require.NotNil(t, createdUser)
	require.Equal(t, createdUser.UserID, userID)
}

func TestUserUseCase_FindByEmail(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userPGRepository := mock.NewMockUserPGRepository(ctrl)
	userRedisRepository := mock.NewMockUserRedisRepository(ctrl)
	apiLogger := logger.NewAPILogger(nil)
	userUC := NewUserUseCase(apiLogger, userPGRepository, userRedisRepository)

	userID := uuid.New()
	mockUser := &models.User{
		UserID:    userID,
		Email:     "email@gmail.com",
		FirstName: "FirstName",
		LastName:  "LastName",
		Role:      "admin",
		Avatar:    nil,
		Password:  "123456",
	}

	ctx := context.Background()

	userPGRepository.EXPECT().FindByEmail(gomock.Any(), mockUser.Email).Return(mockUser, nil)

	user, err := userUC.FindByEmail(ctx, mockUser.Email)
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, user.Email, mockUser.Email)
}

func TestUserUseCase_FindById(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userPGRepository := mock.NewMockUserPGRepository(ctrl)
	userRedisRepository := mock.NewMockUserRedisRepository(ctrl)
	apiLogger := logger.NewAPILogger(nil)
	userUC := NewUserUseCase(apiLogger, userPGRepository, userRedisRepository)

	userID := uuid.New()
	mockUser := &models.User{
		UserID:    userID,
		Email:     "email@gmail.com",
		FirstName: "FirstName",
		LastName:  "LastName",
		Role:      "admin",
		Avatar:    nil,
		Password:  "123456",
	}

	ctx := context.Background()

	userRedisRepository.EXPECT().GetByIDCtx(gomock.Any(), mockUser.UserID.String()).Return(nil, redis.Nil)
	userPGRepository.EXPECT().FindById(gomock.Any(), mockUser.UserID).Return(mockUser, nil)
	userRedisRepository.EXPECT().SetUserCtx(gomock.Any(), mockUser.UserID.String(), 3600, mockUser).Return(nil)

	user, err := userUC.FindById(ctx, mockUser.UserID)
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, user.UserID, mockUser.UserID)
}
