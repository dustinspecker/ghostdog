package resolver

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"
)

func GetBuildFileForPackage(fs afero.Fs, cwd, packagePath string) (string, error) {
	buildFileDirectory := filepath.Join(cwd, packagePath)
	buildFilePath := filepath.Join(buildFileDirectory, "BUILD")

	buildFileExists, err := afero.Exists(fs, buildFilePath)
	if err != nil {
		return "", err
	}

	if !buildFileExists {
		return "", fmt.Errorf("no BUILD file found in %s", buildFileDirectory)
	}

	return buildFilePath, nil
}
