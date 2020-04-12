package graph

import (
	"fmt"
	"io"
	"sort"

	"github.com/apex/log"

	"github.com/dustinspecker/ghostdog/internal/analyze"
	"github.com/dustinspecker/ghostdog/internal/config"
	"github.com/dustinspecker/ghostdog/internal/resolver"
	"github.com/dustinspecker/ghostdog/internal/rule"
)

func GetGraph(config config.Config, buildTarget string, outputFile io.Writer) error {
	buildFileName, targetRule, err := resolver.GetBuildInfoForPackage(config.Fs, config.WorkingDirectory, buildTarget)
	if err != nil {
		config.LogCtx.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("getting build info")
		return err
	}

	config.LogCtx.WithFields(log.Fields{
		"buildFile":  buildFileName,
		"targetRule": targetRule,
	}).Info("build info")

	rules, err := analyze.GetRules(config.LogCtx, config.Fs, buildFileName, targetRule)
	if err != nil {
		return err
	}

	dotGraph := getDotGraph(rules)

	defer config.LogCtx.Trace("writing graph output").Stop(&err)
	_, err = outputFile.Write([]byte(dotGraph))

	return err
}

func getDotGraph(rules []*rule.Rule) string {
	dagString := "digraph g {\n"

	for _, rule := range rules {
		dagString = dagString + getEdges(rule)
	}

	dagString = dagString + "}"

	return dagString
}

func getEdges(r *rule.Rule) string {
	edgesString := ""

	sort.Sort(rule.ByName(r.Children))
	for _, child := range r.Children {
		edgesString = edgesString + getEdges(child) + fmt.Sprintf("  \"%s\" -> \"%s\";\n", r.Name, child.Name)

	}

	return edgesString
}
