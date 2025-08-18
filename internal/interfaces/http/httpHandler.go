package http

import (
	"github.com/danilkompaniets/auth-service/internal/application"
	"github.com/danilkompaniets/auth-service/pkg/api"
	"github.com/danilkompaniets/auth-service/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type HttpHandler struct {
	service *application.AuthService
}

func NewHttpHandler(service *application.AuthService) *HttpHandler {
	return &HttpHandler{
		service: service,
	}
}

// Register godoc
// @Summary      User registration
// @Description  Creates a new user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body registerRequest true "Register request"
// @Success      200  {object} map[string]interfaces{} "userId"
// @Failure      400  {object} map[string]string "bad request"
// @Failure      500  {object} map[string]string "internal error"
// @Router       /auth/register [post]
func (h *HttpHandler) Register(c *gin.Context) {
	var req api.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := model.User{
		Email:     req.Email,
		Password:  req.Password,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	userId, err := h.service.CreateUser(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"userId": userId})
}

// Login godoc
// @Summary      User login
// @Description  Authenticates user and returns tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body loginRequest true "Login request"
// @Success      200  {object} map[string]string "access_token"
// @Failure      400  {object} map[string]string "bad request"
// @Failure      500  {object} map[string]string "internal error"
// @Router       /auth/login [post]
func (h *HttpHandler) Login(c *gin.Context) {
	var req api.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := model.User{
		Email:    req.Email,
		Password: req.Password,
	}

	tokens, err := h.service.LoginUser(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tokens.AccessToken = "Bearer " + tokens.AccessToken
	tokens.RefreshToken = "Bearer " + tokens.RefreshToken

	c.SetCookie("refresh_token", tokens.RefreshToken, 60*60*24*7, "/auth/refresh", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{
		"access_token": tokens.AccessToken,
	})
}

// Logout godoc
// @Summary      User logout
// @Description  Deletes user refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body logoutRequest true "Logout request"
// @Success      200  {object} map[string]string "ok"
// @Failure      400  {object} map[string]string "bad request"
// @Failure      500  {object} map[string]string "internal error"
// @Router       /auth/logout [post]
func (h *HttpHandler) Logout(c *gin.Context) {
	var req api.LogoutRequest
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.DeleteRefreshToken(c, req.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// RefreshTokens godoc
// @Summary      Refresh user tokens
// @Description  Uses refresh_token cookie to issue new tokens
// @Tags         auth
// @Produce      json
// @Success      200  {object} map[string]string "access_token"
// @Failure      500  {object} map[string]string "internal error"
// @Router       /auth/refresh [get]
func (h *HttpHandler) RefreshTokens(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing refresh token"})
		return
	}

	res, err := h.service.RefreshUserTokens(c, refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("refresh_token", res.RefreshToken, 60*60*24*7, "/auth/refresh", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{
		"access_token": res.AccessToken,
	})
}
