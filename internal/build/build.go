package build

import (
	"io"

	"go.starlark.net/starlark"

	"github.com/dustinspecker/ghostdog/internal/builtins"
	"github.com/dustinspecker/ghostdog/internal/dag"
	"github.com/dustinspecker/ghostdog/internal/rule"
)

func RunBuildFile(buildFileName string, buildFileData io.Reader) error {
	thread := &starlark.Thread{Name: "ghostdog-main"}

	rulesDag := dag.NewDag()

	addRule := func(rule rule.Rule) error {
		rulesDag.AddRule(&rule)

		return nil
	}

	nativeFunctions := starlark.StringDict{
		"files": starlark.NewBuiltin("files", builtins.Files(addRule)),
		"rule":  starlark.NewBuiltin("rule", builtins.Rule(addRule)),
	}

	_, err := starlark.ExecFile(thread, buildFileName, buildFileData, nativeFunctions)
	if err != nil {
		return err
	}

	for id, rule := range rulesDag.Rules {
		for _, dependOn := range rule.Sources {
			if err = rulesDag.AddDependency(id, dependOn); err != nil {
				return err
			}
		}
	}

	return nil
}
