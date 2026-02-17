package handlers

import (
	"net/http"

	"github.com/joezmuda/trackside-backend/internal/middleware"
	"github.com/joezmuda/trackside-backend/internal/models"
	"github.com/joezmuda/trackside-backend/internal/repository"
	"github.com/labstack/echo/v4"
)

type LapbookHandler struct {
	lapbookRepo *repository.LapbookRepo
	carRepo     *repository.CarRepo
	trackRepo   *repository.TrackRepo
}

func NewLapbookHandler(lapbookRepo *repository.LapbookRepo, carRepo *repository.CarRepo, trackRepo *repository.TrackRepo) *LapbookHandler {
	return &LapbookHandler{lapbookRepo: lapbookRepo, carRepo: carRepo, trackRepo: trackRepo}
}

// GET /api/lapbook
func (h *LapbookHandler) List(c echo.Context) error {
	userID := middleware.GetUserID(c)
	trackID := c.QueryParam("trackId")
	eventType := c.QueryParam("eventType")
	carID := c.QueryParam("carId")

	records, err := h.lapbookRepo.List(userID, trackID, eventType, carID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch lap records"})
	}
	return c.JSON(http.StatusOK, records)
}

// POST /api/lapbook
func (h *LapbookHandler) Create(c echo.Context) error {
	userID := middleware.GetUserID(c)

	var req models.LapRecordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if req.LapTime == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Lap time is required"})
	}
	if !models.ValidDrivingCondition(req.Conditions) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid conditions"})
	}
	if req.TrackID == "" || req.CarID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Track and car are required"})
	}

	// Verify car belongs to user
	owns, err := h.carRepo.ExistsForUser(req.CarID, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	if !owns {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Car not found"})
	}

	// Verify track exists
	exists, err := h.trackRepo.Exists(req.TrackID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Track not found"})
	}

	record, err := h.lapbookRepo.Create(req, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	return c.JSON(http.StatusCreated, record)
}

// DELETE /api/lapbook/:id
func (h *LapbookHandler) Delete(c echo.Context) error {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	record, err := h.lapbookRepo.FindByIDAndDriver(id, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	if record == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Record not found"})
	}

	if err := h.lapbookRepo.Delete(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	return c.JSON(http.StatusOK, map[string]bool{"success": true})
}
