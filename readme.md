# ghostdog

## Install ghostdog from source

1. Install [Go](https://golang.org/dl/).
1. Add `$GOPATH/bin` to `$PATH` via `export PATH=$(go env GOPATH)/bin:$PATH`.
1. Clone the ghostdog repository via `git clone https://github.com/dustinspecker/ghostdog`.
1. Navigate to the ghostdog local repository via `cd ghostdog`.
1. Install ghostdog via `go install ./...`.
1. ghostdog can then be used via `ghostdog`.


## Using ghostdog

ghostdog uses `build.ghostdog` files written with [Starlark](https://github.com/bazelbuild/starlark) to understand how to build packages. An example `build.ghostdog` file looks like:

```starlark
# `files` functions are used to group files to later be used by `rule` function
files(
  name = "makefile",
  paths = ["Makefile"]
)

files(
  name = "source_code",
  paths = ["pkg/main.go"]
)

files(
  name = "test_code",
  paths = ["pkg/main_test.go"]
)

# this rule is only ran when the makefile or source_code files change
# this rule runs `make build` and expects that command to output a file named main
rule(
  name = "build",
  sources = ["makefile", "source_code"],
  commands = ["make build"],
  outputs = ["main"]
)

# this rule is only ran when the makefile or test_code files change, or when
# the build rule's output changes
rule(
  name = "test",
  sources = ["makefile", "test_code", "build"],
  commands = ["make test"],
  outputs = []
)
```

More examples exist in the [_examples directory](_examples).
