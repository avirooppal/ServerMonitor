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

# Download Binary
# We assume the server hosts the binary at /downloads/agent-linux-amd64
# If running from source (dev), we might copy it. But for production script:
echo "Downloading Agent from $SERVER_URL..."
curl -L "$SERVER_URL/downloads/agent-linux-amd64" -o /usr/local/bin/server-moni-agent

chmod +x /usr/local/bin/server-moni-agent

# Create Systemd Service
echo "Creating Service..."
cat <<EOF > /etc/systemd/system/server-moni-agent.service
[Unit]
Description=Server Monitor Agent
After=network.target

[Service]
Type=simple
User=root
ExecStart=/usr/local/bin/server-moni-agent
Restart=always
Environment="SERVER_URL=$SERVER_URL"
Environment="API_KEY=$API_KEY"

[Install]
WantedBy=multi-user.target
EOF

# Reload and Start
echo "Starting Service..."
systemctl daemon-reload
systemctl enable server-moni-agent
systemctl restart server-moni-agent

echo "Agent installed and connected to $SERVER_URL!"
