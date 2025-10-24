package main

import (
	"strings"
	"testing"
	"time"
)

func TestOrchestratePlan_BasicFlow(t *testing.T) {
	message := "I want to travel to Canada for 7 days with a budget of 100000 baht"

	result := OrchestratePlan(message)

	// Verify basic structure
	if result.Destination == "" {
		t.Error("Destination should not be empty")
	}

	if result.BudgetTHB <= 0 {
		t.Error("BudgetTHB should be positive")
	}

	if result.DurationDays <= 0 {
		t.Error("DurationDays should be positive")
	}

	// Verify budget plan exists
	if result.BudgetPlan.Flight == 0 && result.BudgetPlan.Hotel == 0 {
		t.Error("BudgetPlan should have values")
	}

	// Verify flight information
	if result.Flight.From == "" || result.Flight.To == "" {
		t.Error("Flight should have from and to codes")
	}

	if result.Flight.Date == "" {
		t.Error("Flight date should be set")
	}

	if result.Flight.Price <= 0 {
		t.Error("Flight price should be positive")
	}

	// Verify hotel information
	if result.Hotel.City == "" {
		t.Error("Hotel city should be set")
	}

	if result.Hotel.Nights <= 0 {
		t.Error("Hotel nights should be positive")
	}

	if result.Hotel.TotalPrice <= 0 {
		t.Error("Hotel total price should be positive")
	}

	// Verify weather information
	if result.Weather.City == "" {
		t.Error("Weather city should be set")
	}

	if result.Weather.Month == "" {
		t.Error("Weather month should be set")
	}

	// Verify total cost
	if result.TotalEstimatedCost <= 0 {
		t.Error("Total estimated cost should be positive")
	}

	// Verify message
	if result.Message == "" {
		t.Error("Message should not be empty")
	}

	// Verify total cost calculation
	expectedTotal := result.Flight.Price + result.Hotel.TotalPrice
	if result.TotalEstimatedCost != expectedTotal {
		t.Errorf("Total cost mismatch: got %d, want %d", result.TotalEstimatedCost, expectedTotal)
	}
}

func TestOrchestratePlan_MultipleDestinations(t *testing.T) {
	tests := []struct {
		name        string
		message     string
		wantContain string
	}{
		{
			name:        "Japan",
			message:     "Trip to Japan for 5 days with 80000 baht",
			wantContain: "japan",
		},
		{
			name:        "Singapore",
			message:     "Travel to Singapore for 3 days with 50000 baht",
			wantContain: "singapore",
		},
		{
			name:        "Korea",
			message:     "I want to go to Korea for 7 days with 90000 baht",
			wantContain: "korea",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := OrchestratePlan(tt.message)

			if result.Destination == "" {
				t.Error("Destination should not be empty")
			}

			// Verify basic fields are populated
			if result.TotalEstimatedCost <= 0 {
				t.Error("Total cost should be positive")
			}

			if result.Message == "" {
				t.Error("Message should not be empty")
			}
		})
	}
}

func TestGetAirportCode_ExactMatch(t *testing.T) {
	tests := []struct {
		name        string
		destination string
		want        string
	}{
		{"Canada", "Canada", "YVR"},
		{"Vancouver", "Vancouver", "YVR"},
		{"Japan", "Japan", "NRT"},
		{"Tokyo", "Tokyo", "NRT"},
		{"Korea", "Korea", "ICN"},
		{"Seoul", "Seoul", "ICN"},
		{"Singapore", "Singapore", "SIN"},
		{"Hong Kong", "Hong Kong", "HKG"},
		{"Taipei", "Taipei", "TPE"},
		{"Taiwan", "Taiwan", "TPE"},
		{"Malaysia", "Malaysia", "KUL"},
		{"Kuala Lumpur", "Kuala Lumpur", "KUL"},
		{"Indonesia", "Indonesia", "CGK"},
		{"Jakarta", "Jakarta", "CGK"},
		{"Australia", "Australia", "SYD"},
		{"Sydney", "Sydney", "SYD"},
		{"UK", "UK", "LHR"},
		{"London", "London", "LHR"},
		{"England", "England", "LHR"},
		{"France", "France", "CDG"},
		{"Paris", "Paris", "CDG"},
		{"Germany", "Germany", "FRA"},
		{"Frankfurt", "Frankfurt", "FRA"},
		{"USA", "USA", "LAX"},
		{"America", "America", "LAX"},
		{"Los Angeles", "Los Angeles", "LAX"},
		{"New York", "New York", "JFK"},
		{"UAE", "UAE", "DXB"},
		{"Dubai", "Dubai", "DXB"},
		{"Thailand", "Thailand", "BKK"},
		{"Bangkok", "Bangkok", "BKK"},
		{"Phuket", "Phuket", "HKT"},
		{"Chiang Mai", "Chiang Mai", "CNX"},
		{"Vietnam", "Vietnam", "SGN"},
		{"Hanoi", "Hanoi", "HAN"},
		{"Ho Chi Minh", "Ho Chi Minh", "SGN"},
		{"Philippines", "Philippines", "MNL"},
		{"Manila", "Manila", "MNL"},
		{"India", "India", "DEL"},
		{"Delhi", "Delhi", "DEL"},
		{"China", "China", "PEK"},
		{"Beijing", "Beijing", "PEK"},
		{"Shanghai", "Shanghai", "PVG"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getAirportCode(tt.destination)
			if got != tt.want {
				t.Errorf("getAirportCode(%q) = %q, want %q", tt.destination, got, tt.want)
			}
		})
	}
}

