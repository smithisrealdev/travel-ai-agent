package agents

import (
	"os"
	"testing"
)

func TestGetHotelPrice_WithoutRedis(t *testing.T) {
	// Ensure Redis is not available for this test
	os.Setenv("REDIS_HOST", "invalid-host-123456")
	defer os.Unsetenv("REDIS_HOST")

	price, name := GetHotelPrice("Vancouver", 7)

	// Should return estimated values when Redis is not available
	expectedPrice := 2500 * 7 // Vancouver is 2500 THB/night
	if price != expectedPrice {
		t.Errorf("GetHotelPrice() price = %d, want %d", price, expectedPrice)
	}
	if name == "" {
		t.Errorf("GetHotelPrice() name should not be empty")
	}
}

func TestGetHotelPrice_InvalidInputs(t *testing.T) {
	tests := []struct {
		name        string
		city        string
		nights      int
		wantPrice   int
		description string
	}{
		{
			name:        "Empty city",
			city:        "",
			nights:      5,
			wantPrice:   2500 * 5,
			description: "Should use default price for empty city",
		},
		{
			name:        "Zero nights",
			city:        "Tokyo",
			nights:      0,
			wantPrice:   2500,
			description: "Should use at least 1 night for zero nights",
		},
		{
			name:        "Negative nights",
			city:        "Tokyo",
			nights:      -5,
			wantPrice:   2500,
			description: "Should use at least 1 night for negative nights",
		},
		{
			name:        "Empty city and zero nights",
			city:        "",
			nights:      0,
			wantPrice:   2500,
			description: "Should handle both invalid inputs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure Redis is not available
			os.Setenv("REDIS_HOST", "invalid-host-123456")
			defer os.Unsetenv("REDIS_HOST")

			price, name := GetHotelPrice(tt.city, tt.nights)

			if price != tt.wantPrice {
				t.Errorf("GetHotelPrice(%s, %d) price = %d, want %d: %s", 
					tt.city, tt.nights, price, tt.wantPrice, tt.description)
			}
			if name != "Budget Hotel" {
				t.Errorf("GetHotelPrice(%s, %d) name = %s, want Budget Hotel", 
					tt.city, tt.nights, name)
			}
		})
	}
}

func TestEstimateHotelPricePerNight_KnownCities(t *testing.T) {
	tests := []struct {
		name      string
		city      string
		wantPrice int
	}{
		{name: "Vancouver", city: "Vancouver", wantPrice: 2500},
		{name: "Tokyo", city: "Tokyo", wantPrice: 2000},
		{name: "Seoul", city: "Seoul", wantPrice: 1800},
		{name: "Singapore", city: "Singapore", wantPrice: 2200},
		{name: "Hong Kong", city: "Hong Kong", wantPrice: 2400},
		{name: "Taipei", city: "Taipei", wantPrice: 1600},
		{name: "Kuala Lumpur", city: "Kuala Lumpur", wantPrice: 1200},
		{name: "Jakarta", city: "Jakarta", wantPrice: 1000},
		{name: "Sydney", city: "Sydney", wantPrice: 3000},
		{name: "London", city: "London", wantPrice: 4000},
		{name: "Paris", city: "Paris", wantPrice: 3500},
		{name: "Frankfurt", city: "Frankfurt", wantPrice: 3200},
		{name: "Los Angeles", city: "Los Angeles", wantPrice: 3800},
		{name: "New York", city: "New York", wantPrice: 4500},
		{name: "Dubai", city: "Dubai", wantPrice: 2800},
		{name: "Bangkok", city: "Bangkok", wantPrice: 1000},
		{name: "Phuket", city: "Phuket", wantPrice: 1500},
		{name: "Chiang Mai", city: "Chiang Mai", wantPrice: 800},
		{name: "Pattaya", city: "Pattaya", wantPrice: 1200},
		{name: "Krabi", city: "Krabi", wantPrice: 1400},
		{name: "Osaka", city: "Osaka", wantPrice: 2200},
		{name: "Kyoto", city: "Kyoto", wantPrice: 2400},
		{name: "Busan", city: "Busan", wantPrice: 1600},
		{name: "Bali", city: "Bali", wantPrice: 1300},
		{name: "Hanoi", city: "Hanoi", wantPrice: 900},
		{name: "Ho Chi Minh", city: "Ho Chi Minh", wantPrice: 1100},
		{name: "Phnom Penh", city: "Phnom Penh", wantPrice: 700},
		{name: "Vientiane", city: "Vientiane", wantPrice: 600},
		{name: "Yangon", city: "Yangon", wantPrice: 800},
		{name: "Manila", city: "Manila", wantPrice: 1000},
		{name: "Cebu", city: "Cebu", wantPrice: 900},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			price := estimateHotelPricePerNight(tt.city)
			if price != tt.wantPrice {
				t.Errorf("estimateHotelPricePerNight(%s) = %d, want %d", tt.city, price, tt.wantPrice)
			}
		})
	}
}

func TestEstimateHotelPricePerNight_CaseInsensitive(t *testing.T) {
	tests := []struct {
		name      string
		city      string
		wantPrice int
	}{
		{name: "Lowercase", city: "vancouver", wantPrice: 2500},
		{name: "Uppercase", city: "VANCOUVER", wantPrice: 2500},
		{name: "Mixed case", city: "VaNcOuVeR", wantPrice: 2500},
		{name: "Lowercase Tokyo", city: "tokyo", wantPrice: 2000},
		{name: "Uppercase Tokyo", city: "TOKYO", wantPrice: 2000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			price := estimateHotelPricePerNight(tt.city)
			if price != tt.wantPrice {
				t.Errorf("estimateHotelPricePerNight(%s) = %d, want %d", tt.city, price, tt.wantPrice)
			}
		})
	}
}

