package handlers

import (
	"testing"
)

func TestFormatNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected string
	}{
		{"Small number", 999, "999"},
		{"Thousand", 1000, "1,000"},
		{"Ten thousand", 10000, "10,000"},
		{"Fifty thousand", 50000, "50,000"},
		{"Million", 1234567, "1,234,567"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatNumber(tt.input)
			if result != tt.expected {
				t.Errorf("formatNumber(%d) = %s; want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetStringValue(t *testing.T) {
	m := map[string]interface{}{
		"key1": "value1",
		"key2": "",
		"key3": 123,
	}

	tests := []struct {
		name         string
		key          string
		defaultVal   string
		expected     string
	}{
		{"Existing key", "key1", "default", "value1"},
		{"Empty string key", "key2", "default", "default"},
		{"Non-existent key", "nonexistent", "default", "default"},
		{"Wrong type", "key3", "default", "default"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getStringValue(m, tt.key, tt.defaultVal)
			if result != tt.expected {
				t.Errorf("getStringValue(%s) = %s; want %s", tt.key, result, tt.expected)
			}
		})
	}
}

func TestGetIntValue(t *testing.T) {
	m := map[string]interface{}{
		"key1": float64(42),
		"key2": 100,
		"key3": "not a number",
	}

	tests := []struct {
		name       string
		key        string
		defaultVal int
		expected   int
	}{
		{"Float64 value", "key1", 0, 42},
		{"Int value", "key2", 0, 100},
		{"Wrong type", "key3", 99, 99},
		{"Non-existent key", "nonexistent", 99, 99},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getIntValue(m, tt.key, tt.defaultVal)
			if result != tt.expected {
				t.Errorf("getIntValue(%s) = %d; want %d", tt.key, result, tt.expected)
			}
		})
	}
}

func TestFormatResponseAsMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
	}{
		{
			name:  "Valid trip plan JSON",
			input: `{"destination": "Tokyo", "duration": 5, "total_budget": 80000}`,
			contains: []string{
				"# 🌍 Tokyo Trip Plan",
				"**Duration:** 5 days",
				"**Budget:** 80,000 THB",
				"## 💡 Travel Tips",
			},
		},
		{
			name:  "Non-JSON text",
			input: "This is plain text",
			contains: []string{
				"This is plain text",
			},
		},
		{
			name:  "JSON without destination",
			input: `{"other": "field"}`,
			contains: []string{
				`{"other": "field"}`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatResponseAsMarkdown(tt.input)
			for _, expected := range tt.contains {
				if !contains(result, expected) {
					t.Errorf("formatResponseAsMarkdown() result does not contain %q\nGot: %s", expected, result)
				}
			}
		})
	}
}

func TestFormatTripPlanMarkdown(t *testing.T) {
	plan := map[string]interface{}{
		"destination":  "Chiang Mai",
		"duration":     7,
		"total_budget": 50000.0,
		"itinerary": []interface{}{
			map[string]interface{}{
				"day":    float64(1),
				"budget": 7000.0,
				"activities": []interface{}{
					"Arrival at Chiang Mai Airport",
					"Check into hotel",
				},
			},
		},
		"weather": map[string]interface{}{
			"temperature": 28.5,
			"condition":   "Sunny",
		},
	}

	result := formatTripPlanMarkdown(plan)

	expectedParts := []string{
		"# 🌍 Chiang Mai Trip Plan",
		"**Duration:** 7 days",
		"**Budget:** 50,000 THB",
		"## 📅 Day-by-Day Itinerary",
		"### Day 1 (Budget: 7,000 THB)",
		"- Arrival at Chiang Mai Airport",
		"- Check into hotel",
		"## 🌤️ Weather Forecast",
		"**Current:** 28°C, Sunny",
		"## 💰 Budget Breakdown",
		"## 💡 Travel Tips",
		"*Have a wonderful trip! 🎉*",
	}

	for _, expected := range expectedParts {
		if !contains(result, expected) {
			t.Errorf("formatTripPlanMarkdown() result does not contain %q\nGot: %s", expected, result)
		}
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())
}
