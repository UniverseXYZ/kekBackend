# BarnBridge backend

Getting started as a dev
```shell
# postgres docker command
docker run --name redis -d -p 6379:6379 redis redis-server --appendonly yes
cp config-sample.yml config.yml
# edit config to suit your needs
# set env BB_PG_PASSWORD or db/password in config,yml
go build
make reset
```
