# URL Monitor

A full-stack URL monitoring service built with Go and React. Register URLs, and the system continuously checks their availability every minute — tracking uptime, response times, and HTTP status codes with a live dashboard and Prometheus metrics.

---

## Features

- **Continuous Monitoring** — Background engine checks all registered URLs every 60 seconds using concurrent goroutines
- **Live Dashboard** — React frontend showing URL status, response times, and historical logs
- **Auth** — JWT-based authentication with short-lived access tokens and long-lived refresh tokens
- **Observability** — Prometheus metrics + Grafana dashboards out of the box
- **URL Validation** — Validates and pings a URL before registering it; rejects unreachable or localhost URLs
- **Per-user isolation** — Users can only see and manage their own URLs

---

## Tech Stack

| Layer | Technology |
|---|---|
| Backend | Go 1.24, Chi v5 |
| Database | PostgreSQL, pgx |
| Frontend | React 19, TypeScript, Vite |
| State Management | TanStack React Query |
| Styling | Tailwind CSS |
| Auth | JWT (access + refresh), bcrypt |
| Monitoring | Prometheus, Grafana, Node Exporter |
| Deployment | Docker, Docker Compose |

---

## Architecture

```
┌─────────────────┐         ┌──────────────────────────────────┐
│  React Frontend │──HTTP──▶│            Go API                │
└─────────────────┘         │                                  │
                            │  Handler → Service → Repository  │
                            │                │                  │
                            │         PostgreSQL               │
                            │                                  │
                            │  ┌─────────────────────────┐    │
                            │  │   Monitoring Engine      │    │
                            │  │  (background goroutine)  │    │
                            │  │  - 60s ticker            │    │
                            │  │  - concurrent HTTP checks│    │
                            │  │  - Prometheus metrics    │    │
                            │  └─────────────────────────┘    │
                            └──────────────────────────────────┘
```

The monitoring engine runs as a background goroutine inside the main process. Each tick fetches all active URLs and checks them concurrently using `sync.WaitGroup`, with `sync.Mutex` protecting shared metric state. Results are persisted to PostgreSQL and exported to Prometheus.

---

## Getting Started

### Prerequisites

- Docker and Docker Compose

### Run with Docker

```bash
git clone https://github.com/your-username/url-monitoring-api
cd url-monitoring-api
cp .env.example .env   # fill in values
docker compose up --build
```

| Service | URL |
|---|---|
| Frontend | http://localhost:80 |
| API | http://localhost:3000 |
| Prometheus | http://localhost:9090 |
| Grafana | http://localhost:3001 |

### Run locally (backend only)

```bash
# Start PostgreSQL
docker compose up postgres -d

# Run migrations
psql $DATABASE_URL -f migrations/001_init.sql

# Start API with live reload
air
```

---

## Environment Variables

```env
DATABASE_URL=postgres://user:password@localhost:5432/urlmonitor
JWT_ACCESS_SECRET=your-access-secret
JWT_REFRESH_SECRET=your-refresh-secret
```

---

## API Reference

### Auth

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/auth/signup` | Register a new user |
| `POST` | `/auth/login` | Login and receive tokens |
| `POST` | `/auth/refresh` | Refresh access token |
| `POST` | `/auth/logout` | Clear auth cookies |
| `GET` | `/auth/me` | Get current user |

### URLs (requires auth)

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/urls/register` | Register a URL for monitoring |
| `GET` | `/urls` | List all monitored URLs |
| `DELETE` | `/urls/{id}` | Delete a monitored URL |
| `GET` | `/urls/{id}/logs` | Get check logs for a URL |

### System

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/` | Health check |
| `GET` | `/metrics` | Prometheus metrics |

---

## Prometheus Metrics

| Metric | Labels | Description |
|---|---|---|
| `url_checks_total` | `user_id` | Total checks per user per cycle |
| `url_checks_failed_total` | `user_id` | Failed checks per user per cycle |
| `url_response_time_ms` | `user_id`, `url_id`, `host` | Response time per URL |
| `url_last_status_code` | `user_id`, `url_id`, `host` | Last HTTP status per URL |

---

## Project Structure

```
.
├── cmd/
│   └── main.go              # Entry point, dependency wiring
├── internal/
│   ├── engine/
│   │   └── monitor.go       # Background monitoring engine
│   ├── handler/             # HTTP handlers
│   ├── service/             # Business logic
│   ├── repository/          # Database access (interface + pgx impl)
│   ├── middleware/          # JWT middleware
│   ├── metrics/             # Prometheus metric definitions
│   ├── model/               # Domain models
│   └── router/              # Route registration
├── migrations/              # SQL migration files
├── frontend/                # React + TypeScript SPA
│   └── src/
│       ├── pages/
│       ├── components/
│       ├── hooks/           # TanStack Query hooks
│       └── lib/
├── docker-compose.yml       # Dev stack (includes Prometheus + Grafana)
└── docker-compose.prod.yml  # Production stack
```

---

## Database Schema

```sql
users           -- id, email, password, created_at
monitored_urls  -- id (UUID), url, user_id, is_active, created_at
url_logs        -- id (UUID), url_id, status_code, response_time_ms, checked_at, is_up
```

Foreign keys cascade on delete — removing a URL cleans up all its logs.
