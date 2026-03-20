# Prod infra

## Deploy

./scripts/deploy-stack.sh

The backend container entrypoint runs `hostamat migrate` before `hostamat run`, so the deploy script only needs to refresh Docker configs and update the swarm stack.

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
