package services

import (
	"context"
	"fmt"
	"log"

	"github.com/sashabaranov/go-openai"
	"github.com/smithisrealdev/travel-ai-agent/backend/internal/config"
	"github.com/smithisrealdev/travel-ai-agent/backend/internal/models"
)

// OpenAIService handles OpenAI API interactions
type OpenAIService struct {
	client *openai.Client
	model  string
}

// NewOpenAIService creates a new OpenAI service instance
func NewOpenAIService(cfg *config.Config) *OpenAIService {
	if cfg.OpenAI.APIKey == "" {
		log.Println("Warning: OpenAI API key not configured")
		return nil
	}

	client := openai.NewClient(cfg.OpenAI.APIKey)
	return &OpenAIService{
		client: client,
		model:  cfg.OpenAI.Model,
	}
}

// GenerateTravelRecommendations generates AI-powered travel recommendations
func (s *OpenAIService) GenerateTravelRecommendations(ctx context.Context, req *models.TravelSearchRequest) (string, error) {
	if s == nil || s.client == nil {
		return "", fmt.Errorf("OpenAI service not initialized")
	}

	// Build the prompt for travel recommendations
	prompt := s.buildTravelPrompt(req)

	// Create the completion request
	resp, err := s.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: s.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a helpful travel assistant that provides detailed, personalized travel recommendations. Provide practical advice about destinations, activities, accommodations, and local experiences.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.7,
			MaxTokens:   1000,
		},
	)

	if err != nil {
		return "", fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return resp.Choices[0].Message.Content, nil
}

// buildTravelPrompt constructs a detailed prompt for travel recommendations
func (s *OpenAIService) buildTravelPrompt(req *models.TravelSearchRequest) string {
	prompt := fmt.Sprintf("I'm planning a trip to %s.", req.Destination)

	if req.StartDate != "" && req.EndDate != "" {
		prompt += fmt.Sprintf(" I'll be traveling from %s to %s.", req.StartDate, req.EndDate)
	}

	if req.Budget > 0 {
		prompt += fmt.Sprintf(" My budget is approximately $%.2f.", req.Budget)
	}

	if len(req.Preferences) > 0 {
		prompt += " My preferences include:"
		for key, value := range req.Preferences {
			prompt += fmt.Sprintf(" %s: %v,", key, value)
		}
	}

	prompt += "\n\nPlease provide personalized travel recommendations including:\n"
	prompt += "1. Top attractions and activities\n"
	prompt += "2. Recommended accommodations in different price ranges\n"
	prompt += "3. Local cuisine and dining suggestions\n"
	prompt += "4. Transportation tips\n"
	prompt += "5. Best time to visit specific attractions\n"
	prompt += "6. Cultural tips and local customs\n"
	prompt += "7. Estimated daily budget breakdown\n"

	return prompt
}

// GenerateItinerary generates a day-by-day travel itinerary
func (s *OpenAIService) GenerateItinerary(ctx context.Context, destination string, days int, interests []string) (string, error) {
	if s == nil || s.client == nil {
		return "", fmt.Errorf("OpenAI service not initialized")
	}

	prompt := fmt.Sprintf(
		"Create a detailed %d-day itinerary for %s. Interests: %v. Include specific activities, timing, and locations for each day.",
		days,
		destination,
		interests,
	)

	resp, err := s.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: s.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are an expert travel planner creating detailed day-by-day itineraries.",
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
		return "", fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return resp.Choices[0].Message.Content, nil
}

// HealthCheck verifies the OpenAI service is configured
func (s *OpenAIService) HealthCheck() error {
	if s == nil || s.client == nil {
		return fmt.Errorf("OpenAI service not initialized")
	}
	return nil
}
