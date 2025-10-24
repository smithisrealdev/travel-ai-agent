# Travel AI Agent Setup Guide

## Overview

This document provides detailed setup instructions for the Travel AI Agent monorepo.

## Prerequisites

Before you begin, ensure you have the following installed:

- **Docker** (version 20.10 or higher)
- **Docker Compose** (version 2.0 or higher)
- **Node.js** (version 20.x or higher) - for local frontend development
- **Go** (version 1.21 or higher) - for local backend development

### API Keys Required

You'll need to obtain API keys from the following services:

1. **OpenAI** - https://platform.openai.com/api-keys
   - Used for AI-powered travel recommendations
   - Required for full functionality

2. **OpenWeatherMap** - https://openweathermap.org/api
   - Used for weather information
   - Free tier available (60 calls/minute)

3. **AviationStack** (Optional) - https://aviationstack.com/
   - Used for flight information
   - Free tier available (100 calls/month)

## Quick Start

### 1. Environment Setup

```bash
# Clone the repository
git clone https://github.com/smithisrealdev/travel-ai-agent.git
cd travel-ai-agent

# Copy the environment template
cp .env.example .env

# Edit .env and add your API keys
nano .env  # or use your preferred editor
```

### 2. Configure API Keys

Edit the `.env` file and replace the placeholder values:

```env
# Required
OPENAI_API_KEY=sk-your-actual-openai-key-here
WEATHER_API_KEY=your-weather-api-key-here

# Optional
FLIGHT_API_KEY=your-flight-api-key-here

# Security (change in production)
POSTGRES_PASSWORD=your-secure-database-password
JWT_SECRET=your-secure-jwt-secret
```

### 3. Start All Services

```bash
# Start everything with Docker Compose
docker compose up -d

# Check service status
docker compose ps

# View logs
docker compose logs -f
```

### 4. Access the Application

- **Frontend UI**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **Health Check**: http://localhost:8080/health

## Development Workflow

### Backend Development

```bash
cd backend

# Install Go dependencies
go mod download

# Run locally (requires PostgreSQL and Redis running)
go run cmd/server/main.go

# Build
go build -o main cmd/server/main.go

# Run tests
go test ./...
```

### Frontend Development

```bash
cd frontend

# Install dependencies
npm install

# Run development server with hot reload
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

## Project Structure

```
travel-ai-agent/
├── backend/                      # Go backend service
│   ├── cmd/server/              # Application entry point
│   ├── internal/
│   │   ├── config/              # Configuration management
│   │   ├── database/            # Database connections
│   │   ├── handlers/            # HTTP handlers
│   │   ├── models/              # Data models
│   │   └── services/            # Business logic
│   ├── Dockerfile               # Backend container
│   └── go.mod                   # Go dependencies
│
├── frontend/                     # Nuxt 3 frontend
│   ├── assets/                  # CSS and static assets
│   ├── components/              # Vue components
│   ├── composables/             # Vue composables
│   ├── pages/                   # Page components
│   ├── stores/                  # Pinia stores
│   ├── types/                   # TypeScript types
│   ├── Dockerfile               # Frontend container
│   └── package.json             # Node dependencies
│
├── docker-compose.yml           # Service orchestration
├── .env.example                 # Environment template
└── README.md                    # Project documentation
```

## Service Architecture

### Services

1. **PostgreSQL** (port 5432)
   - Stores travel searches and recommendations
   - Persistent volume for data

2. **Redis** (port 6379)
   - Caches API responses
   - Session storage

3. **Backend** (port 8080)
   - REST API
   - Integrates with OpenAI, Weather, and Flight APIs
   - Manages database operations

4. **Frontend** (port 3000)
   - Nuxt 3 application
   - Server-side rendering
   - Responsive UI

### Data Flow

```
User Browser
    ↓
Frontend (Nuxt 3)
    ↓
Backend API (Go Fiber)
    ↓
├── PostgreSQL (persistence)
├── Redis (caching)
└── External APIs
    ├── OpenAI (recommendations)
    ├── Weather API (weather data)
    └── Flight API (flight info)
```

## API Documentation

### Search Travel

```bash
POST /api/v1/travel/search
Content-Type: application/json

