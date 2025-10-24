# Multi-Agent Travel AI System

This document describes the complete multi-agent architecture implemented for the Travel AI Assistant.

## Architecture Overview

```
User Input → Intent Detection → Orchestrator → Specialized Agents → Response
                                    ├─ PlannerAgent
                                    ├─ FlightAgent
                                    ├─ WeatherAgent
                                    ├─ HotelAgent
                                    ├─ LocalAgent
                                    └─ BudgetAgent
```

## Agents

### 1. IntentAgent (`backend/agents/intent.go`)

**Purpose:** Classifies user messages into intents and extracts entities

**Intent Categories:**
- `plan_trip` - Create new travel plan
- `flight_check` - Check flight status
- `weather_check` - Get weather forecast
- `hotel_search` - Find hotels
- `local_recommendation` - Get nearby places
- `budget_inquiry` - Budget questions
- `plan_update` - Modify existing plan
- `general_chat` - Casual conversation

**Features:**
- Multi-language support (Thai, English)
- Entity extraction (destination, budget, duration, dates, etc.)
- LLM-powered with rule-based fallback
- Context-aware with current time

**Example:**
```go
agent := agents.NewIntentAgent(openaiKey)
result, err := agent.Detect(ctx, "อยากไปเที่ยวแคนาดา 7 วัน งบ 100,000 บาท")
// Result: intent="plan_trip", entities={destination, duration, budget}
```

### 2. PlannerAgent (`backend/agents/planner.go`)

**Purpose:** Creates and updates travel itineraries

**Features:**
- Detailed day-by-day itineraries
- Budget allocation per day
- Activity recommendations
- Plan updates based on conditions (e.g., weather changes)
- Markdown summary generation

**Example:**
```go
agent := agents.NewPlannerAgent(openaiKey)
plan, err := agent.CreatePlan(ctx, "Tokyo", 5, 80000)
// Returns: TripPlan with 5-day itinerary, activities, and budget breakdown
```

### 3. WeatherAgent (`backend/agents/weather.go`)

**Purpose:** Provides weather forecasts and activity suggestions

**Features:**
- 3-day weather forecast
- Rain probability detection (>60% threshold)
- Indoor activity suggestions for rainy weather
- OpenWeatherMap API integration
- Cached results via Redis

**Example:**
```go
agent := agents.NewWeatherAgent(openaiKey, weatherKey)
forecast, err := agent.GetForecast(ctx, "Bangkok")
// Returns: WeatherForecast with 3-day forecast and suggestions
```

### 4. FlightAgent (`backend/agents/flight.go`)

**Purpose:** Tracks flight status and sends notifications

**Features:**
- Flight status checking via AviationStack API
- Delay detection and quantification
- Polite notification generation for delays
- Real-time status updates

**Example:**
```go
agent := agents.NewFlightAgent(openaiKey, flightKey)
status, err := agent.CheckFlight(ctx, "JL708")
// Returns: FlightStatus with status, times, gate, and delay info
```

### 5. LocalAgent (`backend/agents/local.go`)

**Purpose:** Finds nearby places based on interests

**Features:**
- Location-based recommendations
- Interest filtering (restaurants, cafes, etc.)
- Rating and distance information
- Top 3 recommendations within 3 km radius

**Example:**
```go
agent := agents.NewLocalAgent(openaiKey)
places, err := agent.GetRecommendations(ctx, 13.7563, 100.5018, "ramen")
// Returns: []PlaceRecommendation with name, rating, distance, address
```

### 6. HotelAgent (`backend/agents/hotel.go`)

**Purpose:** Searches and recommends hotels

**Features:**
- Budget-based hotel search
- Price per night calculations
- Rating and distance information
- Multiple recommendations

**Example:**
```go
agent := agents.NewHotelAgent(openaiKey, hotelKey)
hotels, err := agent.SearchHotels(ctx, "Tokyo", 2500)
// Returns: []HotelRecommendation within budget
```

### 7. BudgetAgent (`backend/agents/budget.go`)

**Purpose:** Estimates budget breakdown for trips

**Features:**
- Standard allocation percentages:
  - Flights: 45%
  - Hotels: 25%
  - Food: 15%
  - Transport: 10%
  - Miscellaneous: 5%

**Example:**
```go
plan := agents.EstimateBudget(100000)
// Returns: BudgetPlan with allocated amounts
```

