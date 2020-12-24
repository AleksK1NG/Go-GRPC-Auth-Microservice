package service

import (
	"github.com/AleksK1NG/auth-microservice/config"
	"github.com/AleksK1NG/auth-microservice/internal/session"
	"github.com/AleksK1NG/auth-microservice/internal/user"
	"github.com/AleksK1NG/auth-microservice/pkg/logger"
)

type usersService struct {
	logger logger.Logger
	cfg    *config.Config
	userUC user.UserUseCase
	sessUC session.SessionUseCase
}

// Auth service constructor
func NewAuthServerGRPC(logger logger.Logger, cfg *config.Config, userUC user.UserUseCase, sessUC session.SessionUseCase) *usersService {
	return &usersService{logger: logger, cfg: cfg, userUC: userUC, sessUC: sessUC}
}
