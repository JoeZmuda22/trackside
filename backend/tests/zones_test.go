package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestZones_Create(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "User", uniqueEmail("zone"), "password123")
	trackID := app.createTestTrack(t, token)

	body := `{"name":"Turn 1","description":"Tight left-hander","posX":25.5,"posY":50.0}`
	rec := app.doRequest(http.MethodPost, "/api/tracks/"+trackID+"/zones", body, token)
	assert.Equal(t, http.StatusCreated, rec.Code)

	result := parseJSON(t, rec)
	assert.Equal(t, "Turn 1", result["name"])
	assert.Equal(t, 25.5, result["posX"])
	assert.Equal(t, 50.0, result["posY"])
	assert.NotEmpty(t, result["id"])
}

func TestZones_CreateValidation(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "User", uniqueEmail("zone"), "password123")
	trackID := app.createTestTrack(t, token)

	tests := []struct {
		name string
		body string
	}{
		{"missing name", `{"name":"","posX":25,"posY":50}`},
		{"posX out of range", `{"name":"Turn","posX":101,"posY":50}`},
		{"posY out of range", `{"name":"Turn","posX":50,"posY":-1}`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := app.doRequest(http.MethodPost, "/api/tracks/"+trackID+"/zones", tc.body, token)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		})
	}
}

func TestZones_CreateNonExistentTrack(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "User", uniqueEmail("zone"), "password123")

	body := `{"name":"Turn 1","posX":25,"posY":50}`
	rec := app.doRequest(http.MethodPost, "/api/tracks/nonexistent/zones", body, token)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestZones_Update(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "User", uniqueEmail("zone"), "password123")
	trackID := app.createTestTrack(t, token)

	// Create zone
	body := `{"name":"Turn 1","posX":25,"posY":50}`
	rec := app.doRequest(http.MethodPost, "/api/tracks/"+trackID+"/zones", body, token)
	require.Equal(t, http.StatusCreated, rec.Code)
	zoneID := parseJSON(t, rec)["id"].(string)

	// Update zone
	updateBody := `{"name":"Turn 1 - Updated","description":"Better description"}`
	rec = app.doRequest(http.MethodPatch, "/api/tracks/"+trackID+"/zones/"+zoneID, updateBody, token)
	assert.Equal(t, http.StatusOK, rec.Code)

	result := parseJSON(t, rec)
	assert.Equal(t, "Turn 1 - Updated", result["name"])
}

func TestZones_UpdateNonExistentZone(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "User", uniqueEmail("zone"), "password123")
	trackID := app.createTestTrack(t, token)

	body := `{"name":"Updated"}`
	rec := app.doRequest(http.MethodPatch, "/api/tracks/"+trackID+"/zones/nonexistent", body, token)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestZones_Delete(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "User", uniqueEmail("zone"), "password123")
	trackID := app.createTestTrack(t, token)

	body := `{"name":"Turn 1","posX":25,"posY":50}`
	rec := app.doRequest(http.MethodPost, "/api/tracks/"+trackID+"/zones", body, token)
	require.Equal(t, http.StatusCreated, rec.Code)
	zoneID := parseJSON(t, rec)["id"].(string)

	rec = app.doRequest(http.MethodDelete, "/api/tracks/"+trackID+"/zones/"+zoneID, "", token)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestZones_DeleteNonExistent(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "User", uniqueEmail("zone"), "password123")
	trackID := app.createTestTrack(t, token)

	rec := app.doRequest(http.MethodDelete, "/api/tracks/"+trackID+"/zones/nonexistent", "", token)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestZones_VisibleInTrackDetail(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "User", uniqueEmail("zone"), "password123")
	trackID := app.createTestTrack(t, token)

	body := `{"name":"Corkscrew","posX":75,"posY":30}`
	rec := app.doRequest(http.MethodPost, "/api/tracks/"+trackID+"/zones", body, token)
	require.Equal(t, http.StatusCreated, rec.Code)

	// Get track detail
	rec = app.doRequest(http.MethodGet, "/api/tracks/"+trackID, "", "")
	assert.Equal(t, http.StatusOK, rec.Code)

	result := parseJSON(t, rec)
	zones := result["zones"].([]interface{})
	assert.Len(t, zones, 1)
	assert.Equal(t, "Corkscrew", zones[0].(map[string]interface{})["name"])
}

// ─── Zone Tips ──────────────────────────────────────────────────────────────────

func TestZoneTips_Create(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "Tipper", uniqueEmail("tip"), "password123")
	trackID := app.createTestTrack(t, token)

	// Create zone
	zoneBody := `{"name":"Turn 1","posX":25,"posY":50}`
	rec := app.doRequest(http.MethodPost, "/api/tracks/"+trackID+"/zones", zoneBody, token)
	require.Equal(t, http.StatusCreated, rec.Code)
	zoneID := parseJSON(t, rec)["id"].(string)

	// Create tip
	tipBody := `{"content":"Brake before the curb","conditions":"DRY"}`
	rec = app.doRequest(http.MethodPost, "/api/tracks/"+trackID+"/zones/"+zoneID+"/tips", tipBody, token)
	assert.Equal(t, http.StatusCreated, rec.Code)

	result := parseJSON(t, rec)
	assert.Equal(t, "Brake before the curb", result["content"])
}

func TestZoneTips_CreateEmptyContent(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "User", uniqueEmail("tip"), "password123")
	trackID := app.createTestTrack(t, token)

	zoneBody := `{"name":"Turn 1","posX":25,"posY":50}`
	rec := app.doRequest(http.MethodPost, "/api/tracks/"+trackID+"/zones", zoneBody, token)
	require.Equal(t, http.StatusCreated, rec.Code)
	zoneID := parseJSON(t, rec)["id"].(string)

	tipBody := `{"content":""}`
	rec = app.doRequest(http.MethodPost, "/api/tracks/"+trackID+"/zones/"+zoneID+"/tips", tipBody, token)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestZoneTips_InvalidConditions(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "User", uniqueEmail("tip"), "password123")
	trackID := app.createTestTrack(t, token)

	zoneBody := `{"name":"Turn 1","posX":25,"posY":50}`
	rec := app.doRequest(http.MethodPost, "/api/tracks/"+trackID+"/zones", zoneBody, token)
	require.Equal(t, http.StatusCreated, rec.Code)
	zoneID := parseJSON(t, rec)["id"].(string)

	tipBody := `{"content":"A tip","conditions":"SNOWY"}`
	rec = app.doRequest(http.MethodPost, "/api/tracks/"+trackID+"/zones/"+zoneID+"/tips", tipBody, token)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestZoneTips_NonExistentZone(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "User", uniqueEmail("tip"), "password123")
	trackID := app.createTestTrack(t, token)

	tipBody := `{"content":"A tip"}`
	rec := app.doRequest(http.MethodPost, "/api/tracks/"+trackID+"/zones/nonexistent/tips", tipBody, token)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
