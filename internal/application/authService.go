package application

import (
	"context"
	"errors"
	"github.com/danilkompaniets/auth-service/internal/infrastructure/repository"
	"github.com/danilkompaniets/auth-service/pkg/model"
	"github.com/danilkompaniets/go-chat-common/tokens"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo       repository.AuthRepository
	jwtManager tokens.TokenManager
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewAuthService(repo repository.AuthRepository, manager tokens.TokenManager) *AuthService {
	return &AuthService{repo: repo, jwtManager: manager}
}

func (s *AuthService) CreateUser(ctx context.Context, user model.User) (int64, error) {
	if user.Email == "" || user.Password == "" {
		return 0, errors.New("user fields cannot be empty")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	user.Password = string(hash)

	return s.repo.CreateUser(ctx, user)
}

func (s *AuthService) LoginUser(ctx context.Context, user model.User) (*Tokens, error) {
	if user.Email == "" || user.Password == "" {
		return nil, errors.New("user fields cannot be empty")
	}

	userFound, err := s.repo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(user.Password))
	if err != nil {
		return nil, err
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(userFound.Id)
	refreshToken, err := s.jwtManager.GenerateRefreshToken(userFound.Id)

	err = s.repo.SaveRefreshToken(ctx, userFound.Id, refreshToken)
	if err != nil {
		return nil, err
	}

	return &Tokens{RefreshToken: refreshToken, AccessToken: accessToken}, nil
}

func (s *AuthService) RefreshUserTokens(ctx context.Context, refreshToken string) (*Tokens, error) {
	userID, err := s.jwtManager.VerifyRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	newAccessToken, err := s.jwtManager.GenerateAccessToken(userID)
	if err != nil {
		return nil, err
	}
	newRefreshToken, err := s.jwtManager.GenerateRefreshToken(userID)
	if err != nil {
		return nil, err
	}

	err = s.repo.DeleteRefreshToken(ctx, userID)
	if err != nil {
		return nil, err
	}

	err = s.repo.SaveRefreshToken(ctx, userID, newRefreshToken)
	if err != nil {
		return nil, err
	}

	return &Tokens{RefreshToken: newRefreshToken, AccessToken: newAccessToken}, nil
}

func (s *AuthService) GetRefreshToken(ctx context.Context, userID int64) (string, error) {
	if userID == 0 {
		return "", errors.New("userID must not be empty")
	}
	return s.repo.GetRefreshToken(ctx, userID)
}

func (s *AuthService) DeleteRefreshToken(ctx context.Context, userID int64) error {
	if userID == 0 {
		return errors.New("userID must not be empty")
	}
	return s.repo.DeleteRefreshToken(ctx, userID)
}

func (s *AuthService) GetUserByRefreshToken(ctx context.Context, token string) (int64, error) {
	if token == "" {
		return 0, errors.New("token must not be empty")
	}
	userID, err := s.jwtManager.VerifyAccessToken(token)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (s *AuthService) ValidateToken(token string) (int64, error) {
	if token == "" {
		return 0, errors.New("token must not be empty")
	}
	userID, err := s.jwtManager.VerifyAccessToken(token)
	if err != nil {
		return 0, err
	}
	if userID == 0 {
		return 0, errors.New("user not found")
	}

	return userID, nil
}
