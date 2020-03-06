package graph

import (
	"fmt"
	"io"

	"github.com/spf13/afero"

	"github.com/dustinspecker/ghostdog/internal/analyze"
	"github.com/dustinspecker/ghostdog/internal/dag"
	"github.com/dustinspecker/ghostdog/internal/rule"
)

func GetGraph(fs afero.Fs, buildFileName, buildTarget string, outputFile io.Writer) error {
	buildFile, err := fs.Open(buildFileName)
	if err != nil {
		return err
	}

	rulesDag, err := analyze.GetRulesDag(buildFileName, buildFile)
	if err != nil {
		return err
	}

	dotGraph, err := getDotGraph(rulesDag, buildTarget)
	if err != nil {
		return err
	}

	_, err = outputFile.Write([]byte(dotGraph))

	return err
}

func getDotGraph(rulesDag dag.Dag, buildTarget string) (string, error) {
	dagString := "digraph g {\n"

	rules := rulesDag.Rules

	if buildTarget != "all" {
		targetRule, ok := rulesDag.Rules[buildTarget]
		if !ok {
			return "", fmt.Errorf("target %s not found", buildTarget)
		}

		rules = map[string]*rule.Rule{
			targetRule.Name: targetRule,
		}
	}

	for _, rule := range rules {
		for _, child := range rule.Children {
			dagString = dagString + fmt.Sprintf("  \"%s\" -> \"%s\";\n", rule.Name, child.Name)
		}
	}

	dagString = dagString + "}"

	return dagString, nil
}
