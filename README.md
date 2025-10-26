# Travel AI Agent 🌍

AI-powered travel planning assistant MVP that helps you plan your perfect trip with personalized recommendations, budget breakdowns, and real-time information.

## ✨ Features

- ✅ **AI-powered travel planning** - Intelligent trip suggestions powered by OpenAI GPT-4o-mini
- ✅ **Multi-language support** - Supports both Thai and English for seamless communication
- ✅ **Flight search integration** - Real-time flight data and pricing from AviationStack
- ✅ **Hotel recommendations** - Smart hotel suggestions based on your budget and preferences
- ✅ **Weather forecasts** - Current weather and forecasts for your destination
- ✅ **Budget breakdown** - Detailed cost analysis for flights, hotels, food, and activities
- ✅ **Socially popular spots** - Top-rated places from Google Places sorted by review count
- ✅ **Real-time chat interface** - Interactive conversation-based travel planning

## 🏗️ Architecture Overview

The Travel AI Agent uses a **microservices architecture** with specialized AI agents:

- **Intent Agent** - Understands user requests and extracts travel parameters
- **Flight Agent** - Searches and recommends flights based on budget and dates
- **Hotel Agent** - Finds accommodation options matching preferences
- **Weather Agent** - Provides weather forecasts and climate recommendations
- **Budget Agent** - Calculates and breaks down trip costs

**Technology Stack:**
- **Backend**: Go 1.21 + Fiber web framework
- **Frontend**: Nuxt 3 + Vue 3 + Tailwind CSS
- **Database**: PostgreSQL 15 for data persistence
- **Cache**: Redis 7 for performance optimization
- **AI**: OpenAI GPT-4o-mini
- **APIs**: OpenWeatherMap, AviationStack, Booking.com

## 📋 Prerequisites

Before you begin, ensure you have the following:

### Required Software
- **Docker** and **Docker Compose** (for containerized deployment)
- **Node.js 20+** (for local frontend development)
- **Go 1.21+** (for local backend development)

### Required API Keys

You'll need to obtain the following API keys:

1. **OpenAI API Key** (Required)
   - Sign up at https://platform.openai.com/api-keys
   - Create a new API key
   - Used for AI-powered travel planning

2. **OpenWeatherMap API Key** (Required)
   - Sign up at https://openweathermap.org/api
   - Free tier available
   - Used for weather forecasts

3. **AviationStack API Key** (Required)
   - Sign up at https://aviationstack.com
   - Free tier available
   - Used for flight searches

4. **Google Places API Key** (Required for Social Features)
   - Sign up at https://console.cloud.google.com/
   - Enable Places API (New)
   - Used for fetching socially popular spots

5. **Booking.com API Key** (Optional)
   - Used for hotel recommendations
   - Falls back to generic recommendations if not provided

## 🚀 Installation & Setup

### Quick Start (Docker - Recommended)

```bash
# 1. Clone the repository
git clone https://github.com/smithisrealdev/travel-ai-agent.git
cd travel-ai-agent

# 2. Copy environment file
cp .env.example .env

# 3. Edit .env and add your API keys
nano .env  # or use your preferred editor

# 4. Start all services
docker compose up --build
```

That's it! The application will be available at:
- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379

### Environment Variables

Edit the `.env` file and configure these essential variables:

```env
# OpenAI Configuration (Required)
OPENAI_API_KEY=sk-your-openai-api-key-here
OPENAI_MODEL=gpt-4o-mini

# Weather API (Required)
WEATHER_API_KEY=your-openweathermap-api-key-here
WEATHER_API_URL=https://api.openweathermap.org/data/2.5

# Flight API (Required)
FLIGHT_API_KEY=your-aviationstack-api-key-here
FLIGHT_API_URL=https://api.aviationstack.com/v1

# Hotel API (Optional)
HOTEL_API_KEY=your-booking-api-key-here
HOTEL_API_URL=https://api.booking.com

# Database Configuration
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_USER=travelagent
POSTGRES_PASSWORD=your_secure_password_here
POSTGRES_DB=travelagent

# Redis Configuration
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=

# Backend Configuration
BACKEND_PORT=8080
BACKEND_HOST=0.0.0.0
```

### Local Development

If you prefer to run services locally without Docker:

