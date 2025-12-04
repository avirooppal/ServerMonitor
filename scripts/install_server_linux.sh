#!/bin/bash
set -e

# Check for root
if [ "$EUID" -ne 0 ]; then
  echo "Please run as root"
  exit 1
fi

echo "Installing Server Monitor (Server)..."

# Check for Go
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed."
    echo "Please install Go: https://go.dev/doc/install"
    echo "Ubuntu/Debian: sudo apt install golang"
    exit 1
fi

# Build Frontend (if npm exists)
if command -v npm &> /dev/null; then
    echo "Building Frontend..."
    cd web
    npm install
    npm run build
    cd ..
    
    # Copy to dist
    rm -rf cmd/server/dist
    mkdir -p cmd/server/dist
    cp -r web/dist/* cmd/server/dist/
else
    echo "npm not found. Skipping Frontend build (assuming Vercel hosting)."
    # Create dummy dist for embed
    rm -rf cmd/server/dist
    mkdir -p cmd/server/dist
    echo "Frontend hosted externally" > cmd/server/dist/index.html
fi

# Build Backend
echo "Building Backend..."
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
