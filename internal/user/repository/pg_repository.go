package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/AleksK1NG/auth-microservice/internal/models"
)

// User repository
type UserRepository struct {
	db *sqlx.DB
}

// User repository constructor
func NewUserPGRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create new user
func (r *UserRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.Create")
	defer span.Finish()

	createdUser := &models.User{}
	if err := r.db.QueryRowxContext(
		ctx,
		createUserQuery,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
		user.Role,
		user.Avatar,
	).StructScan(createdUser); err != nil {
		return nil, errors.Wrap(err, "Create.QueryRowxContext")
	}

	return createdUser, nil
}

// Find by user email address
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.FindByEmail")
	defer span.Finish()

	user := &models.User{}
	if err := r.db.GetContext(ctx, user, findByEmailQuery, email); err != nil {
		return nil, errors.Wrap(err, "FindByEmail.GetContext")
	}

	return user, nil
}

// Find user by uuid
func (r *UserRepository) FindById(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.FindById")
	defer span.Finish()

	user := &models.User{}
	if err := r.db.GetContext(ctx, user, findByIDQuery, userID); err != nil {
		return nil, errors.Wrap(err, "FindById.GetContext")
	}

	return user, nil
}
