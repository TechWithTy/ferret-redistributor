#!/usr/bin/env bash
set -euo pipefail

if [[ -z "${HETZNER_HOST:-}" ]]; then
  echo "HETZNER_HOST not set" >&2
  exit 1
fi

ssh_opts=(-o StrictHostKeyChecking=no)

ssh "${ssh_opts[@]}" "root@${HETZNER_HOST}" <<'EOF'
set -e
cd /opt/social-scale || mkdir -p /opt/social-scale && cd /opt/social-scale
if [[ -d .git ]]; then
  git fetch origin && git reset --hard origin/main
else
  git clone https://github.com/Deal-Scale/ferret-redistributor.git .
fi
docker compose pull
GIT_SHA=$(git rev-parse --short HEAD) docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --remove-orphans
EOF




