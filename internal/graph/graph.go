package graph

import (
	"fmt"
	"io"

	"github.com/spf13/afero"

	"github.com/dustinspecker/ghostdog/internal/analyze"
	"github.com/dustinspecker/ghostdog/internal/dag"
)

func GetGraph(fs afero.Fs, buildFileName string, outputFile io.Writer) error {
	buildFile, err := fs.Open(buildFileName)
	if err != nil {
		return err
	}

	rulesDag, err := analyze.GetRulesDag(buildFileName, buildFile)
	if err != nil {
		return err
	}

	_, err = outputFile.Write([]byte(getDotGraph(rulesDag)))

	return err
}

func getDotGraph(rulesDag dag.Dag) string {
	dagString := "digraph g {\n"

	for _, rule := range rulesDag.Rules {
		for _, child := range rule.Children {
			dagString = dagString + fmt.Sprintf("  \"%s\" -> \"%s\";\n", rule.Name, child.Name)
		}
	}

	dagString = dagString + "}"

	return dagString
}
