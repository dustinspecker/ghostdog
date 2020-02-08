package build

import (
	"io"

	"go.starlark.net/starlark"

	"github.com/dustinspecker/ghostdog/internal/builtins"
)

func RunBuildFile(buildFileName string, buildFileData io.Reader) error {
	thread := &starlark.Thread{Name: "ghostdog-main"}

	nativeFunctions := starlark.StringDict{
		"files": starlark.NewBuiltin("files", builtins.Files),
	}

	_, err := starlark.ExecFile(thread, buildFileName, buildFileData, nativeFunctions)

	return err
}
