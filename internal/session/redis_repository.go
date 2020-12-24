//go:generate mockgen -source pg_repository.go -destination mock/pg_repository_mock.go -package mock
package session

import (
	"context"
	"github.com/AleksK1NG/auth-microservice/internal/models"
)

// Session repository
type SessRepository interface {
	CreateSession(ctx context.Context, session *models.Session, expire int) (string, error)
	GetSessionByID(ctx context.Context, sessionID string) (*models.Session, error)
	DeleteByID(ctx context.Context, sessionID string) error
}
