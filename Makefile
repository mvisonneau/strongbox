FILES   = $(shell git ls-files '*.go')
VERSION = $(shell git describe --always)
APP     = strongbox

all: test coverage build

build:
	CGO_ENABLED=1 GOOS=linux go build \
	  -ldflags "-extldflags -static -X main.version=$(VERSION)" \
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

test: lint vet imports
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

.PHONY: all build lint vet fm imports test install deps coverage clean prereqs

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
