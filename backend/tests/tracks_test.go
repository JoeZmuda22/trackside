package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTracks_ListEmpty(t *testing.T) {
	app := setupTestApp(t)
	rec := app.doRequest(http.MethodGet, "/api/tracks", "", "")
	assert.Equal(t, http.StatusOK, rec.Code)
	result := parseJSONArray(t, rec)
	assert.Empty(t, result)
}

func TestTracks_ListIsPublic(t *testing.T) {
	app := setupTestApp(t)
	// No token needed
	rec := app.doRequest(http.MethodGet, "/api/tracks", "", "")
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestTracks_Create(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "User", uniqueEmail("tracks"), "password123")

	body := `{"name":"Laguna Seca","location":"Monterey, CA","eventTypes":["ROADCOURSE"]}`
	rec := app.doRequest(http.MethodPost, "/api/tracks", body, token)
	assert.Equal(t, http.StatusCreated, rec.Code)

	result := parseJSON(t, rec)
	assert.Equal(t, "Laguna Seca", result["name"])
	assert.Equal(t, "Monterey, CA", result["location"])
	assert.NotEmpty(t, result["id"])
}

func TestTracks_CreateRequiresAuth(t *testing.T) {
	app := setupTestApp(t)
	body := `{"name":"Test","location":"City","eventTypes":["ROADCOURSE"]}`
	rec := app.doRequest(http.MethodPost, "/api/tracks", body, "")
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestTracks_CreateValidation(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "User", uniqueEmail("tracks"), "password123")

	tests := []struct {
		name string
		body string
	}{
		{"short name", `{"name":"A","location":"City","eventTypes":["ROADCOURSE"]}`},
		{"short location", `{"name":"Track","location":"A","eventTypes":["ROADCOURSE"]}`},
		{"no event types", `{"name":"Track","location":"City","eventTypes":[]}`},
		{"invalid event type", `{"name":"Track","location":"City","eventTypes":["INVALID"]}`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := app.doRequest(http.MethodPost, "/api/tracks", tc.body, token)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		})
	}
}

func TestTracks_GetByID(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "User", uniqueEmail("tracks"), "password123")
	trackID := app.createTestTrack(t, token)

	rec := app.doRequest(http.MethodGet, "/api/tracks/"+trackID, "", "")
	assert.Equal(t, http.StatusOK, rec.Code)

	result := parseJSON(t, rec)
	assert.Equal(t, trackID, result["id"])
	assert.Equal(t, "Test Track", result["name"])
	// Should include nested data
	assert.NotNil(t, result["events"])
	assert.NotNil(t, result["zones"])
	assert.NotNil(t, result["reviews"])
}

func TestTracks_GetByID_NotFound(t *testing.T) {
	app := setupTestApp(t)
	rec := app.doRequest(http.MethodGet, "/api/tracks/nonexistent", "", "")
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestTracks_Update(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "User", uniqueEmail("tracks"), "password123")
	trackID := app.createTestTrack(t, token)

	patchBody := `{"name":"Updated Track","description":"A great track"}`
	rec := app.doRequest(http.MethodPatch, "/api/tracks/"+trackID, patchBody, token)
	assert.Equal(t, http.StatusOK, rec.Code)

	result := parseJSON(t, rec)
	assert.Equal(t, "Updated Track", result["name"])
}

func TestTracks_UpdateNotOwned(t *testing.T) {
	app := setupTestApp(t)
	_, token1 := app.createTestUser(t, "User1", uniqueEmail("tracks1"), "password123")
	_, token2 := app.createTestUser(t, "User2", uniqueEmail("tracks2"), "password123")
	trackID := app.createTestTrack(t, token1)

	patchBody := `{"name":"Hacked Track"}`
	rec := app.doRequest(http.MethodPatch, "/api/tracks/"+trackID, patchBody, token2)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestTracks_ListWithFilter(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "User", uniqueEmail("tracks"), "password123")

	// Create a track
	body := `{"name":"Sebring International","location":"Sebring, FL","eventTypes":["ROADCOURSE"]}`
	rec := app.doRequest(http.MethodPost, "/api/tracks", body, token)
	require.Equal(t, http.StatusCreated, rec.Code)

	// Search by name
	rec = app.doRequest(http.MethodGet, "/api/tracks?search=Sebring", "", "")
	assert.Equal(t, http.StatusOK, rec.Code)
	results := parseJSONArray(t, rec)
	assert.Len(t, results, 1)

	// Search with no match
	rec = app.doRequest(http.MethodGet, "/api/tracks?search=Nonexistent", "", "")
	assert.Equal(t, http.StatusOK, rec.Code)
	results = parseJSONArray(t, rec)
	assert.Empty(t, results)
}
