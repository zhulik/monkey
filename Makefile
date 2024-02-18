ifeq (,$(findstring j,$(MAKEFLAGS)))
	MAKEFLAGS += -j
endif


.PHONY: test lint

GOLANGCI_LINT_VERSION := $(shell cat .golangci-lint-version)

default: test lint

test:
	ginkgo -p -r

lint: bin/golangci-lint
	./bin/golangci-lint run

lint-fix: bin/golangci-lint
	./bin/golangci-lint run --fix

bin/golangci-lint: .golangci-lint-version
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- v$(GOLANGCI_LINT_VERSION)
