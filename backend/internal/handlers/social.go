package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/smithisrealdev/travel-ai-agent/backend/internal/database"
	"github.com/smithisrealdev/travel-ai-agent/backend/internal/models"
	"github.com/smithisrealdev/travel-ai-agent/backend/internal/services"
)

// SocialHandler handles social places-related HTTP requests
type SocialHandler struct {
	redis  *database.RedisCache
	social *services.SocialService
}

// NewSocialHandler creates a new social handler instance
func NewSocialHandler(
	redis *database.RedisCache,
	social *services.SocialService,
) *SocialHandler {
	return &SocialHandler{
		redis:  redis,
		social: social,
	}
}

// GetSocialPlaces handles requests to fetch socially popular places
func (h *SocialHandler) GetSocialPlaces(c *fiber.Ctx) error {
	var req models.SocialPlaceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
			Code:    fiber.StatusBadRequest,
		})
	}

	// Validate required fields
	if req.Keyword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Validation error",
			Message: "Keyword is required",
			Code:    fiber.StatusBadRequest,
		})
	}

	if req.Location == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Validation error",
			Message: "Location is required",
			Code:    fiber.StatusBadRequest,
		})
	}

	// Set default limit
	if req.Limit <= 0 {
		req.Limit = 10
	}

	ctx := context.Background()
	_ = ctx // for future use

	// Check cache first (only if redis is available)
	cacheKey := fmt.Sprintf("social:%s:%s:%d", req.Keyword, req.Location, req.Limit)
	if h.redis != nil {
		if cachedData, err := h.redis.Get(cacheKey); err == nil {
			log.Printf("Returning cached social places for %s in %s", req.Keyword, req.Location)
			var cachedResponse models.SocialPlacesResponse
			if err := json.Unmarshal([]byte(cachedData), &cachedResponse); err == nil {
				return c.JSON(cachedResponse)
			}
		}
	}

	// Fetch places from Google Places API
	if h.social == nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(models.ErrorResponse{
			Error:   "Service unavailable",
			Message: "Social places service is not configured",
			Code:    fiber.StatusServiceUnavailable,
		})
	}

	places, err := h.social.GetTopRatedPlaces(req.Keyword, req.Location, req.Limit)
	if err != nil {
		log.Printf("Error fetching social places: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "Failed to fetch social places",
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		})
	}

	// Build response
	response := models.SocialPlacesResponse{
		Places: places,
		Query:  fmt.Sprintf("%s in %s", req.Keyword, req.Location),
		Count:  len(places),
	}

	// Cache the response for 1 hour (only if redis is available)
	if h.redis != nil {
		if responseJSON, err := json.Marshal(response); err == nil {
			h.redis.Set(cacheKey, responseJSON, time.Hour)
		}
	}

	return c.JSON(response)
}
