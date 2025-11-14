cd ./hostamat/infra/prod
git pull
docker compose pull
docker compose stop
docker compose up --build -d --remove-orphans
