package build

import (
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/spf13/afero"

	"github.com/dustinspecker/ghostdog/internal/analyze"
	"github.com/dustinspecker/ghostdog/internal/cache"
	"github.com/dustinspecker/ghostdog/internal/resolver"
	"github.com/dustinspecker/ghostdog/internal/rule"
)

func RunBuildFile(logCtx *log.Entry, fs afero.Fs, cwd, buildTarget string, cacheDirectory string) error {
	buildFileName, targetRule, err := resolver.GetBuildInfoForPackage(fs, cwd, buildTarget)
	if err != nil {
		logCtx.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("getting build info")
		return err
	}

	rules, err := analyze.GetRules(fs, buildFileName, targetRule)
	if err != nil {
		return err
	}

	for _, rule := range rules {
		if err = run(fs, *rule, cacheDirectory); err != nil {
			return err
		}
	}

	return nil
}

func run(fs afero.Fs, rule rule.Rule, cacheDirectory string) error {
	for _, child := range rule.Children {
		run(fs, *child, cacheDirectory)
	}

	ruleCacheDirectory, err := rule.GetHashDirectory(fs, cacheDirectory)
	if err != nil {
		return err
	}

	_, err = os.Stat(ruleCacheDirectory)
	if err == nil {
		fmt.Println("skipping", rule.Name)

		return cache.CopyRuleCacheToOutputs(fs, rule, ruleCacheDirectory)
	}

	if !os.IsNotExist(err) {
		return err
	}

	fmt.Println("running", rule.Name)
	if err := rule.RunCommand(); err != nil {
		return err
	}

	return cache.CopyOutputsToRuleCache(fs, rule, ruleCacheDirectory)
}
