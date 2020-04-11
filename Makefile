.DEFAULT_GOAL := test

UPDATE_GOLDEN ?= false

ifeq ($(UPDATE_GOLDEN), true)
	_update_arg="-update"
endif

.PHONY: build
build:
	go build -o bin/ghostdog cmd/ghostdog/main.go

.PHONY: test-unit
test-unit:
	go test ./... -cover -coverprofile=cover.out -covermode=count

.PHONY: test-integration
test-integration: build
	EXAMPLES_DIRECTORY=$(realpath ./_examples/) GHOSTDOG_BINARY=$(realpath ./bin/ghostdog) go test -tags=integration ./tests/integration/ $(_update_arg)

.PHONY: test
test: test-unit test-integration
