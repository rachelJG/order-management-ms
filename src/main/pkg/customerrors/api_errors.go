package customerrors

import "net/http"

// API Errors
var (
	// 400 Bad Request - General
	ErrInvalidRequest = &apiError{status: http.StatusBadRequest, code: "INVALID_REQUEST", message: "invalid request"}

	// 400 Bad Request - Order related
	ErrInvalidOrderID     = &apiError{status: http.StatusBadRequest, code: "INVALID_ORDER_ID", message: "invalid order ID"}
	ErrInvalidOrderStatus = &apiError{status: http.StatusBadRequest, code: "INVALID_ORDER_STATUS", message: "invalid order status"}
	ErrInvalidOrder       = &apiError{status: http.StatusBadRequest, code: "INVALID_ORDER", message: "invalid order"}
	ErrInvalidTransition  = &apiError{status: http.StatusBadRequest, code: "INVALID_TRANSITION", message: "invalid status transition"}
	ErrItemsIsRequired    = &apiError{status: http.StatusBadRequest, code: "ITEMS_REQUIRED", message: "items are required"}
	ErrInvalidStatus      = &apiError{status: http.StatusBadRequest, code: "INVALID_STATUS", message: "invalid status"}
	// 401 Unauthorized
	ErrMissingAuthToken = &apiError{status: http.StatusUnauthorized, code: "MISSING_AUTH_TOKEN", message: "missing authorization token"}
	ErrInvalidAuthToken = &apiError{status: http.StatusUnauthorized, code: "INVALID_AUTH_TOKEN", message: "invalid or expired authorization token"}

	// 404 Not Found
	ErrOrderNotFound = &apiError{status: http.StatusNotFound, code: "ORDER_NOT_FOUND", message: "order not found"}

	// 409 Conflict
	ErrOrderAlreadyExists = &apiError{status: http.StatusConflict, code: "ORDER_ALREADY_EXISTS", message: "order already exists"}

	// 500 Internal Server Error
	ErrInternalServer = &apiError{status: http.StatusInternalServerError, code: "INTERNAL_SERVER_ERROR", message: "internal server error"}

	ErrFailedToCreateOrder = &apiError{status: http.StatusInternalServerError, code: "FAILED_TO_CREATE_ORDER", message: "failed to create order"}
)

// apiError implements the Error interface
type apiError struct {
	status  int
	code    string
	message string
}

func (e *apiError) Error() string     { return e.message }
func (e *apiError) StatusCode() int   { return e.status }
func (e *apiError) ErrorCode() string { return e.code }