**Backend:**
```bash
cd backend

# Install dependencies
go mod download

# Run the server
go run cmd/server/main.go
```

**Frontend:**
```bash
cd frontend

# Install dependencies
npm install

# Run development server
npm run dev
```

## 🎯 Running the Application

### With Docker (Recommended)

```bash
# Start all services
docker compose up

# Run in background
docker compose up -d

# Stop services
docker compose down

# View logs
docker compose logs -f backend
docker compose logs -f frontend
```

### Accessing Services

Once running, you can access:
- **Frontend UI**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **Health Check**: http://localhost:8080/health
- **API Root**: http://localhost:8080/api
- **PostgreSQL**: localhost:5432 (user: `travelagent`)
- **Redis**: localhost:6379

## 📡 API Documentation

### Primary Endpoint: `POST /api/plan`

This is the main endpoint for creating travel plans. It accepts natural language input in Thai or English.

#### Request Example (cURL)

```bash
curl -X POST http://localhost:8080/api/plan \
  -H "Content-Type: application/json" \
  -d '{
    "message": "อยากไปเที่ยวแคนาดาในงบ 100,000 บาท"
  }'
```

#### Example Request (Thai)

```json
{
  "message": "อยากไปเที่ยวแคนาดาในงบ 100,000 บาท 7 วัน"
}
```

#### Example Request (English)

```json
{
  "message": "I want to visit Tokyo for 5 days with a budget of 50,000 THB"
}
```

#### Example Response

```json
{
  "success": true,
  "plan": {
    "destination": "Vancouver, Canada",
    "duration": 7,
    "budget": 100000,
    "breakdown": {
      "flights": 35000,
      "hotels": 28000,
      "food": 15000,
      "activities": 17000,
      "misc": 5000
    },
    "weather": {
      "temperature": 18,
      "condition": "Partly Cloudy",
      "forecast": "Pleasant weather for sightseeing"
    },
    "recommendations": [
      "Visit Stanley Park",
      "Explore Granville Island",
      "Capilano Suspension Bridge"
    ],
    "summary": "## 🇨🇦 Your Canada Trip Plan\n\n**Destination:** Vancouver\n**Duration:** 7 days\n**Budget:** 100,000 THB\n\n### 💰 Budget Breakdown\n- Flights: 35,000 THB\n- Hotels: 28,000 THB (4,000 THB/night)\n- Food: 15,000 THB\n- Activities: 17,000 THB\n- Miscellaneous: 5,000 THB\n\n### ✈️ Flights\nRound-trip flights from Bangkok to Vancouver available from 35,000 THB.\n\n### 🏨 Accommodation\nBudget hotels in downtown Vancouver: ~4,000 THB/night\n\n### 🌤️ Weather\nCurrent: 18°C, Partly Cloudy\nGreat weather for exploring!\n\n### 📍 Recommended Activities\n1. Stanley Park - Free\n2. Granville Island Public Market - Budget-friendly\n3. Capilano Suspension Bridge - ~1,000 THB\n4. Vancouver Lookout - ~500 THB\n5. Gastown Walking Tour - Free\n\nHave a wonderful trip! 🎉"
  }
}
```

### Additional Endpoints

#### Social Places

**POST** `/api/social`

Fetches top-rated places from Google Places API sorted by review count. Returns socially popular spots that are integrated into AI agent responses.

**Request:**
```json
{
  "keyword": "restaurants",
  "location": "Tokyo",
  "limit": 10
}
```

