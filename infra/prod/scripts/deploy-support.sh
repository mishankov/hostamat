cd ./hostamat/infra/prod
git pull
docker compose -f support.compose.yml pull
docker compose -f support.compose.yml up -d --build --remove-orphans
