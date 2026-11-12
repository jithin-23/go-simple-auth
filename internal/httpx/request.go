package httpx

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"simple-auth/internal/apperror"
)

// SetupValidator makes validation errors use JSON field names ("email")
// instead of Go field names ("Email"). Call once at startup.
func SetupValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(f reflect.StructField) string {
			name := strings.SplitN(f.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}
}

// BindJSON parses and validates the JSON body into obj (≈ ValidationPipe).
// On failure it writes a 400 and returns false, so the handler just returns.
func BindJSON(c *gin.Context, obj any) bool {
	err := c.ShouldBindJSON(obj)
	if err == nil {
		return true
	}

	var verrs validator.ValidationErrors
	if errors.As(err, &verrs) {
		fields := make(map[string]string, len(verrs))
		for _, fe := range verrs {
			fields[fe.Field()] = validationMessage(fe)
		}
		Error(c, apperror.Validation(fields))
		return false
	}

	Error(c, apperror.BadRequest("invalid JSON body"))
	return false
}

// ParamID reads a positive integer URL parameter (≈ ParseIntPipe).
// On failure it writes a 400 and returns false.
func ParamID(c *gin.Context, name string) (uint, bool) {
	id, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil || id == 0 {
		Error(c, apperror.BadRequest(fmt.Sprintf("invalid %s", name)))
		return 0, false
	}
	return uint(id), true
}

func validationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email"
	case "min":
		return fmt.Sprintf("must be at least %s characters", fe.Param())
	case "max":
		return fmt.Sprintf("must be at most %s characters", fe.Param())
	case "oneof":
		return "must be one of: " + strings.ReplaceAll(fe.Param(), " ", ", ")
	default:
		return "is invalid"
	}
}
