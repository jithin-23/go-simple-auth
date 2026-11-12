package user

import "github.com/gin-gonic/gin"

// RegisterRoutes mounts the user endpoints under rg (e.g. /api/v1 → /api/v1/users).
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/users") // unprotected for now — Module 6 adds auth + roles
	users.POST("", h.Create)
	users.GET("", h.List)
	users.GET("/:id", h.Get)
	users.PATCH("/:id", h.Update)
	users.DELETE("/:id", h.Delete)
}
