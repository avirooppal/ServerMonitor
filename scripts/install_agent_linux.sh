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
      # Assume positional arguments if not a flag
      if [[ "$i" != --* ]]; then
        if [ -z "$API_KEY" ]; then
           API_KEY="$i"
        elif [ -z "$SERVER_URL" ]; then
           SERVER_URL="$i"
        fi
      fi
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
elif command -v systemctl &> /dev/null && systemctl list-units --full -all | grep -Fq "ServerMoniAgent.service"; then
    MODE="systemd"
    SERVICE_NAME="ServerMoniAgent"
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

echo "Detected mode: $MODE"

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

# Install Service using built-in flag (kardianos/service)
# We pass the arguments so they are saved in the service definition
/usr/local/bin/server-moni-agent -server "$SERVER_URL" -token "$API_KEY" --service install

# Start Service
/usr/local/bin/server-moni-agent --service start

echo "Agent installed and connected to $SERVER_URL!"
echo "Use 'server-moni' to manage the service (start/stop/status/logs)."
