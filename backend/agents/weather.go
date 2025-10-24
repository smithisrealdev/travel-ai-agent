package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sashabaranov/go-openai"
)

// DayForecast represents a single day's forecast
type DayForecast struct {
	Date        string  `json:"date"`
	Temperature float64 `json:"temperature"`
	Condition   string  `json:"condition"`
	RainProb    float64 `json:"rain_probability"`
}

// WeatherForecast represents complete weather information with forecast
type WeatherForecast struct {
	City        string        `json:"city"`
	Temperature float64       `json:"temperature"`
	Condition   string        `json:"condition"`
	RainProb    float64       `json:"rain_probability"`
	Forecast    []DayForecast `json:"forecast"`
	Suggestion  string        `json:"suggestion"`
}

// WeatherAgent handles weather forecasting and suggestions
type WeatherAgent struct {
	client *openai.Client
	apiKey string
}

// NewWeatherAgent creates a new weather agent
func NewWeatherAgent(openaiKey, weatherKey string) *WeatherAgent {
	var client *openai.Client
	if openaiKey != "" {
		client = openai.NewClient(openaiKey)
	}
	return &WeatherAgent{
		client: client,
		apiKey: weatherKey,
	}
}

// GetForecast gets 3-day weather forecast and suggestions
func (a *WeatherAgent) GetForecast(ctx context.Context, city string) (*WeatherForecast, error) {
	forecast := &WeatherForecast{
		City:     city,
		Forecast: make([]DayForecast, 0),
	}

	// Try to get forecast from OpenWeatherMap API
	if a.apiKey != "" {
		apiForecasts, avgRainProb := a.fetchForecastFromAPI(city)
		if len(apiForecasts) > 0 {
			forecast.Forecast = apiForecasts
			forecast.RainProb = avgRainProb
			if len(apiForecasts) > 0 {
				forecast.Temperature = apiForecasts[0].Temperature
				forecast.Condition = apiForecasts[0].Condition
			}
		}
	}

	// Fallback to estimated forecast
	if len(forecast.Forecast) == 0 {
		forecast.Forecast = a.estimateForecast(city)
		forecast.RainProb = a.estimateRainProb(city)
		if len(forecast.Forecast) > 0 {
			forecast.Temperature = forecast.Forecast[0].Temperature
			forecast.Condition = forecast.Forecast[0].Condition
		}
	}

	// Generate suggestion if rain probability > 60%
	if forecast.RainProb > 60 {
		forecast.Suggestion = a.generateRainSuggestion(ctx, city, forecast.RainProb)
	} else {
		forecast.Suggestion = "Weather looks good for outdoor activities!"
	}

	log.Printf("WeatherAgent: Forecast for %s - Rain prob: %.1f%%, Suggestion: %s", 
		city, forecast.RainProb, forecast.Suggestion)

	return forecast, nil
}

