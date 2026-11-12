// Package apperror defines the errors the app returns to API clients.
// It knows nothing about HTTP; the httpx package maps codes to status codes.
package apperror

import "errors"

type Code string

const (
	CodeBadRequest   Code = "bad_request"
	CodeValidation   Code = "validation_failed"
	CodeUnauthorized Code = "unauthorized"
	CodeForbidden    Code = "forbidden"
	CodeNotFound     Code = "not_found"
	CodeConflict     Code = "conflict"
	CodeInternal     Code = "internal_error"
)

// Error is an error that is safe to show to the client.
type Error struct {
	Code    Code
	Message string            // shown to the client
	Fields  map[string]string // per-field validation messages (optional)
	Err     error             // underlying cause, for logs only (optional)
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

// Unwrap lets errors.Is / errors.As look at the underlying cause.
func (e *Error) Unwrap() error { return e.Err }

func BadRequest(msg string) *Error   { return &Error{Code: CodeBadRequest, Message: msg} }
func Unauthorized(msg string) *Error { return &Error{Code: CodeUnauthorized, Message: msg} }
func Forbidden(msg string) *Error    { return &Error{Code: CodeForbidden, Message: msg} }
func NotFound(msg string) *Error     { return &Error{Code: CodeNotFound, Message: msg} }
func Conflict(msg string) *Error     { return &Error{Code: CodeConflict, Message: msg} }

func Validation(fields map[string]string) *Error {
	return &Error{Code: CodeValidation, Message: "request validation failed", Fields: fields}
}

// Internal wraps an unexpected error. The cause is logged, never shown to the client.
func Internal(err error) *Error {
	return &Error{Code: CodeInternal, Message: "internal server error", Err: err}
}

// Is reports whether err (or anything it wraps) is an *Error with the given code.
func Is(err error, code Code) bool {
	var e *Error
	return errors.As(err, &e) && e.Code == code
}
