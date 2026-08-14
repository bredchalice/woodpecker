#!/usr/bin/env bash
set -Eeuo pipefail

IMAGE="${1:-woodpecker-server-lha:manual-inputs}"

docker build \
  --pull \
  -f docker/Dockerfile.server.fork \
  -t "$IMAGE" \
  .

echo "$IMAGE"
