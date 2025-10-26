# Social Places API Integration

## Overview

The `/api/social` endpoint fetches top-rated places from Google Places API sorted by review count. These "socially popular spots" are automatically integrated into AI agent responses when users request travel plans.

## Direct API Usage

### Request
```bash
curl -X POST http://localhost:8080/api/social \
  -H "Content-Type: application/json" \
  -d '{
    "keyword": "restaurants",
    "location": "Tokyo",
    "limit": 5
  }'
```

### Response
```json
{
  "places": [
    {
      "placeId": "ChIJ...",
      "name": "Sukiyabashi Jiro",
      "address": "4 Chome-2-15 Ginza, Chuo City, Tokyo",
      "rating": 4.8,
      "reviewCount": 1250,
      "priceLevel": 4,
      "types": ["restaurant", "food"],
      "photoUrl": "https://maps.googleapis.com/...",
      "latitude": 35.6704,
      "longitude": 139.7633,
      "openNow": true
    }
  ],
  "query": "restaurants in Tokyo",
  "count": 5
}
```

## AI Agent Integration

### Example: Trip Planning

When a user asks for a travel plan:

**User Input:**
```
"I want to visit Tokyo for 5 days with a budget of 50,000 THB"
```

**AI Response (includes social spots):**
```markdown
# 5-Day Trip to Tokyo

**Budget:** 50,000 THB

## Itinerary
...

## Recommended Hotels
- Hotel Gracery Shinjuku - 3,500 THB/night (Rating: 4.5★)
- APA Hotel Tokyo - 2,800 THB/night (Rating: 4.3★)

## Socially Popular Spots
*Top-rated places based on reviews*

- **Senso-ji Temple** (4.8★, 12,453 reviews)
- **Tokyo Skytree** (4.7★, 10,892 reviews)
- **Meiji Shrine** (4.6★, 8,234 reviews)
- **Shibuya Crossing** (4.5★, 7,651 reviews)
- **Tsukiji Outer Market** (4.6★, 6,789 reviews)
```

### Example: Local Recommendations

When a user asks for local recommendations:

**User Input:**
```
"Find good cafes in Paris"
```

**AI Response (includes social spots):**
```markdown
# Nearby Cafe Recommendations

1. **Café de Flore**
   - Type: cafe
   - Rating: 4.6★
   - Distance: 1.2 km away
   - Address: 172 Boulevard Saint-Germain

...

## Socially Popular Cafe in Paris
*Top-rated by the community*

4. **Angelina Paris** (4.7★, 8,432 reviews)
   - Address: 226 Rue de Rivoli, 75001 Paris

5. **Café Kitsuné** (4.6★, 6,210 reviews)
   - Address: 51 Galerie de Montpensier, 75001 Paris
```

## Features

- ✅ **Sorted by Review Count**: Results are ordered by the number of reviews (highest first)
- ✅ **Cached Results**: Responses are cached in Redis for 1 hour
- ✅ **Automatic Integration**: No extra work needed - social spots appear automatically in AI responses
- ✅ **Graceful Degradation**: Works even if Google Places API is not configured
- ✅ **Flexible Queries**: Supports any keyword (restaurants, tourist attractions, cafes, etc.)

## Configuration

Add your Google Places API key to `.env`:

```env
GOOGLE_PLACES_API_KEY=your_api_key_here
GOOGLE_PLACES_API_URL=https://maps.googleapis.com/maps/api/place
```

## Rate Limits

- **Free Tier**: $200 credit/month (~28,000 requests)
- **Caching**: 1 hour per query (reduces API calls)
- **Best Practice**: Use reasonable limits (5-10 results) to minimize costs

## Error Handling

The service handles errors gracefully:

1. **Missing API Key**: Service returns "service unavailable" but AI agent continues to work
2. **Invalid Location**: Returns empty results
3. **API Errors**: Falls back gracefully, AI agent continues without social spots
4. **Rate Limits**: Returns cached results or skips social section

## Testing

Run tests:
```bash
cd backend
go test ./internal/handlers -v -run TestSocialHandler
```

Test the endpoint locally:
```bash
# Start the server
go run cmd/server/main.go

# Test with curl
curl -X POST http://localhost:8080/api/social \
  -H "Content-Type: application/json" \
  -d '{"keyword": "museums", "location": "Paris", "limit": 10}'
```
