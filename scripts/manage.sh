#!/bin/bash

# Server Monitor Management Script
# Usage: ./manage.sh [start|stop|restart|status|logs]

SERVICE_NAME="server-moni"
DOCKER_CONTAINER_NAME="server-moni-agent"

command=$1

if [ -z "$command" ]; then
    echo "Usage: $0 [start|stop|restart|status|logs]"
    exit 1
fi

# Detect Mode
MODE=""
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
    echo "Error: Could not detect Server Monitor installation (Systemd service or Docker container)."
    exit 1
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
        echo "Usage: $0 [start|stop|restart|status|logs]"
        exit 1
        ;;
esac
