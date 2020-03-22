package cache

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/spf13/afero"

	"github.com/dustinspecker/ghostdog/internal/rule"
)

func CopyFileToDestination(fs afero.Fs, source, destination string) error {
	fileStat, err := fs.Stat(source)
	if err != nil {
		return err
	}

	if !fileStat.Mode().IsRegular() {
		return fmt.Errorf("expected file %s to be a regular file (not a directory or symlink)", source)
	}

	file, err := fs.Open(source)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := fs.MkdirAll(filepath.Dir(destination), 0755); err != nil {
		return err
	}

	destinationFile, err := fs.Create(destination)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	if _, err := io.Copy(destinationFile, file); err != nil {
		return err
	}

	return nil
}

func CopyOutputsToRuleCache(fs afero.Fs, rule rule.Rule, ruleCacheDirectory string) error {
	// do not copy Files outputs
	if len(rule.Commands) == 0 {
		return nil
	}

	// Always make the ruleCacheDirectory. If a rule doesn't have any outputs
	// then ruleCacheDirectory wouldn't get made otherwise. The existence of
	// a ruleCacheDirectory lets future runs know to now re-run the command.
	if err := fs.MkdirAll(ruleCacheDirectory, 0755); err != nil {
		return err
	}

	for _, output := range rule.Outputs {
		if err := CopyFileToDestination(fs, filepath.Join(rule.WorkingDirectory, output), filepath.Join(ruleCacheDirectory, output)); err != nil {
			return err
		}
	}

	return nil
}

func CopyRuleCacheToOutputs(fs afero.Fs, rule rule.Rule, ruleCacheDirectory string) error {
	// Files do not have a cache directory
	if len(rule.Commands) == 0 {
		return nil
	}

	for _, output := range rule.Outputs {
		if err := CopyFileToDestination(fs, filepath.Join(ruleCacheDirectory, output), filepath.Join(rule.WorkingDirectory, output)); err != nil {
			return err
		}
	}

	return nil
}
