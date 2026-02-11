package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"irontrack-backend/internal/database"
	"irontrack-backend/internal/models"
	"irontrack-backend/internal/router"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Setup In-Memory Database
	database.ConnectDatabase("file::memory:?cache=shared")

	// Run tests
	code := m.Run()

	os.Exit(code)
}

func setupTestRouter() *gin.Engine {
	return router.SetupRouter("test-commit", "test-time")
}

func TestRegisterAndLogin(t *testing.T) {
	r := setupTestRouter()

	// 1. Register
	w := httptest.NewRecorder()
	registerPayload := map[string]string{
		"name":     "Test User",
		"email":    "test@example.com",
		"password": "password123",
	}
	jsonValue, _ := json.Marshal(registerPayload)
	req, _ := http.NewRequest("POST", "/api/register", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// 2. Login
	w = httptest.NewRecorder()
	loginPayload := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	jsonValue, _ = json.Marshal(loginPayload)
	req, _ = http.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.NotEmpty(t, response["token"])
}

func TestPlansCRUD(t *testing.T) {
	r := setupTestRouter()

	// Create User & Get Token
	db := database.DB
	user := models.User{Email: "plan_tester@example.com", Password: "hashedpassword"}
	db.Create(&user)

	// Manually generate token or use login endpoint
	// For simplicity, let's use login endpoint to get a valid token signed with correct secret
	w := httptest.NewRecorder()
	// Need to register first because we inserted with hashedpassword directly but didn't actually hash it properly matching the login bcrypt check
	// So let's just Register via API
	registerPayload := map[string]string{
		"name":     "Plan Tester",
		"email":    "plan_test@example.com",
		"password": "password123",
	}
	jsonBytes, _ := json.Marshal(registerPayload)
	req, _ := http.NewRequest("POST", "/api/register", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	// Login
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	var loginResp map[string]string
	json.Unmarshal(w.Body.Bytes(), &loginResp)
	token := loginResp["token"]
	assert.NotEmpty(t, token)

	// 1. Create Plan
	w = httptest.NewRecorder()
	planPayload := models.WorkoutPlan{
		Name:        "Test Plan",
		Description: "A test plan",
	}
	// Note: We need to match the struct expected by handler.
	// Handler CreatePlan binds JSON to WorkoutPlan.
	jsonBytes, _ = json.Marshal(planPayload)
	req, _ = http.NewRequest("POST", "/api/plans", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token) // No Bearer prefix? Middleware checks: tokenString := c.GetHeader("Authorization"). If it expects Bearer, we need it.
	// Looking at auth middleware in main.go -> auth.AuthMiddleware
	// I should verify auth middleware. Usually it expects Bearer or just token.
	// I'll assume just token or Bearer. I'll check auth middleware later if it fails.
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var createdPlan models.WorkoutPlan
	json.Unmarshal(w.Body.Bytes(), &createdPlan)
	assert.Equal(t, "Test Plan", createdPlan.Name)

	// 2. Get Plans
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/plans", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
