# ChronoGo

Automatic message sending system that sends messages every configured interval from unsent records in the database.

## Features

- 🚀 Automatic message sending with configurable interval and message limit
- 📊 PostgreSQL database for message storage
- 🔄 Redis caching support
- 📝 Swagger API documentation
- 🐳 Docker Compose support with all services containerized
- ⚡ Built with Go 1.25

## Prerequisites

- Docker and Docker Compose (recommended)
- Go 1.25+ (for local development without Docker)

## Quick Start with Docker Compose (Recommended)

### 1. Clone the repository

```bash
git clone <repository-url>
cd ChronoGo
```

### 2. Configure environment variables

Create a `.env` file in the project root:

```env
# Database Configuration
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=chronogo
DB_SSLMODE=disable

# Server Configuration
SERVER_PORT=8080

# Webhook Configuration (Required)
WEBHOOK_URL=https://webhook.site/your-webhook-url
WEBHOOK_AUTH_KEY=your-auth-key-here

# Scheduler Configuration
SCHEDULER_INTERVAL_MINUTES=2
SCHEDULER_MESSAGE_LIMIT=2

# Redis Configuration
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

**Important:** 
- For Docker Compose, `DB_HOST=postgres` and `REDIS_HOST=redis` (container names)
- For local development, use `DB_HOST=localhost` and `REDIS_HOST=localhost`

### 3. Start all services

```bash
docker-compose up --build
```

This will:
- Start PostgreSQL on port 5432
- Start Redis on port 6379
- Build and start the Go application on port 8080
- Run database migrations automatically
- Seed 10 test messages into the database (01-ADANA, 02-ADIYAMAN, etc.)

### 4. Access the application

- **API Server:** `http://localhost:8080`
- **Health Check:** `http://localhost:8080/health`
- **Swagger Documentation:** `http://localhost:8080/swagger/index.html`

### 5. Run in background

```bash
docker-compose up -d --build
```

## Local Setup (without Docker)

### 1. Setup PostgreSQL and Redis

```bash
# PostgreSQL
createdb -U postgres chronogo
psql -U postgres -d chronogo -f migrations/001_create_messages.sql

# Redis (using Homebrew on macOS)
brew install redis
brew services start redis
```

### 2. Configure environment variables

Create a `.env` file with local settings:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=chronogo
DB_SSLMODE=disable

SERVER_PORT=8080

WEBHOOK_URL=https://webhook.site/your-webhook-url
WEBHOOK_AUTH_KEY=your-auth-key-here

SCHEDULER_INTERVAL_MINUTES=2
SCHEDULER_MESSAGE_LIMIT=2

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

### 3. Run seed script

```bash
psql -U postgres -d chronogo -f scripts/seed.sql
```

### 4. Run the application

```bash
go run cmd/server/main.go
```

## Docker Commands

### Start all services
```bash
docker-compose up -d
```

### Stop all services
```bash
docker-compose down
```

