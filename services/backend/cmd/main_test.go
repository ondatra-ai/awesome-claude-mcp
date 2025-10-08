package main

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestORPHAN_VersionEndpoint_Success(t *testing.T) {
	// Arrange
	app := createFiberApp("http://localhost:3000")
	setupRoutes(app)

	// Act
	req := httptest.NewRequest("GET", "/version", nil)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var versionResp VersionResponse
	err = json.NewDecoder(resp.Body).Decode(&versionResp)
	assert.NoError(t, err)
	assert.Equal(t, "1.0.0", versionResp.Version)
}

func TestORPHAN_VersionEndpoint_WrongMethod_MethodNotAllowed(t *testing.T) {
	// Arrange
	app := createFiberApp("http://localhost:3000")
	setupRoutes(app)

	// Act
	req := httptest.NewRequest("POST", "/version", nil)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 405, resp.StatusCode)
}

func TestORPHAN_HealthEndpoint_Success(t *testing.T) {
	// Arrange
	app := createFiberApp("http://localhost:3000")
	setupRoutes(app)

	// Act
	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var healthResp HealthResponse
	err = json.NewDecoder(resp.Body).Decode(&healthResp)
	assert.NoError(t, err)
	assert.Equal(t, "healthy", healthResp.Status)
	assert.Equal(t, "MCP Google Docs Editor - Backend", healthResp.Service)
	assert.NotEmpty(t, healthResp.Timestamp)
}

func TestORPHAN_HealthEndpoint_WrongMethod_MethodNotAllowed(t *testing.T) {
	// Arrange
	app := createFiberApp("http://localhost:3000")
	setupRoutes(app)

	// Act
	req := httptest.NewRequest("DELETE", "/health", nil)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 405, resp.StatusCode)
}

func TestORPHAN_NonExistentEndpoint_NotFound(t *testing.T) {
	// Arrange
	app := createFiberApp("http://localhost:3000")
	setupRoutes(app)

	// Act
	req := httptest.NewRequest("GET", "/nonexistent", nil)
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
}

func TestORPHAN_CreateFiberApp_ReturnsConfiguredApp(t *testing.T) {
	// Act
	app := createFiberApp("http://localhost:3000")

	// Assert
	assert.NotNil(t, app)
	// Test that the app is configured with correct settings
	req := httptest.NewRequest("GET", "/nonexistent", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	// Verify strict routing is enabled (should return 404 for non-existent routes)
	assert.Equal(t, 404, resp.StatusCode)
}
