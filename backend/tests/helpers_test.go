package tests

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joezmuda/trackside-backend/internal/config"
	"github.com/joezmuda/trackside-backend/internal/database"
	mw "github.com/joezmuda/trackside-backend/internal/middleware"
	"github.com/joezmuda/trackside-backend/internal/repository"
	"github.com/joezmuda/trackside-backend/internal/router"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

const testJWTSecret = "test-secret-key-for-testing"

// testApp holds all test dependencies.
type testApp struct {
	e         *echo.Echo
	db        *sql.DB
	cfg       *config.Config
	userRepo  *repository.UserRepo
	carRepo   *repository.CarRepo
	trackRepo *repository.TrackRepo
}

// setupTestApp creates an in-memory SQLite database, runs migrations,
// wires up all routes, and returns a ready-to-use testApp.
func setupTestApp(t *testing.T) *testApp {
	t.Helper()

	// Use a temp file for each test to avoid :memory: connection-pool issues
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := database.Connect(dbPath)
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	err = database.Migrate(db)
	require.NoError(t, err)

	cfg := &config.Config{
		Port:        "0",
		DatabaseURL: ":memory:",
		JWTSecret:   testJWTSecret,
		UploadDir:   t.TempDir(),
		CORSOrigins: []string{"http://localhost:3000"},
		DataDir:     ".",
	}

	e := echo.New()
	router.Setup(e, db, cfg)

	return &testApp{
		e:         e,
		db:        db,
		cfg:       cfg,
		userRepo:  repository.NewUserRepo(db),
		carRepo:   repository.NewCarRepo(db),
		trackRepo: repository.NewTrackRepo(db),
	}
}

// generateToken creates a valid JWT for the given user ID and email.
func generateToken(userID, email string) string {
	claims := &mw.JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Email: email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte(testJWTSecret))
	return tokenStr
}

// authHeader returns the Authorization header value for a given token.
func authHeader(token string) string {
	return "Bearer " + token
}

// doRequest performs an HTTP request against the Echo engine and returns
// the recorder. Body can be nil for GET/DELETE.
func (app *testApp) doRequest(method, path, body string, token string) *httptest.ResponseRecorder {
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	if token != "" {
		req.Header.Set("Authorization", authHeader(token))
	}
	rec := httptest.NewRecorder()
	app.e.ServeHTTP(rec, req)
	return rec
}

// createTestUser creates a user in the database and returns (userID, token).
func (app *testApp) createTestUser(t *testing.T, name, email, password string) (string, string) {
	t.Helper()
	user, err := app.userRepo.Create(name, email, password)
	require.NoError(t, err)
	token := generateToken(user.ID, user.Email)
	return user.ID, token
}

// createTestCar creates a car for a user and returns the car ID.
func (app *testApp) createTestCar(t *testing.T, make, model string, year int, userID string) string {
	t.Helper()
	car, err := app.carRepo.Create(make, model, year, userID)
	require.NoError(t, err)
	return car.ID
}

// createTestTrack creates a track via the API and returns the track ID.
func (app *testApp) createTestTrack(t *testing.T, token string) string {
	t.Helper()
	body := `{"name":"Test Track","location":"Test City, CA","eventTypes":["ROADCOURSE"]}`
	rec := app.doRequest(http.MethodPost, "/api/tracks", body, token)
	require.Equal(t, http.StatusCreated, rec.Code, "create track: %s", rec.Body.String())

	var result map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	require.NoError(t, err)
	return result["id"].(string)
}

// parseJSON unmarshals the response body into a map.
func parseJSON(t *testing.T, rec *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	var result map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	require.NoError(t, err, "failed to parse JSON: %s", rec.Body.String())
	return result
}

// parseJSONArray unmarshals the response body into a slice of maps.
func parseJSONArray(t *testing.T, rec *httptest.ResponseRecorder) []map[string]interface{} {
	t.Helper()
	var result []map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	require.NoError(t, err, "failed to parse JSON array: %s", rec.Body.String())
	return result
}

// uniqueEmail generates a unique email for test isolation.
func uniqueEmail(prefix string) string {
	return fmt.Sprintf("%s-%d@test.com", prefix, time.Now().UnixNano())
}
