package resolver

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

func GetBuildInfoForPackage(fs afero.Fs, cwd, target string) (string, string, error) {
	buildInfos := strings.Split(target, ":")
	if len(buildInfos) != 1 && len(buildInfos) != 2 {
		return "", "", fmt.Errorf("%s is an invalid target", target)
	}

	packagePath := buildInfos[0]

	targetRule := "all"
	if len(buildInfos) == 2 {
		targetRule = buildInfos[1]
	}

	buildFileDirectory := filepath.Join(cwd, packagePath)
	buildFilePath := filepath.Join(buildFileDirectory, "build.ghostdog")

	buildFileExists, err := afero.Exists(fs, buildFilePath)
	if err != nil {
		return "", "", err
	}

	if !buildFileExists {
		return "", "", fmt.Errorf("no build.ghostdog file found in %s", buildFileDirectory)
	}

	return buildFilePath, targetRule, nil
}
