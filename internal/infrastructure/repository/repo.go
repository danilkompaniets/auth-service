package repository

import (
	"context"
	"github.com/danilkompaniets/auth-service/pkg/model"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user model.User) (int64, error)
	DeleteRefreshToken(ctx context.Context, userID int64) error
	SaveRefreshToken(ctx context.Context, userID int64, token string) error
	GetRefreshToken(ctx context.Context, userID int64) (string, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
}
