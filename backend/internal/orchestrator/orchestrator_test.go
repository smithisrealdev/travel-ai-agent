package orchestrator

import (
"context"
"strings"
"testing"
"time"

"github.com/stretchr/testify/assert"
)

func TestOrchestrator_ProcessMessage_PlanTrip(t *testing.T) {
// Test with no API keys (fallback mode)
orch := New("", "", "", "")

ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

message := "I want to visit Vancouver for 7 days with 100,000 THB budget"
response, err := orch.ProcessMessage(ctx, message)

assert.NoError(t, err, "ProcessMessage should not error")
assert.NotEmpty(t, response, "Response should not be empty")
assert.Contains(t, response, "Trip", "Response should mention trip")
}

func TestOrchestrator_ProcessMessage_WeatherCheck(t *testing.T) {
orch := New("", "", "", "")

ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

message := "What's the weather in Bangkok?"
response, err := orch.ProcessMessage(ctx, message)

assert.NoError(t, err, "ProcessMessage should not error")
assert.NotEmpty(t, response, "Response should not be empty")
assert.Contains(t, response, "Weather", "Response should mention weather")
}

func TestOrchestrator_ProcessMessage_FlightCheck(t *testing.T) {
orch := New("", "", "", "")

ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

message := "Is flight JL708 on time?"
response, err := orch.ProcessMessage(ctx, message)

assert.NoError(t, err, "ProcessMessage should not error")
assert.NotEmpty(t, response, "Response should not be empty")
assert.Contains(t, response, "Flight", "Response should mention flight")
}

func TestOrchestrator_ProcessMessage_HotelSearch(t *testing.T) {
orch := New("", "", "", "")

ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

message := "Find hotels in Tokyo"
response, err := orch.ProcessMessage(ctx, message)

assert.NoError(t, err, "ProcessMessage should not error")
assert.NotEmpty(t, response, "Response should not be empty")
assert.Contains(t, response, "Hotels", "Response should mention hotels")
}

func TestOrchestrator_ProcessMessage_LocalRecommendation(t *testing.T) {
orch := New("", "", "", "")

ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

message := "Find good ramen restaurants nearby"
response, err := orch.ProcessMessage(ctx, message)

assert.NoError(t, err, "ProcessMessage should not error")
assert.NotEmpty(t, response, "Response should not be empty")
// Response should mention recommendations or ramen
containsRecommendations := strings.Contains(response, "Recommendations") || strings.Contains(response, "ramen") || strings.Contains(response, "Ramen")
assert.True(t, containsRecommendations, "Response should mention recommendations or ramen")
}

func TestOrchestrator_ProcessMessage_BudgetInquiry(t *testing.T) {
orch := New("", "", "", "")

ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

message := "What can I do with 50000 baht budget?"
response, err := orch.ProcessMessage(ctx, message)

assert.NoError(t, err, "ProcessMessage should not error")
assert.NotEmpty(t, response, "Response should not be empty")
assert.Contains(t, response, "Budget", "Response should mention budget")
}

func TestOrchestrator_ProcessMessage_GeneralChat(t *testing.T) {
orch := New("", "", "", "")

ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

message := "Hello!"
response, err := orch.ProcessMessage(ctx, message)

assert.NoError(t, err, "ProcessMessage should not error")
assert.NotEmpty(t, response, "Response should not be empty")
assert.Contains(t, response, "assistant", "Response should mention assistant")
}

func TestOrchestrator_EntityExtraction(t *testing.T) {
orch := New("", "", "", "")

// Test getStringEntity
entities := map[string]interface{}{
"destination": "Tokyo",
"budget":      100000.0,
"duration":    7.0,
}

dest := orch.getStringEntity(entities, "destination", "Unknown")
assert.Equal(t, "Tokyo", dest, "Should extract destination")

budget := orch.getFloatEntity(entities, "budget", 0)
assert.Equal(t, 100000.0, budget, "Should extract budget")

duration := orch.getIntEntity(entities, "duration", 0)
assert.Equal(t, 7, duration, "Should extract duration")

// Test default values
missing := orch.getStringEntity(entities, "missing", "Default")
assert.Equal(t, "Default", missing, "Should return default for missing entity")
}
