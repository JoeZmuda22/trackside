package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLapbook_ListEmpty(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "User", uniqueEmail("lap"), "password123")

	rec := app.doRequest(http.MethodGet, "/api/lapbook", "", token)
	assert.Equal(t, http.StatusOK, rec.Code)
	result := parseJSONArray(t, rec)
	assert.Empty(t, result)
}

func TestLapbook_Create(t *testing.T) {
	app := setupTestApp(t)
	userID, token := app.createTestUser(t, "User", uniqueEmail("lap"), "password123")
	carID := app.createTestCar(t, "BMW", "M3", 2020, userID)
	trackID := app.createTestTrack(t, token)

	body := `{"lapTime":"1:42.5","conditions":"DRY","trackId":"` + trackID + `","carId":"` + carID + `"}`
	rec := app.doRequest(http.MethodPost, "/api/lapbook", body, token)
	assert.Equal(t, http.StatusCreated, rec.Code)

	result := parseJSON(t, rec)
	assert.Equal(t, "1:42.5", result["lapTime"])
	assert.Equal(t, "DRY", result["conditions"])
}

func TestLapbook_CreateValidation(t *testing.T) {
	app := setupTestApp(t)
	userID, token := app.createTestUser(t, "User", uniqueEmail("lap"), "password123")
	carID := app.createTestCar(t, "BMW", "M3", 2020, userID)
	trackID := app.createTestTrack(t, token)

	tests := []struct {
		name string
		body string
	}{
		{"missing lap time", `{"lapTime":"","conditions":"DRY","trackId":"` + trackID + `","carId":"` + carID + `"}`},
		{"invalid conditions", `{"lapTime":"1:30","conditions":"SNOW","trackId":"` + trackID + `","carId":"` + carID + `"}`},
		{"missing track", `{"lapTime":"1:30","conditions":"DRY","trackId":"","carId":"` + carID + `"}`},
		{"missing car", `{"lapTime":"1:30","conditions":"DRY","trackId":"` + trackID + `","carId":""}`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := app.doRequest(http.MethodPost, "/api/lapbook", tc.body, token)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		})
	}
}

func TestLapbook_CreateNotOwnedCar(t *testing.T) {
	app := setupTestApp(t)
	userID1, _ := app.createTestUser(t, "User1", uniqueEmail("lap1"), "password123")
	_, token2 := app.createTestUser(t, "User2", uniqueEmail("lap2"), "password123")
	carID := app.createTestCar(t, "BMW", "M3", 2020, userID1) // User1's car
	trackID := app.createTestTrack(t, token2)

	// User2 tries to use User1's car
	body := `{"lapTime":"1:42.5","conditions":"DRY","trackId":"` + trackID + `","carId":"` + carID + `"}`
	rec := app.doRequest(http.MethodPost, "/api/lapbook", body, token2)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestLapbook_CreateNonexistentTrack(t *testing.T) {
	app := setupTestApp(t)
	userID, token := app.createTestUser(t, "User", uniqueEmail("lap"), "password123")
	carID := app.createTestCar(t, "BMW", "M3", 2020, userID)

	body := `{"lapTime":"1:42.5","conditions":"DRY","trackId":"nonexistent","carId":"` + carID + `"}`
	rec := app.doRequest(http.MethodPost, "/api/lapbook", body, token)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestLapbook_ListAfterCreate(t *testing.T) {
	app := setupTestApp(t)
	userID, token := app.createTestUser(t, "User", uniqueEmail("lap"), "password123")
	carID := app.createTestCar(t, "BMW", "M3", 2020, userID)
	trackID := app.createTestTrack(t, token)

	body := `{"lapTime":"1:42.5","conditions":"DRY","trackId":"` + trackID + `","carId":"` + carID + `"}`
	rec := app.doRequest(http.MethodPost, "/api/lapbook", body, token)
	require.Equal(t, http.StatusCreated, rec.Code)

	rec = app.doRequest(http.MethodGet, "/api/lapbook", "", token)
	assert.Equal(t, http.StatusOK, rec.Code)
	result := parseJSONArray(t, rec)
	assert.Len(t, result, 1)
	assert.Equal(t, "1:42.5", result[0]["lapTime"])
	// Should include nested track and car
	assert.NotNil(t, result[0]["track"])
	assert.NotNil(t, result[0]["car"])
}

func TestLapbook_Delete(t *testing.T) {
	app := setupTestApp(t)
	userID, token := app.createTestUser(t, "User", uniqueEmail("lap"), "password123")
	carID := app.createTestCar(t, "BMW", "M3", 2020, userID)
	trackID := app.createTestTrack(t, token)

	body := `{"lapTime":"1:42.5","conditions":"DRY","trackId":"` + trackID + `","carId":"` + carID + `"}`
	rec := app.doRequest(http.MethodPost, "/api/lapbook", body, token)
	require.Equal(t, http.StatusCreated, rec.Code)
	recordID := parseJSON(t, rec)["id"].(string)

	rec = app.doRequest(http.MethodDelete, "/api/lapbook/"+recordID, "", token)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify gone
	rec = app.doRequest(http.MethodGet, "/api/lapbook", "", token)
	result := parseJSONArray(t, rec)
	assert.Empty(t, result)
}

func TestLapbook_DeleteNotOwned(t *testing.T) {
	app := setupTestApp(t)
	userID, token1 := app.createTestUser(t, "User1", uniqueEmail("lap1"), "password123")
	_, token2 := app.createTestUser(t, "User2", uniqueEmail("lap2"), "password123")
	carID := app.createTestCar(t, "BMW", "M3", 2020, userID)
	trackID := app.createTestTrack(t, token1)

	body := `{"lapTime":"1:42.5","conditions":"DRY","trackId":"` + trackID + `","carId":"` + carID + `"}`
	rec := app.doRequest(http.MethodPost, "/api/lapbook", body, token1)
	require.Equal(t, http.StatusCreated, rec.Code)
	recordID := parseJSON(t, rec)["id"].(string)

	// User2 tries to delete
	rec = app.doRequest(http.MethodDelete, "/api/lapbook/"+recordID, "", token2)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestLapbook_WithTelemetry(t *testing.T) {
	app := setupTestApp(t)
	userID, token := app.createTestUser(t, "User", uniqueEmail("lap"), "password123")
	carID := app.createTestCar(t, "BMW", "M3", 2020, userID)
	trackID := app.createTestTrack(t, token)

	body := `{
		"lapTime":"1:42.5",
		"conditions":"DRY",
		"trackId":"` + trackID + `",
		"carId":"` + carID + `",
		"tirePressureFL": 32.5,
		"tirePressureFR": 32.5,
		"tirePressureRL": 30.0,
		"tirePressureRR": 30.0,
		"fuelLevel": 75.0,
		"notes": "Best lap of the day"
	}`
	rec := app.doRequest(http.MethodPost, "/api/lapbook", body, token)
	assert.Equal(t, http.StatusCreated, rec.Code)

	result := parseJSON(t, rec)
	assert.Equal(t, 32.5, result["tirePressureFL"])
	assert.Equal(t, 75.0, result["fuelLevel"])
	assert.Equal(t, "Best lap of the day", result["notes"])
}
