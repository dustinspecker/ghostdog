# build

> run rules defined in a build.ghostdog file

## Usage

`ghostdog build BUILD_TARGET`

For information about BUILD_TARGET consult the [docs](../build_targets.md).

The build command will run the rules based on BUILD_TARGET. If the rules have
been previously built with the same content previously then the build command
will skip running the rule.

## high level overview of build

1. ghostdog analyzes the `build.ghostdog`'s [files](../functions/files.md) and [rule](../functions/rule.md)s to determine dependencies
1. ghostdog begins running the desired rule based on the BUILD_TARGET provided
1. ghostdog caches all built artifacts (a rule's outputs) on a successful
 run of the rule
1. ghostdog will use the cached built artifacts if it detects the sources and
 dependencies of the rule have been seen before in previous builds

## Purpose of the cache

ghostdog's cache enables skipping work when nothing has changed from a previous
run. This speeds up CI builds as well as local development. ghostdog only runs
what it has to.

When ghostdog skips running a rule, it will copy the cached built artifact to
the local location. So after running build, the local workspace will always be
the same regardless if the cache was able to be used or not.
