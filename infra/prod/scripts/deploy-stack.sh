cd ./hostamat/infra/prod
git pull

docker stack rm hostamat
docker config rm caddyfile dozzle_users gatus_config

docker config create caddyfile config/Caddyfile
docker config create dozzle_users config/dozzle/users.yml
docker config create gatus_config config/gatus/config.yaml


docker stack config -c compose.yml > swarm.yml
(r=3;while ! docker stack deploy --compose-file swarm.yml --prune --with-registry-auth hostamat ; do ((--r))||exit;sleep 5;done)
