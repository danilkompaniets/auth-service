package sqlRepo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/danilkompaniets/auth-service/internal/infrastructure/repository"
	"github.com/danilkompaniets/auth-service/pkg/model"
)

type Repository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(ctx context.Context, user model.User) (int64, error) {
	query := `
		INSERT INTO users (email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var id int64
	err := r.db.QueryRowContext(ctx, query, user.Email, user.Password, user.CreatedAt, user.UpdatedAt).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repository) DeleteRefreshToken(ctx context.Context, userID int64) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return repository.ErrUserNotFound
	}
	return err
}

func (r *Repository) GetRefreshToken(ctx context.Context, userID int64) (string, error) {
	query := `SELECT token FROM refresh_tokens WHERE user_id = $1`

	row := r.db.QueryRowContext(ctx, query, userID)

	var refreshToken string
	err := row.Scan(&refreshToken)

	if errors.Is(err, sql.ErrNoRows) {
		return "", repository.ErrUserNotFound
	}
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT id, email, password, created_at, updated_at FROM users WHERE email = $1;`

	row := r.db.QueryRowContext(ctx, query, email)

	var user model.User
	err := row.Scan(&user.Id, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, repository.ErrUserNotFound
	}

	return &user, err
}

func (r *Repository) SaveRefreshToken(ctx context.Context, userID int64, token string) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE
		SET token = EXCLUDED.token, created_at = NOW()
	`
	_, err := r.db.ExecContext(ctx, query, userID, token)
	return err
}
