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

// FlightService handles flight API interactions
type FlightService struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// flightAPIResponse represents a flight API response
type flightAPIResponse struct {
	Data []struct {
		FlightNumber string `json:"flight_number"`
		Airline      struct {
			Name string `json:"name"`
		} `json:"airline"`
		Departure struct {
			Airport   string `json:"airport"`
			Scheduled string `json:"scheduled"`
		} `json:"departure"`
		Arrival struct {
			Airport   string `json:"airport"`
			Scheduled string `json:"scheduled"`
		} `json:"arrival"`
		FlightStatus string `json:"flight_status"`
	} `json:"data"`
}

// NewFlightService creates a new Flight service instance
func NewFlightService(cfg *config.Config) *FlightService {
	if cfg.Flight.APIKey == "" {
		log.Println("Warning: Flight API key not configured")
		return nil
	}

	return &FlightService{
		apiKey:  cfg.Flight.APIKey,
		baseURL: cfg.Flight.URL,
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// SearchFlights searches for flights between origin and destination
func (s *FlightService) SearchFlights(origin, destination, date string) ([]models.FlightInfo, error) {
	if s == nil {
		return nil, fmt.Errorf("Flight service not initialized")
	}

	// Build the API URL
	params := url.Values{}
	params.Add("access_key", s.apiKey)
	params.Add("dep_iata", origin)
	params.Add("arr_iata", destination)
	params.Add("flight_date", date)

	apiURL := fmt.Sprintf("%s/flights?%s", s.baseURL, params.Encode())

	// Make the HTTP request
	resp, err := s.client.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch flight data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("flight API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var flightResp flightAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&flightResp); err != nil {
		return nil, fmt.Errorf("failed to parse flight response: %w", err)
	}

	// Convert to our model
	flights := make([]models.FlightInfo, 0, len(flightResp.Data))
	for _, flight := range flightResp.Data {
		departTime, _ := time.Parse(time.RFC3339, flight.Departure.Scheduled)
		arriveTime, _ := time.Parse(time.RFC3339, flight.Arrival.Scheduled)

		duration := arriveTime.Sub(departTime)

		flightInfo := models.FlightInfo{
			FlightNumber: flight.FlightNumber,
			Airline:      flight.Airline.Name,
			Departure:    flight.Departure.Airport,
			Arrival:      flight.Arrival.Airport,
			DepartTime:   departTime,
			ArriveTime:   arriveTime,
			Duration:     formatDuration(duration),
			Price:        0, // Price not provided by basic API
			Stops:        0, // Direct flights only in this basic implementation
		}

		flights = append(flights, flightInfo)
	}

	return flights, nil
}

// GetFlightStatus gets real-time status of a specific flight
func (s *FlightService) GetFlightStatus(flightNumber string) (*models.FlightInfo, error) {
	if s == nil {
		return nil, fmt.Errorf("Flight service not initialized")
	}

	// Build the API URL
	params := url.Values{}
	params.Add("access_key", s.apiKey)
	params.Add("flight_iata", flightNumber)

	apiURL := fmt.Sprintf("%s/flights?%s", s.baseURL, params.Encode())

	// Make the HTTP request
	resp, err := s.client.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch flight status: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("flight API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var flightResp flightAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&flightResp); err != nil {
		return nil, fmt.Errorf("failed to parse flight response: %w", err)
	}

	if len(flightResp.Data) == 0 {
		return nil, fmt.Errorf("flight not found")
	}

	flight := flightResp.Data[0]
	departTime, _ := time.Parse(time.RFC3339, flight.Departure.Scheduled)
	arriveTime, _ := time.Parse(time.RFC3339, flight.Arrival.Scheduled)

	duration := arriveTime.Sub(departTime)

	flightInfo := &models.FlightInfo{
		FlightNumber: flight.FlightNumber,
		Airline:      flight.Airline.Name,
		Departure:    flight.Departure.Airport,
		Arrival:      flight.Arrival.Airport,
		DepartTime:   departTime,
		ArriveTime:   arriveTime,
		Duration:     formatDuration(duration),
		Price:        0,
		Stops:        0,
	}

	return flightInfo, nil
}

// formatDuration formats a duration into a readable string
func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	return fmt.Sprintf("%dh %dm", hours, minutes)
}

// HealthCheck verifies the Flight service is configured
func (s *FlightService) HealthCheck() error {
	if s == nil {
		return fmt.Errorf("Flight service not initialized")
	}
	return nil
}
