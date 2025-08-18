package sqlRepo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/danilkompaniets/auth-service/internal/infrastructure/repository"
)

func (r *Repository) DeleteRefreshToken(ctx context.Context, userID int64) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return repository.ErrUserNotFound
	}
	return err
}
