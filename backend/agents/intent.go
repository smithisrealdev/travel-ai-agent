package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
)

// IntentResult represents the detected intent and extracted entities
type IntentResult struct {
	Intent   string                 `json:"intent"`
	Entities map[string]interface{} `json:"entities"`
}

// IntentAgent handles intent detection and entity extraction
type IntentAgent struct {
	client *openai.Client
}

// NewIntentAgent creates a new intent detection agent
func NewIntentAgent(apiKey string) *IntentAgent {
	if apiKey == "" {
		return &IntentAgent{client: nil}
	}
	return &IntentAgent{
		client: openai.NewClient(apiKey),
	}
}

// Detect analyzes user input and classifies intent with entity extraction
func (a *IntentAgent) Detect(ctx context.Context, userInput string) (*IntentResult, error) {
	if a.client == nil {
		log.Println("IntentAgent: OpenAI client not initialized, using fallback")
		return a.fallbackDetect(userInput), nil
	}

	// Get current time in UTC for context
	currentTime := time.Now().UTC().Format("2006-01-02 15:04:05")

	// Construct the prompt
	prompt := fmt.Sprintf(`Current Date and Time (UTC - YYYY-MM-DD HH:MM:SS formatted): %s
Current User's Login: smithisrealdev

You are an intent detection model for an AI travel assistant.
Classify the user message into one of the following:
[plan_trip, flight_check, weather_check, hotel_search, local_recommendation, budget_inquiry, plan_update, general_chat]

Message: "%s"

Return ONLY valid JSON:
{
  "intent": "one_of_the_intents_above",
  "entities": {
    "destination": "city or country",
    "duration": number_of_days,
    "budget": amount_in_thb,
    "date_from": "YYYY-MM-DD",
    "date_to": "YYYY-MM-DD",
    "travelers": number,
    "interests": ["interest1"],
    "location": {"lat": 0.0, "lng": 0.0},
    "flight_code": "flight number"
  }
}`, currentTime, userInput)

	// Make API request
	resp, err := a.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: os.Getenv("OPENAI_MODEL"),
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are an intent detection assistant. Classify user intents and extract entities. Return ONLY valid JSON.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.3,
			MaxTokens:   300,
		},
	)

	if err != nil {
		log.Printf("IntentAgent: OpenAI API error: %v, using fallback", err)
		return a.fallbackDetect(userInput), nil
	}

	// Parse the response
	if len(resp.Choices) == 0 {
		log.Println("IntentAgent: No response from OpenAI, using fallback")
		return a.fallbackDetect(userInput), nil
	}

	content := resp.Choices[0].Message.Content

	// Try to parse JSON response
	var result IntentResult
	err = json.Unmarshal([]byte(content), &result)
	if err != nil {
		log.Printf("IntentAgent: Failed to parse OpenAI response: %v, using fallback", err)
		return a.fallbackDetect(userInput), nil
	}

	log.Printf("IntentAgent: Detected intent=%s, entities=%v", result.Intent, result.Entities)
	return &result, nil
}

// fallbackDetect provides simple rule-based intent detection as fallback
func (a *IntentAgent) fallbackDetect(userInput string) *IntentResult {
	lowerInput := strings.ToLower(userInput)

	// Simple keyword-based detection
	intent := "general_chat"
	entities := make(map[string]interface{})

	if strings.Contains(lowerInput, "flight") && (strings.Contains(lowerInput, "status") || strings.Contains(lowerInput, "check") || strings.Contains(lowerInput, "on time") || strings.Contains(lowerInput, "is flight")) {
		intent = "flight_check"
		// Try to extract flight code (simple pattern matching)
		words := strings.Fields(userInput)
		for i, word := range words {
			if strings.ToLower(word) == "flight" && i+1 < len(words) {
				possibleCode := strings.TrimSuffix(words[i+1], "?")
				possibleCode = strings.TrimSuffix(possibleCode, ".")
				if len(possibleCode) >= 3 && len(possibleCode) <= 8 {
					entities["flight_code"] = possibleCode
				}
				break
			}
		}
	} else if strings.Contains(lowerInput, "weather") || strings.Contains(lowerInput, "forecast") || strings.Contains(lowerInput, "rain") || strings.Contains(lowerInput, "ฝน") {
		intent = "weather_check"
	} else if strings.Contains(lowerInput, "hotel") || strings.Contains(lowerInput, "accommodation") || strings.Contains(lowerInput, "โรงแรม") {
		intent = "hotel_search"
	} else if strings.Contains(lowerInput, "restaurant") || strings.Contains(lowerInput, "cafe") || strings.Contains(lowerInput, "nearby") || strings.Contains(lowerInput, "ร้านอาหาร") || strings.Contains(lowerInput, "ใกล้") || strings.Contains(lowerInput, "ราเมน") {
		intent = "local_recommendation"
	} else if strings.Contains(lowerInput, "plan") || strings.Contains(lowerInput, "trip") || strings.Contains(lowerInput, "travel") || strings.Contains(lowerInput, "visit") || strings.Contains(lowerInput, "เที่ยว") || strings.Contains(lowerInput, "ไป") {
		intent = "plan_trip"
		// Extract basic entities for plan_trip
		entities["duration"] = 7
		entities["budget"] = 50000
	} else if strings.Contains(lowerInput, "budget") || strings.Contains(lowerInput, "cost") || strings.Contains(lowerInput, "price") {
		// Only set budget_inquiry if not already a trip plan
		if !strings.Contains(lowerInput, "trip") && !strings.Contains(lowerInput, "travel") && !strings.Contains(lowerInput, "เที่ยว") {
			intent = "budget_inquiry"
		} else {
			intent = "plan_trip"
			entities["duration"] = 7
			entities["budget"] = 50000
		}
	} else if strings.Contains(lowerInput, "update") || strings.Contains(lowerInput, "change") || strings.Contains(lowerInput, "modify") || strings.Contains(lowerInput, "เปลี่ยน") {
		intent = "plan_update"
	}

	return &IntentResult{
		Intent:   intent,
		Entities: entities,
	}
}

// Legacy function for backward compatibility
func AnalyzeIntent(message string) (destination string, budgetTHB int, durationDays int) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	agent := NewIntentAgent(apiKey)
	
	ctx := context.Background()
	result, err := agent.Detect(ctx, message)
	if err != nil || result == nil {
		return "Unknown", 50000, 7
	}

	// Extract destination
	if dest, ok := result.Entities["destination"].(string); ok && dest != "" {
		destination = dest
	} else {
		destination = "Unknown"
	}

	// Extract budget
	if budget, ok := result.Entities["budget"].(float64); ok && budget > 0 {
		budgetTHB = int(budget)
	} else {
		budgetTHB = 50000
	}

	// Extract duration
	if duration, ok := result.Entities["duration"].(float64); ok && duration > 0 {
		durationDays = int(duration)
	} else {
		durationDays = 7
	}

	log.Printf("Intent analyzed (legacy): destination=%s, budget=%d THB, duration=%d days",
		destination, budgetTHB, durationDays)

	return destination, budgetTHB, durationDays
}