## Orchestrator (`backend/internal/orchestrator/orchestrator.go`)

**Purpose:** Coordinates multiple agents based on user intent

**Features:**
- Intent-based routing
- Multi-agent coordination
- Error handling and fallbacks
- Natural language response generation

**Workflow:**
1. Receive user message
2. Detect intent using IntentAgent
3. Route to appropriate handler(s)
4. Coordinate multiple agents if needed
5. Generate comprehensive natural language response

**Example:**
```go
orch := orchestrator.New(openaiKey, weatherKey, flightKey, hotelKey)
response, err := orch.ProcessMessage(ctx, "I want to visit Tokyo for 5 days")
// Coordinates PlannerAgent, WeatherAgent, and HotelAgent
```

## API Integration

### Handler (`backend/internal/handlers/plan.go`)

The `/api/plan` endpoint uses the orchestrator:

```go
POST /api/plan
{
  "message": "I want to visit Canada for 7 days with 100,000 baht"
}

Response:
{
  "success": true,
  "response": "# 7-Day Trip to Canada\n\n**Budget:** 100000 THB..."
}
```

## Example User Journeys

### Journey 1: Trip Planning (Thai)
```
Input: "อยากไปเที่ยวแคนาดา 7 วัน งบ 100,000 บาท"

Flow:
1. IntentAgent → plan_trip
2. Extract: destination="Canada", duration=7, budget=100000
3. PlannerAgent creates 7-day itinerary
4. WeatherAgent checks forecast
5. HotelAgent searches hotels
6. Return complete plan with weather info
```

### Journey 2: Weather Check (Thai)
```
Input: "วันนี้ฝนตกที่เกียวโตไหม ถ้าตกช่วยเปลี่ยนกิจกรรมให้หน่อย"

Flow:
1. IntentAgent → weather_check
2. WeatherAgent checks Kyoto weather
3. If rain >60%, suggest indoor activities
4. Return forecast with suggestions
```

### Journey 3: Flight Check (English)
```
Input: "Is flight JL708 on time?"

Flow:
1. IntentAgent → flight_check
2. Extract: flight_code="JL708"
3. FlightAgent checks AviationStack API
4. Return status with notification
```

### Journey 4: Local Recommendations (Thai)
```
Input: "อยากกินราเมนอร่อยใกล้ Shinjuku"

Flow:
1. IntentAgent → local_recommendation
2. Extract: interest="ramen", location="Shinjuku"
3. LocalAgent searches nearby ramen shops
4. Return top 3 recommendations
```

## Multi-Language Support

The system supports both Thai and English:

- Intent detection works in both languages
- Fallback rules handle common Thai keywords
- LLM processes natural language in both languages
- Responses are in the same language as input (when using LLM)

## Configuration

Required environment variables:

```env
# OpenAI for all agents
OPENAI_API_KEY=sk-...
OPENAI_MODEL=gpt-4o-mini

# Weather API
WEATHER_API_KEY=...
WEATHER_API_URL=https://api.openweathermap.org/data/2.5

# Flight API
FLIGHT_API_KEY=...
FLIGHT_API_URL=https://api.aviationstack.com/v1

# Hotel API (optional)
HOTEL_API_KEY=...
HOTEL_API_URL=https://api.booking.com
```

## Testing

All agents have comprehensive tests:

```bash
# Run all tests
go test ./...

# Run specific agent tests
go test ./agents/...

# Run orchestrator tests
go test ./internal/orchestrator/...

# Run integration tests
go test ./internal/orchestrator/integration_test.go
```

Test coverage includes:
- ✅ Thai language input
- ✅ English language input
- ✅ All intent categories
- ✅ Entity extraction
- ✅ Multi-agent coordination
- ✅ Error handling
- ✅ Fallback mechanisms

## Error Handling

The system gracefully handles:
- Missing API keys (uses fallback data)
- API failures (returns estimated data)
- Invalid input (provides helpful messages)
- Timeout scenarios (proper context handling)
- Missing entities (uses sensible defaults)

## Future Enhancements

Potential improvements:
- Google Places API integration for LocalAgent
- Real booking API integration for HotelAgent
- Redis caching for all agent responses
- Plan persistence and retrieval
- User preference learning
- Multi-city itineraries
- Group trip planning
- Real-time currency conversion
