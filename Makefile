
.PHONY: lint
lint:
	$(VERBOSE) go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.0 run

.PHONY: install
install:
	CGO_ENABLED=0 go install -v -ldflags "-s -w" ./...

.PHONY: test
test:
	$(VERBOSE) go run gotest.tools/gotestsum@latest --format testname --junitfile test-results/junit.xml -- -coverprofile cover.out ./...

.PHONY: watch
watch:
	$(VERBOSE) go run gotest.tools/gotestsum@latest --watch --format testname -- -coverprofile cover.out ./...
