ifeq (,$(findstring j,$(MAKEFLAGS)))
	MAKEFLAGS += -j
endif


export GOEXPERIMENT=rangefunc


.PHONY: test lint repl lint-fix bench cpu.prof

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

bench:
	go test -bench=.  -benchtime=10s

bin/golangci-lint: .golangci-lint-version
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- v$(GOLANGCI_LINT_VERSION)

cpuprof: cpu.svg
	xdg-open cpu.svg

cpu.svg: cpu.prof
	go tool pprof -dot cpu.prof | dot -Tsvg > cpu.svg

cpu.prof:
	go test -cpuprofile cpu.prof -bench .
	rm monkey.test
