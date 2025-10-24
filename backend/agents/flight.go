package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
)

// FlightStatus represents complete flight status information
type FlightStatus struct {
	FlightCode    string `json:"flight_code"`
	Status        string `json:"status"`
	DepartureTime string `json:"departure_time"`
	ArrivalTime   string `json:"arrival_time"`
	Gate          string `json:"gate"`
	DelayMinutes  int    `json:"delay_minutes"`
	Notification  string `json:"notification"`
}

// FlightAgent handles flight tracking and status checking
type FlightAgent struct {
	client *openai.Client
	apiKey string
}

// NewFlightAgent creates a new flight agent
func NewFlightAgent(openaiKey, flightKey string) *FlightAgent {
	var client *openai.Client
	if openaiKey != "" {
		client = openai.NewClient(openaiKey)
	}
	return &FlightAgent{
		client: client,
		apiKey: flightKey,
	}
}

// CheckFlight checks flight status and generates notifications
func (a *FlightAgent) CheckFlight(ctx context.Context, flightCode string) (*FlightStatus, error) {
	status := &FlightStatus{
		FlightCode: flightCode,
		Status:     "unknown",
	}

	// Try AviationStack API if available
	if a.apiKey != "" {
		apiStatus := a.fetchFlightStatusFromAPI(flightCode)
		if apiStatus != nil {
			status = apiStatus
		}
	}

	// Fallback to mock data
	if status.Status == "unknown" {
		status = a.mockFlightStatus(flightCode)
	}

	// Generate notification if delayed
	if status.DelayMinutes > 0 {
		status.Notification = a.generateDelayNotification(ctx, status)
	} else {
		status.Notification = fmt.Sprintf("Flight %s is %s.", flightCode, status.Status)
	}

	log.Printf("FlightAgent: Checked %s - Status: %s, Delay: %d min", 
		flightCode, status.Status, status.DelayMinutes)

	return status, nil
}

// fetchFlightStatusFromAPI queries AviationStack API
func (a *FlightAgent) fetchFlightStatusFromAPI(flightCode string) *FlightStatus {
	url := fmt.Sprintf(
		"https://api.aviationstack.com/v1/flights?access_key=%s&flight_iata=%s",
		a.apiKey, flightCode,
	)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("FlightAgent: API request failed: %v", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("FlightAgent: API returned status %d", resp.StatusCode)
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("FlightAgent: Failed to read response: %v", err)
		return nil
	}

	var result struct {
		Data []struct {
			FlightStatus string `json:"flight_status"`
			Departure    struct {
				Scheduled string `json:"scheduled"`
				Actual    string `json:"actual"`
				Gate      string `json:"gate"`
				Delay     int    `json:"delay"`
			} `json:"departure"`
			Arrival struct {
				Scheduled string `json:"scheduled"`
				Actual    string `json:"actual"`
			} `json:"arrival"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("FlightAgent: Failed to parse response: %v", err)
		return nil
	}

	if len(result.Data) == 0 {
		return nil
	}

	flight := result.Data[0]
	status := &FlightStatus{
		FlightCode:    flightCode,
		Status:        flight.FlightStatus,
		DepartureTime: flight.Departure.Scheduled,
		ArrivalTime:   flight.Arrival.Scheduled,
		Gate:          flight.Departure.Gate,
		DelayMinutes:  flight.Departure.Delay,
	}

	if flight.Departure.Delay > 0 {
		status.Status = "delayed"
	}

	return status
}

// mockFlightStatus generates mock flight status for testing
func (a *FlightAgent) mockFlightStatus(flightCode string) *FlightStatus {
	// Simple mock - some flights are on time, some delayed
	delayMinutes := 0
	status := "on-time"
	
	// Hash-based pseudo-random delay
	if len(flightCode) > 0 && flightCode[0] >= 'J' {
		delayMinutes = 30
		status = "delayed"
	}

	departureTime := time.Now().Add(2 * time.Hour).Format("2006-01-02 15:04")
	arrivalTime := time.Now().Add(6 * time.Hour).Format("2006-01-02 15:04")

	return &FlightStatus{
		FlightCode:    flightCode,
		Status:        status,
		DepartureTime: departureTime,
		ArrivalTime:   arrivalTime,
		Gate:          "A12",
		DelayMinutes:  delayMinutes,
	}
}

// generateDelayNotification generates a polite delay notification
func (a *FlightAgent) generateDelayNotification(ctx context.Context, status *FlightStatus) string {
	if a.client == nil {
		return fmt.Sprintf("Your flight %s is delayed by %d minutes. New departure time: %s. Please check the gate information.", 
			status.FlightCode, status.DelayMinutes, status.DepartureTime)
	}

	prompt := fmt.Sprintf(`You are FlightAgent. Flight %s is delayed by %d minutes.
Generate a polite, brief notification message for the passenger. Include:
- Flight code
- Delay duration
- Apology
- Brief advice`, status.FlightCode, status.DelayMinutes)

	resp, err := a.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: "gpt-4o-mini",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a professional airline notification system. Generate polite, brief delay notifications.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.7,
			MaxTokens:   100,
		},
	)

	if err != nil {
		log.Printf("FlightAgent: Failed to generate notification: %v", err)
		return fmt.Sprintf("Flight %s is delayed by %d minutes.", status.FlightCode, status.DelayMinutes)
	}

	if len(resp.Choices) > 0 {
		return resp.Choices[0].Message.Content
	}

	return fmt.Sprintf("Flight %s is delayed by %d minutes.", status.FlightCode, status.DelayMinutes)
}

// GetCheapestFlight searches for the cheapest flight between two cities (Legacy function)
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
