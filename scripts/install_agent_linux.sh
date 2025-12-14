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
command=$1

if [ -z "$command" ]; then
    echo "Usage: server-moni [start|stop|restart|status|logs]"
    exit 1
fi

SVC_SERVER="server-moni"
SVC_AGENT="ServerMoniAgent"
DOCKER_CONTAINER="server-moni-agent"

MODE=""
TARGET=""

if command -v systemctl &> /dev/null; then
    if systemctl is-active --quiet $SVC_SERVER; then
        MODE="systemd"
        TARGET=$SVC_SERVER
    elif systemctl is-active --quiet $SVC_AGENT; then
        MODE="systemd"
        TARGET=$SVC_AGENT
    elif systemctl list-units --full -all | grep -Fq "$SVC_SERVER.service"; then
        MODE="systemd"
        TARGET=$SVC_SERVER
    elif systemctl list-units --full -all | grep -Fq "$SVC_AGENT.service"; then
        MODE="systemd"
        TARGET=$SVC_AGENT
    fi
fi

if [ -z "$MODE" ] && command -v docker &> /dev/null; then
    if docker ps -a --format '{{.Names}}' | grep -Fq "$DOCKER_CONTAINER"; then
        MODE="docker"
        TARGET=$DOCKER_CONTAINER
    fi
fi

if [ -z "$MODE" ]; then
    echo "Error: Could not detect Server Monitor service."
    exit 1
fi

echo "Mode: $MODE"
echo "Target: $TARGET"

case "$command" in
    start)
        [ "$MODE" == "systemd" ] && sudo systemctl start $TARGET || docker start $TARGET
        echo "Started."
        ;;
    stop)
        [ "$MODE" == "systemd" ] && sudo systemctl stop $TARGET || docker stop $TARGET
        echo "Stopped."
        ;;
    restart)
        [ "$MODE" == "systemd" ] && sudo systemctl restart $TARGET || docker restart $TARGET
        echo "Restarted."
        ;;
    status)
        [ "$MODE" == "systemd" ] && sudo systemctl status $TARGET || docker ps -f name=$TARGET
        ;;
    logs)
        [ "$MODE" == "systemd" ] && sudo journalctl -u $TARGET -f || docker logs -f $TARGET
        ;;
    *)
        echo "Invalid command."
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
