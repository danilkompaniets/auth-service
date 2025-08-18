package http

import (
	"github.com/danilkompaniets/auth-service/internal/infrastructure/config"
	"github.com/danilkompaniets/auth-service/internal/interfaces/http"
	"github.com/gin-gonic/gin"
)

type router struct {
	cfg     *config.Config
	handler gin.HandlerFunc
}

func SetupRoutes(handler *http.HttpHandler) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(gin.Logger())

	api := router.Group("api/v1/auth")
	api.POST("/login", handler.Login)
	api.POST("/register", handler.Register)
	api.POST("/refresh-token", handler.RefreshTokens)
	api.POST("/logout", handler.Logout)

	return router
}