### View logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f app
docker-compose logs -f postgres
docker-compose logs -f redis
```

### Restart services
```bash
docker-compose restart
```

### Rebuild and restart
```bash
docker-compose up -d --build
```

### Reset database (remove volumes)
```bash
docker-compose down -v
docker-compose up -d
```

### Stop specific service
```bash
docker-compose stop app
```

## API Endpoints

### Health Check
```
GET /health
```

**Response:**
```json
{
  "status": "ok",
  "message": "Database connection is healthy"
}
```

### List Sent Messages
```
GET /api/messages/sent
```

**Response:**
```json
{
  "messages": [
    {
      "id": 1,
      "to": "+905551111111",
      "content": "01-ADANA",
      "status": "sent",
      "sent_at": "2025-11-02T21:38:05Z",
      "message_id": "uuid-here"
    }
  ],
  "total": 1
}
```

### Toggle Scheduler
```
POST /api/scheduler/toggle
```

**Response:**
```json
{
  "status": "running",
  "message": "Scheduler is now running"
}
```
or
```json
{
  "status": "stopped",
  "message": "Scheduler is now stopped"
}
```

## Project Structure

```
ChronoGo/
├── cmd/
│   └── server/          # Main application entry point
├── internal/
│   ├── database/        # Database connection and config
│   ├── handler/         # HTTP handlers (API endpoints)
│   ├── model/           # Data models
│   ├── queue/           # Scheduler implementation
│   ├── redis/           # Redis connection and caching
│   ├── repository/      # Database operations
│   └── sender/          # Webhook sender
├── migrations/          # Database migrations
├── scripts/             # Seed scripts
├── docs/                # Swagger documentation (auto-generated)
├── docker-compose.yml   # Docker services configuration
└── Dockerfile           # Application container
```

## Environment Variables

| Variable                     | Description                     | Default     | Required |
| ---------------------------- | ------------------------------- | ----------- | -------- |
| `DB_HOST`                    | PostgreSQL host                 | `localhost` | Yes      |
| `DB_PORT`                    | PostgreSQL port                 | `5432`      | Yes      |
| `DB_USER`                    | Database user                   | `postgres`  | Yes      |
| `DB_PASSWORD`                | Database password               | `postgres`  | Yes      |
| `DB_NAME`                    | Database name                   | `chronogo`  | Yes      |
| `DB_SSLMODE`                 | SSL mode                        | `disable`   | No       |
| `WEBHOOK_URL`                | Webhook endpoint URL            | -           | **Yes**  |
| `WEBHOOK_AUTH_KEY`           | Webhook authentication key      | -           | **Yes**  |
| `SCHEDULER_INTERVAL_MINUTES` | Scheduler interval in minutes   | `2`         | No       |
| `SCHEDULER_MESSAGE_LIMIT`    | Number of messages per interval | `2`         | No       |
| `REDIS_HOST`                 | Redis host                      | `localhost` | No       |
| `REDIS_PORT`                 | Redis port                      | `6379`      | No       |
| `REDIS_PASSWORD`             | Redis password                  | -           | No       |
| `REDIS_DB`                   | Redis database number           | `0`         | No       |

**Note:** For Docker Compose, use container names: `DB_HOST=postgres`, `REDIS_HOST=redis`

## How It Works

1. **Scheduler:** Runs automatically when the server starts
   - Fetches unsent messages from the database
   - Sends messages via webhook
   - Updates message status to 'sent'
   - Caches messageId and sent_at to Redis

2. **Message Flow:**
   - Messages are inserted with status 'unsent'
   - Scheduler picks up unsent messages every configured interval
   - Messages are sent to the webhook endpoint
   - Status is updated to 'sent' in the database
   - MessageId and sent_at are cached in Redis (TTL: 24 hours)

3. **Configuration:**
   - Adjust `SCHEDULER_INTERVAL_MINUTES` to change sending frequency
   - Adjust `SCHEDULER_MESSAGE_LIMIT` to change batch size

## Redis Cache

The system caches successfully sent messages in Redis:

- **Key format:** `message:{messageId}`
- **TTL:** 24 hours
- **Cached data:**
  - `message_id`: UUID from webhook response
  - `sent_at`: Timestamp of when message was sent

**Check Redis cache:**
```bash
docker exec chronogo-redis redis-cli KEYS "message:*"
docker exec chronogo-redis redis-cli GET "message:{messageId}"
```

## Development

### Generate Swagger documentation

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/server/main.go -o docs
```

### Build application

```bash
go build -o bin/server cmd/server/main.go
```

### Run tests

```bash
go test ./...
```

## Troubleshooting

### Database connection failed

- Check PostgreSQL is running: `docker ps | grep postgres`
- Verify credentials in `.env` file
- For Docker: Ensure `DB_HOST=postgres` (container name)

### Redis connection failed

- Redis is optional, application continues without cache
- Check Redis is running: `docker ps | grep redis`
- For Docker: Ensure `REDIS_HOST=redis` (container name)

### Webhook errors

- Verify `WEBHOOK_URL` and `WEBHOOK_AUTH_KEY` in `.env`
- Check webhook endpoint is accessible
- For webhook.site: System generates mock messageId if response is not JSON

### Port already in use

```bash
# Find process using port 8080
lsof -ti:8080

# Kill the process
lsof -ti:8080 | xargs kill
```

## License

See LICENSE file for details.
