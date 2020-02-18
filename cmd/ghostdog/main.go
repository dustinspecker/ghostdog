package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/dustinspecker/ghostdog/internal/build"
)

func main() {
	app := &cli.App{
		Name:  "ghostdog",
		Usage: "improve your build process",
		Commands: []*cli.Command{
			{
				Name:  "build",
				Usage: "build projects using BUILD files",
				Action: func(c *cli.Context) error {
					buildFilePath := "BUILD"

					buildFileData, err := os.Open(buildFilePath)
					if err != nil {
						return err
					}

					return build.RunBuildFile(buildFilePath, buildFileData)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
