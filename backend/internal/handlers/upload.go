package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/joezmuda/trackside-backend/internal/middleware"
	"github.com/labstack/echo/v4"
)

type UploadHandler struct {
	uploadDir string
}

func NewUploadHandler(uploadDir string) *UploadHandler {
	return &UploadHandler{uploadDir: uploadDir}
}

// POST /api/upload
func (h *UploadHandler) Upload(c echo.Context) error {
	_ = middleware.GetUserID(c) // Auth is enforced by middleware

	file, err := c.FormFile("file")
	if err != nil || file == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "No file provided"})
	}

	// Validate MIME type
	allowedTypes := map[string]bool{
		"image/jpeg":    true,
		"image/png":     true,
		"image/webp":    true,
		"image/svg+xml": true,
	}
	if !allowedTypes[file.Header.Get("Content-Type")] {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid file type. Allowed: JPEG, PNG, WebP, SVG"})
	}

	// Validate file size (max 10MB)
	if file.Size > 10*1024*1024 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "File too large. Maximum size is 10MB"})
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Upload failed"})
	}
	defer src.Close()

	// Ensure upload directory exists
	tracksDir := filepath.Join(h.uploadDir, "tracks")
	if err := os.MkdirAll(tracksDir, 0755); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Upload failed"})
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	if ext == "" {
		// Fallback based on content type
		switch file.Header.Get("Content-Type") {
		case "image/jpeg":
			ext = ".jpg"
		case "image/png":
			ext = ".png"
		case "image/webp":
			ext = ".webp"
		case "image/svg+xml":
			ext = ".svg"
		}
	}
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	destPath := filepath.Join(tracksDir, filename)

	dst, err := os.Create(destPath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Upload failed"})
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Upload failed"})
	}

	imageURL := fmt.Sprintf("/uploads/tracks/%s", filename)
	return c.JSON(http.StatusCreated, map[string]string{"imageUrl": imageURL})
}

// ServeUploads serves static files from the upload directory.
// This is registered as a group handler for /uploads/*.
func ServeUploads(uploadDir string) echo.HandlerFunc {
	return func(c echo.Context) error {
		// The path param captures everything after /uploads/
		reqPath := c.Param("*")
		// Prevent directory traversal
		reqPath = strings.ReplaceAll(reqPath, "..", "")
		filePath := filepath.Join(uploadDir, reqPath)
		return c.File(filePath)
	}
}
