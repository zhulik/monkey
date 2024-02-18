ifeq (,$(findstring j,$(MAKEFLAGS)))
	MAKEFLAGS += -j
endif


export GOEXPERIMENT=rangefunc


.PHONY: test lint repl

GOLANGCI_LINT_VERSION := $(shell cat .golangci-lint-version)

default: test lint

repl:
	go run ./cmd/monkey/main.go

test:
	ginkgo -p -r

lint: bin/golangci-lint
	./bin/golangci-lint run

lint-fix: bin/golangci-lint
	./bin/golangci-lint run --fix

bin/golangci-lint: .golangci-lint-version
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- v$(GOLANGCI_LINT_VERSION)
