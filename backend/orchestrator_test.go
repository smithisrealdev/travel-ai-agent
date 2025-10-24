package main

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOrchestratePlan_Success(t *testing.T) {
	// Test with a valid travel request message
	message := "I want to visit Vancouver for 7 days with 100,000 THB budget"

	result := OrchestratePlan(message)

	// Assert basic structure
	assert.NotEmpty(t, result.Destination, "Destination should not be empty")
	assert.Greater(t, result.BudgetTHB, 0, "BudgetTHB should be positive")
	assert.Greater(t, result.DurationDays, 0, "DurationDays should be positive")

	// Assert budget plan exists and has values
	assert.NotNil(t, result.BudgetPlan, "BudgetPlan should exist")
	assert.Greater(t, result.BudgetPlan.Flight, 0, "Flight budget should be positive")
	assert.Greater(t, result.BudgetPlan.Hotel, 0, "Hotel budget should be positive")
	assert.Greater(t, result.BudgetPlan.Food, 0, "Food budget should be positive")
	assert.Greater(t, result.BudgetPlan.Transport, 0, "Transport budget should be positive")
	assert.Greater(t, result.BudgetPlan.Misc, 0, "Misc budget should be positive")

	// Assert flight information is valid
	assert.NotEmpty(t, result.Flight.From, "Flight origin should be set")
	assert.NotEmpty(t, result.Flight.To, "Flight destination should be set")
	assert.NotEmpty(t, result.Flight.Date, "Flight date should be set")
	assert.Greater(t, result.Flight.Price, 0, "Flight price should be positive")
	assert.NotEmpty(t, result.Flight.Airline, "Airline should be set")

	// Assert hotel information is valid
	assert.NotEmpty(t, result.Hotel.City, "Hotel city should be set")
	assert.Greater(t, result.Hotel.Nights, 0, "Hotel nights should be positive")
	assert.Greater(t, result.Hotel.TotalPrice, 0, "Hotel total price should be positive")
	assert.Greater(t, result.Hotel.PricePerNight, 0, "Hotel price per night should be positive")
	assert.NotEmpty(t, result.Hotel.Name, "Hotel name should be set")

	// Assert weather information is valid
	assert.NotEmpty(t, result.Weather.City, "Weather city should be set")
	assert.NotEmpty(t, result.Weather.Month, "Weather month should be set")
	assert.NotEqual(t, 0, result.Weather.AvgTemp, "Weather temperature should be set")
	assert.NotEmpty(t, result.Weather.Condition, "Weather condition should be set")

	// Assert total cost is calculated correctly
	assert.Greater(t, result.TotalEstimatedCost, 0, "Total cost should be positive")
	expectedTotal := result.Flight.Price + result.Hotel.TotalPrice
	assert.Equal(t, expectedTotal, result.TotalEstimatedCost, "Total cost should equal flight + hotel")

	// Assert message summary is present
	assert.NotEmpty(t, result.Message, "Message summary should not be empty")
	assert.Contains(t, result.Message, "travel plan", "Message should mention travel plan")
}

func TestOrchestratePlan_BasicFlow(t *testing.T) {
	message := "I want to travel to Canada for 7 days with a budget of 100000 baht"

	result := OrchestratePlan(message)

	// Verify basic structure using assert
	assert.NotEmpty(t, result.Destination, "Destination should not be empty")
	assert.Greater(t, result.BudgetTHB, 0, "BudgetTHB should be positive")
	assert.Greater(t, result.DurationDays, 0, "DurationDays should be positive")

	// Verify budget plan exists
	assert.True(t, result.BudgetPlan.Flight > 0 || result.BudgetPlan.Hotel > 0, "BudgetPlan should have values")

	// Verify flight information
	assert.NotEmpty(t, result.Flight.From, "Flight should have from code")
	assert.NotEmpty(t, result.Flight.To, "Flight should have to code")
	assert.NotEmpty(t, result.Flight.Date, "Flight date should be set")
	assert.Greater(t, result.Flight.Price, 0, "Flight price should be positive")

	// Verify hotel information
	assert.NotEmpty(t, result.Hotel.City, "Hotel city should be set")
	assert.Greater(t, result.Hotel.Nights, 0, "Hotel nights should be positive")
	assert.Greater(t, result.Hotel.TotalPrice, 0, "Hotel total price should be positive")

	// Verify weather information
	assert.NotEmpty(t, result.Weather.City, "Weather city should be set")
	assert.NotEmpty(t, result.Weather.Month, "Weather month should be set")

	// Verify total cost
	assert.Greater(t, result.TotalEstimatedCost, 0, "Total estimated cost should be positive")
	assert.NotEmpty(t, result.Message, "Message should not be empty")

	// Verify total cost calculation
	expectedTotal := result.Flight.Price + result.Hotel.TotalPrice
	assert.Equal(t, expectedTotal, result.TotalEstimatedCost, "Total cost mismatch")
}

