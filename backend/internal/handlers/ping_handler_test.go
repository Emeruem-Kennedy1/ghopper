package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Create a new router
	router := gin.Default()
	router.GET("/ping", Ping())

	// Create a new HTTP request
	req, err := http.NewRequest(http.MethodGet, "/ping", nil)
	assert.NoError(t, err)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check the response body
	expected := `{"message":"pong"}`
	assert.JSONEq(t, expected, rr.Body.String())
}