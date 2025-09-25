package domain

import (
	"errors"
	"fmt"
	"strings"
)

// ValidationError represents field validation failures.
type ValidationError struct {
	ErrorMap map[string][]string
}

// Error implements the error interface for ValidationError.
func (ve *ValidationError) Error() string {
	if len(ve.ErrorMap) == 0 {
		return "validation failed"
	}

	var parts []string
	for field, errs := range ve.ErrorMap {
		parts = append(parts, fmt.Sprintf("%s: %s", field, strings.Join(errs, ", ")))
	}
	return fmt.Sprintf("validation failed: %s", strings.Join(parts, "; "))
}

// Predefined domain errors for consistent error handling.
var (
	BadRequestError     = errors.New("bad request")
	NotFoundError       = errors.New("not found")
	InternalServerError = errors.New("internal server error")
)
