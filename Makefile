.PHONY: help
NUM_OF_VERSIONS=1

ifneq (,$(wildcard ./.env))
    include .env
    export
    ENV_FILE_PARAM = --env-file .env
endif

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

docker-up: requirements ## Start docker containers
	@docker-compose up -d

docker-stop: requirements ## Stop docker containers
	@docker-compose stop

build-and-run: requirements ## Start docker containers and run de service
	go mod download
	go build -o main .
	docker-compose up -d db.digital-account.dev
	./main

mock-generate: ## Generate mocks
	go generate ./...

requirements:
	@docker -v > /dev/null 2>&1 || { echo >&2 "I require docker but it's not installed. See : https://docs.docker.com/engine/install/"; exit 127;}
	@docker-compose -v > /dev/null 2>&1 || { echo >&2 "I require docker-compose but it's not installed. See : https://docs.docker.com/compose/install/"; exit 127;}