func TestGetAirportCode_CaseInsensitive(t *testing.T) {
	tests := []struct {
		name        string
		destination string
		want        string
	}{
		{"Uppercase CANADA", "CANADA", "YVR"},
		{"Lowercase canada", "canada", "YVR"},
		{"Mixed Case CaNaDa", "CaNaDa", "YVR"},
		{"Uppercase JAPAN", "JAPAN", "NRT"},
		{"Mixed Case JaPaN", "JaPaN", "NRT"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getAirportCode(tt.destination)
			if got != tt.want {
				t.Errorf("getAirportCode(%q) = %q, want %q", tt.destination, got, tt.want)
			}
		})
	}
}

func TestGetAirportCode_PartialMatch(t *testing.T) {
	tests := []struct {
		name        string
		destination string
		want        string
	}{
		{"Contains Canada", "I love Canada", "YVR"},
		{"Contains Japan", "Visit Japan now", "NRT"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getAirportCode(tt.destination)
			if got != tt.want {
				t.Errorf("getAirportCode(%q) = %q, want %q", tt.destination, got, tt.want)
			}
		})
	}
}

func TestGetAirportCode_UnknownDestination(t *testing.T) {
	// Unknown destinations should default to SIN (Singapore)
	unknownDestinations := []string{
		"Unknown Place",
		"Mars",
		"Atlantis",
		"",
	}

	for _, destination := range unknownDestinations {
		t.Run(destination, func(t *testing.T) {
			got := getAirportCode(destination)
			if got != "SIN" {
				t.Errorf("getAirportCode(%q) = %q, want SIN (default)", destination, got)
			}
		})
	}
}

func TestGetWeatherMonth_ValidDate(t *testing.T) {
	tests := []struct {
		name      string
		dateStr   string
		wantMonth string
	}{
		{"January date", "2025-01-15", "January"},
		{"February date", "2025-02-20", "February"},
		{"March date", "2025-03-10", "March"},
		{"April date", "2025-04-05", "April"},
		{"May date", "2025-05-12", "May"},
		{"June date", "2025-06-18", "June"},
		{"July date", "2025-07-22", "July"},
		{"August date", "2025-08-30", "August"},
		{"September date", "2025-09-14", "September"},
		{"October date", "2025-10-08", "October"},
		{"November date", "2025-11-23", "November"},
		{"December date", "2025-12-25", "December"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getWeatherMonth(tt.dateStr)
			if got != tt.wantMonth {
				t.Errorf("getWeatherMonth(%q) = %q, want %q", tt.dateStr, got, tt.wantMonth)
			}
		})
	}
}

func TestGetWeatherMonth_InvalidDate(t *testing.T) {
	invalidDates := []string{
		"invalid-date",
		"2025/01/15",
		"15-01-2025",
		"",
		"not a date",
	}

	for _, dateStr := range invalidDates {
		t.Run(dateStr, func(t *testing.T) {
			got := getWeatherMonth(dateStr)
			// Should return current month name
			currentMonth := time.Now().Format("January")
			if got != currentMonth {
				t.Errorf("getWeatherMonth(%q) = %q, want %q (current month)", dateStr, got, currentMonth)
			}
		})
	}
}

