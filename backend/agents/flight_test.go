package agents

import (
	"os"
	"testing"
)

func TestGetCheapestFlight_WithoutAPIKey(t *testing.T) {
	// Ensure FLIGHT_API_KEY is not set
	os.Unsetenv("FLIGHT_API_KEY")

	price, airline := GetCheapestFlight("BKK", "YVR", "2025-12-01")

	// Should return default values when API key is not set
	if price != 38000 {
		t.Errorf("GetCheapestFlight() price = %d, want 38000", price)
	}
	if airline != "EVA Air" {
		t.Errorf("GetCheapestFlight() airline = %s, want EVA Air", airline)
	}
}

func TestGetCheapestFlight_InvalidInputs(t *testing.T) {
	tests := []struct {
		name    string
		from    string
		to      string
		date    string
		wantErr bool
	}{
		{
			name:    "Empty from",
			from:    "",
			to:      "YVR",
			date:    "2025-12-01",
			wantErr: true,
		},
		{
			name:    "Empty to",
			from:    "BKK",
			to:      "",
			date:    "2025-12-01",
			wantErr: true,
		},
		{
			name:    "Empty date",
			from:    "BKK",
			to:      "YVR",
			date:    "",
			wantErr: true,
		},
		{
			name:    "Invalid date format",
			from:    "BKK",
			to:      "YVR",
			date:    "2025/12/01",
			wantErr: true,
		},
		{
			name:    "Invalid date format - no dashes",
			from:    "BKK",
			to:      "YVR",
			date:    "20251201",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			price, airline := GetCheapestFlight(tt.from, tt.to, tt.date)

			// Should return default values for invalid inputs
			if price != 38000 {
				t.Errorf("GetCheapestFlight() price = %d, want 38000 for invalid input", price)
			}
			if airline != "EVA Air" {
				t.Errorf("GetCheapestFlight() airline = %s, want EVA Air for invalid input", airline)
			}
		})
	}
}

func TestEstimateFlightPrice_KnownRoutes(t *testing.T) {
	tests := []struct {
		name      string
		from      string
		to        string
		wantPrice int
	}{
		{
			name:      "Bangkok to Vancouver",
			from:      "BKK",
			to:        "YVR",
			wantPrice: 38000,
		},
		{
			name:      "Bangkok to Tokyo",
			from:      "BKK",
			to:        "NRT",
			wantPrice: 15000,
		},
		{
			name:      "Bangkok to Seoul",
			from:      "BKK",
			to:        "ICN",
			wantPrice: 12000,
		},
		{
			name:      "Bangkok to Singapore",
			from:      "BKK",
			to:        "SIN",
			wantPrice: 5000,
		},
		{
			name:      "Bangkok to Hong Kong",
			from:      "BKK",
			to:        "HKG",
			wantPrice: 6000,
		},
		{
			name:      "Bangkok to Taipei",
			from:      "BKK",
			to:        "TPE",
			wantPrice: 8000,
		},
		{
			name:      "Bangkok to Kuala Lumpur",
			from:      "BKK",
			to:        "KUL",
			wantPrice: 4000,
		},
		{
			name:      "Bangkok to Jakarta",
			from:      "BKK",
			to:        "CGK",
			wantPrice: 7000,
		},
		{
			name:      "Bangkok to Sydney",
			from:      "BKK",
			to:        "SYD",
			wantPrice: 35000,
		},
		{
			name:      "Bangkok to London",
			from:      "BKK",
			to:        "LHR",
			wantPrice: 45000,
		},
		{
			name:      "Bangkok to Paris",
			from:      "BKK",
			to:        "CDG",
			wantPrice: 42000,
		},
		{
			name:      "Bangkok to Frankfurt",
			from:      "BKK",
			to:        "FRA",
			wantPrice: 40000,
		},
		{
			name:      "Bangkok to Los Angeles",
			from:      "BKK",
			to:        "LAX",
			wantPrice: 50000,
		},
		{
			name:      "Bangkok to New York",
			from:      "BKK",
			to:        "JFK",
			wantPrice: 52000,
		},
		{
			name:      "Bangkok to Dubai",
			from:      "BKK",
			to:        "DXB",
			wantPrice: 25000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			price := estimateFlightPrice(tt.from, tt.to)
			if price != tt.wantPrice {
				t.Errorf("estimateFlightPrice(%s, %s) = %d, want %d", tt.from, tt.to, price, tt.wantPrice)
			}
		})
	}
}