**Response:**
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
  "count": 10
}
```

**Parameters:**
- `keyword` (required): Type of place to search for (e.g., "restaurants", "tourist attractions", "cafes")
- `location` (required): City or location name
- `limit` (optional): Maximum number of results (default: 10)

**Features:**
- Results sorted by review count (highest first)
- Cached for 1 hour for performance
- Integrated into AI agent trip planning responses
- Shown as "Socially Popular Spots" in travel plans

#### Travel Search (v1)

**POST** `/api/v1/travel/search`

```json
{
  "destination": "Paris",
  "startDate": "2024-06-01",
  "endDate": "2024-06-07",
  "budget": 2000,
  "preferences": {
    "culture": true,
    "adventure": false
  }
}
```

#### Search History

**GET** `/api/v1/travel/history?userId=user123`

Retrieve user's search history.

#### Health Check

**GET** `/health`

Returns service status and health information.

## 📁 Project Structure

```
travel-ai-agent/
├── backend/                     # Go backend service
│   ├── agents/                  # AI agents (intent, flight, hotel, weather, budget)
│   │   ├── intent.go           # Intent detection agent
│   │   ├── flight.go           # Flight search agent
│   │   ├── hotel.go            # Hotel recommendation agent
│   │   ├── weather.go          # Weather forecast agent
│   │   └── budget.go           # Budget calculation agent
│   ├── cmd/
│   │   └── server/
│   │       └── main.go         # Main server entry point
│   ├── internal/
│   │   ├── config/             # Configuration management
│   │   ├── database/           # Database connections (PostgreSQL, Redis)
│   │   ├── handlers/           # HTTP request handlers
│   │   ├── models/             # Data models
│   │   └── services/           # Business logic services
│   ├── Dockerfile              # Backend Docker configuration
│   └── go.mod                  # Go module dependencies
│
├── frontend/                    # Nuxt 3 frontend
│   ├── assets/                 # Static assets
│   ├── components/             # Vue components
│   │   └── TravelSearch.vue    # Main search component
│   ├── composables/            # Composable functions
│   ├── pages/                  # Nuxt pages
│   │   └── index.vue           # Home page
│   ├── stores/                 # Pinia stores
│   │   └── travel.ts           # Travel state management
│   ├── types/                  # TypeScript definitions
│   ├── Dockerfile              # Frontend Docker configuration
│   ├── nuxt.config.ts          # Nuxt configuration
│   ├── package.json            # NPM dependencies
│   └── tailwind.config.js      # Tailwind CSS configuration
│
├── docker-compose.yml           # Docker services configuration
├── .env.example                # Environment variables template
├── .gitignore                  # Git ignore rules
└── README.md                   # This file
```

## 🛠️ Tech Stack Details

### Backend Stack

- **Language**: Go 1.21
- **Web Framework**: Fiber v2 (Express-like framework for Go)
- **Database**: PostgreSQL 15
- **Cache**: Redis 7
- **AI Engine**: OpenAI GPT-4o-mini
- **HTTP Client**: Standard Go net/http
- **Configuration**: Environment variables

### Frontend Stack

- **Framework**: Nuxt 3 (Vue 3 meta-framework)
- **UI Library**: Vue 3 Composition API
- **Styling**: Tailwind CSS v3
- **HTTP Client**: Axios
- **State Management**: Pinia
- **Language**: TypeScript
- **Package Manager**: npm

### Infrastructure

- **Containerization**: Docker
- **Orchestration**: Docker Compose
- **Reverse Proxy**: Built-in with Fiber
- **Database Migrations**: SQL scripts

## 🧪 Development

### Running Tests

**Backend Tests:**
```bash
cd backend

# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./agents/...
```

**Frontend Tests:**
```bash
cd frontend

# Run tests (if configured)
npm run test
```

### Building for Production

**Using Docker Compose:**
```bash
# Build production images
docker compose -f docker-compose.prod.yml up --build

# Or using the default compose file
docker compose up --build
```

**Manual Build:**
```bash
# Backend
cd backend
go build -o main cmd/server/main.go

# Frontend
cd frontend
npm run build
npm run preview
```

## 🔧 Troubleshooting

### Common Issues and Solutions

#### Docker Credential Errors
```bash
# Remove problematic Docker config
rm ~/.docker/config.json

# Re-run docker compose
docker compose up --build
```

#### Port Already in Use
```bash
# Change ports in .env file
BACKEND_PORT=8081
FRONTEND_PORT=3001

# Or stop conflicting services
lsof -ti:8080 | xargs kill -9  # macOS/Linux
```

#### API Key Errors
```bash
# Verify .env file has valid keys
cat .env | grep API_KEY

# Check if keys are loaded in container
docker compose exec backend env | grep API_KEY
```

#### Frontend Build Errors
```bash
# Ensure you're using Node 20+
node --version

# Clear node_modules and reinstall
cd frontend
rm -rf node_modules package-lock.json
npm install
```

#### Database Connection Issues
```bash
# Check PostgreSQL logs
docker compose logs postgres

# Verify connection
docker compose exec postgres psql -U travelagent -d travelagent