func TestEstimateHotelPricePerNight_UnknownCity(t *testing.T) {
	// Test unknown city returns default price
	price := estimateHotelPricePerNight("UnknownCity123")
	if price != 2500 {
		t.Errorf("estimateHotelPricePerNight(UnknownCity123) = %d, want 2500 (default)", price)
	}
}

func TestEstimateHotelPricePerNight_PartialMatch(t *testing.T) {
	tests := []struct {
		name      string
		city      string
		wantPrice int
	}{
		{name: "Vancouver BC", city: "Vancouver BC", wantPrice: 2500},
		{name: "Downtown Vancouver", city: "Downtown Vancouver", wantPrice: 2500},
		{name: "Tokyo Japan", city: "Tokyo Japan", wantPrice: 2000},
		{name: "Bangkok Thailand", city: "Bangkok Thailand", wantPrice: 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			price := estimateHotelPricePerNight(tt.city)
			if price != tt.wantPrice {
				t.Errorf("estimateHotelPricePerNight(%s) = %d, want %d", tt.city, price, tt.wantPrice)
			}
		})
	}
}

func TestEstimateHotelName_KnownCities(t *testing.T) {
	tests := []struct {
		name string
		city string
	}{
		{name: "Vancouver", city: "Vancouver"},
		{name: "Tokyo", city: "Tokyo"},
		{name: "Seoul", city: "Seoul"},
		{name: "Singapore", city: "Singapore"},
		{name: "Hong Kong", city: "Hong Kong"},
		{name: "Bangkok", city: "Bangkok"},
		{name: "Phuket", city: "Phuket"},
		{name: "London", city: "London"},
		{name: "Paris", city: "Paris"},
		{name: "New York", city: "New York"},
		{name: "Sydney", city: "Sydney"},
		{name: "Dubai", city: "Dubai"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name := estimateHotelName(tt.city)
			if name == "" {
				t.Errorf("estimateHotelName(%s) returned empty name", tt.city)
			}
			// Name should be one of the predefined names or a generated one
		})
	}
}

func TestEstimateHotelName_CaseInsensitive(t *testing.T) {
	tests := []struct {
		name string
		city string
	}{
		{name: "Lowercase", city: "vancouver"},
		{name: "Uppercase", city: "VANCOUVER"},
		{name: "Mixed case", city: "VaNcOuVeR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name := estimateHotelName(tt.city)
			if name == "" {
				t.Errorf("estimateHotelName(%s) returned empty name", tt.city)
			}
		})
	}
}

func TestEstimateHotelName_UnknownCity(t *testing.T) {
	name := estimateHotelName("UnknownCity123")
	// Should return a default hotel name format
	if name == "" {
		t.Errorf("estimateHotelName(UnknownCity123) returned empty name")
	}
	// Check if it contains "Budget Hotel"
	if !contains(name, "Budget Hotel") {
		t.Errorf("estimateHotelName(UnknownCity123) = %s, expected to contain 'Budget Hotel'", name)
	}
}

func TestGetHotelPrice_ValidInputs(t *testing.T) {
	tests := []struct {
		name          string
		city          string
		nights        int
		expectedPrice int
	}{
		{
			name:          "Vancouver 7 nights",
			city:          "Vancouver",
			nights:        7,
			expectedPrice: 2500 * 7,
		},
		{
			name:          "Tokyo 5 nights",
			city:          "Tokyo",
			nights:        5,
			expectedPrice: 2000 * 5,
		},
		{
			name:          "Bangkok 3 nights",
			city:          "Bangkok",
			nights:        3,
			expectedPrice: 1000 * 3,
		},
		{
			name:          "London 4 nights",
			city:          "London",
			nights:        4,
			expectedPrice: 4000 * 4,
		},
		{
			name:          "New York 1 night",
			city:          "New York",
			nights:        1,
			expectedPrice: 4500 * 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ensure Redis is not available
			os.Setenv("REDIS_HOST", "invalid-host-123456")
			defer os.Unsetenv("REDIS_HOST")

			price, name := GetHotelPrice(tt.city, tt.nights)

			if price != tt.expectedPrice {
				t.Errorf("GetHotelPrice(%s, %d) price = %d, want %d", 
					tt.city, tt.nights, price, tt.expectedPrice)
			}
			if name == "" {
				t.Errorf("GetHotelPrice(%s, %d) name should not be empty", tt.city, tt.nights)
			}
		})
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		want int
	}{
		{name: "a > b", a: 5, b: 3, want: 5},
		{name: "b > a", a: 3, b: 5, want: 5},
		{name: "a == b", a: 5, b: 5, want: 5},
		{name: "negative numbers", a: -3, b: -5, want: -3},
		{name: "zero and positive", a: 0, b: 5, want: 5},
		{name: "zero and negative", a: 0, b: -5, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := max(tt.a, tt.b)
			if result != tt.want {
				t.Errorf("max(%d, %d) = %d, want %d", tt.a, tt.b, result, tt.want)
			}
		})
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
