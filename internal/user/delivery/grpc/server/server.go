package server

import (
	"context"
	"github.com/AleksK1NG/auth-microservice/config"
	"github.com/AleksK1NG/auth-microservice/pkg/logger"
	"github.com/AleksK1NG/auth-microservice/proto"
	"github.com/opentracing/opentracing-go"
)

type usersServer struct {
	logger logger.Logger
	cfg    *config.Config
}

// Register new user
func (u *usersServer) Register(ctx context.Context, r *userService.RegisterRequest) (*userService.RegisterResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "user.Register")
	defer span.Finish()

	u.logger.Infof("Get request %s\n", r.String())

	return &userService.RegisterResponse{
		Email:     r.GetEmail(),
		FirstName: r.GetFirstName(),
		LastName:  r.GetLastName(),
		Uid:       1,
	}, nil
}

// Auth server constructor
func NewAuthServerGRPC(logger logger.Logger, cfg *config.Config) *usersServer {
	return &usersServer{logger: logger, cfg: cfg}
}
