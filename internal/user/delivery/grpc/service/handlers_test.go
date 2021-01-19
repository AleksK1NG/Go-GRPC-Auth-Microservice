package service

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/AleksK1NG/auth-microservice/config"
	"github.com/AleksK1NG/auth-microservice/internal/models"
	mockSessUC "github.com/AleksK1NG/auth-microservice/internal/session/mock"
	"github.com/AleksK1NG/auth-microservice/internal/user/mock"
	"github.com/AleksK1NG/auth-microservice/pkg/logger"
	userService "github.com/AleksK1NG/auth-microservice/proto"
)

func TestUsersService_Register(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userUC := mock.NewMockUserUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessionUseCase(ctrl)
	apiLogger := logger.NewAPILogger(nil)
	authServerGRPC := NewAuthServerGRPC(apiLogger, nil, userUC, sessUC)

	reqValue := &userService.RegisterRequest{
		Email:     "email@gmail.com",
		FirstName: "FirstName",
		LastName:  "LastName",
		Password:  "Password",
		Role:      "user",
		Avatar:    "",
	}

	t.Run("Register", func(t *testing.T) {
		t.Parallel()
		userID := uuid.New()
		user := &models.User{
			UserID:    userID,
			Email:     reqValue.Email,
			FirstName: reqValue.FirstName,
			LastName:  reqValue.LastName,
			Password:  reqValue.Password,
			Role:      reqValue.Role,
			Avatar:    nil,
		}

		userUC.EXPECT().Register(gomock.Any(), gomock.Any()).Return(user, nil)

		response, err := authServerGRPC.Register(context.Background(), reqValue)
		require.NoError(t, err)
		require.NotNil(t, response)
		require.Equal(t, reqValue.Email, response.User.Email)
	})
}

func TestUsersService_Login(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userUC := mock.NewMockUserUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessionUseCase(ctrl)
	apiLogger := logger.NewAPILogger(nil)
	cfg := &config.Config{Session: config.Session{
		Expire: 10,
	}}
	authServerGRPC := NewAuthServerGRPC(apiLogger, cfg, userUC, sessUC)

	reqValue := &userService.LoginRequest{
		Email:    "email@gmail.com",
		Password: "Password",
	}

	t.Run("Login", func(t *testing.T) {
		t.Parallel()
		userID := uuid.New()
		session := "session"
		user := &models.User{
			UserID:    userID,
			Email:     "email@gmail.com",
			FirstName: "FirstName",
			LastName:  "LastName",
			Password:  "Password",
			Role:      "user",
			Avatar:    nil,
		}

		userUC.EXPECT().Login(gomock.Any(), reqValue.Email, reqValue.Password).Return(user, nil)
		sessUC.EXPECT().CreateSession(gomock.Any(), &models.Session{
			UserID: user.UserID,
		}, cfg.Session.Expire).Return(session, nil)

		response, err := authServerGRPC.Login(context.Background(), reqValue)
		require.NoError(t, err)
		require.NotNil(t, response)
		require.Equal(t, reqValue.Email, response.User.Email)
	})
}

func TestUsersService_FindByEmail(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userUC := mock.NewMockUserUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessionUseCase(ctrl)
	apiLogger := logger.NewAPILogger(nil)
	cfg := &config.Config{Session: config.Session{
		Expire: 10,
	}}
	authServerGRPC := NewAuthServerGRPC(apiLogger, cfg, userUC, sessUC)

	reqValue := &userService.FindByEmailRequest{
		Email: "email@gmail.com",
	}

	t.Run("FindByEmail", func(t *testing.T) {
		t.Parallel()
		userID := uuid.New()
		user := &models.User{
			UserID:    userID,
			Email:     "email@gmail.com",
			FirstName: "FirstName",
			LastName:  "LastName",
			Password:  "Password",
			Role:      "user",
			Avatar:    nil,
		}

		userUC.EXPECT().FindByEmail(gomock.Any(), reqValue.Email).Return(user, nil)

		response, err := authServerGRPC.FindByEmail(context.Background(), reqValue)
		require.NoError(t, err)
		require.NotNil(t, response)
		require.Equal(t, reqValue.Email, response.User.Email)
	})
}
