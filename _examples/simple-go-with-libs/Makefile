go            ?= docker run --rm -w /app -v ${CURDIR}:/app golang:1.13.10-stretch go
.PHONY: build
build:
	${go} build -o main pkg/main.go

.PHONY: test
test:
	${go} test ./...
