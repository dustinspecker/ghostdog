package graph

import (
	"fmt"
	"io"

	"github.com/spf13/afero"

	"github.com/dustinspecker/ghostdog/internal/analyze"
	"github.com/dustinspecker/ghostdog/internal/resolver"
	"github.com/dustinspecker/ghostdog/internal/rule"
)

func GetGraph(fs afero.Fs, cwd, buildTarget string, outputFile io.Writer) error {
	buildFileName, targetRule, err := resolver.GetBuildInfoForPackage(fs, cwd, buildTarget)
	if err != nil {
		return err
	}

	rules, err := analyze.GetRules(fs, buildFileName, targetRule)
	if err != nil {
		return err
	}

	dotGraph := getDotGraph(rules)

	_, err = outputFile.Write([]byte(dotGraph))

	return err
}

func getDotGraph(rules map[string]*rule.Rule) string {
	dagString := "digraph g {\n"

	for _, rule := range rules {
		for _, child := range rule.Children {
			dagString = dagString + fmt.Sprintf("  \"%s\" -> \"%s\";\n", rule.Name, child.Name)
		}
	}

	dagString = dagString + "}"

	return dagString
}
