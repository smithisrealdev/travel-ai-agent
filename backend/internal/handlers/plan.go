package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/smithisrealdev/travel-ai-agent/backend/internal/models"
	"github.com/smithisrealdev/travel-ai-agent/backend/internal/orchestrator"
	"github.com/smithisrealdev/travel-ai-agent/backend/internal/services"
)

// PlanHandler handles travel plan-related HTTP requests
type PlanHandler struct {
	planService  *services.PlanService
	orchestrator *orchestrator.Orchestrator
}

// NewPlanHandler creates a new plan handler instance
func NewPlanHandler(planService *services.PlanService, orch *orchestrator.Orchestrator) *PlanHandler {
	return &PlanHandler{
		planService:  planService,
		orchestrator: orch,
	}
}

// CreateTravelPlan handles POST /api/plan requests
func (h *PlanHandler) CreateTravelPlan(c *fiber.Ctx) error {
	var req models.PlanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
			Code:    fiber.StatusBadRequest,
		})
	}

	// Validate required fields
	if req.Message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Validation error",
			Message: "message field is required",
			Code:    fiber.StatusBadRequest,
		})
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Use orchestrator if available, otherwise use plan service
	if h.orchestrator != nil {
		log.Printf("Using orchestrator to process message: %s", req.Message)
		rawResponse, err := h.orchestrator.ProcessMessage(ctx, req.Message)
		if err != nil {
			log.Printf("Orchestrator error: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
				Error:   "AI service error",
				Message: fmt.Sprintf("Failed to process request: %v", err),
				Code:    fiber.StatusInternalServerError,
			})
		}

		// Format response as Markdown if it's JSON
		formattedResponse := formatResponseAsMarkdown(rawResponse)

		return c.JSON(fiber.Map{
			"success":  true,
			"response": formattedResponse,
		})
	}

	// Fallback to plan service
	if h.planService == nil {
		log.Println("Plan service not initialized")
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "Service unavailable",
			Message: "OpenAI service is not configured",
			Code:    fiber.StatusInternalServerError,
		})
	}

	// Generate travel plan using OpenAI
	log.Printf("Generating travel plan for message: %s", req.Message)
	planResponse, err := h.planService.GenerateTravelPlan(ctx, req.Message)
	if err != nil {
		log.Printf("Failed to generate travel plan: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "AI service error",
			Message: fmt.Sprintf("Failed to generate travel plan: %v", err),
			Code:    fiber.StatusInternalServerError,
		})
	}

	// Return the plan response
	return c.JSON(planResponse)
}

// formatResponseAsMarkdown converts raw response to beautiful Markdown
func formatResponseAsMarkdown(raw string) string {
	// Try to parse as JSON
	var plan map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &plan); err == nil {
		// Check if it looks like a trip plan
		if _, hasDestination := plan["destination"]; hasDestination {
			return formatTripPlanMarkdown(plan)
		}
	}
	
	// If already markdown or plain text, return as-is
	return raw
}

// formatTripPlanMarkdown formats a trip plan into beautiful Markdown
func formatTripPlanMarkdown(plan map[string]interface{}) string {
	var md strings.Builder
	
	destination := getStringValue(plan, "destination", "Unknown")
	duration := getIntValue(plan, "duration", 0)
	budget := getFloatValue(plan, "total_budget", 0)
	
	// Header with emoji
	md.WriteString(fmt.Sprintf("# 🌍 %s Trip Plan\n\n", destination))
	
	// Duration and budget section
	if duration > 0 && budget > 0 {
		md.WriteString(fmt.Sprintf("**Duration:** %d days | **Budget:** %s THB\n\n",
			duration, formatNumber(int(budget))))
	}
	
	// Itinerary section
	if itinerary, ok := plan["itinerary"].([]interface{}); ok && len(itinerary) > 0 {
		md.WriteString("## 📅 Day-by-Day Itinerary\n\n")
		for _, day := range itinerary {
			if dayMap, ok := day.(map[string]interface{}); ok {
				dayNum := getIntValue(dayMap, "day", 0)
				dailyBudget := getFloatValue(dayMap, "budget", 0)
				
				md.WriteString(fmt.Sprintf("### Day %d", dayNum))
				if dailyBudget > 0 {
					md.WriteString(fmt.Sprintf(" (Budget: %s THB)", formatNumber(int(dailyBudget))))
				}
				md.WriteString("\n\n")
				
				if activities, ok := dayMap["activities"].([]interface{}); ok {
					for _, act := range activities {
						md.WriteString(fmt.Sprintf("- %s\n", act))
					}
					md.WriteString("\n")
				}
			}
		}
	}
	
	// Weather section
	if weather, ok := plan["weather"].(map[string]interface{}); ok {
		md.WriteString("## 🌤️ Weather Forecast\n\n")
		if temp := getFloatValue(weather, "temperature", 0); temp > 0 {
			md.WriteString(fmt.Sprintf("**Current:** %.0f°C, %s\n\n",
				temp, getStringValue(weather, "condition", "Unknown")))
		}
	}
	
	// Budget breakdown
	if budget > 0 {
		md.WriteString("## 💰 Budget Breakdown\n\n")
		md.WriteString(fmt.Sprintf("- **Total Budget:** %s THB\n", formatNumber(int(budget))))
		
		if duration > 0 {
			dailyBudget := budget / float64(duration)
			md.WriteString(fmt.Sprintf("- **Daily Budget:** %s THB\n", formatNumber(int(dailyBudget))))
		}
		md.WriteString("\n")
	}
	
	// Travel tips
	md.WriteString("## 💡 Travel Tips\n\n")
	md.WriteString("- 📅 Book early for best prices\n")
	md.WriteString("- 🏥 Consider travel insurance\n")
	md.WriteString("- 📱 Download offline maps\n")
	md.WriteString("- 💳 Notify your bank about travel plans\n\n")
	
	md.WriteString("---\n\n*Have a wonderful trip! 🎉*\n")
	
	return md.String()
}

// formatNumber adds thousand separators to numbers
func formatNumber(n int) string {
	s := fmt.Sprintf("%d", n)
	if n < 1000 {
		return s
	}
	
	var result strings.Builder
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result.WriteString(",")
		}
		result.WriteRune(c)
	}
	return result.String()
}

// getStringValue safely extracts string from map
func getStringValue(m map[string]interface{}, key, defaultVal string) string {
	if val, ok := m[key].(string); ok && val != "" {
		return val
	}
	return defaultVal
}

// getIntValue safely extracts int from map
func getIntValue(m map[string]interface{}, key string, defaultVal int) int {
	if val, ok := m[key].(float64); ok {
		return int(val)
	}
	if val, ok := m[key].(int); ok {
		return val
	}
	return defaultVal
}

// getFloatValue safely extracts float from map
func getFloatValue(m map[string]interface{}, key string, defaultVal float64) float64 {
	if val, ok := m[key].(float64); ok {
		return val
	}
	if val, ok := m[key].(int); ok {
		return float64(val)
	}
	return defaultVal
}