func TestEstimateFlightPrice_ReverseRoutes(t *testing.T) {
	tests := []struct {
		name      string
		from      string
		to        string
		wantPrice int
	}{
		{
			name:      "Vancouver to Bangkok (reverse)",
			from:      "YVR",
			to:        "BKK",
			wantPrice: 38000,
		},
		{
			name:      "Tokyo to Bangkok (reverse)",
			from:      "NRT",
			to:        "BKK",
			wantPrice: 15000,
		},
		{
			name:      "Singapore to Bangkok (reverse)",
			from:      "SIN",
			to:        "BKK",
			wantPrice: 5000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			price := estimateFlightPrice(tt.from, tt.to)
			if price != tt.wantPrice {
				t.Errorf("estimateFlightPrice(%s, %s) = %d, want %d", tt.from, tt.to, price, tt.wantPrice)
			}
		})
	}
}

func TestEstimateFlightPrice_UnknownRoute(t *testing.T) {
	// Test unknown route returns default price
	price := estimateFlightPrice("XXX", "YYY")
	if price != 38000 {
		t.Errorf("estimateFlightPrice(XXX, YYY) = %d, want 38000 (default)", price)
	}
}

func TestEstimateFlightPrice_CaseInsensitive(t *testing.T) {
	tests := []struct {
		name      string
		from      string
		to        string
		wantPrice int
	}{
		{
			name:      "Lowercase codes",
			from:      "bkk",
			to:        "yvr",
			wantPrice: 38000,
		},
		{
			name:      "Mixed case codes",
			from:      "Bkk",
			to:        "Yvr",
			wantPrice: 38000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			price := estimateFlightPrice(tt.from, tt.to)
			if price != tt.wantPrice {
				t.Errorf("estimateFlightPrice(%s, %s) = %d, want %d", tt.from, tt.to, price, tt.wantPrice)
			}
		})
	}
}

func TestEstimateAirline_KnownRoutes(t *testing.T) {
	tests := []struct {
		name        string
		from        string
		to          string
		wantAirline string
	}{
		{
			name:        "Bangkok to Vancouver",
			from:        "BKK",
			to:          "YVR",
			wantAirline: "EVA Air",
		},
		{
			name:        "Bangkok to Tokyo",
			from:        "BKK",
			to:          "NRT",
			wantAirline: "Thai Airways",
		},
		{
			name:        "Bangkok to Seoul",
			from:        "BKK",
			to:          "ICN",
			wantAirline: "Korean Air",
		},
		{
			name:        "Bangkok to Singapore",
			from:        "BKK",
			to:          "SIN",
			wantAirline: "Singapore Airlines",
		},
		{
			name:        "Bangkok to Hong Kong",
			from:        "BKK",
			to:          "HKG",
			wantAirline: "Cathay Pacific",
		},
		{
			name:        "Bangkok to Kuala Lumpur",
			from:        "BKK",
			to:          "KUL",
			wantAirline: "AirAsia",
		},
		{
			name:        "Bangkok to London",
			from:        "BKK",
			to:          "LHR",
			wantAirline: "British Airways",
		},
		{
			name:        "Bangkok to Dubai",
			from:        "BKK",
			to:          "DXB",
			wantAirline: "Emirates",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			airline := estimateAirline(tt.from, tt.to)
			if airline != tt.wantAirline {
				t.Errorf("estimateAirline(%s, %s) = %s, want %s", tt.from, tt.to, airline, tt.wantAirline)
			}
		})
	}
}

func TestEstimateAirline_ReverseRoutes(t *testing.T) {
	tests := []struct {
		name        string
		from        string
		to          string
		wantAirline string
	}{
		{
			name:        "Vancouver to Bangkok (reverse)",
			from:        "YVR",
			to:          "BKK",
			wantAirline: "EVA Air",
		},
		{
			name:        "Tokyo to Bangkok (reverse)",
			from:        "NRT",
			to:          "BKK",
			wantAirline: "Thai Airways",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			airline := estimateAirline(tt.from, tt.to)
			if airline != tt.wantAirline {
				t.Errorf("estimateAirline(%s, %s) = %s, want %s", tt.from, tt.to, airline, tt.wantAirline)
			}
		})
	}
}

func TestEstimateAirline_UnknownRoute(t *testing.T) {
	// Test unknown route returns default airline
	airline := estimateAirline("XXX", "YYY")
	if airline != "EVA Air" {
		t.Errorf("estimateAirline(XXX, YYY) = %s, want EVA Air (default)", airline)
	}
}

