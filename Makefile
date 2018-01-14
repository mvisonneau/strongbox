FILES   = $(shell git ls-files '*.go')
VERSION = $(shell git describe --always --abbrev=4)
APP     = strongbox

all: lint vet imports test coverage build

build:
	CGO_ENABLED=1 GOOS=linux go build \
	  -ldflags "-linkmode external -extldflags -static -X main.version=$(VERSION)" \
		-o $(APP) \
		main.go $(LIBS)
	strip $(APP)

lint:
	golint -set_exit_status . app config

vet:
	go vet ./...

fmt:
	goimports -w $(FILES)

imports:
	goimports -d $(FILES)

test:
	go test -v ./...

install:
	go install .

deps:
	dep ensure -v

coverage:
	rm -rf *.out
	go test -coverprofile=coverage.out
	@for i in app config; do \
	 	go test -coverprofile=$$i.coverage.out github.com/mvisonneau/$(APP)/$$i; \
		tail -n +2 $$i.coverage.out >> coverage.out; \
	done

clean:
	rm -f $(APP)

prereqs:
	go get -u -v github.com/golang/dep/cmd/dep
	go get -u -v golang.org/x/tools/cmd/goimports
	go get -u -v golang.org/x/tools/cmd/cover
	go get -u -v github.com/golang/lint/golint

dev-env:
	@docker run -d --name vault vault
	@sleep 1
	@docker run -it --rm \
		-e VAULT_ADDR=http://$$(docker inspect vault | jq -r '.[0].NetworkSettings.IPAddress'):8200 \
		-e VAULT_TOKEN=$$(docker logs vault 2>/dev/null | grep 'Root Token' | cut -d' ' -f3) \
		vault mount transit
	@docker run -it --rm \
		-v $(shell pwd):/go/src/github.com/mvisonneau/$(APP) \
		-w /go/src/github.com/mvisonneau/$(APP) \
		-e VAULT_ADDR=http://$$(docker inspect vault | jq -r '.[0].NetworkSettings.IPAddress'):8200 \
		-e VAULT_TOKEN=$$(docker logs vault 2>/dev/null | grep 'Root Token' | cut -d' ' -f3) \
		golang:1.9 \
		/bin/bash -c 'make prereqs; make deps; make install; bash'
	@docker rm vault -f

.PHONY: all build lint vet fm imports test install deps coverage clean prereqs

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
