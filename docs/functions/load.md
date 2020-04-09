# load
> The `load` function allows importing functions and constants from other modules.

`build.ghostdog` may use `load` to import functions and constants from another
module. These other modules are also written in [Starlark](https://github.com/google/starlark-go/blob/master/doc/spec.md). Modules loaded by ghostdog may also call functions that are available in `build.ghostdog`. For a list of functions look at [docs/functions](.). Modules are loaded relatively to the module loading them.

An example of loading modules:

Given the following directory structure:

```
project
├── build.ghostdog
└── libs
    ├── constants.ghostdog
    └── echo_awesome.ghostdog
```

```starlark
# project/libs/constants.ghostdog
MESSAGE = "awesome"
```

```starlark
# project/libs/echo_awesome.ghostdog

# load MESSAGE from constants.ghostdog
# note that the module is loading relative to this file (projects/libs/echo_awesome.ghostdog)
load("constants.ghostdog", "MESSAGE")

# @param filepath (str) - filepath to echo awesome to
def echo_awesome_to_file(filepath):
  rule(
    name = "echo_awesome_to_" + filepath,
    sources = [],
    commands = ["echo " + MESSAGE + " > " + filepath],
    outputs = [filepath]
  )
```

And an example `build.ghostdog` using the above `echo_awesome.ghostdog` module.

```starlark
# project/build.ghostdog

# load echo_awesome_to_file function
# not that the module is loaded relative to this file (projects/build.ghostdog)
load("libs/echo_awesome.ghostdog", "echo_awesome_to_file")

# use the loaded function
echo_awesome_to_file("awesome.txt")
```

It's a convention to use the `.ghostdog` extension when creating modules to be
loaded by `build.ghostdog` files.
