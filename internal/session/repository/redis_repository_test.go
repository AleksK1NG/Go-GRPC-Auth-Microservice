package repository

import (
	"context"
	"log"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/AleksK1NG/auth-microservice/internal/models"
	"github.com/AleksK1NG/auth-microservice/internal/session"
)

func SetupRedis() session.SessRepository {
	mr, err := miniredis.Run()
	if err != nil {
		log.Fatal(err)
	}
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	sessRepository := NewSessionRepository(client, nil)
	return sessRepository
}

func TestCreateSession(t *testing.T) {
	t.Parallel()

	sessRepository := SetupRedis()

	t.Run("CreateSession", func(t *testing.T) {
		sessUUID := uuid.New()
		sess := &models.Session{
			SessionID: sessUUID.String(),
			UserID:    sessUUID,
		}
		s, err := sessRepository.CreateSession(context.Background(), sess, 10)
		require.NoError(t, err)
		require.NotEqual(t, s, "")
	})
}

func TestGetSessionByID(t *testing.T) {
	t.Parallel()

	sessRepository := SetupRedis()

	t.Run("GetSessionByID", func(t *testing.T) {
		sessUUID := uuid.New()
		sess := &models.Session{
			SessionID: sessUUID.String(),
			UserID:    sessUUID,
		}
		createdSess, err := sessRepository.CreateSession(context.Background(), sess, 10)
		require.NoError(t, err)
		require.NotEqual(t, createdSess, "")

		s, err := sessRepository.GetSessionByID(context.Background(), createdSess)
		require.NoError(t, err)
		require.NotEqual(t, s, "")
	})
}

func TestDeleteSession(t *testing.T) {
	t.Parallel()

	sessRepository := SetupRedis()

	t.Run("DeleteByID", func(t *testing.T) {
		sessUUID := uuid.New()
		err := sessRepository.DeleteByID(context.Background(), sessUUID.String())
		require.NoError(t, err)
	})
}
