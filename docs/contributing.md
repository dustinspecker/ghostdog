# contributing

## System Requirements

- Linux or macOS
- `make` is installed (optional, but recommended)
- [`go`](https://golang.org/dl/) is installed

## Project structure

ghostdog follows [golang-standards/project-layout](https://github.com/golang-standards/project-layout)
as much as possible.

## Building

ghostdog is written in go and thus requires go to build. This repository
has a `Makefile` to make this easier.

1. Run `make build` to create `ghostdog` at `bin/ghostdog`

If `make` is not installed, please look at the [Makefile](../Makefile)'s `build`
target to know the `go` command to execute.

## Running unit tests

1. Run `make test-unit` to run ghostdog's unit tests

This also produces a coverage report at `cover.out`. An HTML report may be
open via `go tool cover -html=cover.out`.

## Running integration tests

1. Run `make test-integration` to run ghostdog's integration tests

Some integration tests take advantage of *golden files*, sometimes referred
to as *snapshots*. These files record output from integration tests and compare
output to verify nothing has changed.

## Updating golden files

Sometimes the *golden files* need to be updated because behavior has intentionally
changed.

1. Run `UPDATE_GOLDEN=true make test-integration` to update *golden files*
1. Run `make test-integration` to verify integration tests pass after the previous
updates
