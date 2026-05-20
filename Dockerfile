# ---- Stage 1: Build frontend ----
FROM node:20-alpine AS frontend

WORKDIR /app/web/frontend
COPY web/frontend/package*.json ./
RUN npm ci
COPY web/frontend/ .
RUN npm run build

# ---- Stage 2: Build Go binary ----
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy && go mod download

COPY . .
COPY --from=frontend /app/web/dist ./web/dist

RUN CGO_ENABLED=1 go build -ldflags="-s -w" -o docker-panel .

# ---- Stage 3: Runtime ----
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=builder /app/docker-panel .

EXPOSE 8080

ENTRYPOINT ["./docker-panel"]
