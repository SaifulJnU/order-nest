package custom_error

import "net/http"

type HTTPError struct {
	Message string              `json:"message,omitempty"`
	Type    string              `json:"type,omitempty"`
	Code    int                 `json:"code,omitempty"`
	Errors  map[string][]string `json:"errors,omitempty"`
}

func (e *HTTPError) Error() string {
	return e.Message
}

const (
	TypeError   = "error"
	TypeSuccess = "success"
)

// New returns a new HTTPError with the given message and status code.
func New(message string, code int) HTTPError {
	return HTTPError{Message: message, Type: TypeError, Code: code}
}

// BadRequest returns a 400 error with the given message.
func BadRequest(message string) HTTPError {
	if message == "" {
		message = "bad request"
	}
	return New(message, http.StatusBadRequest)
}

// Internal returns a 500 error with a safe default message.
func Internal() HTTPError {
	return New("internal server error", http.StatusInternalServerError)
}

// Unauthorized returns a 401 error with the given message.
func Unauthorized(message string) HTTPError {
	if message == "" {
		message = "unauthorized"
	}
	return New(message, http.StatusUnauthorized)
}

// WithErrors attaches a validation error map to the base HTTPError.
func WithErrors(base HTTPError, errs map[string][]string) HTTPError {
	base.Errors = errs
	return base
}

// Backward-compatible exported variable used across the codebase.
var (
	Unauthrized = Unauthorized("")
)
