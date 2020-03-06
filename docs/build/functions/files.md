# files
> The `files` function collects a group of files that can later be used by the [rule function](./rule.md)

`files` is a specific type of `rule` that find files. These files are **never** cached by Ghostdog for performance reasons, so `files` is a great fit for finding source files, test files, documentation files, etc., but should **not** be used to find output files that would be created by a `rule`. A good rule of thumb is use `files` to aggregate files that are typically committed.

## Arguments

| Argument Name | Type | Usage | Required |
| ------------- | ---- | ----- | -------- |
| name          | String | This is the name that `rule` will use to reference this file group. Names may only use lowercase letters and underscores (regex used: ^[a-z\_]*$). | true |
| paths         | List of Strings | The filepaths to include in this files group. All paths **must** exist. Paths may **not** start with `/` (absolute paths). Path may **not** use `..` (reference parent directory). | true |

## Example

The `files` function can be used at any place in a `BUILD` file. The order does not matter. And in fact the `rule` function may reference a `files` rule before the `files` rule is even defined.

```python
files(
  name = "source_files",
  paths = ["pkg/main.go"]
)

files(
  name = "test_files",
  paths = ["pkg/main_test.go"]
)
```
