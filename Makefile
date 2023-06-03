.DEFAULT_GOAL := help

export LC_ALL=en_US.UTF-8
export $(shell sed 's/=.*//' .env)
-include .env

.PHONY: help
help: ## Help command
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n"} /^[$$()% a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: run
run: ## Run App
	go run cmd/app/main.go

.PHONY: tests
tests: ## Run Tests
	GOARCH=amd64 go test -v ./... -cover

.PHONY: check-formatting
check-formatting: ## Run Linting
	golangci-lint run

.PHONY: check-formatting
generate: install-dependencies ## Run OpenApi Generator
	go generate ./...

.PHONY: build
build: ## Build app
	go build -o bin/app cmd/app/main.go


migrate: ## run database migrations
	goose -dir migrations up


migrate-rollback: ## run database migrations
	goose -dir migrations down

clean: ## clean all the containers running
