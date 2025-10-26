package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/smithisrealdev/travel-ai-agent/backend/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestSocialHandler_GetSocialPlaces_ValidationError(t *testing.T) {
	// Create a new Fiber app
	app := fiber.New()

	// Create handler without services (will return service unavailable)
	handler := NewSocialHandler(nil, nil)

	// Register route
	app.Post("/api/social", handler.GetSocialPlaces)

	tests := []struct {
		name           string
		requestBody    models.SocialPlaceRequest
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Missing keyword",
			requestBody: models.SocialPlaceRequest{
				Location: "Tokyo",
			},
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  "Keyword is required",
		},
		{
			name: "Missing location",
			requestBody: models.SocialPlaceRequest{
				Keyword: "restaurants",
			},
			expectedStatus: fiber.StatusBadRequest,
			expectedError:  "Location is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request body
			body, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			// Create request
			req := httptest.NewRequest("POST", "/api/social", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			// Perform request
			resp, err := app.Test(req)
			assert.NoError(t, err)

			// Check status code
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			// Parse response
			var errorResp models.ErrorResponse
			err = json.NewDecoder(resp.Body).Decode(&errorResp)
			assert.NoError(t, err)

			// Check error message
			assert.Contains(t, errorResp.Message, tt.expectedError)
		})
	}
}

func TestSocialHandler_GetSocialPlaces_ServiceUnavailable(t *testing.T) {
	// Create a new Fiber app
	app := fiber.New()

	// Create handler without social service
	handler := NewSocialHandler(nil, nil)

	// Register route
	app.Post("/api/social", handler.GetSocialPlaces)

	// Create valid request body
	requestBody := models.SocialPlaceRequest{
		Keyword:  "restaurants",
		Location: "Tokyo",
		Limit:    5,
	}

	body, err := json.Marshal(requestBody)
	assert.NoError(t, err)

	// Create request
	req := httptest.NewRequest("POST", "/api/social", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	resp, err := app.Test(req)
	assert.NoError(t, err)

	// Check status code
	assert.Equal(t, fiber.StatusServiceUnavailable, resp.StatusCode)

	// Parse response
	var errorResp models.ErrorResponse
	err = json.NewDecoder(resp.Body).Decode(&errorResp)
	assert.NoError(t, err)

	// Check error message
	assert.Contains(t, errorResp.Message, "not configured")
}
