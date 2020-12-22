package server

import (
	"context"
	"github.com/AleksK1NG/auth-microservice/config"
	"github.com/AleksK1NG/auth-microservice/internal/models"
	"github.com/AleksK1NG/auth-microservice/internal/user"
	"github.com/AleksK1NG/auth-microservice/pkg/logger"
	"github.com/AleksK1NG/auth-microservice/proto"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type usersServer struct {
	logger logger.Logger
	cfg    *config.Config
	userUC user.UserUseCase
}

// Auth server constructor
func NewAuthServerGRPC(logger logger.Logger, cfg *config.Config, userUC user.UserUseCase) *usersServer {
	return &usersServer{logger: logger, cfg: cfg, userUC: userUC}
}

// Register new user
func (u *usersServer) Register(ctx context.Context, r *userService.RegisterRequest) (*userService.RegisterResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "user.Register")
	defer span.Finish()

	u.logger.Infof("Get request %s\n", r.String())

	createdUser, err := u.userUC.Register(ctx, &models.User{
		Email:     r.GetEmail(),
		FirstName: r.GetFirstName(),
		LastName:  r.GetLastName(),
		Password:  r.GetPassword(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "userUC.Register: %#v", err)
	}

	return &userService.RegisterResponse{
		Email:     createdUser.Email,
		FirstName: createdUser.FirstName,
		LastName:  createdUser.LastName,
		UUID:      uuid.New().String(),
		CreatedAt: timestamppb.New(time.Now()),
	}, nil
}
