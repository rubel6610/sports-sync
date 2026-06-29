package handler

import (
	"net/http"
	"spotsync-api/dto"
	"spotsync-api/service"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type ReservationHandler struct {
	resService service.ReservationService
	validator  *validator.Validate
}

func NewReservationHandler(rs service.ReservationService) *ReservationHandler {
	return &ReservationHandler{resService: rs, validator: validator.New()}
}

func (h *ReservationHandler) Create(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	var req dto.CreateReservationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.APIResponse{Success: false, Message: "Invalid payload"})
	}
	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.APIResponse{Success: false, Message: err.Error()})
	}

	res, err := h.resService.ReserveSpot(userID, req) // এখানে এখন ফিল্টারড ডিটিও আসবে
	if err != nil {
		if err.Error() == "parking zone is full" {
			return c.JSON(http.StatusConflict, dto.APIResponse{Success: false, Message: err.Error()})
		}
		return c.JSON(http.StatusBadRequest, dto.APIResponse{Success: false, Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, dto.APIResponse{
		Success: true,
		Message: "Reservation confirmed successfully",
		Data:    res, // সরাসরি ম্যাপড ডাটা রেসপন্সে চলে যাবে
	})
}

func (h *ReservationHandler) GetMy(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	res, err := h.resService.GetMyReservations(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.APIResponse{Success: false, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, dto.APIResponse{Success: true, Message: "My reservations retrieved successfully", Data: res})
}

func (h *ReservationHandler) Cancel(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	role := c.Get("role").(string)
	id, _ := strconv.Atoi(c.Param("id"))

	err := h.resService.CancelReservation(userID, role, uint(id))
	if err != nil {
		if err.Error() == "not_found" {
			return c.JSON(http.StatusNotFound, dto.APIResponse{Success: false, Message: "Reservation not found"})
		}
		if err.Error() == "forbidden" {
			return c.JSON(http.StatusForbidden, dto.APIResponse{Success: false, Message: "You cannot cancel someone else's reservation"})
		}
		return c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, dto.APIResponse{Success: true, Message: "Reservation cancelled successfully"})
}

func (h *ReservationHandler) GetAll(c echo.Context) error {
	res, err := h.resService.GetAllReservations()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.APIResponse{Success: false, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, dto.APIResponse{Success: true, Message: "All reservations retrieved successfully", Data: res})
}
