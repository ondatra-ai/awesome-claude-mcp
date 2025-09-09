package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetVersion(t *testing.T) {
	// Create test app
	app := fiber.New()
	app.Get("/version", getVersion)

	// Create test request
	req := httptest.NewRequest("GET", "/version", nil)
	resp, err := app.Test(req)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Check response content-type is JSON
	contentType := resp.Header.Get("Content-Type")
	assert.Contains(t, contentType, "application/json")
}

func TestGetHealth(t *testing.T) {
	// Create test app
	app := fiber.New()
	app.Get("/health", getHealth)

	// Create test request
	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Check response content-type is JSON
	contentType := resp.Header.Get("Content-Type")
	assert.Contains(t, contentType, "application/json")
}