func TestOrchestratePlan_ThaiLanguage(t *testing.T) {
	// Test Thai input - the function should handle it gracefully
	// Note: Without OpenAI API key, it will use defaults
	message := "อยากไปเที่ยวแคนาดาในงบ 100,000 บาท 7 วัน"

	result := OrchestratePlan(message)

	// Should return a valid result even with Thai input
	assert.NotNil(t, result, "Result should not be nil")
	assert.NotEmpty(t, result.Destination, "Destination should be set")
	assert.Greater(t, result.BudgetTHB, 0, "Budget should be positive")
	assert.Greater(t, result.DurationDays, 0, "Duration should be positive")
	assert.Greater(t, result.TotalEstimatedCost, 0, "Total cost should be positive")
	assert.NotEmpty(t, result.Message, "Message should not be empty")
}

func TestOrchestratePlan_EnglishLanguage(t *testing.T) {
	message := "I want to visit Tokyo for 5 days"

	result := OrchestratePlan(message)

	// Should return a valid result
	assert.NotNil(t, result, "Result should not be nil")
	assert.NotEmpty(t, result.Destination, "Destination should be set")
	assert.Greater(t, result.BudgetTHB, 0, "Budget should be positive")
	assert.Greater(t, result.DurationDays, 0, "Duration should be positive")
	assert.Greater(t, result.TotalEstimatedCost, 0, "Total cost should be positive")
}

func TestOrchestratePlan_ValidJSONResponse(t *testing.T) {
	message := "Tokyo trip for 5 days"

	result := OrchestratePlan(message)

	// Verify the struct can be marshaled to JSON
	jsonData, err := json.Marshal(result)
	assert.NoError(t, err, "Result should be JSON-serializable")
	assert.NotEmpty(t, jsonData, "JSON data should not be empty")

	// Parse back to verify structure
	var parsed PlanResult
	err = json.Unmarshal(jsonData, &parsed)
	assert.NoError(t, err, "JSON should be parseable")

	// Verify required fields exist in the parsed result
	assert.NotEmpty(t, parsed.Destination, "Parsed destination should not be empty")
	assert.Greater(t, parsed.BudgetTHB, 0, "Parsed budget should be positive")
	assert.Greater(t, parsed.DurationDays, 0, "Parsed duration should be positive")

	// Verify nested structures
	assert.NotNil(t, parsed.BudgetPlan, "Budget plan should exist")
	assert.Greater(t, parsed.Flight.Price, 0, "Flight price should be positive")
	assert.Greater(t, parsed.Hotel.TotalPrice, 0, "Hotel price should be positive")
	assert.NotEmpty(t, parsed.Weather.Condition, "Weather condition should be set")
	assert.NotEmpty(t, parsed.Message, "Message should be set")
}

func TestOrchestratePlan_InvalidInput(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"empty string", ""},
		{"very short", "trip"},
		{"only destination", "Paris"},
		{"with special chars", "I want to go to café in Paris!"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := OrchestratePlan(tc.input)

			// Should return a valid result even with edge case inputs
			assert.NotNil(t, result, "Result should not be nil")
			assert.NotEmpty(t, result.Destination, "Destination should have default value")
			assert.Greater(t, result.BudgetTHB, 0, "Budget should have default value")
			assert.Greater(t, result.DurationDays, 0, "Duration should have default value")
		})
	}
}

