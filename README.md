# Travel AI Agent 🌍✈️

An intelligent, full-stack travel planning application powered by AI. This monorepo includes a Go backend with Fiber framework and a Nuxt 3 frontend with TypeScript.

## Features

### 🤖 AI-Powered Recommendations
- Personalized travel suggestions using OpenAI GPT-4
- Context-aware itinerary generation
- Budget-conscious recommendations

### 🌤️ Real-Time Weather Integration
- Current weather conditions for destinations
- 5-day weather forecasts
- Climate-based travel suggestions

### ✈️ Flight Information
- Real-time flight data integration
- Flight search and status tracking
- Price estimates and duration

### 💾 Data Management
- PostgreSQL for persistent storage
- Redis caching for improved performance
- Travel search history

### 🎨 Modern Frontend
- Responsive design with Tailwind CSS
- TypeScript for type safety
- Pinia state management
- Component-based architecture

## Project Structure

```
travel-ai-agent/
├── backend/                    # Go + Fiber backend
│   ├── cmd/
│   │   └── server/
│   │       └── main.go        # Application entry point
│   ├── internal/
│   │   ├── config/            # Configuration management
│   │   ├── database/          # Database connections (PostgreSQL, Redis)
│   │   ├── handlers/          # HTTP request handlers
│   │   ├── models/            # Data models
│   │   └── services/          # Business logic (OpenAI, Weather, Flight)
│   ├── Dockerfile
│   └── go.mod
│
├── frontend/                   # Nuxt 3 + TypeScript frontend
│   ├── assets/
│   │   └── css/
│   │       └── main.css       # Tailwind styles
│   ├── components/
│   │   └── TravelSearch.vue   # Main search component
│   ├── composables/
│   │   └── useApi.ts          # API integration
│   ├── pages/
│   │   └── index.vue          # Landing page
│   ├── stores/
│   │   └── travel.ts          # Pinia store
│   ├── types/
│   │   └── index.ts           # TypeScript definitions
│   ├── Dockerfile
│   ├── nuxt.config.ts
│   ├── package.json
│   └── tsconfig.json
│
├── docker-compose.yml          # Complete orchestration
├── .env.example               # Environment variables template
├── .gitignore
└── README.md
```

## Prerequisites

- Docker and Docker Compose
- Node.js 20+ (for local development)
- Go 1.21+ (for local development)
- API Keys:
  - OpenAI API key
  - OpenWeatherMap API key
  - AviationStack API key (optional)

## Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/smithisrealdev/travel-ai-agent.git
cd travel-ai-agent
```

### 2. Configure Environment Variables

```bash
cp .env.example .env
```

Edit `.env` and add your API keys:

```env
OPENAI_API_KEY=your_openai_api_key_here
WEATHER_API_KEY=your_weather_api_key_here
FLIGHT_API_KEY=your_flight_api_key_here
POSTGRES_PASSWORD=your_secure_password_here
JWT_SECRET=your_jwt_secret_key_here
```

### 3. Start with Docker Compose

```bash
docker-compose up -d
```

This will start:
- PostgreSQL database on port 5432
- Redis cache on port 6379
- Backend API on port 8080
- Frontend UI on port 3000

### 4. Access the Application

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **Health Check**: http://localhost:8080/health

## API Endpoints

### Travel Search

**POST** `/api/v1/travel/search`

Search for travel recommendations.

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

### Search History

**GET** `/api/v1/travel/history?userId=user123`

Retrieve user's search history.

### Health Check

**GET** `/health`

Check service status.

## Development

### Backend Development

```bash
cd backend

# Install dependencies
go mod download

# Run locally
go run cmd/server/main.go
```

### Frontend Development

```bash
cd frontend

# Install dependencies
npm install

# Run development server
npm run dev
```

The frontend will be available at http://localhost:3000 with hot reload.

### Building for Production

```bash
# Backend
cd backend
go build -o main cmd/server/main.go

# Frontend
cd frontend
npm run build
```

## Environment Variables

### Backend

| Variable | Description | Default |
|----------|-------------|---------|
| `BACKEND_PORT` | Server port | 8080 |
| `POSTGRES_HOST` | PostgreSQL host | localhost |
| `POSTGRES_PORT` | PostgreSQL port | 5432 |
| `POSTGRES_USER` | Database user | travelagent |
| `POSTGRES_PASSWORD` | Database password | - |
| `POSTGRES_DB` | Database name | travelagent |
| `REDIS_HOST` | Redis host | localhost |
| `REDIS_PORT` | Redis port | 6379 |
| `OPENAI_API_KEY` | OpenAI API key | - |
| `WEATHER_API_KEY` | Weather API key | - |
| `FLIGHT_API_KEY` | Flight API key | - |

### Frontend

| Variable | Description | Default |
|----------|-------------|---------|
| `NUXT_PUBLIC_API_BASE` | Backend API URL | http://localhost:8080 |
| `NODE_ENV` | Environment | development |

## Technology Stack

### Backend
- **Framework**: Go Fiber v2
- **Database**: PostgreSQL 15
- **Cache**: Redis 7
- **AI**: OpenAI GPT-4
- **APIs**: OpenWeatherMap, AviationStack

### Frontend
- **Framework**: Nuxt 3
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **State**: Pinia
- **HTTP**: Fetch API

### Infrastructure
- **Containerization**: Docker
- **Orchestration**: Docker Compose
- **Database Migrations**: SQL scripts

## Database Schema

### travel_searches

Stores user travel searches.

| Column | Type | Description |
|--------|------|-------------|
| id | SERIAL | Primary key |
| user_id | VARCHAR(255) | User identifier |
| destination | VARCHAR(255) | Travel destination |
| start_date | DATE | Trip start date |
| end_date | DATE | Trip end date |
| budget | DECIMAL(10,2) | Trip budget |
| preferences | JSONB | User preferences |
| results | JSONB | Search results |
| created_at | TIMESTAMP | Creation time |
| updated_at | TIMESTAMP | Update time |

### travel_recommendations

Stores individual recommendations.

| Column | Type | Description |
|--------|------|-------------|
| id | SERIAL | Primary key |
| search_id | INTEGER | Foreign key to travel_searches |
| recommendation_type | VARCHAR(50) | Type of recommendation |
| title | VARCHAR(255) | Recommendation title |
| description | TEXT | Detailed description |
| price | DECIMAL(10,2) | Price estimate |
| rating | DECIMAL(3,2) | User rating |
| metadata | JSONB | Additional data |
| created_at | TIMESTAMP | Creation time |

## Troubleshooting

### Database Connection Issues

```bash
# Check PostgreSQL logs
docker-compose logs postgres

# Verify connection
docker-compose exec postgres psql -U travelagent -d travelagent
```

### Redis Connection Issues

```bash
# Check Redis logs
docker-compose logs redis

# Test connection
docker-compose exec redis redis-cli ping
```

### Backend Issues

```bash
# Check backend logs
docker-compose logs backend

# Restart backend
docker-compose restart backend
```

### Frontend Issues

```bash
# Check frontend logs
docker-compose logs frontend

# Rebuild frontend
docker-compose up -d --build frontend
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.

## Acknowledgments

- OpenAI for GPT-4 API
- OpenWeatherMap for weather data
- AviationStack for flight information
- Go Fiber community
- Nuxt.js community

## Support

For issues and questions:
- Open an issue on GitHub
- Contact: support@travel-ai-agent.com

---

**Built with ❤️ using Go, Nuxt 3, PostgreSQL, and Redis**