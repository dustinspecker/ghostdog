package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"

	"github.com/dustinspecker/ghostdog/internal/build"
	"github.com/dustinspecker/ghostdog/internal/graph"
)

func main() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	app := &cli.App{
		Name:  "ghostdog",
		Usage: "improve your build process",
		Commands: []*cli.Command{
			{
				Name:  "build",
				Usage: "build projects using BUILD files",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "cache-directory",
						Usage:   "where to write cached results on the system",
						Value:   filepath.Join(userHomeDir, ".cache"),
						EnvVars: []string{"XDG_CACHE_DIR"},
					},
				},
				ArgsUsage: "BUILD_FILE TARGET_RULE",
				Action: func(c *cli.Context) error {
					buildFilePath := c.Args().Get(0)
					buildTarget := c.Args().Get(1)

					return build.RunBuildFile(afero.NewOsFs(), buildFilePath, buildTarget, filepath.Join(c.String("cache-directory"), "ghostdog"))
				},
			},
			{
				Name:      "graph",
				Usage:     "create a graph (DOT) of the build dependencies",
				ArgsUsage: "BUILD_FILE TARGET_RULE",
				Action: func(c *cli.Context) error {
					buildFilePath := c.Args().Get(0)
					buildTarget := c.Args().Get(1)

					return graph.GetGraph(afero.NewOsFs(), buildFilePath, buildTarget, os.Stdout)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
