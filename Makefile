# Define default storage type
STORAGE_TYPE ?= memory
CACHE_TYPE ?= none

# Phony targets
.PHONY: assistance clean-deps format lint-check test-all generate-docs compile execute

build: ## Compile the application binary
	go build -o url-shortener-app ./cmd/main.go

clean-deps: ## Clean up Go module dependencies
	go mod tidy

format: ## Format Go source code
	go fmt ./...

test-all: ## Run all tests and generate coverage report
	go test -cover ./... > coverage.out

generate-docs: ## Generate Swagger API documentation
	swag fmt -d ./cmd/ -d ./internal/handlers/
	swag init -g ./cmd/main.go --parseDependency --parseInternal

execute: ## Run the application with the specified storage type
	go run ./cmd/main.go --storage-type $(STORAGE_TYPE) --cache-type $(CACHE_TYPE)