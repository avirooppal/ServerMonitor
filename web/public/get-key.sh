#!/bin/bash

# Server Moni - Unified Installer & Key Retriever

echo "=================================================="
echo "       SERVER MONI AGENT SETUP"
echo "=================================================="

# 1. Check for Docker
if ! command -v docker &> /dev/null; then
    echo "Error: Docker is not installed."
    echo "Please install Docker first: https://docs.docker.com/get-docker/"
    exit 1
fi

# 2. Check if Agent is running
CONTAINER_NAME=$(docker ps --format '{{.Names}}' | grep -E "server.*agent" | head -n 1)

if [ -z "$CONTAINER_NAME" ]; then
    echo "Agent not found. Installing..."
    
    # Pull Image
    IMAGE_NAME="avirooppal/linux-monitoring-agent:latest"
    echo "Pulling latest agent image: $IMAGE_NAME..."
    docker pull $IMAGE_NAME

    # Remove existing stopped container if any
    if [ "$(docker ps -aq -f name=server-moni-agent)" ]; then
        echo "Removing old agent container..."
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
    
    CONTAINER_NAME="server-moni-agent"
    
    echo "Waiting for agent to initialize..."
    sleep 5
else
    echo "Found running agent: $CONTAINER_NAME"
fi

# 3. Retrieve API Key
API_KEY=$(docker exec $CONTAINER_NAME cat data/api_key.txt 2>/dev/null)

if [ -z "$API_KEY" ]; then
    echo "Error: Could not retrieve API Key. Is the agent initialized?"
    echo "Try running: docker logs $CONTAINER_NAME"
    exit 1
fi

# 4. Detect Public IP & Format URL
PUBLIC_IP=$(curl -s ifconfig.me || curl -s icanhazip.com || echo "<YOUR_SERVER_IP>")

if [[ "$PUBLIC_IP" == *":"* ]]; then
    AGENT_URL="http://[$PUBLIC_IP]:8080"
else
    AGENT_URL="http://$PUBLIC_IP:8080"
fi

# 5. Output Configuration
echo ""
echo "=================================================="
echo "       CONFIGURATION DETAILS"
echo "=================================================="
echo ""
echo "Enter these details in your Cloud Dashboard:"
echo ""
echo "Name:    $(hostname)"
echo "URL:     $AGENT_URL"
echo "API Key: $API_KEY"
echo ""
echo "=================================================="
