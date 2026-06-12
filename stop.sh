#!/bin/bash

IMAGE_NAME="visit-service-go"
CONTAINER_NAME="visit-service-go"

podman stop "$CONTAINER_NAME" 2>/dev/null && echo "Stopped container '$CONTAINER_NAME'" || echo "Container '$CONTAINER_NAME' was not running"
podman rm "$CONTAINER_NAME" 2>/dev/null && echo "Removed container '$CONTAINER_NAME'" || echo "Container '$CONTAINER_NAME' did not exist"
podman rmi "$IMAGE_NAME" 2>/dev/null && echo "Removed image '$IMAGE_NAME'" || echo "Image '$IMAGE_NAME' did not exist"