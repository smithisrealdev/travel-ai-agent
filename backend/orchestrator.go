package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/smithisrealdev/travel-ai-agent/backend/agents"
)

// PlanResult represents a complete travel plan
type PlanResult struct {
	// User Intent
	Destination  string `json:"destination"`
	BudgetTHB    int    `json:"budget_thb"`
	DurationDays int    `json:"duration_days"`

	// Budget Breakdown
	BudgetPlan agents.BudgetPlan `json:"budget_plan"`

	// Flight Information
	Flight struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Date    string `json:"date"`
		Price   int    `json:"price"`
		Airline string `json:"airline"`
	} `json:"flight"`

	// Hotel Information
	Hotel struct {
		City          string `json:"city"`
		Nights        int    `json:"nights"`
		TotalPrice    int    `json:"total_price"`
		PricePerNight int    `json:"price_per_night"`
		Name          string `json:"name"`
	} `json:"hotel"`

	// Weather Information
	Weather struct {
		City      string `json:"city"`
		Month     string `json:"month"`
		AvgTemp   int    `json:"avg_temp"`
		Condition string `json:"condition"`
	} `json:"weather"`

	// Summary
	TotalEstimatedCost int    `json:"total_estimated_cost"`
	Message            string `json:"message"`
}

// OrchestratePlan orchestrates all agents to create a complete travel plan
// Parameters:
//   - message: User's travel request (e.g., "I want to travel to Canada for 7 days with 100000 baht")
//
// Returns:
//   - PlanResult: Complete travel plan with all details
func OrchestratePlan(message string) PlanResult {
	log.Printf("🚀 Orchestrating travel plan for message: %s", message)

	var result PlanResult

	// Step 1: Analyze Intent
	log.Println("Step 1: Analyzing user intent...")
	destination, budgetTHB, durationDays := agents.AnalyzeIntent(message)

	result.Destination = destination
	result.BudgetTHB = budgetTHB
	result.DurationDays = durationDays

	log.Printf("✅ Intent: destination=%s, budget=%d THB, duration=%d days", destination, budgetTHB, durationDays)

	// Step 2: Estimate Budget Breakdown
	log.Println("Step 2: Estimating budget breakdown...")
	budgetPlan := agents.EstimateBudget(budgetTHB)
	result.BudgetPlan = budgetPlan

	log.Printf("✅ Budget Plan: Flight=%d, Hotel=%d, Food=%d, Transport=%d, Misc=%d",
		budgetPlan.Flight, budgetPlan.Hotel, budgetPlan.Food, budgetPlan.Transport, budgetPlan.Misc)

	// Step 3: Get Cheapest Flight (BKK to destination, 30 days from now)
	log.Println("Step 3: Searching for cheapest flight...")
	departureDate := time.Now().AddDate(0, 0, 30).Format("2006-01-02")
	destinationCode := getAirportCode(destination)

	flightPrice, airline := agents.GetCheapestFlight("BKK", destinationCode, departureDate)

	result.Flight.From = "BKK"
	result.Flight.To = destinationCode
	result.Flight.Date = departureDate
	result.Flight.Price = flightPrice
	result.Flight.Airline = airline

	log.Printf("✅ Flight: %s to %s on %s, %s for %d THB", "BKK", destinationCode, departureDate, airline, flightPrice)

	// Step 4: Get Hotel Price
	log.Println("Step 4: Searching for hotel prices...")
	hotelPrice, hotelName := agents.GetHotelPrice(destination, durationDays)

	result.Hotel.City = destination
	result.Hotel.Nights = durationDays
	result.Hotel.TotalPrice = hotelPrice
	result.Hotel.PricePerNight = hotelPrice / max(1, durationDays)
	result.Hotel.Name = hotelName

	log.Printf("✅ Hotel: %s in %s for %d nights, %d THB total (%d THB/night)",
		hotelName, destination, durationDays, hotelPrice, result.Hotel.PricePerNight)

	// Step 5: Get Weather Summary (June or current month + 1)
	log.Println("Step 5: Fetching weather information...")
	weatherMonth := getWeatherMonth(departureDate)
	avgTemp, condition := agents.GetWeatherSummary(destination, weatherMonth)

	result.Weather.City = destination
	result.Weather.Month = weatherMonth
	result.Weather.AvgTemp = avgTemp
	result.Weather.Condition = condition

	log.Printf("✅ Weather: %s in %s - %d°C, %s", destination, weatherMonth, avgTemp, condition)

	// Step 6: Calculate Total Cost
	totalCost := flightPrice + hotelPrice
	result.TotalEstimatedCost = totalCost

	// Step 7: Generate Summary Message
	result.Message = fmt.Sprintf(
		"Complete travel plan for %s: %d days trip costs approximately %d THB (Flight: %d THB, Hotel: %d THB). "+
			"Weather in %s: %d°C, %s. Budget breakdown: Flight %d%%, Hotel %d%%, Food %d%%, Transport %d%%, Misc %d%%.",
		destination, durationDays, totalCost, flightPrice, hotelPrice,
		weatherMonth, avgTemp, condition,
		45, 25, 15, 10, 5,
	)

	log.Printf("🎉 Plan orchestration complete! Total cost: %d THB", totalCost)

	return result
}

// getAirportCode converts city/country name to airport code
func getAirportCode(destination string) string {
	destination = strings.ToLower(destination)

	airportCodes := map[string]string{
		"canada":       "YVR", // Vancouver
		"vancouver":    "YVR",
		"japan":        "NRT", // Tokyo Narita
		"tokyo":        "NRT",
		"korea":        "ICN", // Seoul Incheon
		"seoul":        "ICN",
		"singapore":    "SIN",
		"hong kong":    "HKG",
		"taipei":       "TPE",
		"taiwan":       "TPE",
		"malaysia":     "KUL", // Kuala Lumpur
		"kuala lumpur": "KUL",
		"indonesia":    "CGK", // Jakarta
		"jakarta":      "CGK",
		"australia":    "SYD", // Sydney
		"sydney":       "SYD",
		"uk":           "LHR", // London
		"london":       "LHR",
		"england":      "LHR",
		"france":       "CDG", // Paris
		"paris":        "CDG",
		"germany":      "FRA", // Frankfurt
		"frankfurt":    "FRA",
		"usa":          "LAX", // Los Angeles
		"america":      "LAX",
		"los angeles":  "LAX",
		"new york":     "JFK",
		"uae":          "DXB", // Dubai
		"dubai":        "DXB",
		"thailand":     "BKK", // Bangkok
		"bangkok":      "BKK",
		"phuket":       "HKT",
		"chiang mai":   "CNX",
		"vietnam":      "SGN", // Ho Chi Minh
		"hanoi":        "HAN",
		"ho chi minh":  "SGN",
		"philippines":  "MNL", // Manila
		"manila":       "MNL",
		"india":        "DEL", // Delhi
		"delhi":        "DEL",
		"china":        "PEK", // Beijing
		"beijing":      "PEK",
		"shanghai":     "PVG",
	}

	// Check exact match
	if code, ok := airportCodes[destination]; ok {
		return code
	}

	// Check partial match (only if destination is not empty)
	if destination != "" {
		for key, code := range airportCodes {
			if key != "" && (strings.Contains(destination, key) || strings.Contains(key, destination)) {
				return code
			}
		}
	}

	// Default to major hub
	return "SIN"
}

// getWeatherMonth extracts month from date or returns current month
func getWeatherMonth(dateStr string) string {
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		// Default to current month
		return time.Now().Format("January")
	}

	return parsedDate.Format("January")
}

// Helper functions
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
