package grpc

import (
	"context"
	"github.com/danilkompaniets/auth-service/internal/infrastructure/config"
	grpc2 "github.com/danilkompaniets/auth-service/internal/interfaces/grpc"
	gen_auth "github.com/danilkompaniets/go-chat-common/gen/gen-auth"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"log"
	"net"
)

type GRPCApp struct {
	handler    *grpc2.AuthGRPCHandler
	cfg        config.Config
	grpcServer *grpc.Server
	listener   net.Listener
}

func NewGRPCApp(handler *grpc2.AuthGRPCHandler, cfg config.Config) *GRPCApp {
	return &GRPCApp{
		handler: handler,
		cfg:     cfg,
	}
}

func (a *GRPCApp) Run() error {
	a.grpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
	)

	gen_auth.RegisterAuthServiceServer(a.grpcServer, a.handler)

	grpc_prometheus.Register(a.grpcServer)
	grpc_prometheus.EnableHandlingTimeHistogram() // замер времени запросов

	lis, err := net.Listen("tcp", a.cfg.App.GrpcAddr)
	if err != nil {
		return err
	}
	a.listener = lis

	log.Println("gRPC started on port:", a.cfg.App.GrpcAddr)
	return a.grpcServer.Serve(lis)
}

func (a *GRPCApp) Stop(ctx context.Context) error {
	if a.grpcServer == nil {
		return nil
	}

	stopped := make(chan struct{})
	go func() {
		a.grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-ctx.Done():
		// если не успели — останавливаем принудительно
		a.grpcServer.Stop()
		return ctx.Err()
	case <-stopped:
		log.Println("gRPC server stopped gracefully")
		return nil
	}
}
