package handlers

import (
	"context"
	"fmt"
	"log"
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
		response, err := h.orchestrator.ProcessMessage(ctx, req.Message)
		if err != nil {
			log.Printf("Orchestrator error: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
				Error:   "AI service error",
				Message: fmt.Sprintf("Failed to process request: %v", err),
				Code:    fiber.StatusInternalServerError,
			})
		}

		return c.JSON(fiber.Map{
			"success":  true,
			"response": response,
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

