package application_test

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/danilkompaniets/auth-service/internal/infrastructure/security"
	"testing"
	"time"

	"github.com/danilkompaniets/auth-service/internal/application"
	"github.com/danilkompaniets/auth-service/internal/infrastructure/repository/sqlRepo"
	"github.com/danilkompaniets/auth-service/pkg/model"

	_ "github.com/lib/pq"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var db *sql.DB

func setupTestDB(t *testing.T) func() {
	ctx := context.Background()

	req := tc.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	container, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatal(err)
	}

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "5432")

	dsn := fmt.Sprintf("postgres://testuser:testpass@%s:%s/testdb?sslmode=disable", host, port.Port())
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}

	// Ждём пока база будет готова
	for i := 0; i < 10; i++ {
		if err := db.Ping(); err == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	// Создание таблиц для тестов
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL
	);
	CREATE TABLE IF NOT EXISTS refresh_tokens (
		user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
		token TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT NOW()
	);`
	_, err = db.Exec(schema)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Test DB ready")

	// Возвращаем функцию для остановки контейнера после теста
	return func() {
		db.Close()
		container.Terminate(ctx)
	}
}

func TestAuthServiceIntegration(t *testing.T) {
	cleanup := setupTestDB(t)
	defer cleanup()

	repo := sqlRepo.NewAuthRepository(db)
	jwtManager := security.NewJWTManager("access_secret", "refresh_secret", time.Minute*5, time.Hour*24)
	service := application.NewAuthService(repo, jwtManager)

	ctx := context.Background()

	// 1️⃣ Тест создания пользователя
	user := model.User{
		Email:     "test@example.com",
		Password:  "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	id, err := service.CreateUser(ctx, user)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	if id == 0 {
		t.Fatal("Expected user ID > 0")
	}

	// 2️⃣ Тест логина пользователя
	tokens, err := service.LoginUser(ctx, model.User{Email: user.Email, Password: "password123"})
	if err != nil {
		t.Fatalf("LoginUser failed: %v", err)
	}

	if tokens.AccessToken == "" || tokens.RefreshToken == "" {
		t.Fatal("Expected non-empty tokens")
	}

	// 3️⃣ Тест обновления токенов
	newTokens, err := service.RefreshUserTokens(ctx, tokens.RefreshToken)
	if err != nil {
		t.Fatalf("RefreshUserTokens failed: %v", err)
	}

	if newTokens.AccessToken == "" || newTokens.RefreshToken == "" {
		t.Fatal("Expected non-empty new tokens")
	}

	// 4️⃣ Проверка токена доступа
	userID, err := service.ValidateToken(newTokens.AccessToken)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	if userID != id {
		t.Fatalf("Expected userID %d, got %d", id, userID)
	}

	// 5️⃣ Проверка удаления refresh token
	err = service.DeleteRefreshToken(ctx, id)
	if err != nil {
		t.Fatalf("DeleteRefreshToken failed: %v", err)
	}
}
