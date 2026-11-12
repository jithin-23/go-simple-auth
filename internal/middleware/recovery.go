package middleware

import (
	"fmt"
	"log/slog"
	"runtime/debug"

	"github.com/gin-gonic/gin"

	"simple-auth/internal/apperror"
	"simple-auth/internal/httpx"
)

// Recovery turns a panic into a JSON 500 response and logs the stack trace.
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		slog.Error("panic recovered",
			"panic", recovered, "path", c.Request.URL.Path, "stack", string(debug.Stack()))
		httpx.Error(c, apperror.Internal(fmt.Errorf("panic: %v", recovered)))
	})
}
