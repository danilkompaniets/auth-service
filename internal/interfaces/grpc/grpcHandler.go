package grpc

import (
	"context"
	"github.com/danilkompaniets/auth-service/internal/application"
	gen_auth "github.com/danilkompaniets/go-chat-common/gen/gen-auth"
)

type AuthGRPCHandler struct {
	gen_auth.UnimplementedAuthServiceServer
	service *application.AuthService
}

func NewAuthGRPCHandler(service *application.AuthService) *AuthGRPCHandler {
	return &AuthGRPCHandler{
		service: service,
	}
}

func (h *AuthGRPCHandler) ValidateToken(ctx context.Context, req *gen_auth.ValidateTokenRequest) (*gen_auth.ValidateTokenResponse, error) {
	userId, err := h.service.ValidateToken(req.Token)
	if err != nil {
		return nil, err
	}

	return &gen_auth.ValidateTokenResponse{
		Valid:  userId != 0,
		UserId: userId,
	}, nil
}
