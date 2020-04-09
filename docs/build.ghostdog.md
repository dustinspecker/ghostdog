# build.ghostdog

`build.ghostdog` files are the meat (dogfood?) of using ghostdog. These files instruct ghostdog how
to build projects and what depends on what.

## syntax

`build.ghostdog` files are written using the [Starlark Language](https://github.com/google/starlark-go/blob/master/doc/spec.md).
This language offers a lot of dynamic support over typical configuration files. Starlark is pretty similar to Python, with one
major exception, no side-effects. This means interpreting `build.ghostdog` files will always result in the same outcome.

## functions

More information about each function supported by ghostdog may be found in [functions](functions).

Two main functions need to be known to take advantage of ghostdog.

### files

When a project has files such as source code, test code, or even configuration for build tools like `make` or `webpack` ghostdog
needs a way to know about these files. That's where the [`files`](functions/files.md) function comes in. This function informs
ghostdog of file paths and allows projects to nicely configure these paths into groups. So, typically a `build.ghostdog` file
would have a `files` for souce code and another one test code. They could be combined into a single `files` call, but this will
cause ghostdog to do more work than desired!

ghostdog uses these files to detect changes from prevoius runs. If there have been changes to any of the files within a `files` then
ghostdog will run any `rule` that depends on these files. So, if a single `files` exists for both source code and test code then ghostdog
would always run a build and test `rule` that depend on this `files`. This isn't ideal. There's no point in building a binary of tests don't
depend on source code. By make `files` more granular ghostdog is able to be more intelligent. So having a `files` for souce code and another
one for test code enables ghostdog to skip doing extra work. A win for everyone.

### rule

The [`rule`](functions/rule.md) function tells ghostdog what commands to run, what `files` and other `rule`s it depends on, and the expected
output for the `rule`. Based on this information ghostdog is able to intelligently decide if it needs to run the commands. If ghostdog has ran
the command previously for the given dependencies with the exact same content then there's no need to run the commands. Ghostdog will use the
cached output from a previous run.

Just like `files` a `rule` can depend on a bunch and define a bunch of commands. The more granular `rule`s are the smarter and lazier ghostdog can
be. There could exist a single `rule` to build, run unit tests, and run integration tests. But that means all of these will always be ran when
anything changes. When in reality if only a unit test changed then there's no reason to have to build or run integration tests. So define granular
rules and let ghostdog be lazy.

## load

The Starlark language offers the ability to implement a `load` function. Just so happens ghostdog does implement a `load` function to enable
easy code re-use between projects. More information may be read in [functions/load.md](functions/load.md).
