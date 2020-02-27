package build

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/afero"

	"github.com/dustinspecker/ghostdog/internal/analyze"
	"github.com/dustinspecker/ghostdog/internal/cache"
	"github.com/dustinspecker/ghostdog/internal/rule"
)

func RunBuildFile(fs afero.Fs, buildFileName string, buildFileData io.Reader, cacheDirectory string) error {
	rulesDag, err := analyze.GetRulesDag(buildFileName, buildFileData)
	if err != nil {
		return err
	}

	for _, rule := range rulesDag.GetSources() {
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
