package sqlRepo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/danilkompaniets/auth-service/internal/infrastructure/repository"
	"github.com/danilkompaniets/auth-service/pkg/model"
)

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
