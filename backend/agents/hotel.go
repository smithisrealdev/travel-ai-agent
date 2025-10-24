package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sashabaranov/go-openai"
)

var ctx = context.Background()

// HotelRecommendation represents a hotel search result
type HotelRecommendation struct {
	Name          string  `json:"name"`
	PricePerNight float64 `json:"price_per_night"`
	Rating        float64 `json:"rating"`
	Address       string  `json:"address"`
	Distance      float64 `json:"distance_km"`
}

// HotelAgent handles hotel search and recommendations
type HotelAgent struct {
	client *openai.Client
	apiKey string
}

// NewHotelAgent creates a new hotel agent
func NewHotelAgent(openaiKey, hotelKey string) *HotelAgent {
	var client *openai.Client
	if openaiKey != "" {
		client = openai.NewClient(openaiKey)
	}
	return &HotelAgent{
		client: client,
		apiKey: hotelKey,
	}
}

// SearchHotels searches for hotels within budget
func (a *HotelAgent) SearchHotels(ctx context.Context, destination string, budget float64) ([]HotelRecommendation, error) {
	// Use LLM to generate realistic hotel recommendations
	if a.client != nil {
		recommendations, err := a.searchWithLLM(ctx, destination, budget)
		if err == nil && len(recommendations) > 0 {
			log.Printf("HotelAgent: Found %d hotels in %s within budget %.0f THB", 
				len(recommendations), destination, budget)
			return recommendations, nil
		}
	}

	// Fallback to estimated hotels
	recommendations := a.estimateHotels(destination, budget)
	log.Printf("HotelAgent: Generated %d estimated hotels for %s", len(recommendations), destination)
	return recommendations, nil
}

// searchWithLLM uses LLM to generate hotel recommendations
func (a *HotelAgent) searchWithLLM(ctx context.Context, destination string, budget float64) ([]HotelRecommendation, error) {
	prompt := fmt.Sprintf(`You are HotelAgent. Search for 3 hotels in %s with nightly rate around %.0f THB per night.

Return ONLY valid JSON array:
[
  {"name": "Hotel Name", "price_per_night": 2500.0, "rating": 4.5, "address": "Full address", "distance_km": 1.5},
  ...
]

Generate realistic hotel names, addresses, and ratings for %s.`, destination, budget, destination)

	resp, err := a.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: "gpt-4o-mini",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a hotel search assistant. Generate realistic hotel recommendations and return ONLY valid JSON array.",
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
		return nil, err
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from LLM")
	}

	content := resp.Choices[0].Message.Content

	var recommendations []HotelRecommendation
	err = json.Unmarshal([]byte(content), &recommendations)
	if err != nil {
		return nil, err
	}

	return recommendations, nil
}

// estimateHotels generates estimated hotel recommendations
func (a *HotelAgent) estimateHotels(destination string, budget float64) []HotelRecommendation {
	rand.Seed(time.Now().UnixNano())
	
	basePricePerNight := estimateHotelPricePerNight(destination)
	
	recommendations := make([]HotelRecommendation, 3)
	for i := 0; i < 3; i++ {
		variance := 0.8 + rand.Float64()*0.4 // 0.8 to 1.2
		pricePerNight := float64(basePricePerNight) * variance
		
		recommendations[i] = HotelRecommendation{
			Name:          estimateHotelName(destination),
			PricePerNight: pricePerNight,
			Rating:        4.0 + rand.Float64()*1.0,
			Address:       fmt.Sprintf("%s City Center", destination),
			Distance:      0.5 + rand.Float64()*3.0,
		}
	}

	return recommendations
}

