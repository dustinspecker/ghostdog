package main

import (
	"os"
	"path/filepath"

	"github.com/apex/log"
	apexCli "github.com/apex/log/handlers/cli"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"

	"github.com/dustinspecker/ghostdog/internal/build"
	"github.com/dustinspecker/ghostdog/internal/graph"
)

func getLogCtx(logLevel string) (*log.Entry, error) {
	parsedLogLevel, err := log.ParseLevel(logLevel)
	if err != nil {
		return nil, err
	}

	log.SetLevel(parsedLogLevel)
	log.SetHandler(apexCli.New(os.Stderr))

	logCtx := log.WithFields(log.Fields{
		"app": "ghostdog",
	})

	return logCtx, nil
}

func main() {
	logCtx, err := getLogCtx("error")
	if err != nil {
		panic(err)
	}

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		logCtx.WithError(err).Fatal("getting home directory")
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
					buildTarget := c.Args().Get(0)

					cwd, err := os.Getwd()
					if err != nil {
						return err
					}

					return build.RunBuildFile(afero.NewOsFs(), cwd, buildTarget, filepath.Join(c.String("cache-directory"), "ghostdog"))
				},
			},
			{
				Name:      "graph",
				Usage:     "create a graph (DOT) of the build dependencies",
				ArgsUsage: "BUILD_FILE TARGET_RULE",
				Action: func(c *cli.Context) error {
					buildTarget := c.Args().Get(0)

					cwd, err := os.Getwd()
					if err != nil {
						return err
					}
					return graph.GetGraph(afero.NewOsFs(), cwd, buildTarget, os.Stdout)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		logCtx.WithError(err).Fatal("ran ghostdog")
	}
}
