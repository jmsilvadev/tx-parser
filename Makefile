export GO111MODULE=on
SHELL:=/bin/bash
DOCKER_C := docker-compose
.DEFAULT_GOAL := help
.PHONY: *

build-server: ## Build server component
	go clean -cache
	go mod tidy
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/tx-parser cmd/server/main.go

tests: ## Run unit tests
	go clean -cache
	go test -count=1 -covermode=count -coverprofile=coverage.out github.com/jmsilvadev/tx-parser/...
	go tool cover -func coverage.out

tests-cover: ## Run tests with coverage
	go clean -cache
	go test -count=1 -covermode=count -coverprofile=coverage.out github.com/jmsilvadev/tx-parser/...
	go tool cover -html=coverage.out -o coverage.html
	open coverage.html >&- 2>&- || \
	xdg-open coverage.html >&- 2>&- || \
	gnome-open coverage.html >&- 2>&-

clean: ## Clean all builts
	rm -rf ./bin

clean-tests: ## Clean tests
	go clean -cache
	rm *.out

up: ## Start docker container
	$(DOCKER_C) pull
	$(DOCKER_C) up -d 

up-build: ## Start docker container and rebuild the image
	go mod tidy
	go mod vendor
	$(DOCKER_C) pull
	$(DOCKER_C) up --build -d

down: ## Stop docker container
	$(DOCKER_C) down --remove-orphans

build-image:  ## Build docker image in daemon mode
	go mod tidy
	go mod vendor
	docker build . -t tx-parser
	
logs: ## Watch docker log files
	$(DOCKER_C) logs --tail 100 -f

help:
	@grep -E '^[a-zA-Z._-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
