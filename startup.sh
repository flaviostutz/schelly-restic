#!/bin/bash
set -e
# set -x

echo "Starting Backy2 API..."
backy2-api \
    --listen-port=$LISTEN_PORT \
    --listen-ip=$LISTEN_IP \
    --log-level=$LOG_LEVEL
