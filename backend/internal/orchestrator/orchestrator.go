package orchestrator

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/smithisrealdev/travel-ai-agent/backend/agents"
)

// Orchestrator coordinates multiple agents based on user intent
type Orchestrator struct {
	intentAgent  *agents.IntentAgent
	plannerAgent *agents.PlannerAgent
	weatherAgent *agents.WeatherAgent
	flightAgent  *agents.FlightAgent
	localAgent   *agents.LocalAgent
	hotelAgent   *agents.HotelAgent
	socialService interface {
		GetTopRatedPlaces(keyword, location string, limit int) ([]SocialPlace, error)
	}
}

// SocialPlace represents a socially popular place (imported from models)
type SocialPlace struct {
	PlaceID     string
	Name        string
	Address     string
	Rating      float64
	ReviewCount int
	Types       []string
}

// New creates a new orchestrator with all agents
func New(openaiKey, weatherKey, flightKey, hotelKey string) *Orchestrator {
	return &Orchestrator{
		intentAgent:  agents.NewIntentAgent(openaiKey),
		plannerAgent: agents.NewPlannerAgent(openaiKey),
		weatherAgent: agents.NewWeatherAgent(openaiKey, weatherKey),
		flightAgent:  agents.NewFlightAgent(openaiKey, flightKey),
		localAgent:   agents.NewLocalAgent(openaiKey),
		hotelAgent:   agents.NewHotelAgent(openaiKey, hotelKey),
		socialService: nil, // Will be set via SetSocialService
	}
}

// SetSocialService sets the social service for the orchestrator
func (o *Orchestrator) SetSocialService(service interface {
	GetTopRatedPlaces(keyword, location string, limit int) ([]SocialPlace, error)
}) {
	o.socialService = service
}

// ProcessMessage is the main entry point for handling user messages
func (o *Orchestrator) ProcessMessage(ctx context.Context, userInput string) (string, error) {
	log.Printf("🚀 Orchestrator: Processing message: %s", userInput)

	// Step 1: Detect intent
	intentResult, err := o.intentAgent.Detect(ctx, userInput)
	if err != nil {
		log.Printf("Orchestrator: Intent detection failed: %v", err)
		return "", err
	}

	log.Printf("Orchestrator: Detected intent=%s", intentResult.Intent)

	// Step 2: Route to appropriate handler
	var response string
	switch intentResult.Intent {
	case "plan_trip":
		response, err = o.handlePlanTrip(ctx, intentResult)
	case "weather_check":
		response, err = o.handleWeatherCheck(ctx, intentResult)
	case "flight_check":
		response, err = o.handleFlightCheck(ctx, intentResult)
	case "hotel_search":
		response, err = o.handleHotelSearch(ctx, intentResult)
	case "local_recommendation":
		response, err = o.handleLocalRecommendation(ctx, intentResult)
	case "budget_inquiry":
		response, err = o.handleBudgetInquiry(ctx, intentResult)
	case "plan_update":
		response = "Plan update functionality coming soon!"
	case "general_chat":
		response = "Hello! I'm your AI travel assistant. I can help you plan trips, check weather, find flights, search hotels, and get local recommendations. What would you like to do?"
	default:
		response = "I'm not sure how to help with that. Try asking about planning a trip, checking weather, or finding hotels!"
	}

	if err != nil {
		log.Printf("Orchestrator: Handler error: %v", err)
		return "", err
	}

	log.Printf("✅ Orchestrator: Response generated successfully")
	return response, nil
}

