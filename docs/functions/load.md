# load
> The `load` function allows importing functions and constants from other modules.

`build.ghostdog` may use `load` to import functions and constants from another
module. These other modules are also written in [Starlark](https://github.com/google/starlark-go/blob/master/doc/spec.md). Modules loaded by ghostdog may also call functions that are available in `build.ghostdog`. For a list of functions look at [docs/functions](.).

An example of a module that may be loaded is:

```starlark
# echo_awesome.ghostdog

# @param filepath (str) - filepath to echo awesome to
def echo_awesome_to_file(filepath):
  rule(
    name = "echo_awesome_to_" + filepath,
    sources = [],
    commands = ["echo awesome > " + filepath],
    outputs = [filepath]
  )
```

And an example `build.ghostdog` using the above `echo_awesome.ghostdog` module.

```starlark
# build.ghostdog

# load echo_awesome_to_file function
load("echo_awesome.ghostdog", "echo_awesome_to_file")

# use the loaded function
echo_awesome_to_file("awesome.txt")
```

It's a convention to use the `.ghostdog` extension when creating modules to be
loaded by `build.ghostdog` files.
