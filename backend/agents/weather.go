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
)

// GetWeatherSummary fetches weather information for a city in a specific month
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