# BarnBridge backend

Getting started as a dev
```shell
docker run --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=password -d postgres
docker run --name redis -d -p 6379:6379 redis redis-server
cp config-sample.yml config.yml
# edit config to suit your needs
# set env BB_PG_PASSWORD or db/password in config,yml
go build
make reset
```
