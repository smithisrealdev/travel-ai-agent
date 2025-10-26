package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"time"

	"github.com/smithisrealdev/travel-ai-agent/backend/internal/config"
	"github.com/smithisrealdev/travel-ai-agent/backend/internal/models"
)

// SocialService handles Google Places API interactions
type SocialService struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// googlePlacesTextSearchResponse represents Google Places Text Search API response
type googlePlacesTextSearchResponse struct {
	Results []struct {
		PlaceID         string   `json:"place_id"`
		Name            string   `json:"name"`
		FormattedAddress string  `json:"formatted_address"`
		Rating          float64  `json:"rating"`
		UserRatingsTotal int     `json:"user_ratings_total"`
		PriceLevel      int      `json:"price_level"`
		Types           []string `json:"types"`
		Geometry        struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
		OpeningHours *struct {
			OpenNow bool `json:"open_now"`
		} `json:"opening_hours,omitempty"`
		Photos []struct {
			PhotoReference string `json:"photo_reference"`
		} `json:"photos,omitempty"`
	} `json:"results"`
	Status string `json:"status"`
}

// NewSocialService creates a new Social service instance
func NewSocialService(cfg *config.Config) *SocialService {
	if cfg.GooglePlaces.APIKey == "" {
		log.Println("Warning: Google Places API key not configured")
		return nil
	}

	return &SocialService{
		apiKey:  cfg.GooglePlaces.APIKey,
		baseURL: cfg.GooglePlaces.URL,
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// GetTopRatedPlaces fetches top-rated places from Google Places API
func (s *SocialService) GetTopRatedPlaces(keyword, location string, limit int) ([]models.SocialPlace, error) {
	if s == nil {
		return nil, fmt.Errorf("social service not initialized")
	}

	if limit <= 0 {
		limit = 10
	}

	// Build the query
	query := keyword + " in " + location

	// Build the API URL
	params := url.Values{}
	params.Add("query", query)
	params.Add("key", s.apiKey)

	apiURL := fmt.Sprintf("%s/textsearch/json?%s", s.baseURL, params.Encode())

	// Make the HTTP request
	resp, err := s.client.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch places data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("places API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var placesResp googlePlacesTextSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&placesResp); err != nil {
		return nil, fmt.Errorf("failed to parse places response: %w", err)
	}

	if placesResp.Status != "OK" && placesResp.Status != "ZERO_RESULTS" {
		return nil, fmt.Errorf("places API returned status: %s", placesResp.Status)
	}

	// Convert to our model
	places := make([]models.SocialPlace, 0, len(placesResp.Results))
	for _, result := range placesResp.Results {
		place := models.SocialPlace{
			PlaceID:     result.PlaceID,
			Name:        result.Name,
			Address:     result.FormattedAddress,
			Rating:      result.Rating,
			ReviewCount: result.UserRatingsTotal,
			PriceLevel:  result.PriceLevel,
			Types:       result.Types,
			Latitude:    result.Geometry.Location.Lat,
			Longitude:   result.Geometry.Location.Lng,
		}

		if result.OpeningHours != nil {
			place.OpenNow = result.OpeningHours.OpenNow
		}

		// Get photo URL if available
		if len(result.Photos) > 0 {
			place.PhotoURL = s.getPhotoURL(result.Photos[0].PhotoReference)
		}

		places = append(places, place)
	}

	// Sort by review count (descending)
	sort.Slice(places, func(i, j int) bool {
		return places[i].ReviewCount > places[j].ReviewCount
	})

	// Return top N places
	if len(places) > limit {
		places = places[:limit]
	}

	return places, nil
}

// getPhotoURL constructs the photo URL from a photo reference
func (s *SocialService) getPhotoURL(photoReference string) string {
	if photoReference == "" {
		return ""
	}
	return fmt.Sprintf("%s/photo?maxwidth=400&photoreference=%s&key=%s",
		s.baseURL, photoReference, s.apiKey)
}

// HealthCheck verifies the Social service is configured
func (s *SocialService) HealthCheck() error {
	if s == nil {
		return fmt.Errorf("social service not initialized")
	}
	return nil
}
