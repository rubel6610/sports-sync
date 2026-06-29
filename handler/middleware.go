package handler

import (
	"net/http"
	"spotsync-api/dto"
	"spotsync-api/utils"
	"strings"

	"github.com/labstack/echo/v4"
)

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, dto.APIResponse{Success: false, Message: "Missing auth token"})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, dto.APIResponse{Success: false, Message: "Invalid token format"})
		}

		claims, err := utils.ValidateToken(parts[1])
		if err != nil {
			return c.JSON(http.StatusUnauthorized, dto.APIResponse{Success: false, Message: "Unauthorized access"})
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		return next(c)
	}
}

func RoleBlock(requiredRole string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role := c.Get("role").(string)
			if role != requiredRole {
				return c.JSON(http.StatusForbidden, dto.APIResponse{Success: false, Message: "Access denied. Insufficient privileges."})
			}
			return next(c)
		}
	}
}
