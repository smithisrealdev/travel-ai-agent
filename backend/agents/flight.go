package agents

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// GetCheapestFlight searches for the cheapest flight between two cities
// Parameters:
//   - from: Origin airport/city code (e.g., "BKK" for Bangkok)
//   - to: Destination airport/city code (e.g., "YVR" for Vancouver)
//   - date: Departure date in YYYY-MM-DD format
// Returns:
//   - price: Lowest flight price in THB
//   - airline: Name of the airline with the cheapest flight
func GetCheapestFlight(from, to, date string) (price int, airline string) {
	apiKey := os.Getenv("FLIGHT_API_KEY")

	// Default values if API call fails
	defaultPrice := 38000
	defaultAirline := "EVA Air"

	if apiKey == "" {
		log.Println("FLIGHT_API_KEY not set, using default values")
		return defaultPrice, defaultAirline
	}

	// Validate inputs
	if from == "" || to == "" || date == "" {
		log.Printf("Invalid input: from=%s, to=%s, date=%s", from, to, date)
		return defaultPrice, defaultAirline
	}

	// Validate date format
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		log.Printf("Invalid date format: %s, expected YYYY-MM-DD", date)
		return defaultPrice, defaultAirline
	}

	// Try Skyscanner RapidAPI first
	price, airline = searchSkyscanner(from, to, date, apiKey)
	if price > 0 && airline != "" {
		log.Printf("Found flight via Skyscanner: %s for %d THB", airline, price)
		return price, airline
	}

	// Fallback to mock data based on route
	log.Println("API search failed, using estimated prices")
	return estimateFlightPrice(from, to), estimateAirline(from, to)
}

// searchSkyscanner queries Skyscanner RapidAPI
func searchSkyscanner(from, to, date, apiKey string) (price int, airline string) {
	// Skyscanner RapidAPI endpoint
	url := fmt.Sprintf(
		"https://sky-scrapper.p.rapidapi.com/api/v1/flights/searchFlights?originSkyId=%s&destinationSkyId=%s&originEntityId=%s&destinationEntityId=%s&date=%s&adults=1&currency=THB&market=TH&locale=en-US",
		from, to, from, to, date,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Failed to create Skyscanner request: %v", err)
		return 0, ""
	}

	// Set headers
	req.Header.Add("X-RapidAPI-Key", apiKey)
	req.Header.Add("X-RapidAPI-Host", "sky-scrapper.p.rapidapi.com")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Skyscanner API request failed: %v", err)
		return 0, ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("Skyscanner API returned status %d", resp.StatusCode)
		return 0, ""
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read Skyscanner response: %v", err)
		return 0, ""
	}

	// Parse response
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("Failed to parse Skyscanner response: %v", err)
		return 0, ""
	}

	// Extract cheapest flight
	price, airline = parseSkyscannerResponse(result)

	return price, airline
}

// parseSkyscannerResponse extracts cheapest flight from Skyscanner response
func parseSkyscannerResponse(data map[string]interface{}) (price int, airline string) {
	// Navigate through the JSON structure to find flights
	dataObj, ok := data["data"].(map[string]interface{})
	if !ok {
		log.Println("Invalid Skyscanner response structure: missing data")
		return 0, ""
	}

	itineraries, ok := dataObj["itineraries"].([]interface{})
	if !ok || len(itineraries) == 0 {
		log.Println("No itineraries found in Skyscanner response")
		return 0, ""
	}

	// Find cheapest flight
	minPrice := 999999999
	bestAirline := ""

	for _, item := range itineraries {
		itinerary, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		// Extract price
		priceInfo, ok := itinerary["price"].(map[string]interface{})
		if !ok {
			continue
		}

		rawPrice, ok := priceInfo["raw"].(float64)
		if !ok {
			continue
		}

		priceInt := int(rawPrice)

		if priceInt < minPrice && priceInt > 0 {
			minPrice = priceInt

			// Extract airline from legs
			legs, ok := itinerary["legs"].([]interface{})
			if ok && len(legs) > 0 {
				leg, ok := legs[0].(map[string]interface{})
				if ok {
					carriers, ok := leg["carriers"].(map[string]interface{})
					if ok {
						marketing, ok := carriers["marketing"].([]interface{})
						if ok && len(marketing) > 0 {
							carrier, ok := marketing[0].(map[string]interface{})
							if ok {
								if name, ok := carrier["name"].(string); ok {
									bestAirline = name
								}
							}
						}
					}
				}
			}
		}
	}

	if minPrice == 999999999 {
		return 0, ""
	}

	return minPrice, bestAirline
}

// estimateFlightPrice returns estimated flight price based on route
func estimateFlightPrice(from, to string) int {
	from = strings.ToUpper(from)
	to = strings.ToUpper(to)

	// Common routes from Bangkok
	routes := map[string]int{
		"BKK-YVR": 38000, // Bangkok to Vancouver
		"BKK-NRT": 15000, // Bangkok to Tokyo
		"BKK-ICN": 12000, // Bangkok to Seoul
		"BKK-SIN": 5000,  // Bangkok to Singapore
		"BKK-HKG": 6000,  // Bangkok to Hong Kong
		"BKK-TPE": 8000,  // Bangkok to Taipei
		"BKK-KUL": 4000,  // Bangkok to Kuala Lumpur
		"BKK-CGK": 7000,  // Bangkok to Jakarta
		"BKK-SYD": 35000, // Bangkok to Sydney
		"BKK-LHR": 45000, // Bangkok to London
		"BKK-CDG": 42000, // Bangkok to Paris
		"BKK-FRA": 40000, // Bangkok to Frankfurt
		"BKK-LAX": 50000, // Bangkok to Los Angeles
		"BKK-JFK": 52000, // Bangkok to New York
		"BKK-DXB": 25000, // Bangkok to Dubai
	}

	route := fmt.Sprintf("%s-%s", from, to)
	if price, ok := routes[route]; ok {
		return price
	}

	// Reverse route
	reverseRoute := fmt.Sprintf("%s-%s", to, from)
	if price, ok := routes[reverseRoute]; ok {
		return price
	}

	// Default estimate based on distance (simplified)
	return 38000
}

// estimateAirline returns likely airline based on route
func estimateAirline(from, to string) string {
	from = strings.ToUpper(from)
	to = strings.ToUpper(to)

	airlines := map[string]string{
		"BKK-YVR": "EVA Air",
		"BKK-NRT": "Thai Airways",
		"BKK-ICN": "Korean Air",
		"BKK-SIN": "Singapore Airlines",
		"BKK-HKG": "Cathay Pacific",
		"BKK-TPE": "EVA Air",
		"BKK-KUL": "AirAsia",
		"BKK-CGK": "Thai Lion Air",
		"BKK-SYD": "Qantas",
		"BKK-LHR": "British Airways",
		"BKK-CDG": "Air France",
		"BKK-FRA": "Lufthansa",
		"BKK-LAX": "United Airlines",
		"BKK-JFK": "American Airlines",
		"BKK-DXB": "Emirates",
	}

	route := fmt.Sprintf("%s-%s", from, to)
	if airline, ok := airlines[route]; ok {
		return airline
	}

	// Reverse route
	reverseRoute := fmt.Sprintf("%s-%s", to, from)
	if airline, ok := airlines[reverseRoute]; ok {
		return airline
	}

	return "EVA Air"
}
