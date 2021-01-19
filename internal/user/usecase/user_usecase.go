package usecase

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/AleksK1NG/auth-microservice/internal/models"
	"github.com/AleksK1NG/auth-microservice/internal/user"
	"github.com/AleksK1NG/auth-microservice/pkg/grpc_errors"
	"github.com/AleksK1NG/auth-microservice/pkg/logger"
)

const (
	userByIdCacheDuration = 3600
)

// User UseCase
type userUseCase struct {
	logger     logger.Logger
	userPgRepo user.UserPGRepository
	redisRepo  user.UserRedisRepository
}

// New User UseCase
func NewUserUseCase(logger logger.Logger, userRepo user.UserPGRepository, redisRepo user.UserRedisRepository) *userUseCase {
	return &userUseCase{logger: logger, userPgRepo: userRepo, redisRepo: redisRepo}
}

// Register new user
func (u *userUseCase) Register(ctx context.Context, user *models.User) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserUseCase.Register")
	defer span.Finish()

	existsUser, err := u.userPgRepo.FindByEmail(ctx, user.Email)
	if existsUser != nil || err == nil {
		return nil, grpc_errors.ErrEmailExists
	}

	return u.userPgRepo.Create(ctx, user)
}

// Find use by email address
func (u *userUseCase) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserUseCase.FindByEmail")
	defer span.Finish()

	findByEmail, err := u.userPgRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.Wrap(err, "userPgRepo.FindByEmail")
	}

	findByEmail.SanitizePassword()

	return findByEmail, nil
}

// Find use by uuid
func (u *userUseCase) FindById(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserUseCase.FindById")
	defer span.Finish()

	cachedUser, err := u.redisRepo.GetByIDCtx(ctx, userID.String())
	if err != nil && !errors.Is(err, redis.Nil) {
		u.logger.Errorf("redisRepo.GetByIDCtx", err)
	}
	if cachedUser != nil {
		return cachedUser, nil
	}

	foundUser, err := u.userPgRepo.FindById(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "userPgRepo.FindById")
	}

	if err := u.redisRepo.SetUserCtx(ctx, foundUser.UserID.String(), userByIdCacheDuration, foundUser); err != nil {
		u.logger.Errorf("redisRepo.SetUserCtx", err)
	}

	return foundUser, nil
}

// Login user with email and password
func (u *userUseCase) Login(ctx context.Context, email string, password string) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserUseCase.Login")
	defer span.Finish()

	foundUser, err := u.userPgRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.Wrap(err, "userPgRepo.FindByEmail")
	}

	if err := foundUser.ComparePasswords(password); err != nil {
		return nil, errors.Wrap(err, "user.ComparePasswords")
	}

	return foundUser, err
}
