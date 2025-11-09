.PHONY: help test test-short test-race coverage fmt clean vet build ci-pr ci-main docs

help:
	@echo "Available targets:"
	@echo "  make test       - Run all tests"
	@echo "  make test-race  - Tests with race detector"
	@echo "  make test-short - Quick tests (dev)"
	@echo "  make coverage   - Generate HTML coverage report"
	@echo "  make fmt        - Format code"
	@echo "  make vet        - Run go vet"
	@echo "  make build      - Build application"
	@echo "  make docs       - Generate OpenAPI documentation (JSON and YAML)"
	@echo "  make ci-pr      - Replicate CI checks for PRs"
	@echo "  make ci-main    - Replicate CI checks for main"

test:
	go test ./... -v

test-race:
	go test ./... -race -cover

test-short:
	go test ./... -short -v

coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage generated: coverage.html"

vet:
	go vet ./...

fmt:
	go fmt ./...

build:
	go build -v -o api ./cmd/api

clean:
	rm -f coverage.out coverage.html
	rm -f api api.exe

docs:
	@echo "Generating OpenAPI documentation..."
	swag init -g cmd/api/main.go --output docs --outputTypes json,yaml --parseDependency --parseInternal
	@echo "✅ Documentation generated at docs/doc.json and docs/doc.yaml"

# Replicate exactly the CI checks for PRs
ci-pr: vet test-short build
	@echo "✅ PR checks passed"

# Replicate exactly the CI checks for main
ci-main: vet test-race coverage build
	@echo "✅ Main checks passed"
	@echo "Coverage: $$(go tool cover -func=coverage.out | grep total | awk '{print $$3}')"

