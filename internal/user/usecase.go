//go:generate mockgen -source usecase.go -destination mock/usecase.go -package mock
package user

import (
	"context"

	"github.com/google/uuid"

	"github.com/AleksK1NG/auth-microservice/internal/models"
)

//  User UseCase interface
type UserUseCase interface {
	Register(ctx context.Context, user *models.User) (*models.User, error)
	Login(ctx context.Context, email string, password string) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindById(ctx context.Context, userID uuid.UUID) (*models.User, error)
}
