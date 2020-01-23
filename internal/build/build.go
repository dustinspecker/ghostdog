package build

import (
	"io"

	"go.starlark.net/starlark"
)

func RunBuildFile(buildFileName string, buildFileData io.Reader) error {
	thread := &starlark.Thread{Name: "ghostdog-main"}

	_, err := starlark.ExecFile(thread, buildFileName, buildFileData, nil)

	return err
}
