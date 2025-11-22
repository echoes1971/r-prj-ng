package api

import (
	"encoding/json"
	"net/http"
)

// APIError represents a structured error response with i18n support
type APIError struct {
	Code    string            `json:"code"`             // Error code for frontend i18n
	Message string            `json:"message"`          // Fallback message in English
	Params  map[string]string `json:"params,omitempty"` // Dynamic parameters for interpolation
}

// Error codes - keep in sync with fe/src/locales/*/errors.json
const (
	ErrUserAlreadyExists    = "USER_ALREADY_EXISTS"
	ErrUserNotFound         = "USER_NOT_FOUND"
	ErrGroupAlreadyExists   = "GROUP_ALREADY_EXISTS"
	ErrGroupNotFound        = "GROUP_NOT_FOUND"
	ErrUnauthorized         = "UNAUTHORIZED"
	ErrForbidden            = "FORBIDDEN"
	ErrInvalidRequest       = "INVALID_REQUEST"
	ErrMissingField         = "MISSING_FIELD"
	ErrInternalServer       = "INTERNAL_SERVER_ERROR"
	ErrInvalidToken         = "INVALID_TOKEN"
	ErrMissingAuthorization = "MISSING_AUTHORIZATION"

	ErrObjectNotFound = "OBJECT_NOT_FOUND"
)

// RespondError sends a structured error response
func RespondError(w http.ResponseWriter, code string, message string, params map[string]string, httpStatus int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(APIError{
		Code:    code,
		Message: message,
		Params:  params,
	})
}

// RespondSimpleError sends an error without parameters
func RespondSimpleError(w http.ResponseWriter, code string, message string, httpStatus int) {
	RespondError(w, code, message, nil, httpStatus)
}
