package repository

import (
	"context"
	"github.com/AleksK1NG/auth-microservice/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

// Auth repository
type UserRepository struct {
	db *sqlx.DB
}

// Auth repository constructor
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create new user
func (u *UserRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.Create")
	defer span.Finish()

	query := `INSERT INTO users (first_name, last_name, email, password, role, avatar) VALUES ($1, $2, $3, $4, $5 ,$6) RETURNING *`

	createdUser := &models.User{}
	if err := u.db.QueryRowxContext(ctx, query, user.FirstName, user.LastName, user.Email, user.Password, user.Role, user.Avatar).StructScan(createdUser); err != nil {
		return nil, errors.Wrap(err, "Create.QueryRowxContext")
	}

	return createdUser, nil
}

// Find by user email address
func (u *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.Create")
	defer span.Finish()

	query := `SELECT user_id, email, first_name, last_name, role, avatar FROM users WHERE email = $1`

	user := &models.User{}
	if err := u.db.GetContext(ctx, user, query, email); err != nil {
		return nil, errors.Wrap(err, "FindByEmail.GetContext")
	}

	return user, nil
}

// Find user by uuid
func (u *UserRepository) FindById(ctx context.Context, userID *uuid.UUID) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.Create")
	defer span.Finish()

	query := `SELECT user_id, email, first_name, last_name, role, avatar FROM users WHERE user_id = $1`

	user := &models.User{}
	if err := u.db.GetContext(ctx, user, query, userID); err != nil {
		return nil, errors.Wrap(err, "FindById.GetContext")
	}

	return user, nil
}
