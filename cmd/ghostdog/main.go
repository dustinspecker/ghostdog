package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"

	"github.com/dustinspecker/ghostdog/internal/build"
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
				Action: func(c *cli.Context) error {
					buildFilePath := "BUILD"

					buildFileData, err := os.Open(buildFilePath)
					if err != nil {
						return err
					}

					return build.RunBuildFile(afero.NewOsFs(), buildFilePath, buildFileData, filepath.Join(c.String("cache-directory"), "ghostdog"))
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
