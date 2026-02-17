# Trackside Backend

Go backend for Trackside — a motorsport track day companion app.

## Stack

- **Go** + **Echo v4** web framework
- **SQLite** with WAL mode (via `mattn/go-sqlite3`)
- **JWT** authentication (24h tokens, HS256)
- Raw SQL queries — no ORM

## Quick Start

```bash
# Copy env file and set JWT_SECRET
cp .env.example .env

# Install dependencies
go mod tidy

# Run with demo seed data
make dev-seed

# Or build and run
make seed
```

The server starts on `http://localhost:8080`.

## API Endpoints

### Public
| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/register` | Create account |
| POST | `/api/auth/login` | Login, get JWT |
| GET | `/api/tracks` | List tracks (search, eventType, state filters) |
| GET | `/api/tracks/:id` | Track detail (zones, reviews, events) |
| GET | `/api/tracks/:id/images` | Track images |

### Protected (Bearer token required)
| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/cars` | List user's cars |
| POST | `/api/cars` | Add car |
| PUT | `/api/cars/:id` | Update car |
| DELETE | `/api/cars/:id` | Delete car |
| POST | `/api/cars/:id/mods` | Add car mod |
| DELETE | `/api/cars/:id/mods/:modId` | Remove mod |
| POST | `/api/tracks` | Create track |
| PATCH | `/api/tracks/:id` | Update track (owner only) |
| POST | `/api/tracks/:id/images` | Add image |
| DELETE | `/api/tracks/:id/images?imageId=` | Remove image |
| POST | `/api/tracks/:id/reviews` | Add review |
| POST | `/api/tracks/:id/zones` | Add zone |
| PATCH | `/api/tracks/:id/zones/:zoneId` | Update zone |
| DELETE | `/api/tracks/:id/zones/:zoneId` | Delete zone |
| POST | `/api/tracks/:id/zones/:zoneId/tips` | Add zone tip |
| GET | `/api/lapbook` | List lap records |
| POST | `/api/lapbook` | Add lap record |
| DELETE | `/api/lapbook/:id` | Delete lap record |
| GET | `/api/profile` | Get profile |
| PUT | `/api/profile` | Update profile |
| POST | `/api/upload` | Upload image file |
| POST | `/api/admin/sync-tracks` | Import tracks from data file |

## Testing

```bash
make test          # Run all tests
make test-cover    # Run with coverage report
```

77 integration tests covering all endpoints — auth, validation, ownership, and error cases.

## Project Structure

```
backend/
├── cmd/server/main.go          # Entry point
├── internal/
│   ├── config/                 # Environment config
│   ├── database/               # SQLite connection, migrations, seed
│   ├── handlers/               # HTTP handlers (12 files)
│   ├── middleware/              # JWT auth, CORS
│   ├── models/                 # Structs, enums, DTOs
│   ├── repository/             # Data access layer (raw SQL)
│   └── router/                 # Route registration
├── tests/                      # Integration tests
├── Dockerfile                  # Multi-stage Docker build
└── Makefile                    # Build/run/test commands
```

## Docker

```bash
make docker-build
make docker-run
```

## Demo Credentials

When running with `--seed`:
- Email: `demo@trackside.com`
- Password: `password123`
