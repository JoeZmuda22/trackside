package handlers

import (
	"net/http"

	"github.com/joezmuda/trackside-backend/internal/middleware"
	"github.com/joezmuda/trackside-backend/internal/models"
	"github.com/joezmuda/trackside-backend/internal/repository"
	"github.com/labstack/echo/v4"
)

type TrackHandler struct {
	trackRepo *repository.TrackRepo
}

func NewTrackHandler(trackRepo *repository.TrackRepo) *TrackHandler {
	return &TrackHandler{trackRepo: trackRepo}
}

// GET /api/tracks
func (h *TrackHandler) List(c echo.Context) error {
	search := c.QueryParam("search")
	eventType := c.QueryParam("eventType")
	state := c.QueryParam("state")

	tracks, err := h.trackRepo.List(search, eventType, state)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Failed to fetch tracks",
			"details": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, tracks)
}

// POST /api/tracks
func (h *TrackHandler) Create(c echo.Context) error {
	userID := middleware.GetUserID(c)

	var req models.TrackRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if len(req.Name) < 2 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Track name is required"})
	}
	if len(req.Location) < 2 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Location is required"})
	}
	if len(req.EventTypes) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Select at least one event type"})
	}
	for _, et := range req.EventTypes {
		if !models.ValidEventType(et) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid event type: " + et})
		}
	}

	track, err := h.trackRepo.Create(req, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	return c.JSON(http.StatusCreated, track)
}

// GET /api/tracks/:id
func (h *TrackHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	eventType := c.QueryParam("eventType")

	detail, err := h.trackRepo.GetDetail(id, eventType)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	if detail == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Track not found"})
	}
	return c.JSON(http.StatusOK, detail)
}

// PATCH /api/tracks/:id
func (h *TrackHandler) Update(c echo.Context) error {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	track, err := h.trackRepo.FindByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	if track == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Track not found"})
	}
	if track.UploadedByID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "You can only edit your own tracks"})
	}

	var req models.TrackPatchRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	updated, err := h.trackRepo.Update(id, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	return c.JSON(http.StatusOK, updated)
}
