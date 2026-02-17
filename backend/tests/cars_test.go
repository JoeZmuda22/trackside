package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCars_ListEmpty(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "User", uniqueEmail("cars"), "password123")

	rec := app.doRequest(http.MethodGet, "/api/cars", "", token)
	assert.Equal(t, http.StatusOK, rec.Code)
	result := parseJSONArray(t, rec)
	assert.Empty(t, result)
}

func TestCars_Create(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "User", uniqueEmail("cars"), "password123")

	body := `{"make":"BMW","model":"M3","year":2020}`
	rec := app.doRequest(http.MethodPost, "/api/cars", body, token)
	assert.Equal(t, http.StatusCreated, rec.Code)

	result := parseJSON(t, rec)
	assert.Equal(t, "BMW", result["make"])
	assert.Equal(t, "M3", result["model"])
	assert.Equal(t, float64(2020), result["year"])
	assert.NotEmpty(t, result["id"])
}

func TestCars_CreateValidation(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "User", uniqueEmail("cars"), "password123")

	tests := []struct {
		name string
		body string
	}{
		{"missing make", `{"make":"","model":"M3","year":2020}`},
		{"missing model", `{"make":"BMW","model":"","year":2020}`},
		{"year too low", `{"make":"BMW","model":"M3","year":1800}`},
		{"year too high", `{"make":"BMW","model":"M3","year":2050}`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := app.doRequest(http.MethodPost, "/api/cars", tc.body, token)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		})
	}
}

func TestCars_ListAfterCreate(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "User", uniqueEmail("cars"), "password123")

	body := `{"make":"Honda","model":"Civic","year":2019}`
	rec := app.doRequest(http.MethodPost, "/api/cars", body, token)
	require.Equal(t, http.StatusCreated, rec.Code)

	rec = app.doRequest(http.MethodGet, "/api/cars", "", token)
	assert.Equal(t, http.StatusOK, rec.Code)
	result := parseJSONArray(t, rec)
	assert.Len(t, result, 1)
	assert.Equal(t, "Honda", result[0]["make"])
}

func TestCars_Update(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "User", uniqueEmail("cars"), "password123")

	body := `{"make":"BMW","model":"M3","year":2020}`
	rec := app.doRequest(http.MethodPost, "/api/cars", body, token)
	require.Equal(t, http.StatusCreated, rec.Code)
	created := parseJSON(t, rec)
	carID := created["id"].(string)

	updateBody := `{"make":"BMW","model":"M4","year":2021}`
	rec = app.doRequest(http.MethodPut, "/api/cars/"+carID, updateBody, token)
	assert.Equal(t, http.StatusOK, rec.Code)
	updated := parseJSON(t, rec)
	assert.Equal(t, "M4", updated["model"])
	assert.Equal(t, float64(2021), updated["year"])
}

func TestCars_UpdateNotOwned(t *testing.T) {
	app := setupTestApp(t)
	_, token1 := app.createTestUser(t, "User1", uniqueEmail("cars1"), "password123")
	_, token2 := app.createTestUser(t, "User2", uniqueEmail("cars2"), "password123")

	body := `{"make":"BMW","model":"M3","year":2020}`
	rec := app.doRequest(http.MethodPost, "/api/cars", body, token1)
	require.Equal(t, http.StatusCreated, rec.Code)
	carID := parseJSON(t, rec)["id"].(string)

	updateBody := `{"make":"BMW","model":"M4","year":2021}`
	rec = app.doRequest(http.MethodPut, "/api/cars/"+carID, updateBody, token2)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestCars_Delete(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "User", uniqueEmail("cars"), "password123")

	body := `{"make":"BMW","model":"M3","year":2020}`
	rec := app.doRequest(http.MethodPost, "/api/cars", body, token)
	require.Equal(t, http.StatusCreated, rec.Code)
	carID := parseJSON(t, rec)["id"].(string)

	rec = app.doRequest(http.MethodDelete, "/api/cars/"+carID, "", token)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify deleted
	rec = app.doRequest(http.MethodGet, "/api/cars", "", token)
	result := parseJSONArray(t, rec)
	assert.Empty(t, result)
}

func TestCars_DeleteNotFound(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "User", uniqueEmail("cars"), "password123")

	rec := app.doRequest(http.MethodDelete, "/api/cars/nonexistent", "", token)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

// ─── Car Mods ───────────────────────────────────────────────────────────────────

