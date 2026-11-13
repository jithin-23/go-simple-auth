package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"simple-auth/internal/apperror"
	"simple-auth/internal/httpx"
	"simple-auth/internal/token"
)

const claimsKey = "auth.claims"

// RequireAuth checks the "Authorization: Bearer <token>" header (≈ AuthGuard).
// On success it stores the token's claims in the context for later handlers.
func RequireAuth(tokens *token.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		scheme, tokenString, found := strings.Cut(c.GetHeader("Authorization"), " ")
		if !found || !strings.EqualFold(scheme, "Bearer") || tokenString == "" {
			httpx.Error(c, apperror.Unauthorized("missing or malformed Authorization header"))
			return
		}

		claims, err := tokens.Parse(tokenString)
		if err != nil {
			httpx.Error(c, apperror.Unauthorized("invalid or expired token"))
			return
		}

		c.Set(claimsKey, claims)
		c.Next()
	}
}

// CurrentUser returns the claims stored by RequireAuth (≈ request.user).
// ok is false if the route isn't behind RequireAuth.
func CurrentUser(c *gin.Context) (*token.Claims, bool) {
	value, exists := c.Get(claimsKey)
	if !exists {
		return nil, false
	}
	claims, ok := value.(*token.Claims)
	return claims, ok
}
