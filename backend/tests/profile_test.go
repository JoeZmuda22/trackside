package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProfile_Get(t *testing.T) {
	app := setupTestApp(t)
	email := uniqueEmail("profile")
	_, token := app.createTestUser(t, "Profile User", email, "password123")

	rec := app.doRequest(http.MethodGet, "/api/profile", "", token)
	assert.Equal(t, http.StatusOK, rec.Code)

	result := parseJSON(t, rec)
	assert.Equal(t, "Profile User", result["name"])
	assert.Equal(t, email, result["email"])
	assert.Equal(t, "BEGINNER", result["experience"])
	assert.NotNil(t, result["cars"])
	assert.NotNil(t, result["_count"])
}

func TestProfile_GetRequiresAuth(t *testing.T) {
	app := setupTestApp(t)
	rec := app.doRequest(http.MethodGet, "/api/profile", "", "")
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestProfile_Update(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "User", uniqueEmail("profile"), "password123")

	body := `{"name":"Updated Name","experience":"ADVANCED"}`
	rec := app.doRequest(http.MethodPut, "/api/profile", body, token)
	assert.Equal(t, http.StatusOK, rec.Code)

	result := parseJSON(t, rec)
	assert.Equal(t, "Updated Name", result["name"])
	assert.Equal(t, "ADVANCED", result["experience"])
}

func TestProfile_UpdateValidation(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "User", uniqueEmail("profile"), "password123")

	tests := []struct {
		name string
		body string
	}{
		{"short name", `{"name":"A","experience":"BEGINNER"}`},
		{"invalid experience", `{"name":"Valid Name","experience":"EXPERT"}`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := app.doRequest(http.MethodPut, "/api/profile", tc.body, token)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		})
	}
}

func TestProfile_GetIncludesCarsAndCounts(t *testing.T) {
	app := setupTestApp(t)
	userID, token := app.createTestUser(t, "User", uniqueEmail("profile"), "password123")

	// Add a car
	app.createTestCar(t, "BMW", "M3", 2020, userID)

	// Add a track
	app.createTestTrack(t, token)

	rec := app.doRequest(http.MethodGet, "/api/profile", "", token)
	assert.Equal(t, http.StatusOK, rec.Code)

	result := parseJSON(t, rec)
	cars := result["cars"].([]interface{})
	assert.Len(t, cars, 1)

	counts := result["_count"].(map[string]interface{})
	assert.NotNil(t, counts["tracks"])
	require.Equal(t, float64(1), counts["tracks"])
}
