#!/bin/bash
set -e

IMAGE_NAME="visit-service-go"
CONTAINER_NAME="visit-service-go"

TIMESTAMP=$(date +"%Y-%m-%d_%H-%M-%S")
LOG_FILE="logs/${TIMESTAMP}.log"
mkdir -p logs

podman stop "$CONTAINER_NAME" 2>/dev/null || true
podman rm "$CONTAINER_NAME" 2>/dev/null || true

podman build -t "$IMAGE_NAME" .
podman run -d --name "$CONTAINER_NAME" -p 8088:8088 "$IMAGE_NAME"

podman logs -f "$CONTAINER_NAME" >> "$LOG_FILE" 2>&1 &

echo "Container '$CONTAINER_NAME' is running on port 8088"
echo "Logs: $LOG_FILE"