// fetchForecastFromAPI gets forecast from OpenWeatherMap
func (a *WeatherAgent) fetchForecastFromAPI(city string) ([]DayForecast, float64) {
	url := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/forecast?q=%s&appid=%s&units=metric",
		city, a.apiKey,
	)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("WeatherAgent: API request failed: %v", err)
		return nil, 0
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("WeatherAgent: API returned status %d", resp.StatusCode)
		return nil, 0
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("WeatherAgent: Failed to read response: %v", err)
		return nil, 0
	}

	var result struct {
		List []struct {
			Dt   int64 `json:"dt"`
			Main struct {
				Temp float64 `json:"temp"`
			} `json:"main"`
			Weather []struct {
				Main string `json:"main"`
			} `json:"weather"`
			Pop float64 `json:"pop"` // Probability of precipitation
		} `json:"list"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("WeatherAgent: Failed to parse response: %v", err)
		return nil, 0
	}

	// Group by day and get 3-day forecast
	forecasts := make([]DayForecast, 0, 3)
	rainProbs := make([]float64, 0)
	seenDates := make(map[string]bool)

	for _, item := range result.List {
		date := time.Unix(item.Dt, 0).Format("2006-01-02")
		
		if !seenDates[date] && len(forecasts) < 3 {
			condition := "Clear"
			if len(item.Weather) > 0 {
				condition = item.Weather[0].Main
			}

			forecasts = append(forecasts, DayForecast{
				Date:        date,
				Temperature: item.Main.Temp,
				Condition:   condition,
				RainProb:    item.Pop * 100,
			})
			rainProbs = append(rainProbs, item.Pop*100)
			seenDates[date] = true
		}
	}

	// Calculate average rain probability
	avgRainProb := 0.0
	if len(rainProbs) > 0 {
		sum := 0.0
		for _, prob := range rainProbs {
			sum += prob
		}
		avgRainProb = sum / float64(len(rainProbs))
	}

	return forecasts, avgRainProb
}

// estimateForecast provides estimated 3-day forecast
func (a *WeatherAgent) estimateForecast(city string) []DayForecast {
	baseTemp := estimateTemperature(city, time.Now().Format("January"))
	condition := estimateCondition(city, time.Now().Format("January"))

	forecasts := make([]DayForecast, 3)
	for i := 0; i < 3; i++ {
		date := time.Now().AddDate(0, 0, i).Format("2006-01-02")
		forecasts[i] = DayForecast{
			Date:        date,
			Temperature: float64(baseTemp + (i - 1)), // Slight variation
			Condition:   condition,
			RainProb:    a.estimateRainProb(city),
		}
	}

	return forecasts
}

// estimateRainProb estimates rain probability for a city
func (a *WeatherAgent) estimateRainProb(city string) float64 {
	city = strings.ToLower(city)
	month := strings.ToLower(time.Now().Format("January"))

	rainyMonths := map[string][]string{
		"vancouver":    {"november", "december", "january", "february", "march"},
		"bangkok":      {"may", "june", "july", "august", "september", "october"},
		"tokyo":        {"june", "july", "september"},
		"singapore":    {"november", "december", "january"},
		"kuala lumpur": {"april", "may", "october", "november"},
	}

	if months, ok := rainyMonths[city]; ok {
		for _, m := range months {
			if m == month {
				return 75.0 // High rain probability
			}
		}
	}

	return 20.0 // Low rain probability
}

// generateRainSuggestion generates indoor activity suggestions
func (a *WeatherAgent) generateRainSuggestion(ctx context.Context, city string, rainProb float64) string {
	if a.client == nil {
		return fmt.Sprintf("High chance of rain (%.0f%%). Consider indoor activities like museums, shopping malls, or indoor entertainment.", rainProb)
	}

	prompt := fmt.Sprintf(`You are WeatherAgent. There's a %.0f%% chance of rain in %s.
Recommend 2-3 indoor activities suitable for rainy weather. Keep it brief and friendly.`, rainProb, city)

	resp, err := a.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: "gpt-4o-mini",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a helpful weather advisor. Provide brief, friendly indoor activity suggestions.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.7,
			MaxTokens:   150,
		},
	)

	if err != nil {
		log.Printf("WeatherAgent: Failed to generate suggestion: %v", err)
		return fmt.Sprintf("High chance of rain (%.0f%%). Consider indoor activities like museums, shopping malls, or indoor entertainment.", rainProb)
	}

	if len(resp.Choices) > 0 {
		return resp.Choices[0].Message.Content
	}

	return fmt.Sprintf("High chance of rain (%.0f%%). Consider indoor activities.", rainProb)
}

// GetWeatherSummary fetches weather information for a city in a specific month (Legacy function)
// Parameters:
//   - city: City name (e.g., "Vancouver", "Tokyo")
//   - month: Month name or number (e.g., "December", "12")
// Returns:
//   - avgTemp: Average temperature in Celsius
//   - condition: Main weather condition (e.g., "Sunny", "Rainy", "Cloudy")
func GetWeatherSummary(city, month string) (avgTemp int, condition string) {
	// Default values
	defaultTemp := 15
	defaultCondition := "Sunny"

	if city == "" {
		log.Printf("Invalid input: city is empty")
		return defaultTemp, defaultCondition
	}

	// Normalize month input
	normalizedMonth := normalizeMonth(month)

	// Try to get from Redis cache first
	cachedTemp, cachedCondition := getCachedWeather(city, normalizedMonth)
	if cachedTemp != 0 {
		log.Printf("Cache hit for weather in %s (%s): %d°C, %s", city, normalizedMonth, cachedTemp, cachedCondition)
		return cachedTemp, cachedCondition
	}

	// Try OpenWeatherMap API
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey != "" {
		avgTemp, condition = fetchWeatherFromAPI(city, apiKey)
		if avgTemp != 0 {
			// Cache the result
			cacheWeather(city, normalizedMonth, avgTemp, condition)
			log.Printf("API result for weather in %s: %d°C, %s", city, avgTemp, condition)
			return avgTemp, condition
		}
	}

	// Fallback to estimated weather
	avgTemp = estimateTemperature(city, normalizedMonth)
	condition = estimateCondition(city, normalizedMonth)

	// Cache the estimated result
	cacheWeather(city, normalizedMonth, avgTemp, condition)

	log.Printf("Estimated weather for %s in %s: %d°C, %s", city, normalizedMonth, avgTemp, condition)

	return avgTemp, condition
}

// normalizeMonth converts month input to standard format
func normalizeMonth(month string) string {
	month = strings.ToLower(strings.TrimSpace(month))

	monthMap := map[string]string{
		"1": "january", "jan": "january", "january": "january",
		"2": "february", "feb": "february", "february": "february",
		"3": "march", "mar": "march", "march": "march",
		"4": "april", "apr": "april", "april": "april",
		"5": "may", "may": "may",
		"6": "june", "jun": "june", "june": "june",
		"7": "july", "jul": "july", "july": "july",
		"8": "august", "aug": "august", "august": "august",
		"9": "september", "sep": "september", "september": "september",
		"10": "october", "oct": "october", "october": "october",
		"11": "november", "nov": "november", "november": "november",
		"12": "december", "dec": "december", "december": "december",
	}

	if normalized, ok := monthMap[month]; ok {
		return normalized
	}

	return "january"
}

// getCachedWeather retrieves weather from Redis cache
func getCachedWeather(city, month string) (avgTemp int, condition string) {
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

	ctx := context.Background()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("Redis connection failed: %v", err)
		return 0, ""
	}

	cacheKey := fmt.Sprintf("weather:%s:%s", strings.ToLower(city), strings.ToLower(month))

	val, err := rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return 0, ""
	} else if err != nil {
		log.Printf("Redis get error: %v", err)
		return 0, ""
	}

	var cached struct {
		AvgTemp   int    `json:"avg_temp"`
		Condition string `json:"condition"`
	}

	if err := json.Unmarshal([]byte(val), &cached); err != nil {
		log.Printf("Failed to parse cached weather data: %v", err)
		return 0, ""
	}

	return cached.AvgTemp, cached.Condition
}

// cacheWeather stores weather in Redis cache
func cacheWeather(city, month string, avgTemp int, condition string) {
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

	ctx := context.Background()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("Redis connection failed, skipping cache: %v", err)
		return
	}

	cacheKey := fmt.Sprintf("weather:%s:%s", strings.ToLower(city), strings.ToLower(month))

	data := struct {
		AvgTemp   int    `json:"avg_temp"`
		Condition string `json:"condition"`
	}{
		AvgTemp:   avgTemp,
		Condition: condition,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal weather data: %v", err)
		return
	}

	expiration := 24 * time.Hour

	if err := rdb.Set(ctx, cacheKey, jsonData, expiration).Err(); err != nil {
		log.Printf("Failed to cache weather data: %v", err)
		return
	}

	log.Printf("Cached weather data for %s in %s (key: %s)", city, month, cacheKey)
}

// fetchWeatherFromAPI calls OpenWeatherMap API
func fetchWeatherFromAPI(city, apiKey string) (avgTemp int, condition string) {
	url := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric",
		city, apiKey,
	)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("OpenWeather API request failed: %v", err)
		return 0, ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("OpenWeather API returned status %d", resp.StatusCode)
		return 0, ""
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read OpenWeather response: %v", err)
		return 0, ""
	}

	var result struct {
		Main struct {
			Temp float64 `json:"temp"`
		} `json:"main"`
		Weather []struct {
			Main string `json:"main"`
		} `json:"weather"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("Failed to parse OpenWeather response: %v", err)
		return 0, ""
	}

	avgTemp = int(result.Main.Temp)

	if len(result.Weather) > 0 {
		condition = result.Weather[0].Main
	} else {
		condition = "Clear"
	}

	return avgTemp, condition
}

