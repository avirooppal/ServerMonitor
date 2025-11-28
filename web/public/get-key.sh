#!/bin/bash

# Server Moni - Get Agent Info

CONTAINER_NAME="server-moni-agent"

# Check if container is running
if ! docker ps | grep -q "$CONTAINER_NAME"; then
    echo "Error: Agent container '$CONTAINER_NAME' is not running."
    exit 1
fi

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
