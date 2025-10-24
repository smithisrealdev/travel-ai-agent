package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/smithisrealdev/travel-ai-agent/backend/internal/config"
	"github.com/smithisrealdev/travel-ai-agent/backend/internal/models"
)

// WeatherService handles weather API interactions
type WeatherService struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// weatherAPIResponse represents the OpenWeatherMap API response
type weatherAPIResponse struct {
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
}

// NewWeatherService creates a new Weather service instance
func NewWeatherService(cfg *config.Config) *WeatherService {
	if cfg.Weather.APIKey == "" {
		log.Println("Warning: Weather API key not configured")
		return nil
	}

	return &WeatherService{
		apiKey:  cfg.Weather.APIKey,
		baseURL: cfg.Weather.URL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetWeather fetches current weather for a destination
func (s *WeatherService) GetWeather(destination string) (*models.WeatherInfo, error) {
	if s == nil {
		return nil, fmt.Errorf("Weather service not initialized")
	}

	// Build the API URL
	params := url.Values{}
	params.Add("q", destination)
	params.Add("appid", s.apiKey)
	params.Add("units", "metric") // Use metric units

	apiURL := fmt.Sprintf("%s/weather?%s", s.baseURL, params.Encode())

	// Make the HTTP request
	resp, err := s.client.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("weather API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var weatherResp weatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherResp); err != nil {
		return nil, fmt.Errorf("failed to parse weather response: %w", err)
	}

	// Convert to our model
	weather := &models.WeatherInfo{
		Temperature: weatherResp.Main.Temp,
		Humidity:    weatherResp.Main.Humidity,
		WindSpeed:   weatherResp.Wind.Speed,
	}

	if len(weatherResp.Weather) > 0 {
		weather.Description = weatherResp.Weather[0].Description
		weather.Icon = weatherResp.Weather[0].Icon
	}

	return weather, nil
}

// GetForecast fetches weather forecast for a destination
func (s *WeatherService) GetForecast(destination string, days int) ([]models.DayForecast, error) {
	if s == nil {
		return nil, fmt.Errorf("Weather service not initialized")
	}

	// Build the API URL for forecast
	params := url.Values{}
	params.Add("q", destination)
	params.Add("appid", s.apiKey)
	params.Add("units", "metric")
	params.Add("cnt", fmt.Sprintf("%d", days*8)) // 8 forecasts per day (3-hour intervals)

	apiURL := fmt.Sprintf("%s/forecast?%s", s.baseURL, params.Encode())

	// Make the HTTP request
	resp, err := s.client.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch forecast data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("forecast API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse the forecast response
	var forecastResp struct {
		List []struct {
			Dt   int64 `json:"dt"`
			Main struct {
				TempMin float64 `json:"temp_min"`
				TempMax float64 `json:"temp_max"`
			} `json:"main"`
			Weather []struct {
				Description string `json:"description"`
				Icon        string `json:"icon"`
			} `json:"weather"`
		} `json:"list"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&forecastResp); err != nil {
		return nil, fmt.Errorf("failed to parse forecast response: %w", err)
	}

	// Group forecasts by day
	forecastMap := make(map[string]*models.DayForecast)
	for _, item := range forecastResp.List {
		date := time.Unix(item.Dt, 0).Format("2006-01-02")
		
		if _, exists := forecastMap[date]; !exists {
			forecast := &models.DayForecast{
				Date:    date,
				TempMin: item.Main.TempMin,
				TempMax: item.Main.TempMax,
			}
			if len(item.Weather) > 0 {
				forecast.Description = item.Weather[0].Description
				forecast.Icon = item.Weather[0].Icon
			}
			forecastMap[date] = forecast
		} else {
			// Update min/max temperatures
			if item.Main.TempMin < forecastMap[date].TempMin {
				forecastMap[date].TempMin = item.Main.TempMin
			}
			if item.Main.TempMax > forecastMap[date].TempMax {
				forecastMap[date].TempMax = item.Main.TempMax
			}
		}
	}

	// Convert map to slice
	forecasts := make([]models.DayForecast, 0, len(forecastMap))
	for _, forecast := range forecastMap {
		forecasts = append(forecasts, *forecast)
	}

	return forecasts, nil
}

// HealthCheck verifies the Weather service is configured
func (s *WeatherService) HealthCheck() error {
	if s == nil {
		return fmt.Errorf("Weather service not initialized")
	}
	return nil
}
