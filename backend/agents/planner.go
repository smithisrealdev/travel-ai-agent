package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/sashabaranov/go-openai"
)

// TripPlan represents a complete travel itinerary
type TripPlan struct {
	Destination string         `json:"destination"`
	Duration    int            `json:"duration"`
	TotalBudget float64        `json:"total_budget"`
	Itinerary   []ItineraryDay `json:"itinerary"`
	Summary     string         `json:"summary"`
}

// ItineraryDay represents activities and budget for a single day
type ItineraryDay struct {
	Day        int      `json:"day"`
	Activities []string `json:"activities"`
	Budget     float64  `json:"budget"`
}

// PlannerAgent creates and updates travel itineraries
type PlannerAgent struct {
	client *openai.Client
}

// NewPlannerAgent creates a new planner agent
func NewPlannerAgent(apiKey string) *PlannerAgent {
	if apiKey == "" {
		return &PlannerAgent{client: nil}
	}
	return &PlannerAgent{
		client: openai.NewClient(apiKey),
	}
}

// CreatePlan generates a new travel itinerary
func (a *PlannerAgent) CreatePlan(ctx context.Context, destination string, duration int, budget float64) (*TripPlan, error) {
	if a.client == nil {
		log.Println("PlannerAgent: OpenAI client not initialized, using default plan")
		return a.fallbackPlan(destination, duration, budget), nil
	}

	// Construct the prompt
	prompt := fmt.Sprintf(`You are PlannerAgent, an expert travel planner.

Create a detailed itinerary for:
- Destination: %s
- Duration: %d days
- Budget: %.0f THB

Return ONLY valid JSON:
{
  "destination": "...",
  "duration": 0,
  "total_budget": 0,
  "itinerary": [
    {
      "day": 1,
      "activities": ["Activity 1", "Activity 2"],
      "budget": 15000
    }
  ],
  "summary": "Brief overview in markdown"
}`, destination, duration, budget)

	// Make API request
	resp, err := a.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: "gpt-4o-mini",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are an expert travel planner. Create detailed itineraries and return ONLY valid JSON.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.7,
			MaxTokens:   1500,
		},
	)

	if err != nil {
		log.Printf("PlannerAgent: OpenAI API error: %v, using fallback", err)
		return a.fallbackPlan(destination, duration, budget), nil
	}

	// Parse the response
	if len(resp.Choices) == 0 {
		log.Println("PlannerAgent: No response from OpenAI, using fallback")
		return a.fallbackPlan(destination, duration, budget), nil
	}

	content := resp.Choices[0].Message.Content

	// Try to parse JSON response
	var plan TripPlan
	err = json.Unmarshal([]byte(content), &plan)
	if err != nil {
		log.Printf("PlannerAgent: Failed to parse OpenAI response: %v, using fallback", err)
		return a.fallbackPlan(destination, duration, budget), nil
	}

	log.Printf("PlannerAgent: Created plan for %s, %d days", destination, duration)
	return &plan, nil
}

// UpdatePlan modifies an existing itinerary based on new conditions
func (a *PlannerAgent) UpdatePlan(ctx context.Context, currentPlan *TripPlan, condition string) (*TripPlan, error) {
	if a.client == nil {
		log.Println("PlannerAgent: OpenAI client not initialized, returning current plan")
		return currentPlan, nil
	}

	// Convert current plan to JSON
	planJSON, err := json.Marshal(currentPlan)
	if err != nil {
		log.Printf("PlannerAgent: Failed to marshal current plan: %v", err)
		return currentPlan, err
	}

	// Construct the prompt
	prompt := fmt.Sprintf(`You are PlannerAgent. The user currently has this plan:
%s

Update it according to this condition: %s

Return the revised plan in JSON with activities for each day.`, string(planJSON), condition)

	// Make API request
	resp, err := a.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: "gpt-4o-mini",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are an expert travel planner. Update itineraries based on new conditions and return ONLY valid JSON.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.7,
			MaxTokens:   1500,
		},
	)

	if err != nil {
		log.Printf("PlannerAgent: OpenAI API error: %v, returning current plan", err)
		return currentPlan, err
	}

	// Parse the response
	if len(resp.Choices) == 0 {
		log.Println("PlannerAgent: No response from OpenAI, returning current plan")
		return currentPlan, nil
	}

	content := resp.Choices[0].Message.Content

	// Try to parse JSON response
	var updatedPlan TripPlan
	err = json.Unmarshal([]byte(content), &updatedPlan)
	if err != nil {
		log.Printf("PlannerAgent: Failed to parse OpenAI response: %v, returning current plan", err)
		return currentPlan, err
	}

	log.Printf("PlannerAgent: Updated plan for %s", currentPlan.Destination)
	return &updatedPlan, nil
}

// fallbackPlan generates a simple default itinerary
func (a *PlannerAgent) fallbackPlan(destination string, duration int, budget float64) *TripPlan {
	dailyBudget := budget / float64(duration)
	itinerary := make([]ItineraryDay, duration)

	for i := 0; i < duration; i++ {
		day := i + 1
		activities := []string{
			fmt.Sprintf("Explore %s attractions", destination),
			"Try local cuisine",
			"Visit popular landmarks",
		}
		
		itinerary[i] = ItineraryDay{
			Day:        day,
			Activities: activities,
			Budget:     dailyBudget,
		}
	}

	summary := fmt.Sprintf("## %d-Day Trip to %s\n\nExplore the best of %s with daily activities and local experiences. Budget: %.0f THB", 
		duration, destination, destination, budget)

	return &TripPlan{
		Destination: destination,
		Duration:    duration,
		TotalBudget: budget,
		Itinerary:   itinerary,
		Summary:     summary,
	}
}
