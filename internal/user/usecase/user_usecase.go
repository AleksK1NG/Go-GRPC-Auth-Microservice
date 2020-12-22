package usecase

import (
	"context"
	"github.com/AleksK1NG/auth-microservice/internal/models"
	"github.com/AleksK1NG/auth-microservice/internal/user"
	"github.com/AleksK1NG/auth-microservice/pkg/logger"
	"github.com/opentracing/opentracing-go"
)

// Auth UseCase
type userUseCase struct {
	logger   logger.Logger
	userRepo user.UserRepository
}

// New Auth UseCase
func NewUserUseCase(logger logger.Logger, userRepo user.UserRepository) *userUseCase {
	return &userUseCase{logger: logger, userRepo: userRepo}
}

// Register
func (u userUseCase) Register(ctx context.Context, user *models.User) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserUseCase.Register")
	defer span.Finish()

	return u.userRepo.Register(ctx, user)
}
