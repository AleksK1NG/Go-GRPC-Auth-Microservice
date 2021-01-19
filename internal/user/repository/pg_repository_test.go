package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/AleksK1NG/auth-microservice/internal/models"
)

func TestUserRepository_Create(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	userPGRepository := NewUserPGRepository(sqlxDB)

	columns := []string{"user_id", "first_name", "last_name", "email", "password", "avatar", "role", "created_at", "updated_at"}
	userUUID := uuid.New()
	mockUser := &models.User{
		UserID:    userUUID,
		Email:     "email@gmail.com",
		FirstName: "FirstName",
		LastName:  "LastName",
		Role:      "admin",
		Avatar:    nil,
		Password:  "123456",
	}

	rows := sqlmock.NewRows(columns).AddRow(
		userUUID,
		mockUser.FirstName,
		mockUser.LastName,
		mockUser.Email,
		mockUser.Password,
		mockUser.Avatar,
		mockUser.Role,
		time.Now(),
		time.Now(),
	)

	mock.ExpectQuery(createUserQuery).WithArgs(
		mockUser.FirstName,
		mockUser.LastName,
		mockUser.Email,
		mockUser.Password,
		mockUser.Role,
		mockUser.Avatar,
	).WillReturnRows(rows)

	createdUser, err := userPGRepository.Create(context.Background(), mockUser)
	require.NoError(t, err)
	require.NotNil(t, createdUser)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	userPGRepository := NewUserPGRepository(sqlxDB)

	columns := []string{"user_id", "first_name", "last_name", "email", "password", "avatar", "role", "created_at", "updated_at"}
	userUUID := uuid.New()
	mockUser := &models.User{
		UserID:    userUUID,
		Email:     "email@gmail.com",
		FirstName: "FirstName",
		LastName:  "LastName",
		Role:      "admin",
		Avatar:    nil,
		Password:  "123456",
	}

	rows := sqlmock.NewRows(columns).AddRow(
		userUUID,
		mockUser.FirstName,
		mockUser.LastName,
		mockUser.Email,
		mockUser.Password,
		mockUser.Avatar,
		mockUser.Role,
		time.Now(),
		time.Now(),
	)

	mock.ExpectQuery(findByEmailQuery).WithArgs(mockUser.Email).WillReturnRows(rows)

	foundUser, err := userPGRepository.FindByEmail(context.Background(), mockUser.Email)
	require.NoError(t, err)
	require.NotNil(t, foundUser)
	require.Equal(t, foundUser.Email, mockUser.Email)
}

func TestUserRepository_FindById(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	userPGRepository := NewUserPGRepository(sqlxDB)

	columns := []string{"user_id", "first_name", "last_name", "email", "password", "avatar", "role", "created_at", "updated_at"}
	userUUID := uuid.New()
	mockUser := &models.User{
		UserID:    userUUID,
		Email:     "email@gmail.com",
		FirstName: "FirstName",
		LastName:  "LastName",
		Role:      "admin",
		Avatar:    nil,
		Password:  "123456",
	}

	rows := sqlmock.NewRows(columns).AddRow(
		userUUID,
		mockUser.FirstName,
		mockUser.LastName,
		mockUser.Email,
		mockUser.Password,
		mockUser.Avatar,
		mockUser.Role,
		time.Now(),
		time.Now(),
	)

	mock.ExpectQuery(findByIDQuery).WithArgs(mockUser.UserID).WillReturnRows(rows)

	foundUser, err := userPGRepository.FindById(context.Background(), mockUser.UserID)
	require.NoError(t, err)
	require.NotNil(t, foundUser)
	require.Equal(t, foundUser.UserID, mockUser.UserID)
}
