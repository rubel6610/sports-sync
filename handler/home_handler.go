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

	return c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Welcome to SpotSync API",
	})
}
