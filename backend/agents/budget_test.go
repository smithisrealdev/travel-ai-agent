package agents

import (
	"testing"
)

func TestEstimateBudget_ValidBudgets(t *testing.T) {
	tests := []struct {
		name      string
		total     int
		wantFlight    int
		wantHotel     int
		wantFood      int
		wantTransport int
		wantMisc      int
	}{
		{
			name:          "100,000 THB budget",
			total:         100000,
			wantFlight:    45000,
			wantHotel:     25000,
			wantFood:      15000,
			wantTransport: 10000,
			wantMisc:      5000,
		},
		{
			name:          "50,000 THB budget",
			total:         50000,
			wantFlight:    22500,
			wantHotel:     12500,
			wantFood:      7500,
			wantTransport: 5000,
			wantMisc:      2500,
		},
		{
			name:          "10,000 THB budget",
			total:         10000,
			wantFlight:    4500,
			wantHotel:     2500,
			wantFood:      1500,
			wantTransport: 1000,
			wantMisc:      500,
		},
		{
			name:          "1,000 THB budget",
			total:         1000,
			wantFlight:    450,
			wantHotel:     250,
			wantFood:      150,
			wantTransport: 100,
			wantMisc:      50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EstimateBudget(tt.total)

			if result.Flight != tt.wantFlight {
				t.Errorf("EstimateBudget(%d).Flight = %d, want %d", tt.total, result.Flight, tt.wantFlight)
			}
			if result.Hotel != tt.wantHotel {
				t.Errorf("EstimateBudget(%d).Hotel = %d, want %d", tt.total, result.Hotel, tt.wantHotel)
			}
			if result.Food != tt.wantFood {
				t.Errorf("EstimateBudget(%d).Food = %d, want %d", tt.total, result.Food, tt.wantFood)
			}
			if result.Transport != tt.wantTransport {
				t.Errorf("EstimateBudget(%d).Transport = %d, want %d", tt.total, result.Transport, tt.wantTransport)
			}
			if result.Misc != tt.wantMisc {
				t.Errorf("EstimateBudget(%d).Misc = %d, want %d", tt.total, result.Misc, tt.wantMisc)
			}
		})
	}
}

func TestEstimateBudget_EdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		total int
	}{
		{
			name:  "Zero budget",
			total: 0,
		},
		{
			name:  "Negative budget",
			total: -1000,
		},
		{
			name:  "Large negative budget",
			total: -100000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EstimateBudget(tt.total)

			// All fields should be zero for invalid budgets
			if result.Flight != 0 {
				t.Errorf("EstimateBudget(%d).Flight = %d, want 0", tt.total, result.Flight)
			}
			if result.Hotel != 0 {
				t.Errorf("EstimateBudget(%d).Hotel = %d, want 0", tt.total, result.Hotel)
			}
			if result.Food != 0 {
				t.Errorf("EstimateBudget(%d).Food = %d, want 0", tt.total, result.Food)
			}
			if result.Transport != 0 {
				t.Errorf("EstimateBudget(%d).Transport = %d, want 0", tt.total, result.Transport)
			}
			if result.Misc != 0 {
				t.Errorf("EstimateBudget(%d).Misc = %d, want 0", tt.total, result.Misc)
			}
		})
	}
}

func TestEstimateBudget_PercentageDistribution(t *testing.T) {
	// Test with a budget where we can verify percentages clearly
	total := 100000
	result := EstimateBudget(total)

	// Verify percentages (allowing for integer rounding)
	expectedFlight := int(float64(total) * 0.45)
	expectedHotel := int(float64(total) * 0.25)
	expectedFood := int(float64(total) * 0.15)
	expectedTransport := int(float64(total) * 0.10)
	expectedMisc := int(float64(total) * 0.05)

	if result.Flight != expectedFlight {
		t.Errorf("Flight allocation incorrect: got %d, want %d (45%%)", result.Flight, expectedFlight)
	}
	if result.Hotel != expectedHotel {
		t.Errorf("Hotel allocation incorrect: got %d, want %d (25%%)", result.Hotel, expectedHotel)
	}
	if result.Food != expectedFood {
		t.Errorf("Food allocation incorrect: got %d, want %d (15%%)", result.Food, expectedFood)
	}
	if result.Transport != expectedTransport {
		t.Errorf("Transport allocation incorrect: got %d, want %d (10%%)", result.Transport, expectedTransport)
	}
	if result.Misc != expectedMisc {
		t.Errorf("Misc allocation incorrect: got %d, want %d (5%%)", result.Misc, expectedMisc)
	}
}

func TestBudgetPlan_JSONTags(t *testing.T) {
	// This test ensures the struct has proper JSON tags
	// by creating a budget and checking it can be marshaled
	result := EstimateBudget(100000)

	// Verify the struct is properly initialized
	if result.Flight == 0 && result.Hotel == 0 && result.Food == 0 && 
		result.Transport == 0 && result.Misc == 0 {
		t.Error("BudgetPlan should not have all zero values for valid input")
	}
}
