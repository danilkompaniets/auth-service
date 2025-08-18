package sqlRepo

import (
	"context"
	"github.com/danilkompaniets/auth-service/pkg/model"
)

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
