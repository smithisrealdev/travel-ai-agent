package models

import (
	"time"
)

// TravelSearchRequest represents a travel search request from the client
type TravelSearchRequest struct {
	Destination string                 `json:"destination"`
	StartDate   string                 `json:"startDate,omitempty"`
	EndDate     string                 `json:"endDate,omitempty"`
	Budget      float64                `json:"budget,omitempty"`
	Preferences map[string]interface{} `json:"preferences,omitempty"`
	UserID      string                 `json:"userId,omitempty"`
}

// TravelSearchResponse represents the response to a travel search
type TravelSearchResponse struct {
	SearchID        int                    `json:"searchId"`
	Destination     string                 `json:"destination"`
	Summary         string                 `json:"summary"`
	Recommendations []TravelRecommendation `json:"recommendations"`
	Weather         *WeatherInfo           `json:"weather,omitempty"`
	Flights         []FlightInfo           `json:"flights,omitempty"`
	EstimatedCost   float64                `json:"estimatedCost"`
	CreatedAt       time.Time              `json:"createdAt"`
}

// TravelRecommendation represents a single travel recommendation
type TravelRecommendation struct {
	ID          int                    `json:"id"`
	Type        string                 `json:"type"` // hotel, activity, restaurant, etc.
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Price       float64                `json:"price,omitempty"`
	Rating      float64                `json:"rating,omitempty"`
	Location    string                 `json:"location,omitempty"`
	ImageURL    string                 `json:"imageUrl,omitempty"`
	URL         string                 `json:"url,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// WeatherInfo represents weather information for a destination
type WeatherInfo struct {
	Temperature float64 `json:"temperature"`
	Description string  `json:"description"`
	Humidity    int     `json:"humidity"`
	WindSpeed   float64 `json:"windSpeed"`
	Icon        string  `json:"icon"`
	Forecast    []DayForecast `json:"forecast,omitempty"`
}

// DayForecast represents a daily weather forecast
type DayForecast struct {
	Date        string  `json:"date"`
	TempMin     float64 `json:"tempMin"`
	TempMax     float64 `json:"tempMax"`
	Description string  `json:"description"`
	Icon        string  `json:"icon"`
}

// FlightInfo represents flight information
type FlightInfo struct {
	FlightNumber string    `json:"flightNumber"`
	Airline      string    `json:"airline"`
	Departure    string    `json:"departure"`
	Arrival      string    `json:"arrival"`
	DepartTime   time.Time `json:"departTime"`
	ArriveTime   time.Time `json:"arriveTime"`
	Duration     string    `json:"duration"`
	Price        float64   `json:"price"`
	Stops        int       `json:"stops"`
}

// TravelSearch represents a stored travel search in the database
type TravelSearch struct {
	ID          int                    `json:"id"`
	UserID      string                 `json:"userId"`
	Destination string                 `json:"destination"`
	StartDate   *time.Time             `json:"startDate,omitempty"`
	EndDate     *time.Time             `json:"endDate,omitempty"`
	Budget      float64                `json:"budget"`
	Preferences map[string]interface{} `json:"preferences"`
	Results     map[string]interface{} `json:"results"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
}

// HealthCheckResponse represents the health check response
type HealthCheckResponse struct {
	Status   string            `json:"status"`
	Services map[string]string `json:"services"`
	Time     time.Time         `json:"time"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}
