package tests

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTrackImages_ListEmpty(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "User", uniqueEmail("img"), "password123")
	trackID := app.createTestTrack(t, token)

	rec := app.doRequest(http.MethodGet, "/api/tracks/"+trackID+"/images", "", "")
	assert.Equal(t, http.StatusOK, rec.Code)
	result := parseJSONArray(t, rec)
	assert.Empty(t, result)
}

func TestTrackImages_Create(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "User", uniqueEmail("img"), "password123")
	trackID := app.createTestTrack(t, token)

	body := `{"url":"https://example.com/track.jpg","caption":"Front stretch"}`
	rec := app.doRequest(http.MethodPost, "/api/tracks/"+trackID+"/images", body, token)
	assert.Equal(t, http.StatusCreated, rec.Code)

	result := parseJSON(t, rec)
	assert.Equal(t, "https://example.com/track.jpg", result["url"])
}

func TestTrackImages_CreateInvalidURL(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "User", uniqueEmail("img"), "password123")
	trackID := app.createTestTrack(t, token)

	body := `{"url":"not-a-url"}`
	rec := app.doRequest(http.MethodPost, "/api/tracks/"+trackID+"/images", body, token)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestTrackImages_CreateNonExistentTrack(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "User", uniqueEmail("img"), "password123")

	body := `{"url":"https://example.com/track.jpg"}`
	rec := app.doRequest(http.MethodPost, "/api/tracks/nonexistent/images", body, token)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestTrackImages_Delete(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "User", uniqueEmail("img"), "password123")
	trackID := app.createTestTrack(t, token)

	// Create image
	body := `{"url":"https://example.com/track.jpg"}`
	rec := app.doRequest(http.MethodPost, "/api/tracks/"+trackID+"/images", body, token)
	require.Equal(t, http.StatusCreated, rec.Code)
	imageID := parseJSON(t, rec)["id"].(string)

	// Delete image
	rec = app.doRequest(http.MethodDelete, "/api/tracks/"+trackID+"/images?imageId="+imageID, "", token)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestTrackImages_DeleteNoImageID(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "User", uniqueEmail("img"), "password123")
	trackID := app.createTestTrack(t, token)

	rec := app.doRequest(http.MethodDelete, "/api/tracks/"+trackID+"/images", "", token)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestTrackImages_DeleteForbidden(t *testing.T) {
	app := setupTestApp(t)
	_, token1 := app.createTestUser(t, "User1", uniqueEmail("img1"), "password123")
	_, token2 := app.createTestUser(t, "User2", uniqueEmail("img2"), "password123")
	trackID := app.createTestTrack(t, token1) // User1 owns track

	// User1 creates image
	body := `{"url":"https://example.com/track.jpg"}`
	rec := app.doRequest(http.MethodPost, "/api/tracks/"+trackID+"/images", body, token1)
	require.Equal(t, http.StatusCreated, rec.Code)
	imageID := parseJSON(t, rec)["id"].(string)

	// User2 tries to delete (not track owner, not image uploader)
	rec = app.doRequest(http.MethodDelete, "/api/tracks/"+trackID+"/images?imageId="+imageID, "", token2)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

// ─── File Upload ────────────────────────────────────────────────────────────────

func TestUpload_Success(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "User", uniqueEmail("upload"), "password123")

	// Create multipart form with proper JPEG content type
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="file"; filename="test.jpg"`)
	h.Set("Content-Type", "image/jpeg")
	part, err := writer.CreatePart(h)
	require.NoError(t, err)
	// Write some fake JPEG data
	part.Write([]byte("\xFF\xD8\xFF\xE0test-image-data-here"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/upload", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", authHeader(token))

	rec := httptest.NewRecorder()
	app.e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	result := parseJSON(t, rec)
	assert.Contains(t, result["imageUrl"], "/uploads/tracks/")
	assert.Contains(t, result["imageUrl"], ".jpg")
}

func TestUpload_NoFile(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "User", uniqueEmail("upload"), "password123")

	rec := app.doRequest(http.MethodPost, "/api/upload", "", token)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUpload_RequiresAuth(t *testing.T) {
	app := setupTestApp(t)
	rec := app.doRequest(http.MethodPost, "/api/upload", "", "")
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
