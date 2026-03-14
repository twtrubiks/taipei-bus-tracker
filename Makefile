.PHONY: lint lint-go lint-web test test-go test-web check build build-notify

lint: lint-go lint-web

lint-go:
	go vet ./...
	@if command -v golangci-lint >/dev/null 2>&1; then golangci-lint run; fi

lint-web:
	cd web && npm run lint

test: test-go test-web

test-go:
	go test ./...

test-web:
	cd web && npm run test

check: lint test

build: build-web build-go
	mkdir -p static
	cp -r web/dist/* static/

build-go:
	go build -o taipei-bus ./cmd/server

build-notify:
	go build -o taipei-bus-notify ./cmd/notify

build-web:
	cd web && npm install && npm run build