// GetHotelPrice searches for affordable hotel prices in a city (Legacy function)
// Parameters:
//   - city: Destination city name (e.g., "Vancouver", "Tokyo")
//   - nights: Number of nights to stay
// Returns:
//   - price: Total hotel price in THB for the entire stay
//   - name: Name of the hotel
func GetHotelPrice(city string, nights int) (price int, name string) {
	// Default values
	defaultPricePerNight := 2500
	defaultName := "Budget Hotel"

	if city == "" || nights <= 0 {
		log.Printf("Invalid input: city=%s, nights=%d", city, nights)
		return defaultPricePerNight * max(1, nights), defaultName
	}

	// Try to get from Redis cache first
	cachedPrice, cachedName := getCachedHotelPrice(city, nights)
	if cachedPrice > 0 {
		log.Printf("Cache hit for hotel in %s: %s for %d THB", city, cachedName, cachedPrice)
		return cachedPrice, cachedName
	}

	// Try API call (Booking.com or Trip.com)
	apiKey := os.Getenv("HOTEL_API_KEY")
	if apiKey != "" {
		price, name = searchHotelAPI(city, nights, apiKey)
		if price > 0 {
			// Cache the result
			cacheHotelPrice(city, nights, price, name)
			log.Printf("API result for hotel in %s: %s for %d THB", city, name, price)
			return price, name
		}
	}

	// Fallback to estimated prices
	pricePerNight := estimateHotelPricePerNight(city)
	totalPrice := pricePerNight * nights
	hotelName := estimateHotelName(city)

	// Cache the estimated result
	cacheHotelPrice(city, nights, totalPrice, hotelName)

	log.Printf("Estimated hotel in %s: %s for %d THB (%d nights x %d THB/night)",
		city, hotelName, totalPrice, nights, pricePerNight)

	return totalPrice, hotelName
}

// getCachedHotelPrice retrieves hotel price from Redis cache
func getCachedHotelPrice(city string, nights int) (price int, name string) {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	if redisHost == "" {
		redisHost = "localhost"
	}
	if redisPort == "" {
		redisPort = "6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
		DB:   0,
	})
	defer rdb.Close()

	// Test connection
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("Redis connection failed: %v", err)
		return 0, ""
	}

	// Cache key format: hotel:<city>:<nights>
	cacheKey := fmt.Sprintf("hotel:%s:%d", strings.ToLower(city), nights)

	val, err := rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		// Cache miss
		return 0, ""
	} else if err != nil {
		log.Printf("Redis get error: %v", err)
		return 0, ""
	}

	// Parse cached JSON
	var cached struct {
		Price int    `json:"price"`
		Name  string `json:"name"`
	}

	if err := json.Unmarshal([]byte(val), &cached); err != nil {
		log.Printf("Failed to parse cached hotel data: %v", err)
		return 0, ""
	}

	return cached.Price, cached.Name
}

// cacheHotelPrice stores hotel price in Redis cache
func cacheHotelPrice(city string, nights, price int, name string) {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	if redisHost == "" {
		redisHost = "localhost"
	}
	if redisPort == "" {
		redisPort = "6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
		DB:   0,
	})
	defer rdb.Close()

	// Test connection
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("Redis connection failed, skipping cache: %v", err)
		return
	}

	// Cache key format: hotel:<city>:<nights>
	cacheKey := fmt.Sprintf("hotel:%s:%d", strings.ToLower(city), nights)

	data := struct {
		Price int    `json:"price"`
		Name  string `json:"name"`
	}{
		Price: price,
		Name:  name,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal hotel data: %v", err)
		return
	}

	// Cache for 24 hours
	expiration := 24 * time.Hour

	if err := rdb.Set(ctx, cacheKey, jsonData, expiration).Err(); err != nil {
		log.Printf("Failed to cache hotel data: %v", err)
		return
	}

	log.Printf("Cached hotel data for %s (key: %s)", city, cacheKey)
}

// searchHotelAPI queries Booking.com or Trip.com API
func searchHotelAPI(city string, nights int, apiKey string) (price int, name string) {
	// This is a placeholder for actual API integration
	// In production, implement actual API calls to:
	// - Booking.com Affiliate API
	// - Trip.com Affiliate API
	// - Agoda API
	
	log.Printf("Searching hotels in %s via API (not implemented yet)", city)
	
	// TODO: Implement actual API call
	// Example structure:
	// - Make HTTP request to Booking.com API
	// - Parse response JSON
	// - Extract cheapest hotel
	// - Return price and name
	
	return 0, ""
}

