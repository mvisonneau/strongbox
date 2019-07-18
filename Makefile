NAME          := strongbox
VERSION       := $(shell git describe --tags --abbrev=1)
FILES         := $(shell git ls-files '*.go')
REPOSITORY    := mvisonneau/$(NAME)
VAULT_VERSION := 1.1.3
.DEFAULT_GOAL := help

export GO111MODULE=on

.PHONY: setup
setup: ## Install required libraries/tools for build tasks
	@command -v goveralls 2>&1 >/dev/null  || GO111MODULE=off go get -u -v github.com/mattn/goveralls
	@command -v golint 2>&1 >/dev/null     || GO111MODULE=off go get -u -v golang.org/x/lint/golint
	@command -v cover 2>&1 >/dev/null      || GO111MODULE=off go get -u -v golang.org/x/tools/cmd/cover
	@command -v goimports 2>&1 >/dev/null  || GO111MODULE=off go get -u -v golang.org/x/tools/cmd/goimports

.PHONY: fmt
fmt: setup ## Format source code
	gofmt -s -w $(FILES)
	goimports -w $(FILES)

.PHONY: lint
lint: setup ## Run golint, goimports and go vet against the codebase
	golint -set_exit_status .
	go vet ./...
	goimports -d $(FILES) > goimports.out
	@if [ -s goimports.out ]; then cat goimports.out; rm goimports.out; exit 1; else rm goimports.out; fi

.PHONY: test
test: ## Run the tests against the codebase
	go test -v ./...

.PHONY: install
install: ## Build and install locally the binary (dev purpose)
	go install .

.PHONY: build
build: setup ## Build the binaries
	goreleaser release --snapshot --skip-publish --rm-dist

.PHONY: release
release: setup ## Build & release the binaries
	goreleaser release --rm-dist

.PHONY: publish-coveralls
publish-coveralls: setup ## Publish coverage results on coveralls
	goveralls -service drone.io -coverprofile=coverage.out

.PHONY: clean
clean: ## Remove binary if it exists
	rm -f $(NAME)

.PHONY: coverage
coverage: ## Generates coverage report
	rm -rf *.out
	go test -race -v ./... -coverprofile=coverage.out

.PHONY: dev-env
dev-env: ## Build a local development environment using Docker
	@docker run -d --cap-add IPC_LOCK --name vault vault:$(VAULT_VERSION)
	@sleep 2
	@docker run -it --rm --cap-add IPC_LOCK \
		-e VAULT_ADDR=http://$$(docker inspect vault | jq -r '.[0].NetworkSettings.IPAddress'):8200 \
		-e VAULT_TOKEN=$$(docker logs vault 2>/dev/null | grep 'Root Token' | cut -d' ' -f3 | sed -E "s/[[:cntrl:]]\[[0-9]{1,3}m//g") \
		vault:$(VAULT_VERSION) secrets enable transit
	@docker run -it --rm \
		-v $(shell pwd):/$(NAME) \
		-w /$(NAME) \
		-e VAULT_ADDR=http://$$(docker inspect vault | jq -r '.[0].NetworkSettings.IPAddress'):8200 \
		-e VAULT_TOKEN=$$(docker logs vault 2>/dev/null | grep 'Root Token' | cut -d' ' -f3 | sed -E "s/[[:cntrl:]]\[[0-9]{1,3}m//g") \
		goreleaser/goreleaser:v0.112.2 \
		/bin/bash -c 'make setup; make install; bash'
	@docker kill vault
	@docker rm vault -f

.PHONY: sign-drone
sign-drone: ## Sign Drone CI configuration
	drone sign $(REPOSITORY) --save

.PHONY: all
all: lint test build coverage ## Test, builds and ship package for all supported platforms

.PHONY: help
help: ## Displays this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
