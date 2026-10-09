#!/bin/bash

echo "Stopping and removing vnc-instance containers..."
docker ps -q --filter "name=vnc-instance-" | xargs -r docker stop
docker ps -aq --filter "name=vnc-instance-" | xargs -r docker rm

echo "Stopping main vnc container..."
docker stop vnc-main 2>/dev/null || true
docker rm vnc-main 2>/dev/null || true

echo "Removing vnc-network..."
docker network rm vnc-network 2>/dev/null || true

echo "Cleanup completed!"
