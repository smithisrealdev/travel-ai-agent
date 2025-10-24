package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/sashabaranov/go-openai"
)

type IntentResult struct {
	Destination  string `json:"destination"`
	BudgetTHB    int    `json:"budget_thb"`
	DurationDays int    `json:"duration_days"`
}

func AnalyzeIntent(message string) (destination string, budgetTHB int, durationDays int) {
	// Get API key from environment
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Println("OPENAI_API_KEY not set, using defaults")
		return "Unknown", 50000, 7
	}

	// Create OpenAI client
	client := openai.NewClient(apiKey)
	ctx := context.Background()

	// Construct the prompt
	prompt := fmt.Sprintf(
		"Extract travel destination, budget in THB, and duration in days from: %s\n\n"+
			"Return the information in this exact JSON format:\n"+
			"{\n"+
			"  \"destination\": \"city or country name\",\n"+
			"  \"budget_thb\": number,\n"+
			"  \"duration_days\": number\n"+
			"}\n\n"+
			"If any information is missing, use reasonable defaults:\n"+
			"- destination: \"Unknown\"\n"+
			"- budget_thb: 50000\n"+
			"- duration_days: 7",
		message,
	)

	// Make API request
	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: "gpt-4o-mini",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a travel intent extraction assistant. Extract travel information and return valid JSON only.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.3,
			MaxTokens:   150,
		},
	)

	if err != nil {
		log.Printf("OpenAI API error: %v, using defaults", err)
		return "Unknown", 50000, 7
	}

	// Parse the response
	if len(resp.Choices) == 0 {
		log.Println("No response from OpenAI, using defaults")
		return "Unknown", 50000, 7
	}

	content := resp.Choices[0].Message.Content

	// Try to parse JSON response
	var result IntentResult
	err = json.Unmarshal([]byte(content), &result)
	if err != nil {
		log.Printf("Failed to parse OpenAI response: %v, using defaults", err)
		return "Unknown", 50000, 7
	}

	// Validate and return
	if result.Destination == "" {
		result.Destination = "Unknown"
	}
	if result.BudgetTHB <= 0 {
		result.BudgetTHB = 50000
	}
	if result.DurationDays <= 0 {
		result.DurationDays = 7
	}

	log.Printf("Intent analyzed: destination=%s, budget=%d THB, duration=%d days",
		result.Destination, result.BudgetTHB, result.DurationDays)

	return result.Destination, result.BudgetTHB, result.DurationDays
}
