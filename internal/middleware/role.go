package middleware

import (
	"slices"

	"github.com/gin-gonic/gin"

	"simple-auth/internal/apperror"
	"simple-auth/internal/httpx"
)

// RequireRole allows the request only if the logged-in user has one of the roles
// (≈ RolesGuard + @Roles()). It must run after RequireAuth.
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := CurrentUser(c)
		if !ok {
			httpx.Error(c, apperror.Unauthorized("not authenticated"))
			return
		}
		if !slices.Contains(roles, claims.Role) {
			httpx.Error(c, apperror.Forbidden("you do not have permission to do this"))
			return
		}
		c.Next()
	}
}
