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
cp server-moni /usr/local/bin/server-moni-bin
chmod +x /usr/local/bin/server-moni-bin

# Install Management Tool
echo "Installing management tool..."
cat << 'EOF' > /usr/local/bin/server-moni
#!/bin/bash

# Server Monitor Management Script
SERVICE_NAME="server-moni"
DOCKER_CONTAINER_NAME="server-moni-agent"

command=$1
if [ -z "$command" ]; then
    echo "Usage: server-moni [start|stop|restart|status|logs]"
    exit 1
fi

# Detect Mode
MODE=""
if command -v systemctl &> /dev/null && systemctl list-units --full -all | grep -Fq "$SERVICE_NAME.service"; then
    MODE="systemd"
elif command -v docker &> /dev/null && docker ps -a --format '{{.Names}}' | grep -Fq "$DOCKER_CONTAINER_NAME"; then
    MODE="docker"
else
    # Fallback check for running container if not found by name
    if [ "$command" == "stop" ] || [ "$command" == "restart" ]; then
         if docker ps -q -f name=$DOCKER_CONTAINER_NAME > /dev/null; then
             MODE="docker"
         fi
    fi
    
    if [ -z "$MODE" ]; then
        echo "Error: Could not detect Server Monitor service or container."
        exit 1
    fi
fi

case "$command" in
    start)
        if [ "$MODE" == "systemd" ]; then
            sudo systemctl start $SERVICE_NAME
        else
            docker start $DOCKER_CONTAINER_NAME
        fi
        echo "Started."
        ;;
    stop)
        if [ "$MODE" == "systemd" ]; then
            sudo systemctl stop $SERVICE_NAME
        else
            docker stop $DOCKER_CONTAINER_NAME
        fi
        echo "Stopped."
        ;;
    restart)
        if [ "$MODE" == "systemd" ]; then
            sudo systemctl restart $SERVICE_NAME
        else
            docker restart $DOCKER_CONTAINER_NAME
        fi
        echo "Restarted."
        ;;
    status)
        if [ "$MODE" == "systemd" ]; then
            sudo systemctl status $SERVICE_NAME
        else
            docker ps -f name=$DOCKER_CONTAINER_NAME
        fi
        ;;
    logs)
        if [ "$MODE" == "systemd" ]; then
            sudo journalctl -u $SERVICE_NAME -f
        else
            docker logs -f $DOCKER_CONTAINER_NAME
        fi
        ;;
    *)
        echo "Invalid command: $command"
        echo "Usage: server-moni [start|stop|restart|status|logs]"
        exit 1
        ;;
esac
EOF
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
ExecStart=/usr/local/bin/server-moni-bin
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
