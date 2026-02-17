.PHONY: up down

# Start Go backend + Next.js frontend
up:
	@echo "Building Go backend..."
	@cd backend && CGO_ENABLED=1 go build -o bin/server ./cmd/server/
	@echo "Starting Trackside (backend :8080 | frontend :3000)..."
	@cd backend && ./bin/server --seed & \
	cd trackside && npm run dev & \
	wait

# Stop all background processes
down:
	@pkill -f "bin/server" 2>/dev/null || true
	@pkill -f "next dev" 2>/dev/null || true
	@echo "Stopped."
