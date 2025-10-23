package customerrors

import "net/http"

// ValidationError represents a validation error for a specific field
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrorResponse is the error response for failed validations
type ValidationErrorResponse struct {
	Status  int               `json:"-"`
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Errors  []ValidationError `json:"errors,omitempty"`
}

func (v *ValidationErrorResponse) Error() string     { return v.Message }
func (v *ValidationErrorResponse) StatusCode() int   { return v.Status }
func (v *ValidationErrorResponse) ErrorCode() string { return v.Code }

// NewValidationError creates a new validation error response
func NewValidationError(message string, errors []ValidationError) *ValidationErrorResponse {
	return &ValidationErrorResponse{
		Status:  http.StatusBadRequest,
		Code:    "VALIDATION_ERROR",
		Message: message,
		Errors:  errors,
	}
}
