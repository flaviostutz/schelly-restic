#!/bin/bash
set -e
# set -x

echo "Starting Restic API..."
restic-api \
    --listen-port=$LISTEN_PORT \
    --listen-ip=$LISTEN_IP \
    --log-level=$LOG_LEVEL \
    --repo-dir=/backup-repo \
    --source-path=/backup-source
    --pre-backup-command=$PRE_COMMAND \
    --post-backup-command=$POST_COMMAND

