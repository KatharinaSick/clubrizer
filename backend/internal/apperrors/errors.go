package apperrors

import (
	"errors"
	"fmt"
	"net/http"
)

type Kind string

const (
	NotFound     Kind = "not_found"
	Unauthorized Kind = "Unauthorized"
	BadRequest   Kind = "bad_request"
)

// Error defines a standard application error.
type Error struct {
	kind           Kind
	message        string
	cause          error
	httpStatusCode int
}

func (e *Error) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %s. cause: %s", e.kind, e.message, e.cause)
	}
	return fmt.Sprintf("%s: %s", e.kind, e.message)
}

func Is(err error, kind Kind) bool {
	appErr := &Error{}
	if errors.As(err, &appErr) && appErr.kind == kind {
		return true
	}
	return false
}

func HttpStatusCode(err error) int {
	appErr := &Error{}
	if errors.As(err, &appErr) {
		return appErr.httpStatusCode
	} else {
		return http.StatusInternalServerError
	}
}

func NewNotFound(message string) *Error {
	return &Error{
		kind:           NotFound,
		message:        message,
		httpStatusCode: http.StatusNotFound,
	}
}

func NewUnauthorized(message string, cause error) *Error {
	return &Error{
		kind:           Unauthorized,
		message:        message,
		cause:          cause,
		httpStatusCode: http.StatusUnauthorized,
	}
}

func NewBadRequest(message string, cause error) *Error {
	return &Error{
		kind:           BadRequest,
		message:        message,
		cause:          cause,
		httpStatusCode: http.StatusBadRequest,
	}
}