// estimateTemperature returns estimated temperature
func estimateTemperature(city, month string) int {
	city = strings.ToLower(city)
	month = strings.ToLower(month)

	cityTemps := map[string]map[string]int{
		"vancouver": {
			"january": 6, "february": 7, "march": 9, "april": 12,
			"may": 15, "june": 18, "july": 21, "august": 21,
			"september": 18, "october": 13, "november": 9, "december": 6,
		},
		"tokyo": {
			"january": 6, "february": 7, "march": 11, "april": 16,
			"may": 20, "june": 23, "july": 27, "august": 29,
			"september": 25, "october": 19, "november": 14, "december": 9,
		},
		"bangkok": {
			"january": 27, "february": 29, "march": 30, "april": 31,
			"may": 30, "june": 29, "july": 29, "august": 29,
			"september": 28, "october": 28, "november": 27, "december": 26,
		},
	}

	if temps, ok := cityTemps[city]; ok {
		if temp, ok := temps[month]; ok {
			return temp
		}
	}

	return 15
}

// estimateCondition returns estimated weather condition
func estimateCondition(city, month string) string {
	city = strings.ToLower(city)
	month = strings.ToLower(month)

	rainyMonths := map[string][]string{
		"vancouver": {"november", "december", "january", "february", "march"},
		"bangkok": {"may", "june", "july", "august", "september", "october"},
		"tokyo": {"june", "july", "september"},
	}

	if months, ok := rainyMonths[city]; ok {
		for _, m := range months {
			if m == month {
				return "Rainy"
			}
		}
	}

	return "Sunny"
}