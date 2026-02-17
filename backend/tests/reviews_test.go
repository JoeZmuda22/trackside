package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReviews_Create(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "Reviewer", uniqueEmail("review"), "password123")
	trackID := app.createTestTrack(t, token)

	body := `{"rating":4,"content":"Great track!","conditions":"DRY"}`
	rec := app.doRequest(http.MethodPost, "/api/tracks/"+trackID+"/reviews", body, token)
	assert.Equal(t, http.StatusCreated, rec.Code)

	result := parseJSON(t, rec)
	assert.Equal(t, float64(4), result["rating"])
	assert.Equal(t, "DRY", result["conditions"])
}

func TestReviews_CreateRequiresAuth(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "User", uniqueEmail("review"), "password123")
	trackID := app.createTestTrack(t, token)

	body := `{"rating":4,"content":"Great track!","conditions":"DRY"}`
	rec := app.doRequest(http.MethodPost, "/api/tracks/"+trackID+"/reviews", body, "")
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestReviews_InvalidRating(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "User", uniqueEmail("review"), "password123")
	trackID := app.createTestTrack(t, token)

	tests := []struct {
		name string
		body string
	}{
		{"rating too low", `{"rating":0,"conditions":"DRY"}`},
		{"rating too high", `{"rating":6,"conditions":"DRY"}`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := app.doRequest(http.MethodPost, "/api/tracks/"+trackID+"/reviews", tc.body, token)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		})
	}
}

func TestReviews_InvalidConditions(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "User", uniqueEmail("review"), "password123")
	trackID := app.createTestTrack(t, token)

	body := `{"rating":3,"conditions":"SNOWY"}`
	rec := app.doRequest(http.MethodPost, "/api/tracks/"+trackID+"/reviews", body, token)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestReviews_NonExistentTrack(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "User", uniqueEmail("review"), "password123")

	body := `{"rating":4,"conditions":"DRY"}`
	rec := app.doRequest(http.MethodPost, "/api/tracks/nonexistent/reviews", body, token)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestReviews_VisibleInTrackDetail(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "User", uniqueEmail("review"), "password123")
	trackID := app.createTestTrack(t, token)

	body := `{"rating":5,"content":"Amazing!","conditions":"DRY"}`
	rec := app.doRequest(http.MethodPost, "/api/tracks/"+trackID+"/reviews", body, token)
	require.Equal(t, http.StatusCreated, rec.Code)

	// Get track detail
	rec = app.doRequest(http.MethodGet, "/api/tracks/"+trackID, "", "")
	assert.Equal(t, http.StatusOK, rec.Code)

	result := parseJSON(t, rec)
	reviews := result["reviews"].([]interface{})
	assert.Len(t, reviews, 1)

	review := reviews[0].(map[string]interface{})
	assert.Equal(t, float64(5), review["rating"])
	assert.Equal(t, "Amazing!", review["content"])

	// Should include author info
	author := review["author"].(map[string]interface{})
	assert.NotEmpty(t, author["id"])
}
