// @title Auth Service API
// @version 1.0
// @description API сервис авторизации (JWT + cookies).
// @host localhost:8080
// @BasePath /auth
package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/danilkompaniets/auth-service/internal/application"
	"github.com/danilkompaniets/auth-service/internal/infrastructure/config"
	"github.com/danilkompaniets/auth-service/internal/infrastructure/grpc"
	"github.com/danilkompaniets/auth-service/internal/infrastructure/http"
	sqlRepo "github.com/danilkompaniets/auth-service/internal/infrastructure/repository/sqlRepo"
	grpc2 "github.com/danilkompaniets/auth-service/internal/interfaces/grpc"
	"github.com/danilkompaniets/go-chat-common/tokens"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.MustLoad()

	dbCfg := cfg.App.Database
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.Username,
		dbCfg.Password,
		dbCfg.Database,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("database ping failed: %v", err)
	}

	accessTokenTTL, err := time.ParseDuration(cfg.App.Env.AccessTokenTTL)
	if err != nil {
		log.Fatalf("invalid access token TTL: %v", err)
	}

	refreshTokenTTL, err := time.ParseDuration(cfg.App.Env.RefreshTokenTTL)
	if err != nil {
		log.Fatalf("invalid refresh token TTL: %v", err)
	}

	jwtManager := tokens.NewJWTManager(
		cfg.App.Env.AccessTokenSecret,
		cfg.App.Env.RefreshTokenSecret,
		accessTokenTTL,
		refreshTokenTTL,
	)

	// Сервисы и репо
	repo := sqlRepo.NewAuthRepository(db)
	svc := application.NewAuthService(repo, jwtManager)

	grpcHandler := grpc2.NewAuthGRPCHandler(svc)
	grpcApp := grpc.NewGRPCApp(grpcHandler, *cfg)
	httpApp := http.NewHttpApplication(svc)

	errs := make(chan error, 2)

	go func() {
		log.Println("Starting gRPC server...")
		if err := grpcApp.Run(); err != nil {
			errs <- fmt.Errorf("grpc server error: %w", err)
		}
	}()

	go func() {
		log.Println("Starting HTTP server on :" + cfg.App.HttpAddr)
		if err := httpApp.Run(cfg.App.HttpAddr); err != nil && err != context.Canceled {
			errs <- fmt.Errorf("http server error: %w", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		log.Printf("Shutting down servers due to signal: %v", sig)
	case err := <-errs:
		log.Fatalf("Server error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := grpcApp.Stop(ctx); err != nil {
		log.Printf("Error stopping gRPC server: %v", err)
	}
	if err := httpApp.Shutdown(ctx); err != nil {
		log.Printf("Error when stopping HTTP server: %v", err)
	}

	log.Println("Servers stopped gracefully")
}
