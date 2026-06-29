package handler

import (
	"net/http"
	"spotsync-api/dto"

	"github.com/labstack/echo/v4"
)

type HomeHandler struct{}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (h *HomeHandler) Home(c echo.Context) error {
	data := map[string]interface{}{
		"api_name":    "SpotSync API",
		"version":     "1.0.0",
		"description": "Parking Management System API",
		"endpoints": map[string]interface{}{
			"auth": map[string]string{
				"register": "POST /api/v1/auth/register",
				"login":    "POST /api/v1/auth/login",
			},
			"parking_zones": map[string]string{
				"create":  "POST /api/v1/zones (admin only)",
				"get_all": "GET /api/v1/zones",
				"get_one": "GET /api/v1/zones/:id",
			},
			"reservations": map[string]string{
				"create":  "POST /api/v1/reservations",
				"get_my":  "GET /api/v1/reservations/my-reservations",
				"cancel":  "DELETE /api/v1/reservations/:id",
				"get_all": "GET /api/v1/reservations (admin only)",
			},
		},
	}

	return c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Welcome to SpotSync API",
		Data:    data,
	})
}
