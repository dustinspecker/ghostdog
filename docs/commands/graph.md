# graph

> create a graph representation of a project's build dependencies

## Usage

`ghostdog graph BUILD_TARGET`

For information about BUILD_TARGET consult the [docs](../build_targets.md).

This command will print a graph in the [DOT language](https://www.graphviz.org/doc/info/lang.html).
This output may be used by the `graphviz` package to create an image.

## Installing graphviz

Installing the `graphviz` package installs the `dot` CLI, which can transform a `.dot` file into an
image.

### debian

On debian systems `graphviz` may be intalled by running `sudo apt install graphviz`.

### macOS

On macOS, [brew](https://brew.sh/) may be used to install `graphviz` via `brew install graphviz`.

## Creating a graph image with graphviz

1. Run `ghostdog graph BUILD_TARGET | dot -Tpng > graph.png`
1. open `graph.png` in an image viewer
