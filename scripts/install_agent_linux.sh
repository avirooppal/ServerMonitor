#!/bin/bash
set -e

# Usage: ./install_agent_linux.sh --server=http://your-server.com --token=YOUR_API_KEY

SERVER_URL=""
API_KEY=""

for i in "$@"; do
  case $i in
    --server=*)
      SERVER_URL="${i#*=}"
      shift
      ;;
    --token=*)
      API_KEY="${i#*=}"
      shift
      ;;
    *)
      ;;
  esac
done

if [ -z "$SERVER_URL" ] || [ -z "$API_KEY" ]; then
  echo "Usage: $0 --server=http://your-server.com --token=YOUR_API_KEY"
  exit 1
fi

# Check for root
if [ "$EUID" -ne 0 ]; then
  echo "Please run as root"
  exit 1
fi

echo "Installing Server Monitor Agent..."

# Stop existing service if running
if [ -f "/usr/local/bin/server-moni-agent" ]; then
  echo "Stopping existing agent..."
  /usr/local/bin/server-moni-agent --service stop || true
  /usr/local/bin/server-moni-agent --service uninstall || true
fi

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
