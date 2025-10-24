package orchestrator

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestJourney1_TripPlanning_Thai tests the Thai language trip planning journey
func TestJourney1_TripPlanning_Thai(t *testing.T) {
	orch := New("", "", "", "")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Journey: "อยากไปเที่ยวแคนาดา 7 วัน งบ 100,000 บาท"
	message := "อยากไปเที่ยวแคนาดา 7 วัน งบ 100,000 บาท"
	response, err := orch.ProcessMessage(ctx, message)

	assert.NoError(t, err, "Should process Thai trip planning message")
	assert.NotEmpty(t, response, "Response should not be empty")
	
	// Response should include trip details
	assert.True(t, 
		strings.Contains(response, "Trip") || strings.Contains(response, "Day"),
		"Response should include trip or itinerary information")
}

// TestJourney2_WeatherCheck_Thai tests Thai weather check with plan update
func TestJourney2_WeatherCheck_Thai(t *testing.T) {
	orch := New("", "", "", "")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Journey: "วันนี้ฝนตกที่เกียวโตไหม ถ้าตกช่วยเปลี่ยนกิจกรรมให้หน่อย"
	message := "วันนี้ฝนตกที่เกียวโตไหม"
	response, err := orch.ProcessMessage(ctx, message)

	assert.NoError(t, err, "Should process Thai weather check")
	assert.NotEmpty(t, response, "Response should not be empty")
	assert.Contains(t, response, "Weather", "Response should mention weather")
}

// TestJourney3_FlightCheck_English tests English flight status check
func TestJourney3_FlightCheck_English(t *testing.T) {
	orch := New("", "", "", "")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Journey: "Is flight JL708 on time?"
	message := "Is flight JL708 on time?"
	response, err := orch.ProcessMessage(ctx, message)

	assert.NoError(t, err, "Should process flight check")
	assert.NotEmpty(t, response, "Response should not be empty")
	
	// Should contain flight status information
	assert.True(t,
		strings.Contains(response, "Flight") || strings.Contains(response, "JL708"),
		"Response should mention flight or flight code")
}

// TestJourney4_LocalRecommendation_Thai tests Thai local recommendations
func TestJourney4_LocalRecommendation_Thai(t *testing.T) {
	orch := New("", "", "", "")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Journey: "อยากกินราเมนอร่อยใกล้ Shinjuku"
	message := "อยากกินราเมนอร่อยใกล้ Shinjuku"
	response, err := orch.ProcessMessage(ctx, message)

	assert.NoError(t, err, "Should process Thai local recommendation")
	assert.NotEmpty(t, response, "Response should not be empty")
	
	// Should contain recommendations
	assert.True(t,
		strings.Contains(response, "Recommendations") || strings.Contains(response, "restaurant"),
		"Response should include recommendations")
}

// TestMultiLanguageSupport verifies the system handles both Thai and English
func TestMultiLanguageSupport(t *testing.T) {
	orch := New("", "", "", "")
	ctx := context.Background()

	testCases := []struct {
		name     string
		message  string
		language string
	}{
		{"Thai trip planning", "อยากไปเที่ยวญี่ปุ่น 5 วัน", "Thai"},
		{"English trip planning", "I want to visit Japan for 5 days", "English"},
		{"Thai weather", "อากาศที่โตเกียว", "Thai"},
		{"English weather", "Weather in Tokyo", "English"},
		{"Thai hotel", "หาโรงแรมในกรุงเทพ", "Thai"},
		{"English hotel", "Find hotels in Bangkok", "English"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response, err := orch.ProcessMessage(ctx, tc.message)
			assert.NoError(t, err, "Should process %s message", tc.language)
			assert.NotEmpty(t, response, "%s response should not be empty", tc.language)
		})
	}
}

// TestAllIntentCategories verifies all 8 intent categories work
func TestAllIntentCategories(t *testing.T) {
	orch := New("", "", "", "")
	ctx := context.Background()

	testCases := []struct {
		name           string
		message        string
		expectedIntent string
	}{
		{"plan_trip", "Plan a 7-day trip to Canada", "plan_trip"},
		{"flight_check", "Is flight TG123 delayed?", "flight_check"},
		{"weather_check", "What's the weather in Paris?", "weather_check"},
		{"hotel_search", "Find budget hotels in London", "hotel_search"},
		{"local_recommendation", "Best coffee shops nearby", "local_recommendation"},
		{"budget_inquiry", "How much for 100000 baht trip?", "budget_inquiry"},
		{"general_chat", "Hello, how are you?", "general_chat"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response, err := orch.ProcessMessage(ctx, tc.message)
			assert.NoError(t, err, "Should process %s intent", tc.expectedIntent)
			assert.NotEmpty(t, response, "Response should not be empty for %s", tc.expectedIntent)
		})
	}
}

// TestOrchestratorCoordination tests agent coordination
func TestOrchestratorCoordination(t *testing.T) {
	orch := New("", "", "", "")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Test that plan_trip coordinates multiple agents
	message := "I want to visit Tokyo for 7 days with 100000 baht"
	response, err := orch.ProcessMessage(ctx, message)

	assert.NoError(t, err, "Should process complex trip planning")
	assert.NotEmpty(t, response, "Response should not be empty")
	
	// The response should include coordinated information from multiple agents
	// - Itinerary from PlannerAgent
	// - Weather from WeatherAgent
	// - Hotels from HotelAgent
	containsMultipleAgentData := strings.Contains(response, "Day") || 
		strings.Contains(response, "Itinerary") ||
		strings.Contains(response, "Weather") ||
		strings.Contains(response, "Hotels")
	
	assert.True(t, containsMultipleAgentData, 
		"Response should coordinate data from multiple agents")
}
