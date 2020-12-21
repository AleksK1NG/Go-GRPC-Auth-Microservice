package user

import (
	"context"
	"github.com/AleksK1NG/auth-microservice/internal/models"
)

//  User UseCase interface
type UserUseCase interface {
	Register(ctx context.Context, user *models.User) (*models.User, error)
}
