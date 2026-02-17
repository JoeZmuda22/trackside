package handlers

import (
	"net/http"

	"github.com/joezmuda/trackside-backend/internal/middleware"
	"github.com/joezmuda/trackside-backend/internal/models"
	"github.com/joezmuda/trackside-backend/internal/repository"
	"github.com/labstack/echo/v4"
)

type ProfileHandler struct {
	userRepo *repository.UserRepo
}

func NewProfileHandler(userRepo *repository.UserRepo) *ProfileHandler {
	return &ProfileHandler{userRepo: userRepo}
}

// GET /api/profile
func (h *ProfileHandler) Get(c echo.Context) error {
	userID := middleware.GetUserID(c)

	profile, err := h.userRepo.GetProfile(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	if profile == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}
	return c.JSON(http.StatusOK, profile)
}

// PUT /api/profile
func (h *ProfileHandler) Update(c echo.Context) error {
	userID := middleware.GetUserID(c)

	var req models.ProfileUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if len(req.Name) < 2 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Name must be at least 2 characters"})
	}
	if !models.ValidExperienceLevel(req.Experience) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid experience level"})
	}

	user, err := h.userRepo.UpdateProfile(userID, req.Name, req.Experience)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"experience": user.Experience,
	})
}
