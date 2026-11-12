// Package httpx holds HTTP helpers shared by every feature's handlers.
package httpx

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"simple-auth/internal/apperror"
)

type errorBody struct {
	Code    apperror.Code     `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

var statusByCode = map[apperror.Code]int{
	apperror.CodeBadRequest:   http.StatusBadRequest,
	apperror.CodeValidation:   http.StatusBadRequest,
	apperror.CodeUnauthorized: http.StatusUnauthorized,
	apperror.CodeForbidden:    http.StatusForbidden,
	apperror.CodeNotFound:     http.StatusNotFound,
	apperror.CodeConflict:     http.StatusConflict,
	apperror.CodeInternal:     http.StatusInternalServerError,
}

// Error writes err as a JSON error response and stops the handler chain.
// Equivalent of the HttpExceptionFilter in the NestJS app.
func Error(c *gin.Context, err error) {
	var appErr *apperror.Error
	if !errors.As(err, &appErr) {
		appErr = apperror.Internal(err) // anything unknown is a 500
	}

	status, ok := statusByCode[appErr.Code]
	if !ok {
		status = http.StatusInternalServerError
	}

	if status >= 500 {
		slog.Error("request failed",
			"method", c.Request.Method, "path", c.Request.URL.Path, "error", err)
	}

	c.AbortWithStatusJSON(status, gin.H{"error": errorBody{
		Code:    appErr.Code,
		Message: appErr.Message,
		Fields:  appErr.Fields,
	}})
}
