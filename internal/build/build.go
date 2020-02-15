package build

import (
	"io"

	"go.starlark.net/starlark"

	"github.com/dustinspecker/ghostdog/internal/builtins"
	"github.com/dustinspecker/ghostdog/internal/rule"
)

func RunBuildFile(buildFileName string, buildFileData io.Reader) error {
	thread := &starlark.Thread{Name: "ghostdog-main"}

	addRule := func(rule rule.Rule) error {
		return nil
	}

	nativeFunctions := starlark.StringDict{
		"files": starlark.NewBuiltin("files", builtins.Files(addRule)),
	}

	_, err := starlark.ExecFile(thread, buildFileName, buildFileData, nativeFunctions)

	return err
}
