.PHONY: build
build:
	go build -o bin/ghostdog cmd/ghostdog/main.go

.PHONY: test-unit
test-unit:
	go test ./... -cover -coverprofile=cover.out

.PHONY: test-integration
test-integration: build
	EXAMPLES_DIRECTORY=$(realpath ./_examples/) GHOSTDOG_BINARY=$(realpath ./bin/ghostdog) go test -tags=integration ./tests/integration/

.PHONY: update-golden-files
update-golden-files:
	EXAMPLES_DIRECTORY=$(realpath ./_examples/) GHOSTDOG_BINARY=$(realpath ./bin/ghostdog) go test -tags=integration ./tests/integration/ -update

.PHONY: test
test: test-unit test-integration
