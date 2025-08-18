package sqlRepo

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/danilkompaniets/auth-service/internal/infrastructure/repository"
	"github.com/danilkompaniets/auth-service/pkg/model"
	"github.com/stretchr/testify/assert"
)

func setupMockDB(t *testing.T) (repository.AuthRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	repo := NewAuthRepository(db)

	return repo, mock, func() { db.Close() }
}

func TestCreateUser(t *testing.T) {
	repo, mock, closeDB := setupMockDB(t)
	defer closeDB()

	user := model.User{
		Email:     "test@example.com",
		Password:  "hashedPassword",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO users (email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id`)).
		WithArgs(user.Email, user.Password, user.CreatedAt, user.UpdatedAt).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	id, err := repo.CreateUser(context.Background(), user)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDeleteRefreshToken(t *testing.T) {
	repo, mock, closeDB := setupMockDB(t)
	defer closeDB()

	userID := int64(1)
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM refresh_tokens WHERE user_id = $1`)).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.DeleteRefreshToken(context.Background(), userID)
	assert.NoError(t, err)
}

func TestGetRefreshToken(t *testing.T) {
	repo, mock, closeDB := setupMockDB(t)
	defer closeDB()

	userID := int64(1)
	expectedToken := "refresh-token-123"

	// Тест успешного запроса
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT token FROM refresh_tokens WHERE user_id = $1`)).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"token"}).AddRow(expectedToken))

	token, err := repo.GetRefreshToken(context.Background(), userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedToken, token)

	// Тест ошибки: нет записи
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT token FROM refresh_tokens WHERE user_id = $1`)).
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	_, err = repo.GetRefreshToken(context.Background(), userID)
	assert.Equal(t, repository.ErrUserNotFound, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetUserByEmail(t *testing.T) {
	repo, mock, closeDB := setupMockDB(t)
	defer closeDB()

	email := "test@example.com"
	now := time.Now()

	// Успешный кейс
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, email, password, created_at, updated_at FROM users WHERE email = $1;`)).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "created_at", "updated_at"}).
			AddRow(1, email, "hashedPassword", now, now))

	user, err := repo.GetUserByEmail(context.Background(), email)
	assert.NoError(t, err)
	assert.Equal(t, email, user.Email)

	// Ошибка: не найден
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, email, password, created_at, updated_at FROM users WHERE email = $1;`)).
		WithArgs(email).
		WillReturnError(sql.ErrNoRows)

	user, err = repo.GetUserByEmail(context.Background(), email)
	assert.Nil(t, user)
	assert.Equal(t, repository.ErrUserNotFound, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestSaveRefreshToken(t *testing.T) {
	repo, mock, closeDB := setupMockDB(t)
	defer closeDB()

	userID := int64(1)
	token := "refresh-token-123"

	mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO refresh_tokens (user_id, token)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE
		SET token = EXCLUDED.token, created_at = NOW()`)).
		WithArgs(userID, token).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.SaveRefreshToken(context.Background(), userID, token)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