# Reset database
docker compose down -v
docker compose up -d
```

#### Redis Connection Issues
```bash
# Check Redis logs
docker compose logs redis

# Test connection
docker compose exec redis redis-cli ping

# Should return: PONG
```

## 🔑 API Keys Setup Guide

### 1. OpenAI API Key

1. Go to https://platform.openai.com/api-keys
2. Sign in or create an account
3. Click "Create new secret key"
4. Copy the key (starts with `sk-`)
5. Add to `.env`: `OPENAI_API_KEY=sk-...`

**Cost**: Pay-per-use, GPT-4o-mini is very affordable (~$0.15/1M tokens)

### 2. OpenWeatherMap API Key

1. Sign up at https://openweathermap.org/api
2. Go to API keys section
3. Generate a new API key
4. Add to `.env`: `WEATHER_API_KEY=...`

**Cost**: Free tier includes 1,000 calls/day

### 3. AviationStack API Key

1. Sign up at https://aviationstack.com
2. Choose a plan (free tier available)
3. Get your API access key
4. Add to `.env`: `FLIGHT_API_KEY=...`

**Cost**: Free tier includes 100 requests/month

### 4. Google Places API Key (Required for Social Features)

1. Go to https://console.cloud.google.com/
2. Create a new project or select an existing one
3. Enable the "Places API" (New)
4. Go to "Credentials" and create an API key
5. Add to `.env`: `GOOGLE_PLACES_API_KEY=...`

**Cost**: Free tier includes $200 credit/month (~28,000 requests)

### 5. Booking.com API Key (Optional)

1. Apply for API access at https://www.booking.com/affiliate
2. Wait for approval
3. Get API credentials
4. Add to `.env`: `HOTEL_API_KEY=...`

**Note**: This is optional; the system works without it using fallback recommendations.

## 🤝 Contributing

Contributions are welcome! Here's how you can help:

1. **Fork the repository**
2. **Create a feature branch**
   ```bash
   git checkout -b feature/amazing-feature
   ```
3. **Make your changes**
4. **Commit your changes**
   ```bash
   git commit -m 'Add some amazing feature'
   ```
5. **Push to the branch**
   ```bash
   git push origin feature/amazing-feature
   ```
6. **Open a Pull Request**

Please ensure your PR:
- Follows the existing code style
- Includes tests for new features
- Updates documentation as needed
- Has a clear description of changes

## 📄 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## 🎨 Screenshots

> **Note**: Add screenshots here showing:
> - Chat interface with travel planning conversation
> - Travel plan results with budget breakdown
> - Weather and flight information display
> - Multi-language support demonstration

## 🗺️ Roadmap

Future features and improvements planned:

- [ ] **User authentication** - Login/signup with JWT
- [ ] **Save/share trip plans** - Bookmark and share itineraries
- [ ] **Multi-destination trips** - Plan trips to multiple cities
- [ ] **Hotel booking integration** - Direct booking through the platform
- [ ] **Flight booking integration** - Complete booking flow
- [ ] **Mobile app** - React Native mobile application
- [ ] **Email notifications** - Price alerts and trip reminders
- [ ] **Social features** - Share trips with friends
- [ ] **Trip reviews** - Rate and review destinations
- [ ] **Advanced filters** - More detailed search criteria
- [ ] **Currency converter** - Support for multiple currencies
- [ ] **Offline mode** - Access saved plans without internet

## 👨‍💻 Author

Created by [@smithisrealdev](https://github.com/smithisrealdev)

## 🙏 Acknowledgments

Special thanks to:
- **OpenAI** for providing the GPT-4o-mini API
- **OpenWeatherMap** for weather data
- **AviationStack** for flight information
- **Go Fiber** community for the excellent web framework
- **Nuxt.js** community for the amazing Vue framework
- All open-source contributors

## 📞 Support

Need help? Here's how to get support:

- **Issues**: Open an issue on [GitHub Issues](https://github.com/smithisrealdev/travel-ai-agent/issues)
- **Discussions**: Join [GitHub Discussions](https://github.com/smithisrealdev/travel-ai-agent/discussions)
- **Email**: support@travel-ai-agent.com

---

**Built with ❤️ using Go, Nuxt 3, PostgreSQL, Redis, and AI**