package main

import (
	"log"
	"os"

	"github.com/dustinspecker/ghostdog/internal/build"
)

func main() {
	buildFilePath := "BUILD"

	buildFileData, err := os.Open(buildFilePath)
	if err != nil {
		log.Fatal(err)
	}

	if err = build.RunBuildFile(buildFilePath, buildFileData); err != nil {
		log.Fatal(err)
	}
}
