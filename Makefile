.DEFAULT_GOAL := help
.PHONY: help build vet test cover fmt fmt-check lint lint-fix check clean dev

# Note: `./...` relies on `ignore ./node_modules` in go.mod — the Vercel CLI's
# npm packages ship .go files that would otherwise break package matching.
PKGS    := ./...
LINT    := go tool -modfile=tools/go.mod golangci-lint
TESTFLAGS ?= -v

help: ## Show this help
	@awk 'BEGIN {FS = ":.*## "} /^[a-zA-Z_-]+:.*## / {printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Compile every package
	go build $(PKGS)

vet: ## Run go vet
	go vet $(PKGS)

test: ## Run tests (TESTFLAGS=-v by default; `make test TESTFLAGS=` for quiet)
	go test $(TESTFLAGS) $(PKGS)

cover: ## Run tests with coverage and print the total
	go test -coverprofile=coverage.out $(PKGS)
	go tool cover -func=coverage.out | tail -1

fmt: ## Format code in place
	go fmt $(PKGS)

fmt-check: ## Fail if any file needs formatting
	@out=$$(gofmt -l $$(go list -f '{{.Dir}}' $(PKGS))); \
	if [ -n "$$out" ]; then echo "Needs gofmt:"; echo "$$out"; exit 1; fi

lint: ## Run golangci-lint (report only)
	$(LINT) run $(PKGS)

lint-fix: ## Run golangci-lint and apply auto-fixes
	$(LINT) run --fix $(PKGS)

check: fmt-check vet lint test ## Everything CI should run

clean: ## Clear test cache and coverage output
	go clean -testcache
	rm -f coverage.out coverage.html

dev: ## Run `vercel dev` with the global CLI (npm i -g vercel@latest)
	vercel dev
