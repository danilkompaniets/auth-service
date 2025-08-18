package sqlRepo

import (
	"context"
)

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
