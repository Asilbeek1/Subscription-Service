# Subscription Service

REST API for managing user subscriptions, built with Go.

## Stack

- **Go** — Chi router, slog logger
- **PostgreSQL** — data storage
- **Docker** — containerization
- **Swagger** — API docs

## Project Structure
cmd/api          → entrypoint
internal/
config         → config loading
database       → postgres + repository
domain         → interfaces & models
service        → business logic
transport/http → handlers, DTOs, router
migrations/      → SQL migrations
pkg/ctxutil      → request ID helper

## Getting Started

**1. Clone**
```bash
git clone https://github.com/Asilbeek1/Subscription-Service
cd Subscription-Service
```

**2. Configure**
```bash
cp .env.example .env  # edit DB credentials inside
```

**3. Run**
```bash
docker compose up --build
```

## API

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/subscriptions` | Create subscription |
| GET | `/subscriptions` | List all |
| GET | `/subscriptions/{id}` | Get by ID |
| DELETE | `/subscriptions/{id}` | Delete |
| GET | `/subscriptions/total` | Calculate total cost |

### Total query params
| Param | Type | Required |
|-------|------|----------|
| `user_id` | uuid | No |
| `service_name` | string | No |
| `from` | MM-YYYY | No |
| `to` | MM-YYYY | No |

## Swagger
http://localhost:8080/swagger/index.html

## Migrations

```bash
make migrate-up
make migrate-down
```