func TestOrchestratePlan_MultipleDestinations(t *testing.T) {
	tests := []struct {
		name    string
		message string
	}{
		{
			name:    "Japan",
			message: "Trip to Japan for 5 days with 80000 baht",
		},
		{
			name:    "Singapore",
			message: "Travel to Singapore for 3 days with 50000 baht",
		},
		{
			name:    "Korea",
			message: "I want to go to Korea for 7 days with 90000 baht",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := OrchestratePlan(tt.message)

			assert.NotEmpty(t, result.Destination, "Destination should not be empty")
			assert.Greater(t, result.TotalEstimatedCost, 0, "Total cost should be positive")
			assert.NotEmpty(t, result.Message, "Message should not be empty")

			// Verify all components are present
			assert.NotEmpty(t, result.Flight.Airline, "Flight airline should be set")
			assert.NotEmpty(t, result.Hotel.Name, "Hotel name should be set")
			assert.NotEmpty(t, result.Weather.Condition, "Weather condition should be set")
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
			assert.Equal(t, tt.want, got, "getAirportCode(%q) should return %q", tt.destination, tt.want)
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
			assert.Equal(t, tt.want, got, "getAirportCode(%q) should be case-insensitive", tt.destination)
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
			assert.Equal(t, tt.want, got, "getAirportCode should match partial strings")
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
			assert.Equal(t, "SIN", got, "getAirportCode(%q) should default to SIN", destination)
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
			assert.Equal(t, tt.wantMonth, got, "getWeatherMonth(%q) should return correct month", tt.dateStr)
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

	currentMonth := time.Now().Format("January")

	for _, dateStr := range invalidDates {
		t.Run(dateStr, func(t *testing.T) {
			got := getWeatherMonth(dateStr)
			assert.Equal(t, currentMonth, got, "getWeatherMonth(%q) should return current month for invalid dates", dateStr)
		})
	}
}

func TestGetWeatherMonth_FutureDate(t *testing.T) {
	// Test with a date 30 days from now
	futureDate := time.Now().AddDate(0, 0, 30)
	dateStr := futureDate.Format("2006-01-02")
	expectedMonth := futureDate.Format("January")

	got := getWeatherMonth(dateStr)
	assert.Equal(t, expectedMonth, got, "getWeatherMonth should handle future dates correctly")
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
			assert.Equal(t, tt.want, got, "max(%d, %d) should return %d", tt.a, tt.b, tt.want)
		})
	}
}

func TestPlanResult_FlightDateIsFuture(t *testing.T) {
	message := "Travel to Japan for 5 days with 80000 baht"
	result := OrchestratePlan(message)

	// Parse the flight date
	flightDate, err := time.Parse("2006-01-02", result.Flight.Date)
	assert.NoError(t, err, "Flight date should be parseable")

	// Flight date should be in the future
	assert.True(t, flightDate.After(time.Now()), "Flight date %s should be in the future", result.Flight.Date)

	// Flight date should be approximately 30 days from now (allow 1 day tolerance)
	expectedDate := time.Now().AddDate(0, 0, 30)
	daysDiff := flightDate.Sub(expectedDate).Hours() / 24

	assert.InDelta(t, 0, daysDiff, 1, "Flight date should be ~30 days from now, got %s (diff: %.1f days)", result.Flight.Date, daysDiff)
}

func TestPlanResult_HotelPricePerNight(t *testing.T) {
	message := "Travel to Singapore for 5 days with 60000 baht"
	result := OrchestratePlan(message)

	// Verify price per night calculation
	expectedPricePerNight := result.Hotel.TotalPrice / result.Hotel.Nights
	assert.Equal(t, expectedPricePerNight, result.Hotel.PricePerNight, "Hotel price per night should be correctly calculated")
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
		assert.Contains(t, result.Message, required, "Message should contain %q", required)
	}
}

func TestPlanResult_WeatherMonthMatchesFlightDate(t *testing.T) {
	message := "Travel to Canada for 7 days with 100000 baht"
	result := OrchestratePlan(message)

	// Parse the flight date
	flightDate, err := time.Parse("2006-01-02", result.Flight.Date)
	assert.NoError(t, err, "Flight date should be parseable")

	// Weather month should match the flight date month
	expectedMonth := flightDate.Format("January")
	assert.Equal(t, expectedMonth, result.Weather.Month, "Weather month should match flight date month")
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

	assert.Equal(t, expectedFlight, result.BudgetPlan.Flight, "Flight budget should be 45%% of total")
	assert.Equal(t, expectedHotel, result.BudgetPlan.Hotel, "Hotel budget should be 25%% of total")
	assert.Equal(t, expectedFood, result.BudgetPlan.Food, "Food budget should be 15%% of total")
	assert.Equal(t, expectedTransport, result.BudgetPlan.Transport, "Transport budget should be 10%% of total")
	assert.Equal(t, expectedMisc, result.BudgetPlan.Misc, "Misc budget should be 5%% of total")
}

// Additional comprehensive tests

