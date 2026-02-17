package handlers

import (
	"net/http"

	"github.com/joezmuda/trackside-backend/internal/middleware"
	"github.com/joezmuda/trackside-backend/internal/models"
	"github.com/joezmuda/trackside-backend/internal/repository"
	"github.com/labstack/echo/v4"
)

type TrackReviewHandler struct {
	trackRepo  *repository.TrackRepo
	reviewRepo *repository.ReviewRepo
}

func NewTrackReviewHandler(trackRepo *repository.TrackRepo, reviewRepo *repository.ReviewRepo) *TrackReviewHandler {
	return &TrackReviewHandler{trackRepo: trackRepo, reviewRepo: reviewRepo}
}

// POST /api/tracks/:id/reviews
func (h *TrackReviewHandler) Create(c echo.Context) error {
	userID := middleware.GetUserID(c)
	trackID := c.Param("id")

	exists, err := h.trackRepo.Exists(trackID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Track not found"})
	}

	var req models.TrackReviewRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if req.Rating < 1 || req.Rating > 5 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Rating must be between 1 and 5"})
	}
	if !models.ValidDrivingCondition(req.Conditions) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid conditions"})
	}

	review, err := h.reviewRepo.Create(req.Rating, req.Content, req.Conditions, trackID, req.TrackEventID, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	return c.JSON(http.StatusCreated, review)
}
