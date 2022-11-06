# GC is the go compiler.
GC = go

export GO111MODULE=on

default:
	chmod +rwx scripts/*.sh
	$(GC) build -o gidari-cli cmd/gidari/main.go

# containers build the docker containers for performing integration tests.
.PHONY: containers
containers:
	scripts/build-storage.sh

# test runs all of the unit tests locally. Each test is run 5 times to minimize flakiness.
.PHONY: tests
tests:
	$(GC) clean -testcache && go test -v -count=5 -tags=utests ./...

# ci runs the tests with a ci container
.PHONY: ci
ci:
	chmod +rwx scripts/*.sh
	$(GC) clean -testcache
	./scripts/run-ci-tests.sh

# fmt runs the formatter.
.PHONY: fmt
fmt:
	gofumpt -l -w .

# lint runs the linter.
.PHONY: lint
lint:
	golangci-lint run --fix
	golangci-lint run --config .golangci.yml

# add-license adds the license to all the top of all the .go files.
.PHONY: add-license
add-license:
	./scripts/add-license.sh