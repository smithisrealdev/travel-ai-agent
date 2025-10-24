package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/sashabaranov/go-openai"
)

// PlaceRecommendation represents a nearby place recommendation
type PlaceRecommendation struct {
	Name       string  `json:"name"`
	Type       string  `json:"type"`
	Rating     float64 `json:"rating"`
	DistanceKm float64 `json:"distance_km"`
	Address    string  `json:"address"`
}

// LocalAgent finds nearby places based on interests
type LocalAgent struct {
	client *openai.Client
}

// NewLocalAgent creates a new local recommendations agent
func NewLocalAgent(apiKey string) *LocalAgent {
	if apiKey == "" {
		return &LocalAgent{client: nil}
	}
	return &LocalAgent{
		client: openai.NewClient(apiKey),
	}
}

// GetRecommendations finds nearby places based on location and interest
func (a *LocalAgent) GetRecommendations(ctx context.Context, lat, lng float64, interest string) ([]PlaceRecommendation, error) {
	if a.client == nil {
		log.Println("LocalAgent: OpenAI client not initialized, using fallback recommendations")
		return a.fallbackRecommendations(interest), nil
	}

	// Construct the prompt
	prompt := fmt.Sprintf(`You are LocalAgent.
Given current location (lat: %.6f, lng: %.6f) and preference: %s,
recommend 3 options within 3 km.

Format (return ONLY valid JSON array):
[
  {"name": "...", "type": "cafe", "rating": 4.6, "distance_km": 1.2, "address": "..."},
  {"name": "...", "type": "restaurant", "rating": 4.5, "distance_km": 0.8, "address": "..."},
  {"name": "...", "type": "cafe", "rating": 4.7, "distance_km": 2.1, "address": "..."}
]`, lat, lng, interest)

	// Make API request
	resp, err := a.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: "gpt-4o-mini",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a local recommendations expert. Provide realistic place recommendations and return ONLY valid JSON array.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.7,
			MaxTokens:   500,
		},
	)

	if err != nil {
		log.Printf("LocalAgent: OpenAI API error: %v, using fallback", err)
		return a.fallbackRecommendations(interest), nil
	}

	// Parse the response
	if len(resp.Choices) == 0 {
		log.Println("LocalAgent: No response from OpenAI, using fallback")
		return a.fallbackRecommendations(interest), nil
	}

	content := resp.Choices[0].Message.Content

	// Try to parse JSON response
	var recommendations []PlaceRecommendation
	err = json.Unmarshal([]byte(content), &recommendations)
	if err != nil {
		log.Printf("LocalAgent: Failed to parse OpenAI response: %v, using fallback", err)
		return a.fallbackRecommendations(interest), nil
	}

	log.Printf("LocalAgent: Found %d recommendations for %s", len(recommendations), interest)
	return recommendations, nil
}

// fallbackRecommendations provides default recommendations
func (a *LocalAgent) fallbackRecommendations(interest string) []PlaceRecommendation {
	rand.Seed(time.Now().UnixNano())
	
	recommendations := []PlaceRecommendation{
		{
			Name:       fmt.Sprintf("Popular %s Spot #1", interest),
			Type:       interest,
			Rating:     4.5 + rand.Float64()*0.5,
			DistanceKm: 0.5 + rand.Float64()*2.0,
			Address:    "City Center Area",
		},
		{
			Name:       fmt.Sprintf("Local %s Place #2", interest),
			Type:       interest,
			Rating:     4.3 + rand.Float64()*0.5,
			DistanceKm: 0.8 + rand.Float64()*1.5,
			Address:    "Downtown District",
		},
		{
			Name:       fmt.Sprintf("Best %s #3", interest),
			Type:       interest,
			Rating:     4.6 + rand.Float64()*0.4,
			DistanceKm: 1.2 + rand.Float64()*1.8,
			Address:    "Tourist Area",
		},
	}

	return recommendations
}
