package builtins

import (
	"errors"

	"go.starlark.net/starlark"

	"github.com/dustinspecker/ghostdog/internal/rule"
	rulesUtils "github.com/dustinspecker/ghostdog/internal/utils/rules"
	starlarkUtils "github.com/dustinspecker/ghostdog/internal/utils/starlark"
)

var (
	ErrNoCommands = errors.New("no commands were provided")
)

func Rule(addRule func(rule rule.Rule) error) func(thread *starlark.Thread, builtin *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return func(thread *starlark.Thread, builtin *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var name string
		var sources *starlark.List
		var commands *starlark.List
		var outputs *starlark.List
		if err := starlark.UnpackArgs(builtin.Name(), args, kwargs, "name", &name, "sources", &sources, "commands", &commands, "outputs", &outputs); err != nil {
			return nil, err
		}

		if err := rulesUtils.ValidateName(name); err != nil {
			return nil, err
		}

		convertedCommands := starlarkUtils.GetStringSlice(*commands)
		convertedSources := starlarkUtils.GetStringSlice(*sources)

		// todo validate sources rule names

		if len(convertedCommands) == 0 {
			return nil, ErrNoCommands
		}

		newRule := rule.Rule{
			Name:     name,
			Sources:  convertedSources,
			Commands: convertedCommands,
			Outputs:  starlarkUtils.GetStringSlice(*outputs),
		}
		if err := addRule(newRule); err != nil {
			return nil, err
		}

		return starlark.None, nil
	}
}
