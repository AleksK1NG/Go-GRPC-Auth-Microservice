package user

import (
	"context"
	"github.com/AleksK1NG/auth-microservice/internal/models"
	"github.com/google/uuid"
)

// User pg repository
type UserRepository interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindById(ctx context.Context, userID uuid.UUID) (*models.User, error)
}
