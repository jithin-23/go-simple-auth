package user

import "github.com/gin-gonic/gin"

// RegisterRoutes mounts the user endpoints under rg (e.g. /api/v1 → /api/v1/users).
// The router passes in the middleware, so this package decides WHICH routes need
// auth/admin without knowing HOW tokens are checked.
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup, requireAuth, requireAdmin gin.HandlerFunc) {
	users := rg.Group("/users", requireAuth) // every /users route needs a valid token

	users.GET("", h.List)         // any logged-in user
	users.GET("/:id", h.Get)      // any logged-in user
	users.PATCH("/:id", h.Update) // self or admin (checked in the service)

	users.POST("", requireAdmin, h.Create)
	users.DELETE("/:id", requireAdmin, h.Delete)
	users.PATCH("/:id/role", requireAdmin, h.ChangeRole)
}
