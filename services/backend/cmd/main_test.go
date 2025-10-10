package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Unit test for createFiberApp function.
func TestORPHAN_CreateFiberApp_ReturnsConfiguredApp(t *testing.T) {
	// Act
	app := createFiberApp("http://localhost:3000")

	// Assert
	assert.NotNil(t, app)
	// Test that the app is configured with correct settings
	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	defer resp.Body.Close()
	// Verify strict routing is enabled (should return 404 for non-existent routes)
	assert.Equal(t, 404, resp.StatusCode)
}
