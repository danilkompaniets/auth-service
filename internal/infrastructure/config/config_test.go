package config

import (
	"os"
	"testing"
)

func TestMustLoad_ConfigPathNotSet(t *testing.T) {
	old := os.Getenv("CONFIG_PATH")
	defer os.Setenv("CONFIG_PATH", old)

	os.Unsetenv("CONFIG_PATH")

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic when CONFIG_PATH is not set")
		}
	}()

	_ = MustLoad()
}

func TestMustLoad_FileDoesNotExist(t *testing.T) {
	old := os.Getenv("CONFIG_PATH")
	defer os.Setenv("CONFIG_PATH", old)

	os.Setenv("CONFIG_PATH", "nonexistent.yml")

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic when config file does not exist")
		}
	}()

	_ = MustLoad()
}

func TestMustLoad_Success(t *testing.T) {
	tmpFile := "tmp_config.yml"
	yamlContent := `application:
  grpc_addr: "localhost:9090"
  http_addr: "localhost:8080"
  database:
    host: "dbhost"
    port: "5432"
    username: "user"
    password: "pass"
    database: "mydb"
  environment:
    accessTokenSecret: ""
    refreshTokenSecret: ""
    accessTokenTTL: "15m"
    refreshTokenTTL: "720h"
`
	err := os.WriteFile(tmpFile, []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("failed to create temp config file: %v", err)
	}
	defer os.Remove(tmpFile)

	os.Setenv("CONFIG_PATH", tmpFile)
	os.Setenv("ACCESS_TOKEN_SECRET", "access123")
	os.Setenv("REFRESH_TOKEN_SECRET", "refresh123")
	os.Setenv("GRPC_ADDR", "127.0.0.1:9090")
	os.Setenv("HTTP_ADDR", "127.0.0.1:8080")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USERNAME", "admin")
	os.Setenv("DB_PASSWORD", "adminpass")
	os.Setenv("DB_NAME", "authdb")

	cfg := MustLoad()

	if cfg.App.Env.AccessTokenSecret != "access123" {
		t.Errorf("expected AccessTokenSecret to be 'access123', got %s", cfg.App.Env.AccessTokenSecret)
	}
	if cfg.App.Env.RefreshTokenSecret != "refresh123" {
		t.Errorf("expected RefreshTokenSecret to be 'refresh123', got %s", cfg.App.Env.RefreshTokenSecret)
	}
	if cfg.App.GrpcAddr != "127.0.0.1:9090" {
		t.Errorf("expected GrpcAddr to be '127.0.0.1:9090', got %s", cfg.App.GrpcAddr)
	}
	if cfg.App.Database.Username != "admin" {
		t.Errorf("expected DB username to be 'admin', got %s", cfg.App.Database.Username)
	}
}
