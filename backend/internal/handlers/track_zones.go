package handlers

import (
	"net/http"

	"github.com/joezmuda/trackside-backend/internal/middleware"
	"github.com/joezmuda/trackside-backend/internal/models"
	"github.com/joezmuda/trackside-backend/internal/repository"
	"github.com/labstack/echo/v4"
)

type TrackZoneHandler struct {
	trackRepo *repository.TrackRepo
	zoneRepo  *repository.ZoneRepo
}

func NewTrackZoneHandler(trackRepo *repository.TrackRepo, zoneRepo *repository.ZoneRepo) *TrackZoneHandler {
	return &TrackZoneHandler{trackRepo: trackRepo, zoneRepo: zoneRepo}
}

// POST /api/tracks/:id/zones
func (h *TrackZoneHandler) Create(c echo.Context) error {
	userID := middleware.GetUserID(c)
	_ = userID
	trackID := c.Param("id")

	exists, err := h.trackRepo.Exists(trackID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Track not found"})
	}

	var req models.TrackZoneRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if req.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Zone name is required"})
	}
	if req.PosX < 0 || req.PosX > 100 || req.PosY < 0 || req.PosY > 100 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Position must be between 0 and 100"})
	}

	zone, err := h.zoneRepo.Create(req.Name, req.Description, req.PosX, req.PosY, trackID, req.EventType)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	return c.JSON(http.StatusCreated, zone)
}

// PATCH /api/tracks/:id/zones/:zoneId
func (h *TrackZoneHandler) Update(c echo.Context) error {
	trackID := c.Param("id")
	zoneID := c.Param("zoneId")

	exists, err := h.trackRepo.Exists(trackID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Track not found"})
	}

	zone, err := h.zoneRepo.FindByID(zoneID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	if zone == nil || zone.TrackID != trackID {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Zone not found"})
	}

	var req models.ZoneUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	updated, err := h.zoneRepo.Update(zoneID, req.Name, req.Description)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	return c.JSON(http.StatusOK, updated)
}

// DELETE /api/tracks/:id/zones/:zoneId
func (h *TrackZoneHandler) Delete(c echo.Context) error {
	trackID := c.Param("id")
	zoneID := c.Param("zoneId")

	exists, err := h.trackRepo.Exists(trackID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Track not found"})
	}

	zone, err := h.zoneRepo.FindByID(zoneID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	if zone == nil || zone.TrackID != trackID {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Zone not found"})
	}

	if err := h.zoneRepo.Delete(zoneID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	return c.JSON(http.StatusOK, map[string]bool{"success": true})
}
