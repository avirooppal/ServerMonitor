#!/bin/bash

# Server Moni - Get Agent Info

# Find container name (handles docker-compose naming)
CONTAINER_NAME=$(docker ps --format '{{.Names}}' | grep -E "server.*agent" | head -n 1)

if [ -z "$CONTAINER_NAME" ]; then
    echo "Error: Could not find a running agent container."
    echo "Expected a container name containing 'server...agent'"
    exit 1
fi

echo "Found agent container: $CONTAINER_NAME"

# Retrieve API Key
API_KEY=$(docker exec $CONTAINER_NAME cat data/api_key.txt 2>/dev/null)

if [ -z "$API_KEY" ]; then
    echo "Error: Could not retrieve API Key. Is the agent initialized?"
    exit 1
fi

# Detect Public IP
PUBLIC_IP=$(curl -s ifconfig.me || curl -s icanhazip.com || echo "<YOUR_SERVER_IP>")

echo "=================================================="
echo "       SERVER MONI AGENT CONFIGURATION"
echo "=================================================="
echo ""
echo "Enter these details in your Cloud Dashboard:"
echo ""
echo "Name:    $(hostname)"
echo "URL:     http://$PUBLIC_IP:8080"
echo "API Key: $API_KEY"
echo ""
echo "=================================================="
