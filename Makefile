VERSION := "$(shell git describe --abbrev=0 --tags 2> /dev/null || echo 'v0.0.0')+$(shell git rev-parse --short HEAD)"

build:
	go build -ldflags "-X main.buildVersion=$(VERSION)"

run:
	go run main.go

gen:
	go generate ./...

reset:
	./barnbridge-backend reset --force
	./barnbridge-backend migrate
	./barnbridge-backend dev-setup

clean-run: build reset
	./barnbridge-backend scrape --vv
