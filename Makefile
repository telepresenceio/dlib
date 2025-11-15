#
# Intro

help:
	@echo 'Usage:'
	@echo '  make help'
	@echo '  make test'
	@echo '  make generate'
	@echo '  make dlib.cov.html'
	@echo '  make lint'
.PHONY: help

.SECONDARY:
.PHONY: FORCE
SHELL = bash

#
# Test

dlib.cov: test
	test -e $@
	touch $@
test:
	GOCOVERDIR=. go test -count=1 -coverprofile=dlib.cov -coverpkg=./... -race ./...
.PHONY: test

%.cov.html: %.cov
	go tool cover -html=$< -o=$@

#
# Generate

generate-clean:
	rm -f dlog/convenience.go
.PHONY: generate-clean

generate:
	go generate ./...
.PHONY: generate

#
# Lint

GOLANGCI_VERSION=v2.6.1

lint:
	docker run -e GOOS=linux --rm -v $$(pwd):/app -v ~/.cache/golangci-lint/$(GOLANGCI_VERSION):/root/.cache -w /app golangci/golangci-lint:$(GOLANGCI_VERSION) golangci-lint run ./...
	docker run -e GOOS=darwin --rm -v $$(pwd):/app -v ~/.cache/golangci-lint/$(GOLANGCI_VERSION):/root/.cache -w /app golangci/golangci-lint:$(GOLANGCI_VERSION) golangci-lint run ./...
	docker run -e GOOS=windows --rm -v $$(pwd):/app -v ~/.cache/golangci-lint/$(GOLANGCI_VERSION):/root/.cache -w /app golangci/golangci-lint:$(GOLANGCI_VERSION) golangci-lint run ./...
.PHONY: lint
