package sqlRepo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/danilkompaniets/auth-service/internal/infrastructure/repository"
)

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
