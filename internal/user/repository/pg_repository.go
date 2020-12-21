package repository

import (
	"context"
	"github.com/AleksK1NG/auth-microservice/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
)

// Auth repository
type UserRepository struct {
	db *sqlx.DB
}

func (u *UserRepository) Register(ctx context.Context, user *models.User) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.Register")
	defer span.Finish()

	return user, nil
}

// Auth repository constructor
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}
