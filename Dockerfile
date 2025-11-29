# ==========================================
# Stage 1: Build Frontend (React/Vite)
# ==========================================
FROM node:18-alpine AS frontend-builder
WORKDIR /app/web
COPY web/package*.json ./
RUN npm install
COPY web/ .
RUN npm run build

# ==========================================
# Stage 2: Build Backend (Go)
# ==========================================
FROM golang:1.24-alpine AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Embed the frontend build into the Go binary
COPY --from=frontend-builder /app/web/dist ./cmd/server/dist
RUN go build -o server-moni ./cmd/server

# ==========================================
# Stage 3: Final Production Image
# ==========================================
FROM alpine:latest
WORKDIR /app
COPY --from=backend-builder /app/server-moni .
# Expose the application port
EXPOSE 8080
CMD ["./server-moni"]
