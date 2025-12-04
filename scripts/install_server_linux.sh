#!/bin/bash
set -e

# Check for root
if [ "$EUID" -ne 0 ]; then
  echo "Please run as root"
  exit 1
fi

echo "Installing Server Monitor (Server)..."

# Build Frontend
echo "Building Frontend..."
cd web
npm install
npm run build
cd ..

# Build Backend
echo "Building Backend..."
# Ensure we embed the frontend
rm -rf cmd/server/dist
mkdir -p cmd/server/dist
cp -r web/dist/* cmd/server/dist/
go build -o server-moni ./cmd/server

# Build Agents (for downloads)
echo "Building Agents..."
GOOS=linux GOARCH=amd64 go build -o agent-linux-amd64 ./cmd/agent
GOOS=windows GOARCH=amd64 go build -o agent-windows-amd64.exe ./cmd/agent

# Install Binary
echo "Installing Binary..."
cp server-moni /usr/local/bin/server-moni
chmod +x /usr/local/bin/server-moni

# Create Data Directory
mkdir -p /var/lib/server-moni
mkdir -p /var/lib/server-moni/downloads

# Copy Agents to Downloads
cp agent-linux-amd64 /var/lib/server-moni/downloads/
cp agent-windows-amd64.exe /var/lib/server-moni/downloads/

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
