package http

import (
	"context"
	"github.com/danilkompaniets/auth-service/internal/application"
	http2 "github.com/danilkompaniets/auth-service/internal/interfaces/http"
	"net/http"
	"time"
)

type HttpApplication struct {
	server  *http.Server
	service *application.AuthService
}

func NewHttpApplication(service *application.AuthService) *HttpApplication {
	handler := http2.NewHttpHandler(service)
	r := SetupRoutes(handler)

	return &HttpApplication{
		service: service,
		server: &http.Server{
			Handler:      r,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
	}
}

func (app *HttpApplication) Run(port string) error {
	app.server.Addr = port
	return app.server.ListenAndServe()
}

func (app *HttpApplication) Shutdown(ctx context.Context) error {
	return app.server.Shutdown(ctx)
}
