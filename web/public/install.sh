#!/bin/bash

# Server Moni Agent Installer

echo "=================================================="
echo "Installing Server Moni Agent..."
echo "=================================================="

# Check for Docker
if ! command -v docker &> /dev/null; then
    echo "Error: Docker is not installed."
    echo "Please install Docker to run the agent."
    exit 1
fi

# Pull Image
IMAGE_NAME="avirooppal/linux-monitoring-agent:latest"
echo "Pulling latest agent image: $IMAGE_NAME..."
docker pull $IMAGE_NAME

# Remove existing container if any
if [ "$(docker ps -aq -f name=server-moni-agent)" ]; then
    echo "Removing existing agent container..."
    docker rm -f server-moni-agent
fi

# Run Agent
echo "Starting Agent container..."
docker run -d \
    --name server-moni-agent \
    --restart unless-stopped \
    --network host \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v server-moni-data:/app/data \
    -e API_PORT=8080 \
    $IMAGE_NAME

# Wait for startup and key generation
echo "Waiting for agent to initialize..."
sleep 5

# Retrieve API Key
API_KEY=$(docker exec server-moni-agent cat data/api_key.txt 2>/dev/null)

if [ -z "$API_KEY" ]; then
    echo "Warning: Could not retrieve API Key automatically."
    echo "Please run: docker exec server-moni-agent cat data/api_key.txt"
else
    # Detect Public IP
    PUBLIC_IP=$(curl -s ifconfig.me || curl -s icanhazip.com || echo "<YOUR_SERVER_IP>")

    echo ""
    echo "Agent is running!"
    echo "=================================================="
    echo "SETUP INSTRUCTIONS:"
    echo "1. Go to your Cloud Dashboard"
    echo "2. Click 'Add System'"
    echo "3. Enter these details:"
    echo ""
    echo "   Name:    $(hostname)"
    echo "   URL:     http://$PUBLIC_IP:8080"
    echo "   API Key: $API_KEY"
    echo ""
    echo "=================================================="
fi
