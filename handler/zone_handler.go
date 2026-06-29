package handler

import (
	"net/http"
	"strconv"
	"spotsync-api/dto"
	"spotsync-api/service"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type ZoneHandler struct {
	zoneService service.ZoneService
	validator   *validator.Validate
}

func NewZoneHandler(zs service.ZoneService) *ZoneHandler {
	return &ZoneHandler{zoneService: zs, validator: validator.New()}
}

func (h *ZoneHandler) Create(c echo.Context) error {
	var req dto.CreateZoneRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.APIResponse{Success: false, Message: "Invalid payload"})
	}
	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.APIResponse{Success: false, Message: err.Error()})
	}
	res, err := h.zoneService.CreateZone(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.APIResponse{Success: false, Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, dto.APIResponse{Success: true, Message: "Parking zone created successfully", Data: res})
}

func (h *ZoneHandler) GetAll(c echo.Context) error {
	res, err := h.zoneService.GetAllZones()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.APIResponse{Success: false, Message: err.Error()})
	}
	return c.JSON(http.StatusOK, dto.APIResponse{Success: true, Message: "Parking zones retrieved successfully", Data: res})
}

func (h *ZoneHandler) GetOne(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	res, err := h.zoneService.GetZoneByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, dto.APIResponse{Success: false, Message: "Zone not found"})
	}
	return c.JSON(http.StatusOK, dto.APIResponse{Success: true, Message: "Parking zone retrieved successfully", Data: res})
}