package errors

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

type AppError struct {
	Success    bool   `json:"success"`             // Always false for errors
	Message    string `json:"message"`             // Human-readable error message
	ErrorCode  string `json:"errorCode"`           // Machine-readable error code
	Details    string `json:"details,omitempty"`   // Additional details
	Path       string `json:"path,omitempty"`      // Request path
	Timestamp  int64  `json:"timestamp,omitempty"` // Unix timestamp
	HTTPStatus int    `json:"-"`                   // HTTP status code
	Internal   error  `json:"-"`                   // Internal error (not exposed to client)
}

func (e *AppError) Error() string {
	if e.Internal != nil {
		return fmt.Sprintf("[%s] %s: %s (internal: %v)", e.ErrorCode, e.Message, e.Details, e.Internal)
	}
	return fmt.Sprintf("[%s] %s: %s", e.ErrorCode, e.Message, e.Details)
}

func NewAppError(code, message string, httpStatus int) *AppError {
	return &AppError{
		Success:    false,
		ErrorCode:  code,
		Message:    message,
		HTTPStatus: httpStatus,
		Timestamp:  time.Now().Unix(),
	}
}

func (e *AppError) WithDetails(details string) *AppError {
	e.Details = details
	return e
}

func (e *AppError) WithInternal(err error) *AppError {
	e.Internal = err
	if e.Details == "" && err != nil {
		e.Details = err.Error()
	}
	return e
}

func (e *AppError) WithPath(path string) *AppError {
	e.Path = path
	return e
}

const (
	ErrCodeBadRequest         = "BAD_REQUEST"
	ErrCodeUnauthorized       = "UNAUTHORIZED"
	ErrCodeForbidden          = "FORBIDDEN"
	ErrCodeNotFound           = "NOT_FOUND"
	ErrCodeConflict           = "CONFLICT"
	ErrCodeValidationFailed   = "VALIDATION_FAILED"
	ErrCodeInternalServer     = "INTERNAL_SERVER_ERROR"
	ErrCodeServiceUnavailable = "SERVICE_UNAVAILABLE"
)

const (
	ErrCodePluginNotFound         = "PLUGIN_NOT_FOUND"
	ErrCodePluginAlreadyExists    = "PLUGIN_ALREADY_EXISTS"
	ErrCodePluginInvalid          = "PLUGIN_INVALID"
	ErrCodePluginInstallFailed    = "PLUGIN_INSTALL_FAILED"
	ErrCodePluginActivateFailed   = "PLUGIN_ACTIVATE_FAILED"
	ErrCodePluginDeactivateFailed = "PLUGIN_DEACTIVATE_FAILED"
	ErrCodePluginUninstallFailed  = "PLUGIN_UNINSTALL_FAILED"
	ErrCodePluginCallFailed       = "PLUGIN_CALL_FAILED"
)

// Predefined errors
var (
	ErrBadRequest         = NewAppError(ErrCodeBadRequest, "Bad request", http.StatusBadRequest)
	ErrUnauthorized       = NewAppError(ErrCodeUnauthorized, "Unauthorized", http.StatusUnauthorized)
	ErrForbidden          = NewAppError(ErrCodeForbidden, "Forbidden", http.StatusForbidden)
	ErrNotFound           = NewAppError(ErrCodeNotFound, "Resource not found", http.StatusNotFound)
	ErrConflict           = NewAppError(ErrCodeConflict, "Resource conflict", http.StatusConflict)
	ErrValidationFailed   = NewAppError(ErrCodeValidationFailed, "Validation failed", http.StatusBadRequest)
	ErrInternalServer     = NewAppError(ErrCodeInternalServer, "Internal server error", http.StatusInternalServerError)
	ErrServiceUnavailable = NewAppError(ErrCodeServiceUnavailable, "Service unavailable", http.StatusServiceUnavailable)

	ErrPluginNotFound         = NewAppError(ErrCodePluginNotFound, "Plugin not found", http.StatusNotFound)
	ErrPluginAlreadyExists    = NewAppError(ErrCodePluginAlreadyExists, "Plugin already exists", http.StatusConflict)
	ErrPluginInvalid          = NewAppError(ErrCodePluginInvalid, "Invalid plugin", http.StatusBadRequest)
	ErrPluginInstallFailed    = NewAppError(ErrCodePluginInstallFailed, "Failed to install plugin", http.StatusInternalServerError)
	ErrPluginActivateFailed   = NewAppError(ErrCodePluginActivateFailed, "Failed to activate plugin", http.StatusInternalServerError)
	ErrPluginDeactivateFailed = NewAppError(ErrCodePluginDeactivateFailed, "Failed to deactivate plugin", http.StatusInternalServerError)
	ErrPluginUninstallFailed  = NewAppError(ErrCodePluginUninstallFailed, "Failed to uninstall plugin", http.StatusInternalServerError)
	ErrPluginCallFailed       = NewAppError(ErrCodePluginCallFailed, "Failed to call plugin", http.StatusInternalServerError)
)

func APIErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	var appErr *AppError

	if e, ok := err.(*AppError); ok {
		appErr = e
		appErr.Path = lo.Ternary(appErr.Path == "", c.Request().URL.Path, appErr.Path)
	} else if he, ok := err.(*echo.HTTPError); ok {
		// Handle Echo's HTTPError
		appErr = &AppError{
			Success:    false,
			ErrorCode:  mapHTTPStatusToErrorCode(he.Code),
			Message:    fmt.Sprintf("%v", he.Message),
			Path:       c.Request().URL.Path,
			Timestamp:  time.Now().Unix(),
			HTTPStatus: he.Code,
		}
	} else {
		// Generic error
		appErr = &AppError{
			Success:    false,
			ErrorCode:  ErrCodeInternalServer,
			Message:    "Internal server error",
			Details:    err.Error(),
			Path:       c.Request().URL.Path,
			Timestamp:  time.Now().Unix(),
			HTTPStatus: http.StatusInternalServerError,
		}
	}

	c.Logger().Error(err)

	if err := c.JSON(appErr.HTTPStatus, appErr); err != nil {
		c.Logger().Error(err)
	}
}

// mapHTTPStatusToErrorCode maps HTTP status codes to error codes
func mapHTTPStatusToErrorCode(status int) string {
	switch status {
	case http.StatusBadRequest:
		return ErrCodeBadRequest
	case http.StatusUnauthorized:
		return ErrCodeUnauthorized
	case http.StatusForbidden:
		return ErrCodeForbidden
	case http.StatusNotFound:
		return ErrCodeNotFound
	case http.StatusConflict:
		return ErrCodeConflict
	case http.StatusInternalServerError:
		return ErrCodeInternalServer
	case http.StatusServiceUnavailable:
		return ErrCodeServiceUnavailable
	default:
		return ErrCodeInternalServer
	}
}
