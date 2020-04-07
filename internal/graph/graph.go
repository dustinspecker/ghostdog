package graph

import (
	"fmt"
	"io"

	"github.com/apex/log"
	"github.com/spf13/afero"

	"github.com/dustinspecker/ghostdog/internal/analyze"
	"github.com/dustinspecker/ghostdog/internal/resolver"
	"github.com/dustinspecker/ghostdog/internal/rule"
)

func GetGraph(logCtx *log.Entry, fs afero.Fs, cwd, buildTarget string, outputFile io.Writer) error {
	buildFileName, targetRule, err := resolver.GetBuildInfoForPackage(fs, cwd, buildTarget)
	if err != nil {
		logCtx.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("getting build info")
		return err
	}

	logCtx.WithFields(log.Fields{
		"buildFile":  buildFileName,
		"targetRule": targetRule,
	}).Info("build info")

	rules, err := analyze.GetRules(logCtx, fs, buildFileName, targetRule)
	if err != nil {
		return err
	}

	dotGraph := getDotGraph(rules)

	defer logCtx.Trace("writing graph output").Stop(&err)
	_, err = outputFile.Write([]byte(dotGraph))

	return err
}

func getDotGraph(rules map[string]*rule.Rule) string {
	dagString := "digraph g {\n"

	for _, rule := range rules {
		dagString = dagString + getEdges(rule)
	}

	dagString = dagString + "}"

	return dagString
}

func getEdges(rule *rule.Rule) string {
	edgesString := ""

	for _, child := range rule.Children {
		edgesString = edgesString + getEdges(child) + fmt.Sprintf("  \"%s\" -> \"%s\";\n", rule.Name, child.Name)

	}

	return edgesString
}
