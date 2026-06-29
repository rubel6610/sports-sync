package handler

import (
	"net/http"
	"spotsync-api/dto"
	"spotsync-api/service"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService service.AuthService
	validator   *validator.Validate
}

func NewAuthHandler(as service.AuthService) *AuthHandler {
	return &AuthHandler{authService: as, validator: validator.New()}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.APIResponse{Success: false, Message: "Invalid request payload"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.APIResponse{Success: false, Message: "Validation failed", Errors: err.Error()})
	}

	res, err := h.authService.Register(req)
	if err != nil {
		return c.JSON(http.StatusConflict, dto.APIResponse{Success: false, Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, dto.APIResponse{Success: true, Message: "User registered successfully", Data: res})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.APIResponse{Success: false, Message: "Invalid payload"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.APIResponse{Success: false, Message: err.Error()})
	}

	res, err := h.authService.Login(req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, dto.APIResponse{Success: false, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.APIResponse{Success: true, Message: "Login successful", Data: res})
}