package handlers

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/smithisrealdev/travel-ai-agent/backend/internal/database"
	"github.com/smithisrealdev/travel-ai-agent/backend/internal/models"
	"github.com/smithisrealdev/travel-ai-agent/backend/internal/services"
)

// TravelHandler handles travel-related HTTP requests
type TravelHandler struct {
	db      *database.PostgresDB
	redis   *database.RedisCache
	openai  *services.OpenAIService
	weather *services.WeatherService
	flight  *services.FlightService
}

// NewTravelHandler creates a new travel handler instance
func NewTravelHandler(
	db *database.PostgresDB,
	redis *database.RedisCache,
	openai *services.OpenAIService,
	weather *services.WeatherService,
	flight *services.FlightService,
) *TravelHandler {
	return &TravelHandler{
		db:      db,
		redis:   redis,
		openai:  openai,
		weather: weather,
		flight:  flight,
	}
}

// SearchTravel handles travel search requests
func (h *TravelHandler) SearchTravel(c *fiber.Ctx) error {
	var req models.TravelSearchRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
			Code:    fiber.StatusBadRequest,
		})
	}

	// Validate required fields
	if req.Destination == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Validation error",
			Message: "Destination is required",
			Code:    fiber.StatusBadRequest,
		})
	}

	ctx := context.Background()

	// Check cache first
	cacheKey := "travel:" + req.Destination
	if cachedData, err := h.redis.Get(cacheKey); err == nil {
		log.Println("Returning cached travel data for", req.Destination)
		var cachedResponse models.TravelSearchResponse
		if err := json.Unmarshal([]byte(cachedData), &cachedResponse); err == nil {
			return c.JSON(cachedResponse)
		}
	}

	// Generate AI recommendations
	var aiRecommendations string
	if h.openai != nil {
		recommendations, err := h.openai.GenerateTravelRecommendations(ctx, &req)
		if err != nil {
			log.Printf("Warning: Failed to get AI recommendations: %v", err)
		} else {
			aiRecommendations = recommendations
		}
	}

	// Fetch weather data
	var weatherInfo *models.WeatherInfo
	if h.weather != nil {
		weather, err := h.weather.GetWeather(req.Destination)
		if err != nil {
			log.Printf("Warning: Failed to fetch weather: %v", err)
		} else {
			weatherInfo = weather
			// Fetch forecast
			forecast, err := h.weather.GetForecast(req.Destination, 5)
			if err == nil {
				weatherInfo.Forecast = forecast
			}
		}
	}

	// Store in database
	searchID, err := h.storeSearch(&req, aiRecommendations)
	if err != nil {
		log.Printf("Warning: Failed to store search: %v", err)
	}

	// Build response
	response := models.TravelSearchResponse{
		SearchID:    searchID,
		Destination: req.Destination,
		Summary:     aiRecommendations,
		Weather:     weatherInfo,
		Recommendations: []models.TravelRecommendation{
			{
				Type:        "hotel",
				Title:       "Luxury Hotel Recommendation",
				Description: "Experience comfort in the heart of " + req.Destination,
				Price:       150.00,
				Rating:      4.5,
			},
			{
				Type:        "activity",
				Title:       "City Tour",
				Description: "Discover the best of " + req.Destination,
				Price:       50.00,
				Rating:      4.8,
			},
		},
		EstimatedCost: req.Budget,
		CreatedAt:     time.Now(),
	}

	// Cache the response for 1 hour
	if responseJSON, err := json.Marshal(response); err == nil {
		h.redis.Set(cacheKey, responseJSON, time.Hour)
	}

	return c.JSON(response)
}

// GetSearchHistory retrieves user's search history
func (h *TravelHandler) GetSearchHistory(c *fiber.Ctx) error {
	userID := c.Query("userId")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Validation error",
			Message: "userId is required",
			Code:    fiber.StatusBadRequest,
		})
	}

	query := `
		SELECT id, user_id, destination, start_date, end_date, budget, 
		       preferences, results, created_at, updated_at
		FROM travel_searches
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 20
	`

	rows, err := h.db.DB.Query(query, userID)
	if err != nil {
		log.Printf("Database error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "Database error",
			Message: "Failed to retrieve search history",
			Code:    fiber.StatusInternalServerError,
		})
	}
	defer rows.Close()

	searches := []models.TravelSearch{}
	for rows.Next() {
		var search models.TravelSearch
		var preferencesJSON, resultsJSON []byte

		err := rows.Scan(
			&search.ID,
			&search.UserID,
			&search.Destination,
			&search.StartDate,
			&search.EndDate,
			&search.Budget,
			&preferencesJSON,
			&resultsJSON,
			&search.CreatedAt,
			&search.UpdatedAt,
		)
		if err != nil {
			log.Printf("Row scan error: %v", err)
			continue
		}

		// Parse JSON fields
		if len(preferencesJSON) > 0 {
			json.Unmarshal(preferencesJSON, &search.Preferences)
		}
		if len(resultsJSON) > 0 {
			json.Unmarshal(resultsJSON, &search.Results)
		}

		searches = append(searches, search)
	}

	return c.JSON(searches)
}

// HealthCheck handles health check requests
func (h *TravelHandler) HealthCheck(c *fiber.Ctx) error {
	services := make(map[string]string)

	// Check database
	if err := h.db.HealthCheck(); err != nil {
		services["database"] = "unhealthy: " + err.Error()
	} else {
		services["database"] = "healthy"
	}

	// Check Redis
	if err := h.redis.HealthCheck(); err != nil {
		services["redis"] = "unhealthy: " + err.Error()
	} else {
		services["redis"] = "healthy"
	}

	// Check OpenAI
	if h.openai != nil {
		if err := h.openai.HealthCheck(); err != nil {
			services["openai"] = "unhealthy: " + err.Error()
		} else {
			services["openai"] = "healthy"
		}
	} else {
		services["openai"] = "not configured"
	}

	// Check Weather
	if h.weather != nil {
		if err := h.weather.HealthCheck(); err != nil {
			services["weather"] = "unhealthy: " + err.Error()
		} else {
			services["weather"] = "healthy"
		}
	} else {
		services["weather"] = "not configured"
	}

	// Check Flight
	if h.flight != nil {
		if err := h.flight.HealthCheck(); err != nil {
			services["flight"] = "unhealthy: " + err.Error()
		} else {
			services["flight"] = "healthy"
		}
	} else {
		services["flight"] = "not configured"
	}

	response := models.HealthCheckResponse{
		Status:   "healthy",
		Services: services,
		Time:     time.Now(),
	}

	return c.JSON(response)
}

// storeSearch stores a travel search in the database
func (h *TravelHandler) storeSearch(req *models.TravelSearchRequest, results string) (int, error) {
	var searchID int

	preferencesJSON, _ := json.Marshal(req.Preferences)
	resultsMap := map[string]interface{}{
		"recommendations": results,
	}
	resultsJSON, _ := json.Marshal(resultsMap)

	query := `
		INSERT INTO travel_searches (user_id, destination, budget, preferences, results)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	err := h.db.DB.QueryRow(
		query,
		req.UserID,
		req.Destination,
		req.Budget,
		preferencesJSON,
		resultsJSON,
	).Scan(&searchID)

	if err != nil {
		return 0, err
	}

	return searchID, nil
}
