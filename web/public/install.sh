#!/bin/bash

# Server Moni Agent Installer

SERVER_URL=""

# Parse args
while [[ "$#" -gt 0 ]]; do
    case $1 in
        -server) SERVER_URL="$2"; shift ;;
        *) echo "Unknown parameter passed: $1"; exit 1 ;;
    esac
    shift
done

if [ -z "$SERVER_URL" ]; then
    echo "Usage: $0 -server <URL>"
    exit 1
fi

# Generate a random Agent Token
if command -v uuidgen &> /dev/null; then
    AGENT_TOKEN=$(uuidgen)
else
    AGENT_TOKEN=$(cat /proc/sys/kernel/random/uuid 2>/dev/null || date +%s | md5sum | head -c 32)
fi

echo "=================================================="
echo "Installing Server Moni Agent..."
echo "Generated Agent Token: $AGENT_TOKEN"
echo "=================================================="

# Check for Docker
if command -v docker &> /dev/null; then
    echo "Docker found. Running Agent container..."
    
    docker run -d \
        --name server-moni-agent \
        --restart unless-stopped \
        --network host \
        -v /var/run/docker.sock:/var/run/docker.sock \
        -e SERVER_URL="$SERVER_URL" \
        -e API_KEY="$AGENT_TOKEN" \
        server-moni-agent
        
elif command -v go &> /dev/null; then
    echo "Go found. Building Agent from source..."
    go build -o agent cmd/agent/main.go
    
    echo "Starting Agent in background..."
    nohup ./agent -server "$SERVER_URL" -key "$AGENT_TOKEN" > agent.log 2>&1 &
    
else
    echo "Error: Neither Docker nor Go found."
    echo "Please install Docker or Go to run the agent."
    exit 1
fi

echo ""
echo "Agent is running!"
echo "PLEASE COPY THE TOKEN BELOW AND PASTE IT IN THE DASHBOARD:"
echo ""
echo "   $AGENT_TOKEN"
echo ""
echo "=================================================="
