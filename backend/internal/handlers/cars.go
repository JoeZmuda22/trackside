package handlers

import (
	"net/http"

	"github.com/joezmuda/trackside-backend/internal/middleware"
	"github.com/joezmuda/trackside-backend/internal/models"
	"github.com/joezmuda/trackside-backend/internal/repository"
	"github.com/labstack/echo/v4"
)

type CarHandler struct {
	carRepo *repository.CarRepo
}

func NewCarHandler(carRepo *repository.CarRepo) *CarHandler {
	return &CarHandler{carRepo: carRepo}
}

// GET /api/cars
func (h *CarHandler) List(c echo.Context) error {
	userID := middleware.GetUserID(c)
	cars, err := h.carRepo.FindByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch cars"})
	}
	return c.JSON(http.StatusOK, cars)
}

// POST /api/cars
func (h *CarHandler) Create(c echo.Context) error {
	userID := middleware.GetUserID(c)
	var req models.CarRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if req.Make == "" || req.Model == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Validation failed"})
	}
	if req.Year < 1900 || req.Year > 2030 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Validation failed"})
	}

	car, err := h.carRepo.Create(req.Make, req.Model, req.Year, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	return c.JSON(http.StatusCreated, car)
}

// PUT /api/cars/:id
func (h *CarHandler) Update(c echo.Context) error {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	car, err := h.carRepo.FindByIDAndUser(id, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	if car == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Car not found"})
	}

	var req models.CarRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
	if req.Make == "" || req.Model == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Validation failed"})
	}
	if req.Year < 1900 || req.Year > 2030 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Validation failed"})
	}

	updated, err := h.carRepo.Update(id, req.Make, req.Model, req.Year)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	return c.JSON(http.StatusOK, updated)
}

// DELETE /api/cars/:id
func (h *CarHandler) Delete(c echo.Context) error {
	userID := middleware.GetUserID(c)
	id := c.Param("id")

	car, err := h.carRepo.FindByIDAndUser(id, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	if car == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Car not found"})
	}

	if err := h.carRepo.Delete(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	return c.JSON(http.StatusOK, map[string]bool{"success": true})
}
