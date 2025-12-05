#!/bin/bash
set -e

# Check for root
if [ "$EUID" -ne 0 ]; then
  echo "Please run as root"
  exit 1
fi

echo "Installing Server Monitor (Server)..."

# Skip Frontend Build (Hosted on Vercel)
echo "Skipping Frontend Build..."

# Prepare Dist for Embed (Minimal)
rm -rf cmd/server/dist
mkdir -p cmd/server/dist
# Copy get-key.sh which is needed for the install script
cp web/public/get-key.sh cmd/server/dist/
# Create a placeholder index.html
echo "Server Monitor API is running." > cmd/server/dist/index.html

# Build Backend
echo "Building Backend..."
go build -o server-moni ./cmd/server

# Build Agent (Linux AMD64)
echo "Building Agent..."
GOOS=linux GOARCH=amd64 go build -o agent-linux-amd64 ./cmd/agent

# Install Binary
echo "Installing Binary..."
systemctl stop server-moni || true
cp server-moni /usr/local/bin/server-moni
chmod +x /usr/local/bin/server-moni

# Create Data Directory
mkdir -p /var/lib/server-moni
mkdir -p /var/lib/server-moni/downloads
cp agent-linux-amd64 /var/lib/server-moni/downloads/agent-linux-amd64
cp scripts/install_agent_linux.sh /var/lib/server-moni/downloads/install_agent_linux.sh

# Create Systemd Service
echo "Creating Service..."
cat <<EOF > /etc/systemd/system/server-moni.service
[Unit]
Description=Server Monitor Dashboard
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/var/lib/server-moni
ExecStart=/usr/local/bin/server-moni
Restart=always
Environment="PORT=8080"

[Install]
WantedBy=multi-user.target
EOF

# Reload and Start
echo "Starting Service..."
systemctl daemon-reload
systemctl enable server-moni
systemctl restart server-moni

echo "Server Monitor installed and started on port 8080!"