// estimateHotelPricePerNight returns estimated hotel price per night based on city
func estimateHotelPricePerNight(city string) int {
	city = strings.ToLower(city)

	// Price per night in THB for budget hotels
	cityPrices := map[string]int{
		"vancouver":     2500,
		"tokyo":         2000,
		"seoul":         1800,
		"singapore":     2200,
		"hong kong":     2400,
		"taipei":        1600,
		"kuala lumpur":  1200,
		"jakarta":       1000,
		"sydney":        3000,
		"london":        4000,
		"paris":         3500,
		"frankfurt":     3200,
		"los angeles":   3800,
		"new york":      4500,
		"dubai":         2800,
		"bangkok":       1000,
		"phuket":        1500,
		"chiang mai":    800,
		"pattaya":       1200,
		"krabi":         1400,
		"osaka":         2200,
		"kyoto":         2400,
		"busan":         1600,
		"bali":          1300,
		"hanoi":         900,
		"ho chi minh":   1100,
		"phnom penh":    700,
		"vientiane":     600,
		"yangon":        800,
		"manila":        1000,
		"cebu":          900,
	}

	// Check for exact match
	if price, ok := cityPrices[city]; ok {
		return price
	}

	// Check for partial match
	for key, price := range cityPrices {
		if strings.Contains(city, key) || strings.Contains(key, city) {
			return price
		}
	}

	// Default price for unknown cities
	return 2500
}

// estimateHotelName returns a realistic hotel name based on city
func estimateHotelName(city string) string {
	city = strings.ToLower(city)

	hotelNames := map[string][]string{
		"vancouver":     {"Comfort Inn Downtown", "Budget Hotel Vancouver", "City Center Inn"},
		"tokyo":         {"Tokyo Budget Hotel", "Shinjuku Comfort Inn", "Asakusa Guesthouse"},
		"seoul":         {"Seoul Budget Hotel", "Gangnam Inn", "Myeongdong Guesthouse"},
		"singapore":     {"Budget Hotel Singapore", "Chinatown Inn", "Little India Hotel"},
		"hong kong":     {"Hong Kong Budget Inn", "Tsim Sha Tsui Hotel", "Kowloon Guesthouse"},
		"taipei":        {"Taipei Budget Hotel", "Ximending Inn", "Da'an Guesthouse"},
		"kuala lumpur":  {"KL Budget Hotel", "Bukit Bintang Inn", "KLCC Guesthouse"},
		"bangkok":       {"Bangkok Budget Inn", "Sukhumvit Hotel", "Silom Guesthouse"},
		"phuket":        {"Patong Beach Hotel", "Phuket Budget Inn", "Kata Guesthouse"},
		"london":        {"London Budget Hotel", "Westminster Inn", "Camden Guesthouse"},
		"paris":         {"Paris Budget Hotel", "Marais Inn", "Montmartre Guesthouse"},
		"new york":      {"NYC Budget Hotel", "Manhattan Inn", "Brooklyn Guesthouse"},
		"sydney":        {"Sydney Budget Inn", "Darling Harbour Hotel", "Bondi Guesthouse"},
		"dubai":         {"Dubai Budget Hotel", "Deira Inn", "Downtown Guesthouse"},
	}

	// Check for exact match
	if names, ok := hotelNames[city]; ok {
		// Return random name from list
		rand.Seed(time.Now().UnixNano())
		return names[rand.Intn(len(names))]
	}

	// Check for partial match
	for key, names := range hotelNames {
		if strings.Contains(city, key) || strings.Contains(key, city) {
			rand.Seed(time.Now().UnixNano())
			return names[rand.Intn(len(names))]
		}
	}

	// Default hotel name
	return fmt.Sprintf("%s Budget Hotel", strings.Title(city))
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
