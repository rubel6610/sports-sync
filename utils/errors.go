package utils

import (
	"net/http"
	"spotsync-api/dto"

	"github.com/labstack/echo/v4"
)

// CustomHTTPErrorHandler handles all errors globally and prevents internal server/GORM details from leaking
func CustomHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	message := "An unexpected server error occurred"
	var errorsDetail interface{} = nil

	// Check if it's an Echo HTTP Error (like 404 Not Found, 405 Method Not Allowed)
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		if m, ok := he.Message.(string); ok {
			message = m
		} else {
			message = "Resource or request error"
		}
		errorsDetail = he.Internal
	}

	// Masking raw internal errors (GORM/Database structural leaks)
	// If the status code is still 500, we keep the message generic for security
	if code == http.StatusInternalServerError {
		// Log the actual error internally for debugging (Optional)
		c.Logger().Error(err)
	} else {
		message = err.Error()
	}

	// Send standard uniform error response structure matching specification
	_ = c.JSON(code, dto.APIResponse{
		Success: false,
		Message: message,
		Errors:  errorsDetail,
	})
}