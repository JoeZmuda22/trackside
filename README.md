# Trackside

A mobile-first web app for motorsports drivers to review tracks, log laps, and manage their garage.

## What It Does

- **Track Browser** — Search and filter tracks by location and event type (drift, drag, grip)
- **Track Reviews** — Rate tracks with wet/dry condition tags
- **Interactive Zones** — Tap zones on track maps to view or add driving tips
- **Garage** — Manage cars and modification lists
- **Lap Book** — Log lap times with telemetry (tire pressures, fuel, alignment)
- **Track Uploads** — Contribute new tracks with layout images

## Tech Stack

| Layer | Tech |
|-------|------|
| Frontend | Next.js 15, React 19, TypeScript, Tailwind CSS |
| Backend | Go, Echo v4, raw SQL |
| Database | SQLite (WAL mode) |
| Auth | JWT (HS256, 24h tokens) |

## Project Structure

```
trackside-app/
├── backend/          # Go API server
│   ├── cmd/server/   # Entry point
│   ├── internal/     # Config, handlers, middleware, models, repository, router
│   └── tests/        # 77 integration tests
├── trackside/        # Next.js frontend
│   └── src/
│       ├── app/      # Pages (tracks, garage, lapbook, profile, login, register)
│       ├── components/
│       └── lib/      # API client, auth, utils
└── Makefile          # Start/stop everything
```

## Getting Started

### Prerequisites

- Go 1.25+
- Node.js 18+
- npm

### Run Locally

```bash
# Clone the repo
git clone https://github.com/JoeZmuda22/trackside.git
cd trackside

# Start both backend and frontend
make up
```

This builds the Go backend, seeds demo data, and starts the Next.js dev server.

- **Frontend:** http://localhost:3000
- **Backend API:** http://localhost:8080

### Stop

```bash
make down
```

### Demo Login

- **Email:** `demo@trackside.com`
- **Password:** `password123`

## Environment Variables

### Backend (`backend/.env`)

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `DATABASE_URL` | `./trackside.db` | SQLite database path |
| `JWT_SECRET` | `change-me-in-production` | JWT signing key |
| `UPLOAD_DIR` | `./uploads` | Image upload directory |
| `CORS_ORIGINS` | `http://localhost:3000` | Allowed CORS origins |
| `DATA_DIR` | `../trackside/data` | Track data directory |

### Frontend (`trackside/.env.local`)

| Variable | Default | Description |
|----------|---------|-------------|
| `NEXT_PUBLIC_API_URL` | `http://localhost:8080` | Backend API URL |

## API Overview

See [backend/README.md](backend/README.md) for the full endpoint reference.

**Public:** Register, login, browse tracks and details.

**Protected (JWT required):** Create/edit tracks, reviews, zones, tips, cars, mods, lap records, profile, image uploads.

## Testing

```bash
cd backend
make test          # Run all 77 integration tests
make test-cover    # Run with coverage report
```
