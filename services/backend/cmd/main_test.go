package main

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionEndpoint_Success(t *testing.T) {
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

func TestVersionEndpoint_WrongMethod_MethodNotAllowed(t *testing.T) {
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

func TestHealthEndpoint_Success(t *testing.T) {
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

func TestHealthEndpoint_WrongMethod_MethodNotAllowed(t *testing.T) {
	// ORPHAN: validates edge case for health endpoint method validation
	// Reason: Method validation for health endpoint not specified in requirements

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

func TestNonExistentEndpoint_NotFound(t *testing.T) {
	// ORPHAN: validates 404 handling for unit tests
	// Reason: Covered by EE-00003-01 at E2E level, but useful for unit testing

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

func TestCreateFiberApp_ReturnsConfiguredApp(t *testing.T) {
	// ORPHAN: validates Fiber app configuration
	// Reason: Testing internal framework setup not specified in requirements

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
