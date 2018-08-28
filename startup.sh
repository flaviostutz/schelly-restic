#!/bin/bash
set -e
# set -x

echo "Starting Restic API..."
schelly-restic \
    --listen-ip=$LISTEN_IP \
    --listen-port=$LISTEN_PORT \
    --log-level=$LOG_LEVEL \
    --repo-dir=/backup-repo \
    --source-path=/backup-source \
    --pre-backup-command=$PRE_COMMAND \
    --post-backup-command=$POST_COMMAND

