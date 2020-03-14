.PHONY: build
build:
	go build -o bin/ghostdog cmd/ghostdog/main.go

.PHONY: test-unit
test-unit:
	go test ./... -cover -coverprofile=cover.out
