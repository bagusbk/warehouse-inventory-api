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

func setupBarangRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	handler := handlers.NewBarangHandler()

	api := r.Group("/api")
	{
		api.GET("/barang", handler.GetAll)
		api.GET("/barang/stok", handler.GetAllWithStok)
		api.GET("/barang/:id", handler.GetByID)
		api.POST("/barang", handler.Create)
		api.PUT("/barang/:id", handler.Update)
		api.DELETE("/barang/:id", handler.Delete)
	}

	return r
}

func TestGetAllBarang(t *testing.T) {
	router := setupBarangRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/barang?page=1&limit=10", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
}

func TestGetAllBarangWithStok(t *testing.T) {
	router := setupBarangRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/barang/stok", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
}

func TestGetBarangByIDNotFound(t *testing.T) {
	router := setupBarangRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/barang/99999", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreateBarangValidation(t *testing.T) {
	router := setupBarangRouter()

	invalidBarang := models.BarangRequest{
		NamaBarang: "Test",
	}

	jsonBody, _ := json.Marshal(invalidBarang)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/barang", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestSearchBarang(t *testing.T) {
	router := setupBarangRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/barang?search=laptop", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
}

func TestPaginationParams(t *testing.T) {
	router := setupBarangRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/barang?page=2&limit=5", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Meta)
	assert.Equal(t, 2, response.Meta.Page)
	assert.Equal(t, 5, response.Meta.Limit)
}
