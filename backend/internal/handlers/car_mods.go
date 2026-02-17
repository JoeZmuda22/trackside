package handlers

import (
	"net/http"

	"github.com/joezmuda/trackside-backend/internal/middleware"
	"github.com/joezmuda/trackside-backend/internal/models"
	"github.com/joezmuda/trackside-backend/internal/repository"
	"github.com/labstack/echo/v4"
)

type CarModHandler struct {
	carRepo *repository.CarRepo
}

func NewCarModHandler(carRepo *repository.CarRepo) *CarModHandler {
	return &CarModHandler{carRepo: carRepo}
}

// POST /api/cars/:id/mods
func (h *CarModHandler) Create(c echo.Context) error {
	userID := middleware.GetUserID(c)
	carID := c.Param("id")

	car, err := h.carRepo.FindByIDAndUser(carID, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	if car == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Car not found"})
	}

	var req models.CarModRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if req.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Mod name is required"})
	}
	if !models.ValidModCategory(req.Category) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid mod category"})
	}

	mod, err := h.carRepo.CreateMod(req.Name, req.Category, req.Notes, carID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	return c.JSON(http.StatusCreated, mod)
}

// DELETE /api/cars/:id/mods/:modId
func (h *CarModHandler) Delete(c echo.Context) error {
	userID := middleware.GetUserID(c)
	carID := c.Param("id")
	modID := c.Param("modId")

	car, err := h.carRepo.FindByIDAndUser(carID, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	if car == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Car not found"})
	}

	mod, err := h.carRepo.FindMod(modID, carID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	if mod == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Mod not found"})
	}

	if err := h.carRepo.DeleteMod(modID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	return c.JSON(http.StatusOK, map[string]bool{"success": true})
}
