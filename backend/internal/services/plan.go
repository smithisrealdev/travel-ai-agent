package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/sashabaranov/go-openai"
	"github.com/smithisrealdev/travel-ai-agent/backend/internal/config"
	"github.com/smithisrealdev/travel-ai-agent/backend/internal/models"
)

// PlanService handles travel plan generation using OpenAI
type PlanService struct {
	client *openai.Client
	model  string
}

// NewPlanService creates a new plan service instance
func NewPlanService(cfg *config.Config) *PlanService {
	if cfg.OpenAI.APIKey == "" {
		log.Println("Warning: OpenAI API key not configured")
		return nil
	}

	client := openai.NewClient(cfg.OpenAI.APIKey)
	return &PlanService{
		client: client,
		model:  cfg.OpenAI.Model,
	}
}

// GenerateTravelPlan generates a comprehensive travel plan based on user message
func (s *PlanService) GenerateTravelPlan(ctx context.Context, message string) (*models.PlanResponse, error) {
	if s == nil || s.client == nil {
		return nil, fmt.Errorf("plan service not initialized")
	}

	// Build the system prompt for structured JSON output
	systemPrompt := `You are a travel planning assistant. Analyze the user's travel request and provide a detailed travel plan in JSON format.

Extract the following information:
- destination: The travel destination mentioned by the user
- budget: The budget in Thai Baht (THB). If another currency is mentioned, convert it to THB
- duration_days: The number of days for the trip
- itinerary: An array of daily activities with "day" (integer) and "activity" (string) fields
- weather: Object with "avg_temp" (float, in Celsius) and "condition" (string) for the destination
- flight_price: Estimated round-trip flight price in THB
- hotel_price: Estimated hotel price per night in THB

Return ONLY valid JSON matching this exact structure. Do not include any markdown formatting or additional text.

Example format:
{
  "destination": "Canada",
  "budget": 100000,
  "duration_days": 7,
  "itinerary": [
    {"day": 1, "activity": "Arrive Vancouver"},
    {"day": 2, "activity": "Stanley Park + Aquarium"}
  ],
  "weather": {"avg_temp": 15, "condition": "Sunny"},
  "flight_price": 38000,
  "hotel_price": 2500
}`

	// Create the completion request
	resp, err := s.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: s.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: message,
				},
			},
			Temperature: 0.7,
			MaxTokens:   2000,
		},
	)

	if err != nil {
		log.Printf("OpenAI API error: %v", err)
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	// Parse the JSON response
	content := resp.Choices[0].Message.Content
	log.Printf("OpenAI response: %s", content)

	var planResponse models.PlanResponse
	if err := json.Unmarshal([]byte(content), &planResponse); err != nil {
		log.Printf("Failed to parse OpenAI response: %v", err)
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	return &planResponse, nil
}

// HealthCheck verifies the plan service is configured
func (s *PlanService) HealthCheck() error {
	if s == nil || s.client == nil {
		return fmt.Errorf("plan service not initialized")
	}
	return nil
}
