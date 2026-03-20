#!/bin/sh
set -eu

cd ./hostamat/infra/prod

# Pulling latest git and remembering the hash
git pull --ff-only
REV=$(git rev-parse --short=12 HEAD)

# Exporting unique config namse to use in compose file
export CADDYFILE_CONFIG="caddyfile-$REV"
export DOZZLE_USERS_CONFIG="dozzle_users-$REV"
export GATUS_CONFIG="gatus_config-$REV"
export PROMETHEUS_CONFIG="prometheus_config-$REV"

# Creating new configs
docker config inspect "$CADDYFILE_CONFIG" >/dev/null 2>&1 ||
  docker config create "$CADDYFILE_CONFIG" config/Caddyfile

docker config inspect "$DOZZLE_USERS_CONFIG" >/dev/null 2>&1 ||
  docker config create "$DOZZLE_USERS_CONFIG" config/dozzle/users.yml

docker config inspect "$GATUS_CONFIG" >/dev/null 2>&1 ||
  docker config create "$GATUS_CONFIG" config/gatus/config.yaml

docker config inspect "$PROMETHEUS_CONFIG" >/dev/null 2>&1 ||
  docker config create "$PROMETHEUS_CONFIG" config/prometheus.yaml

# Generating random temporary swarm.yml file
TMP_SWARM_YML=$(mktemp)
trap 'rm -f "$TMP_SWARM_YML"' EXIT INT TERM HUP

docker stack config -c compose.yml > "$TMP_SWARM_YML"

# Deploying stack using the temporary swarm.yml file
docker stack deploy --compose-file "$TMP_SWARM_YML" --prune --with-registry-auth --detach=true hostamat
