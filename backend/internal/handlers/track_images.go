package handlers

import (
	"net/http"
	"net/url"

	"github.com/joezmuda/trackside-backend/internal/middleware"
	"github.com/joezmuda/trackside-backend/internal/models"
	"github.com/joezmuda/trackside-backend/internal/repository"
	"github.com/labstack/echo/v4"
)

type TrackImageHandler struct {
	trackRepo *repository.TrackRepo
}

func NewTrackImageHandler(trackRepo *repository.TrackRepo) *TrackImageHandler {
	return &TrackImageHandler{trackRepo: trackRepo}
}

// GET /api/tracks/:id/images
func (h *TrackImageHandler) List(c echo.Context) error {
	trackID := c.Param("id")
	images, err := h.trackRepo.GetImages(trackID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	return c.JSON(http.StatusOK, images)
}

// POST /api/tracks/:id/images
func (h *TrackImageHandler) Create(c echo.Context) error {
	userID := middleware.GetUserID(c)
	trackID := c.Param("id")

	var req models.TrackImageRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Validate URL
	if _, err := url.ParseRequestURI(req.URL); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid image URL"})
	}

	exists, err := h.trackRepo.Exists(trackID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Track not found"})
	}

	img, err := h.trackRepo.CreateImage(req.URL, req.Caption, trackID, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	return c.JSON(http.StatusCreated, img)
}

// DELETE /api/tracks/:id/images
func (h *TrackImageHandler) Delete(c echo.Context) error {
	userID := middleware.GetUserID(c)
	trackID := c.Param("id")
	imageID := c.QueryParam("imageId")

	if imageID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Image ID is required"})
	}

	image, err := h.trackRepo.FindImage(imageID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	if image == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Image not found"})
	}

	track, err := h.trackRepo.FindByID(trackID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	isTrackOwner := track != nil && track.UploadedByID == userID
	isImageUploader := image.UploadedByID == userID

	if !isTrackOwner && !isImageUploader {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Forbidden"})
	}

	if err := h.trackRepo.DeleteImage(imageID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	return c.JSON(http.StatusOK, map[string]bool{"success": true})
}
