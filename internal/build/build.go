package build

import (
	"io"

	"github.com/dustinspecker/ghostdog/internal/analyze"
	"github.com/dustinspecker/ghostdog/internal/rule"
)

func RunBuildFile(buildFileName string, buildFileData io.Reader) error {
	rulesDag, err := analyze.GetRulesDag(buildFileName, buildFileData)
	if err != nil {
		return err
	}

	for _, rule := range rulesDag.GetSources() {
		if err = run(rule); err != nil {
			return err
		}
	}

	return nil
}

func run(rule *rule.Rule) error {
	for _, child := range rule.Children {
		run(child)
	}

	if err := rule.RunCommand(); err != nil {
		return err
	}

	return nil
}
