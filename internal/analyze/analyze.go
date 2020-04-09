package analyze

import (
	"fmt"
	"path/filepath"

	"github.com/apex/log"
	"github.com/spf13/afero"
	"go.starlark.net/starlark"

	"github.com/dustinspecker/ghostdog/internal/builtins"
	"github.com/dustinspecker/ghostdog/internal/dag"
	"github.com/dustinspecker/ghostdog/internal/rule"
)

func GetRules(logCtx *log.Entry, fs afero.Fs, buildFileName string, buildTarget string) ([]*rule.Rule, error) {
	buildFileData, err := fs.Open(buildFileName)
	if err != nil {
		return nil, err
	}

	workingDirectory := filepath.Dir(buildFileName)

	rulesDag := dag.NewDag()

	addRule := func(rule rule.Rule) error {
		logCtx.WithFields(log.Fields{
			"ruleName": rule.Name,
		}).Info("adding rule")

		rule.WorkingDirectory = workingDirectory

		return rulesDag.AddRule(&rule)
	}

	nativeFunctions := starlark.StringDict{
		"files": starlark.NewBuiltin("files", builtins.Files(addRule)),
		"rule":  starlark.NewBuiltin("rule", builtins.Rule(addRule)),
	}

	var load func(directoryOfModuleCallingLoad string) func(thread *starlark.Thread, module string) (starlark.StringDict, error)
	load = func(directoryOfModuleCallingLoad string) func(thread *starlark.Thread, module string) (starlark.StringDict, error) {
		return func(thread *starlark.Thread, module string) (starlark.StringDict, error) {
			moduleFilePath := filepath.Join(directoryOfModuleCallingLoad, module)
			moduleData, err := fs.Open(moduleFilePath)
			if err != nil {
				return nil, err
			}

			loadThread := &starlark.Thread{Name: "load-" + module, Load: load(filepath.Dir(moduleFilePath))}

			return starlark.ExecFile(loadThread, module, moduleData, nativeFunctions)
		}
	}

	thread := &starlark.Thread{Name: "ghostdog-main", Load: load(workingDirectory)}

	_, err = starlark.ExecFile(thread, buildFileName, buildFileData, nativeFunctions)
	if err != nil {
		return nil, err
	}

	for id, rule := range rulesDag.Rules {
		for _, source := range rule.Sources {
			logCtx.WithFields(log.Fields{
				"parent": id,
				"child":  source,
			}).Info("adding rule dependency")

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

	return []*rule.Rule{
		rulesDag.Rules[buildTarget],
	}, nil
}
