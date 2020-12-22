package usecase

import (
	"context"
	"github.com/AleksK1NG/auth-microservice/internal/models"
	"github.com/AleksK1NG/auth-microservice/internal/user"
	"github.com/AleksK1NG/auth-microservice/pkg/logger"
	"github.com/google/uuid"
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

// Register new user
func (u userUseCase) Register(ctx context.Context, user *models.User) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserUseCase.Create")
	defer span.Finish()

	return u.userRepo.Create(ctx, user)
}

// Find use by email address
func (u userUseCase) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserUseCase.FindByEmail")
	defer span.Finish()

	return u.userRepo.FindByEmail(ctx, email)
}

// Find use by uuid
func (u userUseCase) FindById(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserUseCase.FindById")
	defer span.Finish()

	return u.userRepo.FindById(ctx, userID)
}