func TestEstimateAirline_CaseInsensitive(t *testing.T) {
	tests := []struct {
		name        string
		from        string
		to          string
		wantAirline string
	}{
		{
			name:        "Lowercase codes",
			from:        "bkk",
			to:          "yvr",
			wantAirline: "EVA Air",
		},
		{
			name:        "Mixed case codes",
			from:        "Bkk",
			to:          "Yvr",
			wantAirline: "EVA Air",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			airline := estimateAirline(tt.from, tt.to)
			if airline != tt.wantAirline {
				t.Errorf("estimateAirline(%s, %s) = %s, want %s", tt.from, tt.to, airline, tt.wantAirline)
			}
		})
	}
}

func TestParseSkyscannerResponse_ValidResponse(t *testing.T) {
	// Test with a valid Skyscanner response structure
	response := map[string]interface{}{
		"data": map[string]interface{}{
			"itineraries": []interface{}{
				map[string]interface{}{
					"price": map[string]interface{}{
						"raw": float64(15000),
					},
					"legs": []interface{}{
						map[string]interface{}{
							"carriers": map[string]interface{}{
								"marketing": []interface{}{
									map[string]interface{}{
										"name": "Thai Airways",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	price, airline := parseSkyscannerResponse(response)

	if price != 15000 {
		t.Errorf("parseSkyscannerResponse() price = %d, want 15000", price)
	}
	if airline != "Thai Airways" {
		t.Errorf("parseSkyscannerResponse() airline = %s, want Thai Airways", airline)
	}
}

func TestParseSkyscannerResponse_MultipleFlights(t *testing.T) {
	// Test with multiple flights, should return cheapest
	response := map[string]interface{}{
		"data": map[string]interface{}{
			"itineraries": []interface{}{
				map[string]interface{}{
					"price": map[string]interface{}{
						"raw": float64(25000),
					},
					"legs": []interface{}{
						map[string]interface{}{
							"carriers": map[string]interface{}{
								"marketing": []interface{}{
									map[string]interface{}{
										"name": "Expensive Air",
									},
								},
							},
						},
					},
				},
				map[string]interface{}{
					"price": map[string]interface{}{
						"raw": float64(15000),
					},
					"legs": []interface{}{
						map[string]interface{}{
							"carriers": map[string]interface{}{
								"marketing": []interface{}{
									map[string]interface{}{
										"name": "Budget Air",
									},
								},
							},
						},
					},
				},
				map[string]interface{}{
					"price": map[string]interface{}{
						"raw": float64(20000),
					},
					"legs": []interface{}{
						map[string]interface{}{
							"carriers": map[string]interface{}{
								"marketing": []interface{}{
									map[string]interface{}{
										"name": "Mid-range Air",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	price, airline := parseSkyscannerResponse(response)

	if price != 15000 {
		t.Errorf("parseSkyscannerResponse() price = %d, want 15000 (cheapest)", price)
	}
	if airline != "Budget Air" {
		t.Errorf("parseSkyscannerResponse() airline = %s, want Budget Air", airline)
	}
}

func TestParseSkyscannerResponse_InvalidResponse(t *testing.T) {
	tests := []struct {
		name     string
		response map[string]interface{}
	}{
		{
			name:     "Empty response",
			response: map[string]interface{}{},
		},
		{
			name: "Missing data field",
			response: map[string]interface{}{
				"error": "not found",
			},
		},
		{
			name: "Missing itineraries",
			response: map[string]interface{}{
				"data": map[string]interface{}{},
			},
		},
		{
			name: "Empty itineraries",
			response: map[string]interface{}{
				"data": map[string]interface{}{
					"itineraries": []interface{}{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			price, airline := parseSkyscannerResponse(tt.response)

			if price != 0 {
				t.Errorf("parseSkyscannerResponse() price = %d, want 0 for invalid response", price)
			}
			if airline != "" {
				t.Errorf("parseSkyscannerResponse() airline = %s, want empty string for invalid response", airline)
			}
		})
	}
}

func TestGetCheapestFlight_ValidDate(t *testing.T) {
	// Test with valid date formats
	tests := []struct {
		name string
		date string
	}{
		{
			name: "Future date",
			date: "2025-12-01",
		},
		{
			name: "Near future",
			date: "2025-11-15",
		},
		{
			name: "Far future",
			date: "2026-01-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Without API key, should use default values
			os.Unsetenv("FLIGHT_API_KEY")
			price, airline := GetCheapestFlight("BKK", "NRT", tt.date)

			// Should return default values when API key is not set
			if price != 38000 {
				t.Errorf("GetCheapestFlight() price = %d, want 38000 (default)", price)
			}
			if airline != "EVA Air" {
				t.Errorf("GetCheapestFlight() airline = %s, want EVA Air (default)", airline)
			}
		})
	}
}
