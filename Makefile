SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

DATE 		:= $(shell date +"%a %b %d %T %Y")
UNAME_S 	:= $(shell uname -s | tr A-Z a-z)

##@ Development

.PHONY: build
build: ## Build geniveev (output: build/gen)
	go build -o build/gen cmd/geniveev/main.go

.PHONY: test
test: ## Run tests
	go test ./...

.PHONY: watch
watch: ##
	find . *.go | entr ${MAKE} build

##@ Common

.PHONY: help
help: ## Show this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Other

.PHONY: release
release: ## Push a new release to github
	rm -rf dist/
	goreleaser release
