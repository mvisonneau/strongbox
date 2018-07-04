NAME          := strongbox
VERSION       := $(shell git describe --tags --abbrev=1)
FILES         := $(shell git ls-files '*.go')
LDFLAGS       := -linkmode external -extldflags -static -X 'main.version=$(VERSION)'
VAULT_VERSION := 0.9.3
.DEFAULT_GOAL := help

.PHONY: setup
setup: ## Install required libraries/tools
	go get -u -v github.com/golang/dep/cmd/dep
	go get -u -v golang.org/x/tools/cmd/goimports
	go get -u -v golang.org/x/tools/cmd/cover
	go get -u -v github.com/golang/lint/golint

.PHONY: fmt
fmt: ## Format source code
	goimports -w $(FILES)

.PHONY: lint
lint: ## Run golint and go vet against the codebase
	golint -set_exit_status . app config rand
	go vet ./...

.PHONY: test
test: ## Run the tests against the codebase
	go test -v ./...

.PHONY: install
install: ## Build and install locally the binary (dev purpose)
	go install .

.PHONY: build
build: ## Build the binary
	CGO_ENABLED=1 go build -ldflags "$(LDFLAGS)" -o $(NAME) main.go
	strip $(NAME)

.PHONY: deps
deps: ## Fetch all dependencies
	dep ensure -v

.PHONY: imports
imports: ## Fixes the syntax (linting) of the codebase
	goimports -d $(FILES)

.PHONY: clean
clean: ## Remove binary if it exists
	rm -f $(NAME)

.PHONY: coverage
coverage: ## Generates coverage report
	rm -rf *.out
	go test -coverprofile=coverage.out
	@for i in app config; do \
	 	go test -coverprofile=$$i.coverage.out github.com/mvisonneau/$(NAME)/$$i; \
		tail -n +2 $$i.coverage.out >> coverage.out; \
	done

.PHONY: dev-env
dev-env: ## Build a local development environment using Docker
	@docker run -d --cap-add IPC_LOCK --name vault vault:$(VAULT_VERSION)
	@sleep 1
	@docker run -it --rm \
		-e VAULT_ADDR=http://$$(docker inspect vault | jq -r '.[0].NetworkSettings.IPAddress'):8200 \
		-e VAULT_TOKEN=$$(docker logs vault 2>/dev/null | grep 'Root Token' | cut -d' ' -f3 | sed -E "s/[[:cntrl:]]\[[0-9]{1,3}m//g") \
		vault:$(VAULT_VERSION) secrets enable transit
	@docker run -it --rm \
		-v $(shell pwd):/go/src/github.com/mvisonneau/$(NAME) \
		-w /go/src/github.com/mvisonneau/$(NAME) \
		-e VAULT_ADDR=http://$$(docker inspect vault | jq -r '.[0].NetworkSettings.IPAddress'):8200 \
		-e VAULT_TOKEN=$$(docker logs vault 2>/dev/null | grep 'Root Token' | cut -d' ' -f3 | sed -E "s/[[:cntrl:]]\[[0-9]{1,3}m//g") \
		golang:1.9 \
		/bin/bash -c 'make setup; make deps; make install; bash'
	@docker kill vault
	@docker rm vault -f

.PHONY: all
all: lint imports test coverage build ## Test, builds and ship package for all supported platforms

.PHONY: help
help: ## Displays this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
