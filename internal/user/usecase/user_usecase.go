package usecase

import (
	"context"
	"github.com/AleksK1NG/auth-microservice/internal/models"
	"github.com/AleksK1NG/auth-microservice/internal/user"
	"github.com/AleksK1NG/auth-microservice/pkg/logger"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
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

	findByEmail, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.Wrap(err, "userRepo.FindByEmail")
	}

	findByEmail.SanitizePassword()

	return findByEmail, nil
}

// Find use by uuid
func (u userUseCase) FindById(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserUseCase.FindById")
	defer span.Finish()

	return u.userRepo.FindById(ctx, userID)
}

// Login user with email and password
func (u userUseCase) Login(ctx context.Context, email string, password string) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserUseCase.FindById")
	defer span.Finish()

	foundUser, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if err := foundUser.ComparePasswords(password); err != nil {
		return nil, errors.Wrap(err, "user.ComparePasswords")
	}

	return foundUser, err
}
