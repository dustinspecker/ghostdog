package analyze

import (
	"fmt"
	"io"

	"go.starlark.net/starlark"

	"github.com/dustinspecker/ghostdog/internal/builtins"
	"github.com/dustinspecker/ghostdog/internal/dag"
	"github.com/dustinspecker/ghostdog/internal/rule"
)

func GetRules(buildFileName string, buildFileData io.Reader, buildTarget string) (map[string]*rule.Rule, error) {
	thread := &starlark.Thread{Name: "ghostdog-main"}

	rulesDag := dag.NewDag()

	addRule := func(rule rule.Rule) error {
		if err := rulesDag.AddRule(&rule); err != nil {
			return err
		}

		return nil
	}

	nativeFunctions := starlark.StringDict{
		"files": starlark.NewBuiltin("files", builtins.Files(addRule)),
		"rule":  starlark.NewBuiltin("rule", builtins.Rule(addRule)),
	}

	_, err := starlark.ExecFile(thread, buildFileName, buildFileData, nativeFunctions)
	if err != nil {
		return nil, err
	}

	for id, rule := range rulesDag.Rules {
		for _, source := range rule.Sources {
			if err = rulesDag.AddDependency(id, source); err != nil {
				return nil, err
			}
		}
	}

	if buildTarget == "all" {
		return rulesDag.GetSources(), nil
	}

	if _, ok := rulesDag.Rules[buildTarget]; !ok {
		return nil, fmt.Errorf("target %s not found", buildTarget)
	}

	return map[string]*rule.Rule{
		buildTarget: rulesDag.Rules[buildTarget],
	}, nil
}
