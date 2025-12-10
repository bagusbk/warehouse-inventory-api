package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"warehouse/handlers"
	"warehouse/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupAuthRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	handler := handlers.NewUserHandler()

	api := r.Group("/api")
	{
		api.POST("/register", handler.Register)
		api.POST("/login", handler.Login)
	}

	return r
}

func TestLoginSuccess(t *testing.T) {
	router := setupAuthRouter()

	loginReq := models.LoginRequest{
		Username: "admin",
		Password: "password123",
	}

	jsonBody, _ := json.Marshal(loginReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
}

func TestLoginInvalidCredentials(t *testing.T) {
	router := setupAuthRouter()

	loginReq := models.LoginRequest{
		Username: "admin",
		Password: "wrongpassword",
	}

	jsonBody, _ := json.Marshal(loginReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLoginMissingFields(t *testing.T) {
	router := setupAuthRouter()

	loginReq := map[string]string{
		"username": "admin",
	}

	jsonBody, _ := json.Marshal(loginReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestRegisterValidation(t *testing.T) {
	router := setupAuthRouter()

	invalidReq := map[string]string{
		"username": "ab",
		"password": "123",
	}

	jsonBody, _ := json.Marshal(invalidReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestRegisterDuplicateUsername(t *testing.T) {
	router := setupAuthRouter()

	registerReq := models.RegisterRequest{
		Username: "admin",
		Password: "password123",
		Email:    "admin@warehouse.com",
		FullName: "Admin",
	}

	jsonBody, _ := json.Marshal(registerReq)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
