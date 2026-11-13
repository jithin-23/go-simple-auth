package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"simple-auth/internal/apperror"
	"simple-auth/internal/httpx"
	"simple-auth/internal/middleware"
	"simple-auth/internal/user"
)

// Handler is the HTTP layer for auth (≈ AuthController in NestJS).
type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes mounts /auth routes. signup and login are public; me needs a token.
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup, requireAuth gin.HandlerFunc) {
	a := rg.Group("/auth")
	a.POST("/signup", h.Signup)
	a.POST("/login", h.Login)
	a.GET("/me", requireAuth, h.Me)
}

// POST /auth/signup
func (h *Handler) Signup(c *gin.Context) {
	var req SignupRequest
	if !httpx.BindJSON(c, &req) {
		return
	}

	u, err := h.svc.Signup(c.Request.Context(), req)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	c.JSON(http.StatusCreated, user.ToResponse(u))
}

// POST /auth/login
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if !httpx.BindJSON(c, &req) {
		return
	}

	tokens, err := h.svc.Login(c.Request.Context(), req)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, tokens)
}

// GET /auth/me
func (h *Handler) Me(c *gin.Context) {
	claims, ok := middleware.CurrentUser(c)
	if !ok {
		httpx.Error(c, apperror.Unauthorized("not authenticated"))
		return
	}

	u, err := h.svc.Me(c.Request.Context(), claims.UserID)
	if err != nil {
		httpx.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, user.ToResponse(u))
}
