//go:generate mockgen -source usecase.go -destination mock/usecase.go -package mock
package session

import (
	"context"

	"github.com/AleksK1NG/auth-microservice/internal/models"
)

// Session UseCase
type SessionUseCase interface {
	CreateSession(ctx context.Context, session *models.Session, expire int) (string, error)
	GetSessionByID(ctx context.Context, sessionID string) (*models.Session, error)
	DeleteByID(ctx context.Context, sessionID string) error
}
