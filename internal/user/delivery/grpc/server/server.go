package server

import (
	"github.com/AleksK1NG/auth-microservice/config"
	"github.com/AleksK1NG/auth-microservice/internal/user"
	"github.com/AleksK1NG/auth-microservice/pkg/logger"
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