{
  "destination": "Paris",
  "startDate": "2024-06-01",
  "endDate": "2024-06-07",
  "budget": 2000,
  "preferences": {
    "culture": true,
    "adventure": false,
    "relaxation": true,
    "food": true
  }
}
```

### Get Search History

```bash
GET /api/v1/travel/history?userId=user123
```

### Health Check

```bash
GET /health
```

## Docker Commands

### Start Services

```bash
# Start all services
docker compose up -d

# Start specific service
docker compose up -d backend

# Build and start
docker compose up --build -d
```

### View Logs

```bash
# All services
docker compose logs -f

# Specific service
docker compose logs -f backend

# Last 100 lines
docker compose logs --tail 100 backend
```

### Stop Services

```bash
# Stop all services
docker compose stop

# Stop specific service
docker compose stop backend

# Stop and remove
docker compose down

# Stop and remove volumes
docker compose down -v
```

### Database Operations

```bash
# Connect to PostgreSQL
docker compose exec postgres psql -U travelagent -d travelagent

# Connect to Redis
docker compose exec redis redis-cli
```

## Troubleshooting

### Backend Won't Start

1. Check PostgreSQL is running:
   ```bash
   docker compose ps postgres
   docker compose logs postgres
   ```

2. Verify environment variables:
   ```bash
   docker compose exec backend env | grep POSTGRES
   ```

3. Check backend logs:
   ```bash
   docker compose logs backend
   ```

### Frontend Won't Start

1. Check if backend is accessible:
   ```bash
   curl http://localhost:8080/health
   ```

2. Verify frontend environment:
   ```bash
   docker compose exec frontend env | grep NUXT
   ```

3. Check frontend logs:
   ```bash
   docker compose logs frontend
   ```

### Database Connection Issues

1. Ensure PostgreSQL is healthy:
   ```bash
   docker compose exec postgres pg_isready -U travelagent
   ```

2. Check network connectivity:
   ```bash
   docker compose exec backend ping postgres
   ```

3. Verify database exists:
   ```bash
   docker compose exec postgres psql -U travelagent -l
   ```

### Redis Connection Issues

1. Test Redis connection:
   ```bash
   docker compose exec redis redis-cli ping
   ```

2. Check Redis logs:
   ```bash
   docker compose logs redis
   ```

## Security Considerations

### Production Deployment

1. **Change Default Passwords**
   - Update `POSTGRES_PASSWORD`
   - Update `JWT_SECRET`

2. **Enable HTTPS**
   - Use a reverse proxy (nginx/traefik)
   - Obtain SSL certificates

3. **Secure API Keys**
   - Never commit `.env` to version control
   - Use secrets management (e.g., Vault, AWS Secrets Manager)

4. **Database Security**
   - Limit network exposure
   - Use strong passwords
   - Enable SSL connections

5. **Rate Limiting**
   - Implement API rate limiting
   - Protect against DDoS

### Environment Variables

Never commit these to version control:
- API keys
- Database passwords
- JWT secrets
- Any sensitive configuration

## Performance Optimization

### Backend

- Connection pooling is pre-configured
- Redis caching reduces API calls
- Health checks ensure service availability

### Frontend

- Server-side rendering for better SEO
- Static asset optimization
- Lazy loading of components

### Database

- Indexes on frequently queried fields
- JSONB for flexible data storage
- Regular vacuum and analyze

## Monitoring

### Health Checks

```bash
# Application health
curl http://localhost:8080/health

# Service status
docker compose ps
```

### Logs

```bash
# Real-time logs
docker compose logs -f

# Service-specific logs
docker compose logs -f backend frontend
```

### Resource Usage

```bash
# Container stats
docker stats

# Disk usage
docker system df
```

## Backup and Recovery

### Database Backup

```bash
# Backup PostgreSQL
docker compose exec postgres pg_dump -U travelagent travelagent > backup.sql

# Restore
docker compose exec -T postgres psql -U travelagent travelagent < backup.sql
```

### Redis Backup

```bash
# Trigger save
docker compose exec redis redis-cli SAVE

# Copy RDB file
docker compose cp redis:/data/dump.rdb ./redis-backup.rdb
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## Support

For issues and questions:
- GitHub Issues: https://github.com/smithisrealdev/travel-ai-agent/issues
- Email: support@travel-ai-agent.com

## License

MIT License - see LICENSE file for details
