.PHONY: dev api desktop web build build-api build-desktop build-web \
       test test-v test-cover test-web clean clean-web clean-bin \
       lint lint-web install-web tidy env help

# Load .env if it exists (values become available as env vars to all targets)
-include .env
export

# --- Development (run) ---

dev: ## Run API server + web dev server concurrently
	@echo "Starting API on :$${PORT:-8080} and Web on :3000..."
	@trap 'kill 0' EXIT; \
		$(MAKE) api & \
		$(MAKE) web & \
		wait

api: ## Run the Go API server
	go run ./cmd/api

desktop: ## Run the Fyne desktop app
	go run ./cmd/desktop

web: ## Run the React dev server (requires API running)
	cd web && bun start

# --- Build ---

build: build-api build-desktop build-web ## Build everything

build-api: ## Build the API server binary
	go build -o bin/api ./cmd/api

build-desktop: ## Build the desktop app binary
	go build -o bin/desktop ./cmd/desktop

build-web: ## Build the React app for production
	cd web && npx react-scripts build

# --- Test ---

test: ## Run all Go tests
	go test ./internal/...

test-cover: ## Run Go tests with coverage report
	go test -coverprofile=coverage.out ./internal/...
	go tool cover -func=coverage.out
	@rm -f coverage.out

test-web: ## Run React tests (TODO: no test cases written yet)
	cd web && bun test -- --watchAll=false

test-all: ## Run all tests (Go + coverage + React), stops on first failure
	go test ./internal/... && \
	go test -coverprofile=coverage.out ./internal/... && \
	go tool cover -func=coverage.out && \
	rm -f coverage.out && \
	cd web && bun test -- --watchAll=false

# --- Lint / Tidy ---

lint: ## Vet Go code
	go vet ./...

lint-web: ## Lint React code
	cd web && npm run lint

tidy: ## Tidy Go modules
	go mod tidy

# --- Setup ---

env: ## Create .env files from examples (won't overwrite existing)
	@test -f .env || (cp .env.example .env && echo "Created .env")
	@test -f web/.env || (cp web/.env.example web/.env && echo "Created web/.env")
	@test -f .env && test -f web/.env && echo "Environment files ready."

install-web: ## Install web dependencies
	cd web && bun install

setup: env tidy install-web ## Full project setup (env + deps)

# --- Clean ---

clean: clean-bin clean-web ## Clean all build artifacts

clean-bin: ## Remove Go binaries
	rm -rf bin/

clean-web: ## Remove web build output
	rm -rf web/build

# --- Help ---

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
