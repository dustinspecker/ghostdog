package main

import (
	"os"
	"path/filepath"

	"github.com/apex/log"
	"github.com/urfave/cli/v2"

	"github.com/dustinspecker/ghostdog/internal/build"
	"github.com/dustinspecker/ghostdog/internal/config"
	"github.com/dustinspecker/ghostdog/internal/graph"
)

func main() {
	appConfig, err := config.New("error")
	if err != nil {
		panic(err)
	}

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		appConfig.LogCtx.WithError(err).Fatal("getting home directory")
	}

	app := &cli.App{
		Name:  "ghostdog",
		Usage: "improve your build process",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "log-level",
				Usage: "level of logs to write (debug, info, warn, error, fatal)",
				Value: "error",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "build",
				Usage: "build projects using build.ghostdog files",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "cache-directory",
						Usage:   "where to write cached results on the system",
						Value:   filepath.Join(userHomeDir, ".cache"),
						EnvVars: []string{"XDG_CACHE_DIR"},
					},
				},
				ArgsUsage: "build.ghostdog_FILE TARGET_RULE",
				Action: func(c *cli.Context) error {
					userConfig, err := config.New(c.String("log-level"))
					if err != nil {
						appConfig.LogCtx.WithFields(log.Fields{
							"error": err.Error(),
						}).Fatal("creating userConfig")
					}
					userConfig.LogCtx = userConfig.LogCtx.WithFields(log.Fields{
						"subcommand": "build",
					})

					buildTarget := c.Args().Get(0)

					return build.RunBuildFile(userConfig.LogCtx, userConfig.Fs, userConfig.WorkingDirectory, buildTarget, filepath.Join(c.String("cache-directory"), "ghostdog"))
				},
			},
			{
				Name:      "graph",
				Usage:     "create a graph (DOT) of the build dependencies",
				ArgsUsage: "build.ghostdog_FILE TARGET_RULE",
				Action: func(c *cli.Context) error {
					userConfig, err := config.New(c.String("log-level"))
					if err != nil {
						appConfig.LogCtx.WithFields(log.Fields{
							"error": err.Error(),
						}).Fatal("creating userConfig")
					}
					userConfig.LogCtx = userConfig.LogCtx.WithFields(log.Fields{
						"subcommand": "graph",
					})

					buildTarget := c.Args().Get(0)

					return graph.GetGraph(userConfig.LogCtx, userConfig.Fs, userConfig.WorkingDirectory, buildTarget, os.Stdout)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		appConfig.LogCtx.WithError(err).Fatal("ran ghostdog")
	}
}
