package models

// PlanRequest represents a request to create a travel plan
type PlanRequest struct {
	Message string `json:"message" validate:"required"`
}

// PlanResponse represents a comprehensive travel plan response
type PlanResponse struct {
	Destination  string          `json:"destination"`
	Budget       float64         `json:"budget"`
	DurationDays int             `json:"duration_days"`
	Itinerary    []ItineraryDay  `json:"itinerary"`
	Weather      PlanWeatherInfo `json:"weather"`
	FlightPrice  float64         `json:"flight_price"`
	HotelPrice   float64         `json:"hotel_price"`
}

// ItineraryDay represents a single day in the travel itinerary
type ItineraryDay struct {
	Day      int    `json:"day"`
	Activity string `json:"activity"`
}

// PlanWeatherInfo represents weather information for the destination in plan response
type PlanWeatherInfo struct {
	AvgTemp   float64 `json:"avg_temp"`
	Condition string  `json:"condition"`
}