// handlePlanTrip creates a complete travel plan
func (o *Orchestrator) handlePlanTrip(ctx context.Context, intent *agents.IntentResult) (string, error) {
	// Extract entities
	destination := o.getStringEntity(intent.Entities, "destination", "Unknown")
	duration := o.getIntEntity(intent.Entities, "duration", 7)
	budget := o.getFloatEntity(intent.Entities, "budget", 50000)

	log.Printf("Creating plan: destination=%s, duration=%d days, budget=%.0f THB", 
		destination, duration, budget)

	// Create itinerary
	plan, err := o.plannerAgent.CreatePlan(ctx, destination, duration, budget)
	if err != nil {
		return "", err
	}

	// Get weather forecast
	weather, err := o.weatherAgent.GetForecast(ctx, destination)
	if err != nil {
		log.Printf("Weather check failed: %v", err)
	}

	// Search for hotels
	hotels, err := o.hotelAgent.SearchHotels(ctx, destination, budget/float64(duration))
	if err != nil {
		log.Printf("Hotel search failed: %v", err)
	}

	// Build response
	response := fmt.Sprintf("# %d-Day Trip to %s\n\n", duration, destination)
	response += fmt.Sprintf("**Budget:** %.0f THB\n\n", budget)
	
	response += "## Itinerary\n"
	for _, day := range plan.Itinerary {
		response += fmt.Sprintf("\n**Day %d:**\n", day.Day)
		for _, activity := range day.Activities {
			response += fmt.Sprintf("- %s\n", activity)
		}
		response += fmt.Sprintf("*Daily Budget: %.0f THB*\n", day.Budget)
	}

	if weather != nil {
		response += fmt.Sprintf("\n## Weather Forecast\n")
		response += fmt.Sprintf("Current: %.0f°C, %s\n", weather.Temperature, weather.Condition)
		if weather.RainProb > 60 {
			response += fmt.Sprintf("\n⚠️ %s\n", weather.Suggestion)
		}
	}

	if len(hotels) > 0 {
		response += fmt.Sprintf("\n## Recommended Hotels\n")
		for i, hotel := range hotels {
			if i < 3 {
				response += fmt.Sprintf("- **%s** - %.0f THB/night (Rating: %.1f★)\n", 
					hotel.Name, hotel.PricePerNight, hotel.Rating)
			}
		}
	}

	// Get socially popular spots
	if o.socialService != nil {
		socialPlaces, err := o.socialService.GetTopRatedPlaces("tourist attractions", destination, 5)
		if err == nil && len(socialPlaces) > 0 {
			response += fmt.Sprintf("\n## Socially Popular Spots\n")
			response += "*Top-rated places based on reviews*\n\n"
			for i, place := range socialPlaces {
				if i < 5 {
					response += fmt.Sprintf("- **%s** (%.1f★, %d reviews)\n", 
						place.Name, place.Rating, place.ReviewCount)
				}
			}
		}
	}

	response += fmt.Sprintf("\n%s", plan.Summary)

	return response, nil
}

// handleWeatherCheck gets weather forecast for a city
func (o *Orchestrator) handleWeatherCheck(ctx context.Context, intent *agents.IntentResult) (string, error) {
	city := o.getStringEntity(intent.Entities, "destination", "Bangkok")

	log.Printf("Checking weather for: %s", city)

	forecast, err := o.weatherAgent.GetForecast(ctx, city)
	if err != nil {
		return "", err
	}

	response := fmt.Sprintf("# Weather Forecast for %s\n\n", city)
	response += fmt.Sprintf("**Current:** %.0f°C, %s\n\n", forecast.Temperature, forecast.Condition)
	response += "## 3-Day Forecast\n"
	
	for _, day := range forecast.Forecast {
		response += fmt.Sprintf("- %s: %.0f°C, %s (Rain: %.0f%%)\n", 
			day.Date, day.Temperature, day.Condition, day.RainProb)
	}

	if forecast.RainProb > 60 {
		response += fmt.Sprintf("\n⚠️ **Rain Alert:** %s\n", forecast.Suggestion)
	}

	return response, nil
}

// handleFlightCheck checks flight status
func (o *Orchestrator) handleFlightCheck(ctx context.Context, intent *agents.IntentResult) (string, error) {
	flightCode := o.getStringEntity(intent.Entities, "flight_code", "")
	
	if flightCode == "" {
		return "Please provide a flight code (e.g., 'Is flight JL708 on time?')", nil
	}

	log.Printf("Checking flight: %s", flightCode)

	status, err := o.flightAgent.CheckFlight(ctx, flightCode)
	if err != nil {
		return "", err
	}

	response := fmt.Sprintf("# Flight %s Status\n\n", status.FlightCode)
	response += fmt.Sprintf("**Status:** %s\n", strings.Title(status.Status))
	response += fmt.Sprintf("**Departure:** %s\n", status.DepartureTime)
	response += fmt.Sprintf("**Arrival:** %s\n", status.ArrivalTime)
	
	if status.Gate != "" {
		response += fmt.Sprintf("**Gate:** %s\n", status.Gate)
	}

	if status.DelayMinutes > 0 {
		response += fmt.Sprintf("\n⚠️ **Delayed by %d minutes**\n\n", status.DelayMinutes)
	}

	response += fmt.Sprintf("\n%s", status.Notification)

	return response, nil
}

// handleHotelSearch searches for hotels
func (o *Orchestrator) handleHotelSearch(ctx context.Context, intent *agents.IntentResult) (string, error) {
	destination := o.getStringEntity(intent.Entities, "destination", "Bangkok")
	budget := o.getFloatEntity(intent.Entities, "budget", 3000)

	log.Printf("Searching hotels in %s, budget: %.0f THB/night", destination, budget)

	hotels, err := o.hotelAgent.SearchHotels(ctx, destination, budget)
	if err != nil {
		return "", err
	}

	response := fmt.Sprintf("# Hotels in %s\n\n", destination)
	response += fmt.Sprintf("Budget: Up to %.0f THB per night\n\n", budget)

	for i, hotel := range hotels {
		response += fmt.Sprintf("%d. **%s**\n", i+1, hotel.Name)
		response += fmt.Sprintf("   - Price: %.0f THB/night\n", hotel.PricePerNight)
		response += fmt.Sprintf("   - Rating: %.1f★\n", hotel.Rating)
		response += fmt.Sprintf("   - Distance: %.1f km from center\n", hotel.Distance)
		response += fmt.Sprintf("   - Address: %s\n\n", hotel.Address)
	}

	return response, nil
}

