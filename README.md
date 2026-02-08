# Trackside - Motorsports Track Review & Driver App

A mobile-first web application for motorsports drivers to review tracks, manage vehicles, log lap data with telemetry, and browse tracks by event type.

## Features

- **Track Browser** — Search and filter tracks by location and event type (drift, drag, grip)
- **Track Upload** — Any driver can contribute tracks with layout images
- **Interactive Zone Map** — Tap zones on track images to view/add tips and comments
- **Track Reviews** — Rate tracks with wet/dry condition tags per event type
- **Driver Profile** — Manage experience level, multiple cars, and modification lists
- **My Lap Book** — Log lap times with full telemetry: tire pressures, fuel, alignment (camber, caster, toe)
- **Event Types** — Tracks support multiple event offerings: drift, drag, and grip racing

## Tech Stack

- **Framework:** Next.js 15 (App Router)
- **Language:** TypeScript
- **Styling:** Tailwind CSS (mobile-first)
- **Database:** PostgreSQL
- **ORM:** Prisma
- **Auth:** NextAuth.js v5
- **Validation:** Zod

## Getting Started

```bash
# Install dependencies
npm install

# Set up environment variables
cp .env.example .env
# Edit .env with your DATABASE_URL and NEXTAUTH_SECRET

# Generate Prisma client and push schema
npm run db:generate
npm run db:push

# Seed the database (optional)
npm run db:seed

# Start development server
npm run dev
```

Open [http://localhost:3000](http://localhost:3000) in your browser.

## Environment Variables

```
DATABASE_URL="postgresql://user:password@localhost:5432/trackside"
NEXTAUTH_SECRET="your-secret-here"
NEXTAUTH_URL="http://localhost:3000"
```
