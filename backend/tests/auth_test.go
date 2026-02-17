package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegister_Success(t *testing.T) {
	app := setupTestApp(t)
	email := uniqueEmail("register")

	body := `{"name":"Test User","email":"` + email + `","password":"password123","confirmPassword":"password123"}`
	rec := app.doRequest(http.MethodPost, "/api/register", body, "")

	assert.Equal(t, http.StatusCreated, rec.Code)
	result := parseJSON(t, rec)
	assert.Equal(t, "Test User", result["name"])
	assert.Equal(t, email, result["email"])
	assert.NotEmpty(t, result["id"])
}

func TestRegister_ShortName(t *testing.T) {
	app := setupTestApp(t)
	body := `{"name":"A","email":"x@test.com","password":"password123","confirmPassword":"password123"}`
	rec := app.doRequest(http.MethodPost, "/api/register", body, "")
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	result := parseJSON(t, rec)
	assert.Contains(t, result["error"], "Name must be at least")
}

func TestRegister_ShortPassword(t *testing.T) {
	app := setupTestApp(t)
	body := `{"name":"Test","email":"x@test.com","password":"short","confirmPassword":"short"}`
	rec := app.doRequest(http.MethodPost, "/api/register", body, "")
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestRegister_PasswordMismatch(t *testing.T) {
	app := setupTestApp(t)
	body := `{"name":"Test","email":"x@test.com","password":"password123","confirmPassword":"different123"}`
	rec := app.doRequest(http.MethodPost, "/api/register", body, "")
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	result := parseJSON(t, rec)
	assert.Contains(t, result["error"], "Passwords do not match")
}

func TestRegister_DuplicateEmail(t *testing.T) {
	app := setupTestApp(t)
	email := uniqueEmail("dup")

	body := `{"name":"User1","email":"` + email + `","password":"password123","confirmPassword":"password123"}`
	rec := app.doRequest(http.MethodPost, "/api/register", body, "")
	require.Equal(t, http.StatusCreated, rec.Code)

	// Same email again
	rec = app.doRequest(http.MethodPost, "/api/register", body, "")
	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestRegister_MissingEmail(t *testing.T) {
	app := setupTestApp(t)
	body := `{"name":"Test","email":"","password":"password123","confirmPassword":"password123"}`
	rec := app.doRequest(http.MethodPost, "/api/register", body, "")
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestLogin_Success(t *testing.T) {
	app := setupTestApp(t)
	email := uniqueEmail("login")

	// Register first
	regBody := `{"name":"Login User","email":"` + email + `","password":"password123","confirmPassword":"password123"}`
	rec := app.doRequest(http.MethodPost, "/api/register", regBody, "")
	require.Equal(t, http.StatusCreated, rec.Code)

	// Login
	loginBody := `{"email":"` + email + `","password":"password123"}`
	rec = app.doRequest(http.MethodPost, "/api/auth/login", loginBody, "")
	assert.Equal(t, http.StatusOK, rec.Code)

	result := parseJSON(t, rec)
	assert.NotEmpty(t, result["token"])
	user := result["user"].(map[string]interface{})
	assert.Equal(t, email, user["email"])
	assert.Equal(t, "Login User", user["name"])
}

func TestLogin_WrongPassword(t *testing.T) {
	app := setupTestApp(t)
	email := uniqueEmail("loginwrong")

	regBody := `{"name":"User","email":"` + email + `","password":"password123","confirmPassword":"password123"}`
	app.doRequest(http.MethodPost, "/api/register", regBody, "")

	loginBody := `{"email":"` + email + `","password":"wrongpassword"}`
	rec := app.doRequest(http.MethodPost, "/api/auth/login", loginBody, "")
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestLogin_NonExistentEmail(t *testing.T) {
	app := setupTestApp(t)
	body := `{"email":"nobody@test.com","password":"password123"}`
	rec := app.doRequest(http.MethodPost, "/api/auth/login", body, "")
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestLogin_MissingFields(t *testing.T) {
	app := setupTestApp(t)
	body := `{"email":"","password":""}`
	rec := app.doRequest(http.MethodPost, "/api/auth/login", body, "")
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestProtectedRoute_NoToken(t *testing.T) {
	app := setupTestApp(t)
	rec := app.doRequest(http.MethodGet, "/api/cars", "", "")
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestProtectedRoute_InvalidToken(t *testing.T) {
	app := setupTestApp(t)
	rec := app.doRequest(http.MethodGet, "/api/cars", "", "invalid-token")
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