// handleLocalRecommendation finds nearby places
func (o *Orchestrator) handleLocalRecommendation(ctx context.Context, intent *agents.IntentResult) (string, error) {
	interest := o.getStringEntity(intent.Entities, "interests", "restaurant")
	destination := o.getStringEntity(intent.Entities, "destination", "")
	
	// Get location from entities
	lat := 13.7563 // Default Bangkok
	lng := 100.5018
	
	if loc, ok := intent.Entities["location"].(map[string]interface{}); ok {
		if latVal, ok := loc["lat"].(float64); ok {
			lat = latVal
		}
		if lngVal, ok := loc["lng"].(float64); ok {
			lng = lngVal
		}
	}

	log.Printf("Finding %s near (%.4f, %.4f)", interest, lat, lng)

	places, err := o.localAgent.GetRecommendations(ctx, lat, lng, interest)
	if err != nil {
		return "", err
	}

	response := fmt.Sprintf("# Nearby %s Recommendations\n\n", strings.Title(interest))

	// Add AI-generated recommendations
	for i, place := range places {
		response += fmt.Sprintf("%d. **%s**\n", i+1, place.Name)
		response += fmt.Sprintf("   - Type: %s\n", place.Type)
		response += fmt.Sprintf("   - Rating: %.1f★\n", place.Rating)
		response += fmt.Sprintf("   - Distance: %.1f km away\n", place.DistanceKm)
		response += fmt.Sprintf("   - Address: %s\n\n", place.Address)
	}

	// Add socially popular spots if available and destination is provided
	if o.socialService != nil && destination != "" {
		socialPlaces, err := o.socialService.GetTopRatedPlaces(interest, destination, 3)
		if err == nil && len(socialPlaces) > 0 {
			response += fmt.Sprintf("\n## Socially Popular %s in %s\n", strings.Title(interest), destination)
			response += "*Top-rated by the community*\n\n"
			for i, place := range socialPlaces {
				response += fmt.Sprintf("%d. **%s** (%.1f★, %d reviews)\n", 
					len(places)+i+1, place.Name, place.Rating, place.ReviewCount)
				response += fmt.Sprintf("   - Address: %s\n\n", place.Address)
			}
		}
	}

	return response, nil
}

// handleBudgetInquiry provides budget breakdown
func (o *Orchestrator) handleBudgetInquiry(ctx context.Context, intent *agents.IntentResult) (string, error) {
	budget := o.getFloatEntity(intent.Entities, "budget", 50000)

	budgetPlan := agents.EstimateBudget(int(budget))

	response := fmt.Sprintf("# Budget Breakdown for %.0f THB\n\n", budget)
	response += fmt.Sprintf("- **Flights:** %d THB (45%%)\n", budgetPlan.Flight)
	response += fmt.Sprintf("- **Hotels:** %d THB (25%%)\n", budgetPlan.Hotel)
	response += fmt.Sprintf("- **Food:** %d THB (15%%)\n", budgetPlan.Food)
	response += fmt.Sprintf("- **Transport:** %d THB (10%%)\n", budgetPlan.Transport)
	response += fmt.Sprintf("- **Miscellaneous:** %d THB (5%%)\n", budgetPlan.Misc)

	total := budgetPlan.Flight + budgetPlan.Hotel + budgetPlan.Food + budgetPlan.Transport + budgetPlan.Misc
	response += fmt.Sprintf("\n**Total:** %d THB\n", total)

	return response, nil
}

// Helper methods to extract entities safely
func (o *Orchestrator) getStringEntity(entities map[string]interface{}, key, defaultValue string) string {
	if val, ok := entities[key].(string); ok && val != "" {
		return val
	}
	return defaultValue
}

func (o *Orchestrator) getIntEntity(entities map[string]interface{}, key string, defaultValue int) int {
	if val, ok := entities[key].(float64); ok {
		return int(val)
	}
	if val, ok := entities[key].(int); ok {
		return val
	}
	return defaultValue
}

func (o *Orchestrator) getFloatEntity(entities map[string]interface{}, key string, defaultValue float64) float64 {
	if val, ok := entities[key].(float64); ok {
		return val
	}
	if val, ok := entities[key].(int); ok {
		return float64(val)
	}
	return defaultValue
}
