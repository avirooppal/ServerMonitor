# 🐳 Docker Deployment Guide

This guide details how to build, publish, and run the Server Monitor components using Docker.

## 🏗️ Architecture

The system consists of two main Docker images:

1.  **Dashboard (`server-moni`)**: The central web interface and API server.
2.  **Agent (`server-moni-agent`)**: The lightweight collector running on target servers.

## 📦 Building Images

### 1. Build the Dashboard
The dashboard image includes both the Go backend and the React frontend (embedded).

```bash
# Build locally
docker build -t avirooppal/server-moni:latest .
```

### 2. Build the Agent
The agent image is a minimal Alpine-based image for metrics collection.

```bash
# Build locally
docker build -t avirooppal/linux-monitoring-agent:latest -f Dockerfile.agent .
```

---

## ☁️ Publishing to Docker Hub

To make the **One-Click Installer** work for everyone, the Agent image must be hosted on a public registry (like Docker Hub).

```bash
# 1. Login to Docker Hub
docker login

# 2. Tag and Push the Agent
docker tag avirooppal/linux-monitoring-agent:latest avirooppal/linux-monitoring-agent:latest
docker push avirooppal/linux-monitoring-agent:latest

# 3. (Optional) Push the Dashboard
docker tag avirooppal/server-moni:latest avirooppal/server-moni:latest
docker push avirooppal/server-moni:latest
```

---

## 🚀 Running Manually

If you prefer not to use the helper scripts, you can run containers manually.

### Running the Dashboard
```bash
docker run -d \
  --name dashboard \
  -p 8082:8080 \
  -v $(pwd)/data:/app/data \
  avirooppal/server-moni:latest
```

### Running the Agent
```bash
docker run -d \
  --name agent \
  --network host \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v agent_data:/app/data \
  -e API_PORT=8080 \
  avirooppal/linux-monitoring-agent:latest
```
*Note: `--network host` is recommended for the agent to accurately report host network metrics.*

---

## 🐙 Docker Compose (Local Dev)

For local development, use the provided `docker-compose.yml`:

```bash
docker-compose up -d --build
```

This starts:
- **Dashboard** on `http://localhost:8082`
- **Agent** on `http://localhost:8081`
