package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/smithisrealdev/travel-ai-agent/backend/internal/config"
	"github.com/smithisrealdev/travel-ai-agent/backend/internal/database"
	"github.com/smithisrealdev/travel-ai-agent/backend/internal/handlers"
	"github.com/smithisrealdev/travel-ai-agent/backend/internal/orchestrator"
	"github.com/smithisrealdev/travel-ai-agent/backend/internal/services"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connections
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close()

	redis, err := database.NewRedisCache(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redis.Close()

	// Initialize services
	openaiService := services.NewOpenAIService(cfg)
	weatherService := services.NewWeatherService(cfg)
	flightService := services.NewFlightService(cfg)
	planService := services.NewPlanService(cfg)

	// Initialize orchestrator
	orch := orchestrator.New(
		cfg.OpenAI.APIKey,
		cfg.Weather.APIKey,
		cfg.Flight.APIKey,
		cfg.Hotel.APIKey,
	)

	// Initialize handlers
	travelHandler := handlers.NewTravelHandler(
		db,
		redis,
		openaiService,
		weatherService,
		flightService,
	)
	planHandler := handlers.NewPlanHandler(planService, orch)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Travel AI Agent API",
		ServerHeader: "Travel-AI-Agent",
		ErrorHandler: customErrorHandler,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} - ${latency}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: false,
	}))

	// API Routes
	api := app.Group("/api")

	// Plan endpoint
	api.Post("/plan", planHandler.CreateTravelPlan)

	// API v1 Routes
	apiv1 := app.Group("/api/v1")

	// Travel endpoints
	apiv1.Post("/travel/search", travelHandler.SearchTravel)
	apiv1.Get("/travel/history", travelHandler.GetSearchHistory)

	// Health check endpoint
	app.Get("/health", travelHandler.HealthCheck)

	// Root endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service": "Travel AI Agent API",
			"version": "1.0.0",
			"status":  "running",
		})
	})

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Gracefully shutting down...")
		app.Shutdown()
	}()

	// Start server
	address := fmt.Sprintf("%s:%s", cfg.Backend.Host, cfg.Backend.Port)
	log.Printf("Starting server on %s", address)

	if err := app.Listen(address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// customErrorHandler handles application errors
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"error":   true,
		"message": err.Error(),
		"code":    code,
	})
}
