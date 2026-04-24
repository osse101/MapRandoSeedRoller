.PHONY: test lint clean fmt

LINT    := go run -modfile=tools/go.mod github.com/golangci/golangci-lint/cmd/golangci-lint

test:
	go test ./... -v

lint:
	$(LINT) run --fix ./...

fmt:
	go fmt ./...

clean:
	go clean -testcache
