# rule
> The `rule` function is the meat of a `build.ghostdog` file and informs Ghostdog what do.

`rule` functions explain what do to with files (and other `rule`s' output files) by listing commands to run. `rule` functions also list their expected output files. All of this information is used by Ghostdog to know what to run. More importantly, Ghostdog leverages this information to decide if it even needs to run anything! If Ghostdog detects that it has already ran the commands for same inputs, then it knows it can use the outputs it cached from a previous run!

## Arguments

| Argument Name | Type | Usage | Required |
| ------------- | ---- | ----- | -------- |
| name          | String | This is the name that `rule` that other rules may reference this `rule` by to depend on its outputs. Names may only use lowercase letters, numbers, and underscores (regex used: ^[a-z0-9\_]*$). | true |
| sources | List of Strings | The name of the `rules` this `rule` depends on. Ghostdog won't run this `rule` until the `rules` in `sources` are all completed.| true |
| commands | List of Strings | A list of commands to run as part of this `rule`. The commands are ran sequentially. If any command fails then the remaining commands are skipped. | true |
| outputs | List of Strings | The list of filepaths that are created by this `rule`. It's okay to have zero outputs (must provide an empty list) as this sometimes makes sense like linting for example. | true |

## Example

The `rule` function can be used at any place in a `build.ghostdog` file. The order does not matter. And in fact the `rule` function may reference another `rule` before the `rule` is even defined.

```python
files(
  name = "source_files",
  paths = ["pkg/main.go"]
)

files(
  name = "test_files",
  paths = ["pkg/main_test.go"]
)

rule(
  name = "build",
  sources = ["source_files"],
  commands = ["go build pkg/main.go"],
  outputs = ["main"]
)

rule(
  name = "test",
  sources = ["source_files", "test_files"],
  commands = ["go test ./pkg/"],
  outputs = []
)

rule(
  name = "lint",
  sources = ["source_files", "test_files"],
  commands = ["go fmt -l -s **/*.go"],
  outputs = []
)
```
