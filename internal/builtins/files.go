package builtins

import (
	"go.starlark.net/starlark"

	"github.com/dustinspecker/ghostdog/internal/rule"
	rulesUtils "github.com/dustinspecker/ghostdog/internal/utils/rules"
	starlarkUtils "github.com/dustinspecker/ghostdog/internal/utils/starlark"
)

func Files(addRule func(rule rule.Rule) error) func(thread *starlark.Thread, builtin *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {

	return func(thread *starlark.Thread, builtin *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var name string
		var paths *starlark.List
		if err := starlark.UnpackArgs(builtin.Name(), args, kwargs, "name", &name, "paths", &paths); err != nil {
			return nil, err
		}

		if err := rulesUtils.ValidateName(name); err != nil {
			return nil, err
		}

		convertedPaths := starlarkUtils.GetStringSlice(*paths)

		if err := rulesUtils.ValidatePaths(convertedPaths); err != nil {
			return nil, err
		}

		newRule := rule.Rule{
			Name:    name,
			Outputs: convertedPaths,
		}
		if err := addRule(newRule); err != nil {
			return nil, err
		}

		return starlark.None, nil
	}
}
