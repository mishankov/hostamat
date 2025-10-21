git pull
docker compose pull
docker compose up --build -d --remove-orphans
docker compose -f support.compose.yml up -d
docker compose logs hostamat -f
