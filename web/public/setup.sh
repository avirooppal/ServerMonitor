#!/bin/bash

# Server Moni Agent Setup Script
# Supports systemd-based Linux systems

function exit_with_error() {
  echo "Error: $1"
  echo "Please check the logs or contact support."
  exit 1
}

function detect_system() {
  if command -v systemctl &> /dev/null && [ -d "/etc/systemd/system" ]; then
    echo "systemd"
    return
  fi
  echo "unknown"
}

function setup_systemd() {
  echo "Detected systemd system - installing service..."

  local API_URL="$1"
  local API_TOKEN="$2"
  local AGENT_PATH="/opt/server-moni/agent"
  local SERVICE_FILE="/etc/systemd/system/server-moni-agent.service"

  # Create service file
  cat <<EOF > "$SERVICE_FILE"
[Unit]
Description=Server Moni Agent
After=network.target

[Service]
Type=simple
User=server-moni
Group=server-moni
WorkingDirectory=/opt/server-moni
ExecStart=$AGENT_PATH
Restart=always
RestartSec=5
Environment="API_URL=$API_URL"
Environment="API_TOKEN=$API_TOKEN"

[Install]
WantedBy=multi-user.target
EOF

  # Reload systemd
  systemctl daemon-reload
  
  # Enable and start service
  systemctl enable server-moni-agent.service
  systemctl start server-moni-agent.service

  if [ $? -eq 0 ]; then
    echo "Service installed and started successfully."
  else
    exit_with_error "Failed to start service."
  fi
}

# Main Execution

if [ "$EUID" -ne 0 ]; then
  exit_with_error "This script must be run as root."
fi

if [ $# -lt 2 ]; then
  echo "Usage: ./setup.sh <API_URL> <API_TOKEN>"
  exit 1
fi

API_URL="$1"
API_TOKEN="$2"

SYSTEM_TYPE=$(detect_system)
if [ "$SYSTEM_TYPE" != "systemd" ]; then
  exit_with_error "Unsupported system type. Only systemd is supported."
fi

# Create user
if ! id -u server-moni >/dev/null 2>&1; then
  echo "Creating system user server-moni..."
  useradd --system --user-group --shell /bin/false --create-home -b /opt/ server-moni
fi

# Prepare directory
mkdir -p /opt/server-moni/data
chown -R server-moni:server-moni /opt/server-moni

# Download Agent
ARCH=$(uname -m)
AGENT_URL=""

# TODO: Replace with actual release URL pattern
BASE_URL="https://github.com/avirooppal/ServerMonitor/releases/latest/download"

if [ "$ARCH" == "x86_64" ]; then
  AGENT_URL="$BASE_URL/agent-linux-amd64"
elif [ "$ARCH" == "aarch64" ]; then
  AGENT_URL="$BASE_URL/agent-linux-arm64"
else
  exit_with_error "Unsupported architecture: $ARCH"
fi

echo "Downloading agent from $AGENT_URL..."
wget -q --show-progress -O /opt/server-moni/agent "$AGENT_URL" || exit_with_error "Failed to download agent."
chmod +x /opt/server-moni/agent
chown server-moni:server-moni /opt/server-moni/agent

# Setup Service
setup_systemd "$API_URL" "$API_TOKEN"

echo ""
echo "=========================================="
echo "Server Moni Agent setup complete!"
echo "=========================================="
echo "Logs: journalctl -u server-moni-agent -f"
echo "Status: systemctl status server-moni-agent"
