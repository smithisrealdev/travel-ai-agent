package agents

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIntentAgent_Detect_Thai(t *testing.T) {
	agent := NewIntentAgent("")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tests := []struct {
		name           string
		input          string
		expectedIntent string
	}{
		{
			name:           "Thai trip planning",
			input:          "อยากไปเที่ยวแคนาดา 7 วัน งบ 100,000 บาท",
			expectedIntent: "plan_trip",
		},
		{
			name:           "Thai weather check",
			input:          "วันนี้ฝนตกที่เกียวโตไหม",
			expectedIntent: "weather_check",
		},
		{
			name:           "Thai hotel search",
			input:          "หาโรงแรมในโตเกียว",
			expectedIntent: "hotel_search",
		},
		{
			name:           "Thai local recommendation",
			input:          "อยากกินราเมนอร่อยใกล้ Shinjuku",
			expectedIntent: "local_recommendation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := agent.Detect(ctx, tt.input)
			assert.NoError(t, err, "Detect should not error")
			assert.NotNil(t, result, "Result should not be nil")
			assert.Equal(t, tt.expectedIntent, result.Intent, "Intent should match")
		})
	}
}

func TestIntentAgent_Detect_English(t *testing.T) {
	agent := NewIntentAgent("")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tests := []struct {
		name           string
		input          string
		expectedIntent string
	}{
		{
			name:           "English trip planning",
			input:          "I want to visit Canada for 7 days with 100,000 baht",
			expectedIntent: "plan_trip",
		},
		{
			name:           "English flight check",
			input:          "Is flight JL708 on time?",
			expectedIntent: "flight_check",
		},
		{
			name:           "English weather check",
			input:          "What's the weather in Tokyo?",
			expectedIntent: "weather_check",
		},
		{
			name:           "English hotel search",
			input:          "Find hotels in Bangkok",
			expectedIntent: "hotel_search",
		},
		{
			name:           "English local recommendation",
			input:          "Looking for good ramen restaurants nearby",
			expectedIntent: "local_recommendation",
		},
		{
			name:           "English budget inquiry",
			input:          "What can I do with 50000 baht budget?",
			expectedIntent: "budget_inquiry",
		},
		{
			name:           "General greeting",
			input:          "Hello!",
			expectedIntent: "general_chat",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := agent.Detect(ctx, tt.input)
			assert.NoError(t, err, "Detect should not error")
			assert.NotNil(t, result, "Result should not be nil")
			assert.Equal(t, tt.expectedIntent, result.Intent, "Intent should match")
		})
	}
}

func TestIntentAgent_EntityExtraction(t *testing.T) {
	agent := NewIntentAgent("")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test plan_trip entity extraction
	result, err := agent.Detect(ctx, "I want to visit Tokyo for 5 days with 80000 baht")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.Entities, "Entities should not be nil")
	
	// In fallback mode, it should at least have some default entities
	assert.NotEmpty(t, result.Entities, "Entities should have some values")
}

func TestIntentAgent_LegacyFunction(t *testing.T) {
	// Test backward compatibility with AnalyzeIntent
	destination, budget, duration := AnalyzeIntent("I want to visit Tokyo for 5 days")
	
	assert.NotEmpty(t, destination, "Destination should be set")
	assert.Greater(t, budget, 0, "Budget should be positive")
	assert.Greater(t, duration, 0, "Duration should be positive")
}
