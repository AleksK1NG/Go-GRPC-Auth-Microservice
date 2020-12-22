package repository

import (
	"context"
	"github.com/AleksK1NG/auth-microservice/internal/models"
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
func (u *UserRepository) Register(ctx context.Context, user *models.User) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.Register")
	defer span.Finish()

	query := `INSERT INTO users (first_name, last_name, email, password, role, avatar) VALUES ($1, $2, $3, $4, $5 ,$6) RETURNING *`

	createdUser := &models.User{}
	if err := u.db.QueryRowxContext(ctx, query, user.FirstName, user.LastName, user.Email, user.Password, user.Role, user.Avatar).StructScan(createdUser); err != nil {
		return nil, errors.Wrap(err, "Register.QueryRowxContext")
	}

	return createdUser, nil
}
