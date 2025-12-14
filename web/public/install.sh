#!/bin/bash

# Server Moni Agent Installer

echo "=================================================="
echo "Installing Server Moni Agent..."
echo "=================================================="

# Check for Docker
if ! command -v docker &> /dev/null; then
    echo "Error: Docker is not installed."
    echo "Please install Docker to run the agent."
    exit 1
fi

# Pull Image
IMAGE_NAME="avirooppal/linux-monitoring-agent:latest"
echo "Pulling latest agent image: $IMAGE_NAME..."
docker pull $IMAGE_NAME

# Remove existing container if any
if [ "$(docker ps -aq -f name=server-moni-agent)" ]; then
    echo "Removing existing agent container..."
    docker rm -f server-moni-agent
fi

# Run Agent
echo "Starting Agent container..."
docker run -d \
    --name server-moni-agent \
    --restart unless-stopped \
    --network host \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v server-moni-data:/app/data \
    -e API_PORT=8080 \
    $IMAGE_NAME

# Wait for startup and key generation
echo "Waiting for agent to initialize..."
sleep 5

# Retrieve API Key
API_KEY=$(docker exec server-moni-agent cat data/api_key.txt 2>/dev/null)

if [ -z "$API_KEY" ]; then
    echo "Warning: Could not retrieve API Key automatically."
    echo "Please run: docker exec server-moni-agent cat data/api_key.txt"
else
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
    echo "Management tool installed as 'server-moni'"

    # Detect Public IP
    PUBLIC_IP=$(curl -s ifconfig.me || curl -s icanhazip.com || echo "<YOUR_SERVER_IP>")

    echo ""
    echo "Agent is running!"
    echo "=================================================="
    echo "SETUP INSTRUCTIONS:"
    echo "1. Go to your Cloud Dashboard"
    echo "2. Click 'Add System'"
    echo "3. Enter these details:"
    echo ""
    echo "   Name:    $(hostname)"
    echo "   URL:     http://$PUBLIC_IP:8080"
    echo "   API Key: $API_KEY"
    echo ""
    echo "MANAGEMENT COMMANDS:"
    echo "   Use 'server-moni' to control the agent:"
    echo "   - server-moni stop"
    echo "   - server-moni start"
    echo "   - server-moni restart"
    echo ""
    echo "=================================================="
fi
