# Prod infra

## Deploy

docker config create caddyfile config/Caddyfile && \
docker config create dozzle_users config/dozzle/users.yml && \
docker config create gatus_config config/gatus/config.yaml && \
docker stack config -c compose.yml > swarm.yml && docker stack deploy --compose-file swarm.yml --prune --with-registry-auth hostamat

### Nuke

docker stack rm hostamat && docker config rm caddyfile dozzle_users gatus_config

## Links

- [Hostamat](https://hostamat.ru/)
- [Portainer](https://5.129.199.174:8080)
- [Dozzle](http://5.129.199.174:8081)
- [Gatus](http://5.129.199.174:8082)
- [Dashboard](http://5.129.199.174:8083/dashboard)

## Настройка сервера

Настройка адреса в приватной сети для VPS - https://timeweb.cloud/docs/vpc/managing-bgp-networks/configuring-private-ip#ubuntu
