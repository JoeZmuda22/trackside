package handlers

import (
	"net/http"

	"github.com/joezmuda/trackside-backend/internal/middleware"
	"github.com/joezmuda/trackside-backend/internal/models"
	"github.com/joezmuda/trackside-backend/internal/repository"
	"github.com/labstack/echo/v4"
)

type ZoneTipHandler struct {
	zoneRepo *repository.ZoneRepo
}

func NewZoneTipHandler(zoneRepo *repository.ZoneRepo) *ZoneTipHandler {
	return &ZoneTipHandler{zoneRepo: zoneRepo}
}

// POST /api/tracks/:id/zones/:zoneId/tips
func (h *ZoneTipHandler) Create(c echo.Context) error {
	userID := middleware.GetUserID(c)
	trackID := c.Param("id")
	zoneID := c.Param("zoneId")

	zone, err := h.zoneRepo.FindZoneForTrack(zoneID, trackID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	if zone == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Zone not found"})
	}

	var req models.ZoneTipRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if req.Content == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Tip content is required"})
	}
	if req.Conditions != nil && !models.ValidDrivingCondition(*req.Conditions) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid conditions"})
	}

	tip, err := h.zoneRepo.CreateTip(req.Content, req.Conditions, zoneID, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	return c.JSON(http.StatusCreated, tip)
}
