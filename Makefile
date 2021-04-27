VERSION := "$(shell git describe --abbrev=0 --tags 2> /dev/null || echo 'v0.0.0')+$(shell git rev-parse --short HEAD)"

build:
	go build -ldflags "-X main.buildVersion=$(VERSION)"

run:
	go run main.go

gen:
	go generate ./...

reset:
	./kekBackend reset --force
	./kekBackend migrate

clean-run: build reset
	./kekBackend sync-sy-pools
	./kekBackend scrape --vv
