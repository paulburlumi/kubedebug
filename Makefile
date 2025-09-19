
.PHONY: lint
lint:
	$(VERBOSE) golangci-lint run

.PHONY: build
build:
	CGO_ENABLED=0 go build -v -ldflags "-s -w" -o bin/ ./...

.PHONY: test
test:
	$(VERBOSE) go run gotest.tools/gotestsum@latest --format testname --junitfile test-results/junit.xml -- -coverprofile cover.out ./...

.PHONY: watch
watch:
	$(VERBOSE) go run gotest.tools/gotestsum@latest --watch --format testname -- -coverprofile cover.out ./...
