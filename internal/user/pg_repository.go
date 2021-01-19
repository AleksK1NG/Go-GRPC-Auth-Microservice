//go:generate mockgen -source pg_repository.go -destination mock/pg_repository.go -package mock
package user

import (
	"context"

	"github.com/google/uuid"

	"github.com/AleksK1NG/auth-microservice/internal/models"
)

// User pg repository
type UserPGRepository interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindById(ctx context.Context, userID uuid.UUID) (*models.User, error)
}
