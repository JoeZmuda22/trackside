package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/joezmuda/trackside-backend/internal/middleware"
	"github.com/joezmuda/trackside-backend/internal/models"
	"github.com/joezmuda/trackside-backend/internal/repository"
	"github.com/labstack/echo/v4"
)

type AdminHandler struct {
	trackRepo *repository.TrackRepo
	userRepo  *repository.UserRepo
	dataDir   string
}

func NewAdminHandler(trackRepo *repository.TrackRepo, userRepo *repository.UserRepo, dataDir string) *AdminHandler {
	return &AdminHandler{trackRepo: trackRepo, userRepo: userRepo, dataDir: dataDir}
}

// POST /api/admin/sync-tracks
func (h *AdminHandler) SyncTracks(c echo.Context) error {
	email := middleware.GetUserEmail(c)

	if !strings.Contains(email, "admin") && email != "system@trackside.local" {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Admin access required"})
	}

	tracksPath := filepath.Join(h.dataDir, "usa-tracks.json")
	data, err := os.ReadFile(tracksPath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Sync failed",
			"details": "Could not read tracks data file",
		})
	}

	var file struct {
		Tracks []models.ImportedTrack `json:"tracks"`
	}
	if err := json.Unmarshal(data, &file); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Sync failed",
			"details": "Invalid JSON in tracks data file",
		})
	}

	systemUser, err := h.userRepo.FindOrCreateSystem()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Sync failed",
			"details": err.Error(),
		})
	}

	var created, updated, failed int
	var errors []string

	for _, trackData := range file.Tracks {
		existing, _ := h.trackRepo.FindByNameAndLocation(trackData.Name, trackData.Location)
		err := h.trackRepo.UpsertImported(trackData, systemUser.ID)
		if err != nil {
			failed++
			errors = append(errors, trackData.Name+": "+err.Error())
		} else if existing != nil {
			updated++
		} else {
			created++
		}
	}

	resp := models.SyncTracksResponse{
		Status: "success",
	}
	resp.Summary.Total = len(file.Tracks)
	resp.Summary.Created = created
	resp.Summary.Updated = updated
	resp.Summary.Failed = failed
	if len(errors) > 0 {
		resp.Errors = errors
	}

	return c.JSON(http.StatusOK, resp)
}