func TestGetWeatherMonth_FutureDate(t *testing.T) {
	// Test with a date 30 days from now
	futureDate := time.Now().AddDate(0, 0, 30)
	dateStr := futureDate.Format("2006-01-02")
	expectedMonth := futureDate.Format("January")

	got := getWeatherMonth(dateStr)
	if got != expectedMonth {
		t.Errorf("getWeatherMonth(%q) = %q, want %q", dateStr, got, expectedMonth)
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		want int
	}{
		{"a greater", 10, 5, 10},
		{"b greater", 5, 10, 10},
		{"equal", 7, 7, 7},
		{"negative numbers", -5, -10, -5},
		{"zero and positive", 0, 5, 5},
		{"zero and negative", 0, -5, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := max(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("max(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestPlanResult_FlightDateIsFuture(t *testing.T) {
	message := "Travel to Japan for 5 days with 80000 baht"
	result := OrchestratePlan(message)

	// Parse the flight date
	flightDate, err := time.Parse("2006-01-02", result.Flight.Date)
	if err != nil {
		t.Fatalf("Failed to parse flight date: %v", err)
	}

	// Flight date should be in the future
	if !flightDate.After(time.Now()) {
		t.Errorf("Flight date %s should be in the future", result.Flight.Date)
	}

	// Flight date should be approximately 30 days from now (allow 1 day tolerance)
	expectedDate := time.Now().AddDate(0, 0, 30)
	daysDiff := flightDate.Sub(expectedDate).Hours() / 24

	if daysDiff < -1 || daysDiff > 1 {
		t.Errorf("Flight date should be ~30 days from now, got %s (diff: %.1f days)", result.Flight.Date, daysDiff)
	}
}

func TestPlanResult_HotelPricePerNight(t *testing.T) {
	message := "Travel to Singapore for 5 days with 60000 baht"
	result := OrchestratePlan(message)

	// Verify price per night calculation
	expectedPricePerNight := result.Hotel.TotalPrice / result.Hotel.Nights
	if result.Hotel.PricePerNight != expectedPricePerNight {
		t.Errorf("Hotel price per night mismatch: got %d, want %d",
			result.Hotel.PricePerNight, expectedPricePerNight)
	}
}

func TestPlanResult_MessageFormat(t *testing.T) {
	message := "Travel to Korea for 7 days with 90000 baht"
	result := OrchestratePlan(message)

	// Verify message contains key information
	requiredStrings := []string{
		"Complete travel plan",
		"days trip",
		"THB",
		"Flight:",
		"Hotel:",
		"Weather",
		"Budget breakdown",
	}

	for _, required := range requiredStrings {
		if !strings.Contains(result.Message, required) {
			t.Errorf("Message should contain %q, got: %s", required, result.Message)
		}
	}
}

func TestPlanResult_WeatherMonthMatchesFlightDate(t *testing.T) {
	message := "Travel to Canada for 7 days with 100000 baht"
	result := OrchestratePlan(message)

	// Parse the flight date
	flightDate, err := time.Parse("2006-01-02", result.Flight.Date)
	if err != nil {
		t.Fatalf("Failed to parse flight date: %v", err)
	}

	// Weather month should match the flight date month
	expectedMonth := flightDate.Format("January")
	if result.Weather.Month != expectedMonth {
		t.Errorf("Weather month should be %s to match flight date, got %s",
			expectedMonth, result.Weather.Month)
	}
}

func TestPlanResult_BudgetBreakdownPercentages(t *testing.T) {
	message := "Travel to Japan for 5 days with 100000 baht"
	result := OrchestratePlan(message)

	// Verify budget breakdown follows the 45/25/15/10/5 split
	budget := result.BudgetTHB

	expectedFlight := int(float64(budget) * 0.45)
	expectedHotel := int(float64(budget) * 0.25)
	expectedFood := int(float64(budget) * 0.15)
	expectedTransport := int(float64(budget) * 0.10)
	expectedMisc := int(float64(budget) * 0.05)

	if result.BudgetPlan.Flight != expectedFlight {
		t.Errorf("Flight budget: got %d, want %d", result.BudgetPlan.Flight, expectedFlight)
	}
	if result.BudgetPlan.Hotel != expectedHotel {
		t.Errorf("Hotel budget: got %d, want %d", result.BudgetPlan.Hotel, expectedHotel)
	}
	if result.BudgetPlan.Food != expectedFood {
		t.Errorf("Food budget: got %d, want %d", result.BudgetPlan.Food, expectedFood)
	}
	if result.BudgetPlan.Transport != expectedTransport {
		t.Errorf("Transport budget: got %d, want %d", result.BudgetPlan.Transport, expectedTransport)
	}
	if result.BudgetPlan.Misc != expectedMisc {
		t.Errorf("Misc budget: got %d, want %d", result.BudgetPlan.Misc, expectedMisc)
	}
}
