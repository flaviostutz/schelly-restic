#!/bin/bash
set -e
# set -x

if [ "$RESTIC_PASSWORD" == "" ]; then
    echo "ENV RESTIC_PASSWORD is mandatory"
    exit 1
fi

export RESTIC_REPOSITORY=/backup-repo

if [ ! -d "$RESTIC_REPOSITORY/initialized" ]; then
    echo "Initializing new local path repository..."
    restic init
    touch "$RESTIC_REPOSITORY/initialized"
    echo "Repo initialized"
fi

echo "Starting Restic API..."
restic-api \
    --listen-port=$LISTEN_PORT \
    --listen-ip=$LISTEN_IP \
    --log-level=$LOG_LEVEL
