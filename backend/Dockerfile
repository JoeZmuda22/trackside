# Build stage
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO is required for go-sqlite3
RUN CGO_ENABLED=1 GOOS=linux go build -o server ./cmd/server/

# Runtime stage
FROM alpine:3.20

RUN apk add --no-cache ca-certificates sqlite-libs

WORKDIR /app

COPY --from=builder /app/server .

RUN mkdir -p uploads data

EXPOSE 8080

ENTRYPOINT ["./server"]
