#!/bin/bash

# Server Monitor Management Script
# Usage: ./manage.sh [start|stop|restart|status|logs]

command=$1

if [ -z "$command" ]; then
    echo "Usage: $0 [start|stop|restart|status|logs]"
    exit 1
fi

# Service Names
SVC_SERVER="server-moni"
SVC_AGENT="ServerMoniAgent"
DOCKER_CONTAINER="server-moni-agent"

# Detect Mode and Target
MODE=""
TARGET=""

if command -v systemctl &> /dev/null; then
    # Prioritize active services
    if systemctl is-active --quiet $SVC_SERVER; then
        MODE="systemd"
        TARGET=$SVC_SERVER
    elif systemctl is-active --quiet $SVC_AGENT; then
        MODE="systemd"
        TARGET=$SVC_AGENT
    # If neither is active, check which one is installed
    elif systemctl list-units --full -all | grep -Fq "$SVC_SERVER.service"; then
        MODE="systemd"
        TARGET=$SVC_SERVER
    elif systemctl list-units --full -all | grep -Fq "$SVC_AGENT.service"; then
        MODE="systemd"
        TARGET=$SVC_AGENT
    fi
fi

# If systemd didn't match, try Docker
if [ -z "$MODE" ]; then
    if command -v docker &> /dev/null; then
        if docker ps -a --format '{{.Names}}' | grep -Fq "$DOCKER_CONTAINER"; then
            MODE="docker"
            TARGET=$DOCKER_CONTAINER
        fi
    fi
fi

if [ -z "$MODE" ]; then
    echo "Error: Could not detect Server Monitor service or container."
    exit 1
fi

echo "Mode: $MODE"
echo "Target: $TARGET"

case "$command" in
    start)
        if [ "$MODE" == "systemd" ]; then
            sudo systemctl start $TARGET
        else
            docker start $TARGET
        fi
        echo "Started."
        ;;
    stop)
        if [ "$MODE" == "systemd" ]; then
            sudo systemctl stop $TARGET
        else
            docker stop $TARGET
        fi
        echo "Stopped."
        ;;
    restart)
        if [ "$MODE" == "systemd" ]; then
            sudo systemctl restart $TARGET
        else
            docker restart $TARGET
        fi
        echo "Restarted."
        ;;
    status)
        if [ "$MODE" == "systemd" ]; then
            sudo systemctl status $TARGET
        else
            docker ps -f name=$TARGET
        fi
        ;;
    logs)
        if [ "$MODE" == "systemd" ]; then
            sudo journalctl -u $TARGET -f
        else
            docker logs -f $TARGET
        fi
        ;;
    *)
        echo "Invalid command: $command"
        echo "Usage: $0 [start|stop|restart|status|logs]"
        exit 1
        ;;
esac