func TestCarMods_Create(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "User", uniqueEmail("mods"), "password123")

	// Create car first
	carBody := `{"make":"BMW","model":"M3","year":2020}`
	rec := app.doRequest(http.MethodPost, "/api/cars", carBody, token)
	require.Equal(t, http.StatusCreated, rec.Code)
	carID := parseJSON(t, rec)["id"].(string)

	// Create mod
	modBody := `{"name":"Turbo Kit","category":"ENGINE","notes":"Stage 2"}`
	rec = app.doRequest(http.MethodPost, "/api/cars/"+carID+"/mods", modBody, token)
	assert.Equal(t, http.StatusCreated, rec.Code)

	result := parseJSON(t, rec)
	assert.Equal(t, "Turbo Kit", result["name"])
	assert.Equal(t, "ENGINE", result["category"])
	assert.NotEmpty(t, result["id"])
}

func TestCarMods_InvalidCategory(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "User", uniqueEmail("mods"), "password123")

	carBody := `{"make":"BMW","model":"M3","year":2020}`
	rec := app.doRequest(http.MethodPost, "/api/cars", carBody, token)
	require.Equal(t, http.StatusCreated, rec.Code)
	carID := parseJSON(t, rec)["id"].(string)

	modBody := `{"name":"Thing","category":"INVALID"}`
	rec = app.doRequest(http.MethodPost, "/api/cars/"+carID+"/mods", modBody, token)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCarMods_Delete(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "User", uniqueEmail("mods"), "password123")

	carBody := `{"make":"BMW","model":"M3","year":2020}`
	rec := app.doRequest(http.MethodPost, "/api/cars", carBody, token)
	require.Equal(t, http.StatusCreated, rec.Code)
	carID := parseJSON(t, rec)["id"].(string)

	modBody := `{"name":"Coilovers","category":"SUSPENSION"}`
	rec = app.doRequest(http.MethodPost, "/api/cars/"+carID+"/mods", modBody, token)
	require.Equal(t, http.StatusCreated, rec.Code)
	modID := parseJSON(t, rec)["id"].(string)

	rec = app.doRequest(http.MethodDelete, "/api/cars/"+carID+"/mods/"+modID, "", token)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCarMods_DeleteNotOwned(t *testing.T) {
	app := setupTestApp(t)
	_, token1 := app.createTestUser(t, "User1", uniqueEmail("mods1"), "password123")
	_, token2 := app.createTestUser(t, "User2", uniqueEmail("mods2"), "password123")

	carBody := `{"make":"BMW","model":"M3","year":2020}`
	rec := app.doRequest(http.MethodPost, "/api/cars", carBody, token1)
	require.Equal(t, http.StatusCreated, rec.Code)
	carID := parseJSON(t, rec)["id"].(string)

	modBody := `{"name":"Coilovers","category":"SUSPENSION"}`
	rec = app.doRequest(http.MethodPost, "/api/cars/"+carID+"/mods", modBody, token1)
	require.Equal(t, http.StatusCreated, rec.Code)
	modID := parseJSON(t, rec)["id"].(string)

	// User2 tries to delete
	rec = app.doRequest(http.MethodDelete, "/api/cars/"+carID+"/mods/"+modID, "", token2)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestCarMods_VisibleInCarList(t *testing.T) {
	app := setupTestApp(t)
	_, token := app.createTestUser(t, "User", uniqueEmail("mods"), "password123")

	carBody := `{"make":"BMW","model":"M3","year":2020}`
	rec := app.doRequest(http.MethodPost, "/api/cars", carBody, token)
	require.Equal(t, http.StatusCreated, rec.Code)
	carID := parseJSON(t, rec)["id"].(string)

	modBody := `{"name":"Big Brake Kit","category":"BRAKES"}`
	rec = app.doRequest(http.MethodPost, "/api/cars/"+carID+"/mods", modBody, token)
	require.Equal(t, http.StatusCreated, rec.Code)

	// List cars and verify mods are included
	rec = app.doRequest(http.MethodGet, "/api/cars", "", token)
	assert.Equal(t, http.StatusOK, rec.Code)

	var cars []json.RawMessage
	err := json.Unmarshal(rec.Body.Bytes(), &cars)
	require.NoError(t, err)
	require.Len(t, cars, 1)

	var car map[string]interface{}
	json.Unmarshal(cars[0], &car)
	mods := car["mods"].([]interface{})
	assert.Len(t, mods, 1)
	assert.Equal(t, "Big Brake Kit", mods[0].(map[string]interface{})["name"])
}
