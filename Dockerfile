# Build Frontend
FROM node:18-alpine AS frontend-builder
WORKDIR /app/web
COPY web/package*.json ./
RUN npm install
COPY web/ .
RUN npm run build

# Build Backend
FROM golang:1.24-alpine AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Copy frontend build to embed location
COPY --from=frontend-builder /app/web/dist ./cmd/server/dist
RUN go build -o server-moni ./cmd/server

# Final Image
FROM alpine:latest
WORKDIR /app
COPY --from=backend-builder /app/server-moni .
# Create data directory if needed, though sqlite file is created in CWD
EXPOSE 8080
CMD ["./server-moni"]
