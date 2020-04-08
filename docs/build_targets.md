# build targets

Some of ghostdog's commands take a build target. A build target instructs ghostdog where to
find a `build.ghostdog` file and which rule to run within the `build.ghostdog` file. A build
target is typically in the structure of `path/to/package:rule_name_to_build`.

Given the following directory structure:

```
home
└── cool_projects
    └── project
        ├── Makefile
        ├── build.ghostdog
        └── pkg
            ├── main.go
            └── main_test.go
```

With the following contents for `home/cool_projects/project/build.ghostdog`:

```starlark
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

rule(
  name = "build",
  sources = ["makefile", "source_code"],
  commands = ["make build"],
  outputs = ["main"]
)

rule(
  name = "test",
  sources = ["makefile", "test_code", "build"],
  commands = ["make test"],
  outputs = []
)
```

Then the following table demonstrates how to specify a build target based on the current
working directory.

<table>
  <thead>
    <tr>
      <th>working directory</th>
      <th>build target</th>
      <th>result</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>home</td>
      <td>cool_projects/project:test</td>
      <td rowspan=4>runs test rule from <code>home/cool_projects/project/build.ghostdog</code> and its dependencies</td>
    </tr>
    <tr>
      <td>home/cool_projects</td>
      <td>project:test</td>
    </tr>
    <tr>
      <td rowspan=2>home/cool_projects/project</td>
      <td>.:test</td>
    </tr>
    <tr>
      <td>:test</td>
    </tr>
    <tr>
      <td rowspan=2>home</td>
      <td>cool_projects/project:all</td>
      <td rowspan=8>runs all rules in <code>home/cool_projects/project/build.ghostdog</code></td>
    </tr>
    <tr>
      <td>cool_projects/project</td>
    </tr>
    <tr>
      <td rowspan=2>home/cool_projects</td>
      <td>project:all</td>
    </tr>
    <tr>
      <td>project</td>
    </tr>
    <tr>
      <td rowspan=4>home/cool_projects/project</td>
      <td>.:all</td>
    </tr>
    <tr>
      <td>:all</td>
    </tr>
    <tr>
      <td>.</td>
    </tr>
    <tr>
      <td>&nbsp;</td>
    </tr>
  </tbody>
</table>
