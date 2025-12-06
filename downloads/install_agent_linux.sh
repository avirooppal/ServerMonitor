#!/bin/bash
set -e

# Usage: ./install_agent_linux.sh --server=http://your-server.com --token=YOUR_API_KEY

SERVER_URL="__SERVER_URL__"
API_KEY="$1"
URL_ARG="$2"

if [ -z "$API_KEY" ]; then
  echo "Usage: $0 <API_KEY> [SERVER_URL]"
  exit 1
fi

# Fallback if placeholder isn't replaced (e.g. manual download or GitHub raw)
if [ "$SERVER_URL" = "__SERVER_URL__" ]; then
    if [ -n "$URL_ARG" ]; then
        SERVER_URL="$URL_ARG"
    else
        echo "Error: SERVER_URL not configured. Pass it as the second argument:"
        echo "Usage: $0 <API_KEY> <SERVER_URL>"
        exit 1
    fi
fi

# Check for root
if [ "$EUID" -ne 0 ]; then
  echo "Please run as root"
  exit 1
fi

echo "Installing Server Monitor Agent..."

# Download Binary
echo "Downloading Agent from $SERVER_URL..."
curl -L "$SERVER_URL/downloads/agent-linux-amd64" -o /usr/local/bin/server-moni-agent
chmod +x /usr/local/bin/server-moni-agent

# Install Service using built-in flag (kardianos/service)
# We pass the arguments so they are saved in the service definition
/usr/local/bin/server-moni-agent -server "$SERVER_URL" -token "$API_KEY" --service install

# Start Service
/usr/local/bin/server-moni-agent --service start

echo "Agent installed and connected to $SERVER_URL!"
