#!/bin/bash

docker network create vnc-network 2>/dev/null || true

docker build ./ --tag vnc-main:latest
docker run -p 8080:8080 -v /var/run/docker.sock:/var/run/docker.sock --network vnc-network vnc-main:latest