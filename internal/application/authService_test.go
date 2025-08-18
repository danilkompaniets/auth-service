package application

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"testing"

	"github.com/danilkompaniets/auth-service/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) CreateUser(ctx context.Context, user model.User) (int64, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockRepo) SaveRefreshToken(ctx context.Context, userID int64, token string) error {
	args := m.Called(ctx, userID, token)
	return args.Error(0)
}

func (m *MockRepo) DeleteRefreshToken(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockRepo) GetRefreshToken(ctx context.Context, userID int64) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

type MockJWT struct {
	mock.Mock
}

func (m *MockJWT) GenerateAccessToken(userID int64) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *MockJWT) GenerateRefreshToken(userID int64) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *MockJWT) VerifyRefreshToken(token string) (int64, error) {
	args := m.Called(token)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockJWT) VerifyAccessToken(token string) (int64, error) {
	args := m.Called(token)
	return args.Get(0).(int64), args.Error(1)
}

// --- Тесты ---

func TestCreateUser_Success(t *testing.T) {
	repo := new(MockRepo)
	jwt := new(MockJWT)
	service := NewAuthService(repo, jwt)

	user := model.User{Email: "test@test.com", Password: "123456"}
	repo.On("CreateUser", mock.Anything, mock.Anything).Return(int64(1), nil)

	id, err := service.CreateUser(context.Background(), user)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)
}

func TestCreateUser_EmptyFields(t *testing.T) {
	service := NewAuthService(nil, nil)

	_, err := service.CreateUser(context.Background(), model.User{})
	assert.Error(t, err)
}

func TestLoginUser_Success(t *testing.T) {
	repo := new(MockRepo)
	jwt := new(MockJWT)
	service := NewAuthService(repo, jwt)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	user := &model.User{Id: 1, Email: "test@test.com", Password: string(hashedPassword)}

	repo.On("GetUserByEmail", mock.Anything, "test@test.com").Return(user, nil)
	jwt.On("GenerateAccessToken", int64(1)).Return("access", nil)
	jwt.On("GenerateRefreshToken", int64(1)).Return("refresh", nil)
	repo.On("SaveRefreshToken", mock.Anything, int64(1), "refresh").Return(nil)

	tokens, err := service.LoginUser(context.Background(), model.User{Email: "test@test.com", Password: "123456"})
	assert.NoError(t, err)
	assert.Equal(t, "access", tokens.AccessToken)
	assert.Equal(t, "refresh", tokens.RefreshToken)
}

func TestLoginUser_WrongPassword(t *testing.T) {
	repo := new(MockRepo)
	jwt := new(MockJWT)
	service := NewAuthService(repo, jwt)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	user := &model.User{Id: 1, Email: "test@test.com", Password: string(hashedPassword)}

	repo.On("GetUserByEmail", mock.Anything, "test@test.com").Return(user, nil)

	_, err := service.LoginUser(context.Background(), model.User{Email: "test@test.com", Password: "wrong"})
	assert.Error(t, err)
}

func TestRefreshUserTokens_Success(t *testing.T) {
	repo := new(MockRepo)
	jwt := new(MockJWT)
	service := NewAuthService(repo, jwt)

	jwt.On("VerifyRefreshToken", "oldToken").Return(int64(1), nil)
	jwt.On("GenerateAccessToken", int64(1)).Return("newAccess", nil)
	jwt.On("GenerateRefreshToken", int64(1)).Return("newRefresh", nil)
	repo.On("DeleteRefreshToken", mock.Anything, int64(1)).Return(nil)
	repo.On("SaveRefreshToken", mock.Anything, int64(1), "newRefresh").Return(nil)

	tokens, err := service.RefreshUserTokens(context.Background(), "oldToken")
	assert.NoError(t, err)
	assert.Equal(t, "newAccess", tokens.AccessToken)
	assert.Equal(t, "newRefresh", tokens.RefreshToken)
}

func TestGetRefreshToken_Success(t *testing.T) {
	repo := new(MockRepo)
	service := NewAuthService(repo, nil)

	repo.On("GetRefreshToken", mock.Anything, int64(1)).Return("refreshToken", nil)

	token, err := service.GetRefreshToken(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, "refreshToken", token)
}

func TestDeleteRefreshToken_Success(t *testing.T) {
	repo := new(MockRepo)
	service := NewAuthService(repo, nil)

	repo.On("DeleteRefreshToken", mock.Anything, int64(1)).Return(nil)

	err := service.DeleteRefreshToken(context.Background(), 1)
	assert.NoError(t, err)
}
