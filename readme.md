# ghostdog

## Install ghostdog

1. Install [Go](https://golang.org/dl/).
1. Add `$GOPATH/bin` to `$PATH` via `export PATH=$(go env GOPATH)/bin:$PATH`.
1. Clone the ghostdog repository via `git clone https://github.com/dustinspecker/ghostdog`.
1. Navigate to the ghostdog local repository via `cd ghostdog`.
1. Install ghostdog via `go install ./...`.
1. ghostdog can then be used via `ghostdog`.


## Using ghostdog

ghostdog uses `BUILD` files written with [Starlark](https://github.com/bazelbuild/starlark) to understand how to build packages. An example `BUILD` file looks like:

```starlark
files(
  name = "makefile",
  paths = ["Makefile"]
)

files(
  name = "source_code",
  paths = ["pkg/main.go"]
)

rule(
  name = "build",
  sources = ["makefile", "source_code"],
  commands = ["make build"],
  outputs = ["main"]
)
```

More examples exist in the [examples directory](examples).
