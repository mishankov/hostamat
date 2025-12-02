# Prod infra

## Deploy

docker config create caddyfile config/Caddyfile
docker config create dozzle_users config/dozzle/users.yml
docker stack config -c compose.yml > swarm.yml
docker stack deploy --compose-file swarm.yml --prune hostamat

## Links

- [Hostamat](https://hostamat.ru/)
- [Portainer](https://5.129.199.174:8080)
- [Dozzle](http://5.129.199.174:8081)
- [Gatus](http://147.45.109.74:8082)
- [Dashboard](http://147.45.109.74:8083/dashboard)

## Настройка сервера

Настройка адреса в приватной сети для VPS - https://timeweb.cloud/docs/vpc/managing-bgp-networks/configuring-private-ip#ubuntu
