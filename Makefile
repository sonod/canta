TEST ?= $(shell go list ./...)
VERSION = $(shell cat version.go |grep Version |cut -f 2 -d "\"")
REVISION = $(shell git describe --always)

INFO_COLOR=\033[1;34m
RESET=\033[0m
BOLD=\033[1m

REVISION=$(shell git describe --always)
GOVERSION=$(shell go version)
BUILDDATE=$(shell date '+%Y/%m/%d %H:%M:%S %Z')
DIST ?= darwin
BUILD=pkg
ME=$(shell whoami)
default: build

GO ?= GO111MODULE=on go

ci: depsdev test lint ## Run test and more...

depsdev: ## Installing dependencies for development
	$(GO) get -u github.com/golang/lint/golint
	$(GO) get -u github.com/tcnksm/ghr
	GO111MODULE=off go get github.com/motemen/gobump/cmd/gobump

test: ## Run test
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Testing$(RESET)"
	$(GO) test -v $(TEST) -timeout=30s -parallel=4
	$(GO) test -race $(TEST)

lint: ## Exec golint
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Linting$(RESET)"
	golint -min_confidence 1.1 -set_exit_status $(TEST)


build: ## Build as linux binary
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Building$(RESET)"
	BUILD=$(BUILD) DIST=$(DIST) misc/build

ghr: ## Upload to Github releases without token check
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Releasing for Github$(RESET)"
	$(eval ver = v$(shell gobump show -r ./))
	ghr -u sonod ${ver} pkg

.PHONY: default test