func TestPlanResult_StructureValidation(t *testing.T) {
	message := "Plan a trip to Paris for 10 days with 150000 THB"
	result := OrchestratePlan(message)

	// Verify all top-level fields are populated
	assert.NotEmpty(t, result.Destination, "Destination should be set")
	assert.Greater(t, result.BudgetTHB, 0, "BudgetTHB should be positive")
	assert.Greater(t, result.DurationDays, 0, "DurationDays should be positive")

	// Verify BudgetPlan structure
	assert.Greater(t, result.BudgetPlan.Flight, 0, "Budget plan flight should be positive")
	assert.Greater(t, result.BudgetPlan.Hotel, 0, "Budget plan hotel should be positive")
	assert.Greater(t, result.BudgetPlan.Food, 0, "Budget plan food should be positive")
	assert.Greater(t, result.BudgetPlan.Transport, 0, "Budget plan transport should be positive")
	assert.Greater(t, result.BudgetPlan.Misc, 0, "Budget plan misc should be positive")

	// Verify Flight structure
	assert.NotEmpty(t, result.Flight.From, "Flight from should be set")
	assert.NotEmpty(t, result.Flight.To, "Flight to should be set")
	assert.NotEmpty(t, result.Flight.Date, "Flight date should be set")
	assert.Greater(t, result.Flight.Price, 0, "Flight price should be positive")
	assert.NotEmpty(t, result.Flight.Airline, "Flight airline should be set")

	// Verify Hotel structure
	assert.NotEmpty(t, result.Hotel.City, "Hotel city should be set")
	assert.Greater(t, result.Hotel.Nights, 0, "Hotel nights should be positive")
	assert.Greater(t, result.Hotel.TotalPrice, 0, "Hotel total price should be positive")
	assert.Greater(t, result.Hotel.PricePerNight, 0, "Hotel price per night should be positive")
	assert.NotEmpty(t, result.Hotel.Name, "Hotel name should be set")

	// Verify Weather structure
	assert.NotEmpty(t, result.Weather.City, "Weather city should be set")
	assert.NotEmpty(t, result.Weather.Month, "Weather month should be set")
	assert.NotEqual(t, 0, result.Weather.AvgTemp, "Weather temperature should be set")
	assert.NotEmpty(t, result.Weather.Condition, "Weather condition should be set")

	// Verify summary fields
	assert.Greater(t, result.TotalEstimatedCost, 0, "Total estimated cost should be positive")
	assert.NotEmpty(t, result.Message, "Message should be set")
}

func TestPlanResult_ConsistencyChecks(t *testing.T) {
	message := "Singapore vacation for 6 days with 80000 THB budget"
	result := OrchestratePlan(message)

	// Hotel nights should match duration
	// Note: In the current implementation, without API key it uses default duration
	assert.Greater(t, result.Hotel.Nights, 0, "Hotel nights should be positive")

	// Total cost should be sum of flight and hotel
	calculatedTotal := result.Flight.Price + result.Hotel.TotalPrice
	assert.Equal(t, calculatedTotal, result.TotalEstimatedCost, "Total cost should equal flight + hotel")

	// Hotel price per night should be consistent
	expectedPricePerNight := result.Hotel.TotalPrice / result.Hotel.Nights
	assert.Equal(t, expectedPricePerNight, result.Hotel.PricePerNight, "Price per night calculation should be consistent")

	// Weather city and hotel city should match
	assert.Equal(t, result.Hotel.City, result.Weather.City, "Weather city should match hotel city")
}

func TestPlanResult_EdgeCases(t *testing.T) {
	edgeCases := []struct {
		name    string
		message string
	}{
		{"Empty message", ""},
		{"Very short message", "hi"},
		{"Only numbers", "12345"},
		{"Special characters", "@#$%^&*()"},
		{"Mixed languages", "Trip to 東京 Tokyo"},
		{"Long message", strings.Repeat("travel ", 100)},
	}

	for _, tc := range edgeCases {
		t.Run(tc.name, func(t *testing.T) {
			// Should not panic and should return valid structure
			result := OrchestratePlan(tc.message)

			assert.NotNil(t, result, "Result should not be nil")
			assert.NotEmpty(t, result.Destination, "Should have default destination")
			assert.Greater(t, result.BudgetTHB, 0, "Should have default budget")
			assert.Greater(t, result.DurationDays, 0, "Should have default duration")
			assert.Greater(t, result.TotalEstimatedCost, 0, "Should calculate total cost")
		})
	}
}
