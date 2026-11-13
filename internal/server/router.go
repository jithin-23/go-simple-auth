package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"simple-auth/internal/apperror"
	"simple-auth/internal/auth"
	"simple-auth/internal/httpx"
	"simple-auth/internal/middleware"
	"simple-auth/internal/token"
	"simple-auth/internal/user"
)

// Handlers groups every feature's handler so NewRouter has one parameter per concern.
type Handlers struct {
	Auth *auth.Handler
	User *user.Handler
}

// NewRouter builds the Gin engine with global middleware and all routes.
func NewRouter(db *gorm.DB, tokens *token.Manager, h Handlers) *gin.Engine {
	httpx.SetupValidator() // validation errors use JSON field names

	r := gin.New() // empty engine: we choose the middleware ourselves
	r.Use(middleware.Logger(), middleware.Recovery())

	// We are not behind a proxy, so don't trust X-Forwarded-For headers.
	_ = r.SetTrustedProxies(nil)

	// Unknown routes get the same JSON error shape as everything else
	r.NoRoute(func(c *gin.Context) {
		httpx.Error(c, apperror.NotFound("route not found"))
	})

	r.GET("/health", healthHandler(db))

	api := r.Group("/api/v1")
	requireAuth := middleware.RequireAuth(tokens)

	h.Auth.RegisterRoutes(api, requireAuth)
	h.User.RegisterRoutes(api)

	return r
}

// healthHandler reports whether the API and the database are reachable.
func healthHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sqlDB, err := db.DB()
		if err == nil {
			err = sqlDB.PingContext(c.Request.Context())
		}
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable", "database": "down"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok", "database": "up"})
	}
}
