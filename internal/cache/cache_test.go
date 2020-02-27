package cache

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/afero"

	"github.com/dustinspecker/ghostdog/internal/rule"
)

func TestCopyFileToDestination(t *testing.T) {
	fs := afero.NewMemMapFs()

	if err := fs.MkdirAll("cache", 0755); err != nil {
		t.Fatalf("got error while creating cache: %w", err)
	}

	if err := fs.MkdirAll(filepath.Join("build", "output"), 0755); err != nil {
		t.Fatalf("got error while creating build/output/dir: %w", err)
	}

	if err := afero.WriteFile(fs, filepath.Join("build", "output", "file"), []byte("hey"), 0644); err != nil {
		t.Fatalf("got error writing build/output/file: %w", err)
	}

	if err := CopyFileToDestination(fs, filepath.Join("build", "output", "file"), filepath.Join("cache", "build", "output", "file")); err != nil {
		t.Fatalf("failed to copy files: %w", err)
	}

	if _, err := fs.Stat(filepath.Join("cache/build/output/file")); err != nil {
		t.Errorf("expected build/output/file to be copied: %w", err)
	}
}

func TestCopyFileToDestinationReturnsErrorWhenFilepathDoesntExist(t *testing.T) {
	fs := afero.NewMemMapFs()

	err := CopyFileToDestination(fs, "some_file", "cache")
	if err == nil {
		t.Fatalf("expected CopyFileToDestination to return an error when file doesn't exist")
	}

	if !strings.Contains(err.Error(), "some_file") {
		t.Errorf("error message should contain which file didn't exist: %w", err)
	}
}

func TestCopyFilePathsToDirectoryReturnsErrorWhenFilepathIsNotRegularFile(t *testing.T) {
	fs := afero.NewMemMapFs()

	if err := fs.MkdirAll("build", 0755); err != nil {
		t.Fatalf("got an error while creating build: %w", err)
	}

	err := CopyFileToDestination(fs, "build", "cache")
	if err == nil {
		t.Fatalf("expected an error when trying to copy a non-regular file")
	}

	if !strings.Contains(err.Error(), "build") {
		t.Errorf("error message should contain filepath that isn't regular file: %w", err)
	}
}

func TestCopyOutputsToRuleCache(t *testing.T) {
	fs := afero.NewMemMapFs()

	if err := afero.WriteFile(fs, filepath.Join("some", "file"), []byte("hey"), 0644); err != nil {
		t.Fatalf("expected no error creating some/file, but got: %w", err)
	}

	rule := rule.Rule{
		Commands: []string{"echo hey"},
		Outputs:  []string{"some/file"},
	}

	if err := CopyOutputsToRuleCache(fs, rule, "rule-cache"); err != nil {
		t.Fatalf("expected no error when calling CopyOutputsToRuleCache, but got: %w", err)
	}

	if _, err := fs.Stat(filepath.Join("rule-cache", "some", "file")); err != nil {
		t.Errorf("expected rule-cache/some/file to exist: %w", err)
	}
}

func TestCopyOutputsToRuleCacheDoesNothingWhenNoCommands(t *testing.T) {
	if err := CopyOutputsToRuleCache(afero.NewMemMapFs(), rule.Rule{}, "cache"); err != nil {
		t.Errorf("expected no error when rule had zero commands, but got: %w", err)
	}
}

func TestCopyOutputsToRuleCacheReturnsErrorWhenOutputDoesntExist(t *testing.T) {
	rule := rule.Rule{
		Commands: []string{"echo hey"},
		Outputs:  []string{"doesnt_exist"},
	}
	err := CopyOutputsToRuleCache(afero.NewMemMapFs(), rule, "cache")
	if err == nil {
		t.Error("expected error when output doesn't exist")
	}
}

func TestCopyOutputsToRuleCacheCreatesRuleCacheDirectoryEvenWhenZeroOutputs(t *testing.T) {
	fs := afero.NewMemMapFs()

	rule := rule.Rule{
		Commands: []string{"echo hey"},
	}

	if err := CopyOutputsToRuleCache(fs, rule, "cache/rule-cache"); err != nil {
		t.Errorf("expected no error when rule had zero outputs, but got: %w", err)
	}

	if _, err := fs.Stat("cache/rule-cache"); err != nil {
		t.Errorf("expected cache/rule-cache to be created even when no outputs: %w", err)
	}
}

func TestCopyRuleCacheToOutputs(t *testing.T) {
	fs := afero.NewMemMapFs()

	rule := rule.Rule{
		Commands: []string{"echo hey"},
		Outputs:  []string{"dist/file"},
	}

	if err := fs.MkdirAll(filepath.Join("rule-cache", "dist"), 0755); err != nil {
		t.Fatalf("expected no error creating rule-cache/dist: %w", err)
	}

	if err := afero.WriteFile(fs, filepath.Join("rule-cache", "dist", "file"), []byte("hey"), 0644); err != nil {
		t.Fatalf("expected no error creating rule-cache/dist/file: %w", err)
	}

	if err := CopyRuleCacheToOutputs(fs, rule, "rule-cache"); err != nil {
		t.Fatalf("expected no error calling CopyRuleCacheToOutputs: %w", err)
	}

	if _, err := fs.Stat(filepath.Join("dist", "file")); err != nil {
		t.Errorf("expected dist/file to be created: %w", err)
	}
}

func TestCopyRuleCacheToOutputsDoesNothingWhenNoCommands(t *testing.T) {
	if err := CopyRuleCacheToOutputs(afero.NewMemMapFs(), rule.Rule{}, "cache"); err != nil {
		t.Errorf("expected no error when rule had zero commands, but got: %w", err)
	}
}

func TestCopyRuleCacheToOutputsReturnsErrorWhenRuleCacheOutputDoesntExist(t *testing.T) {
	fs := afero.NewMemMapFs()

	rule := rule.Rule{
		Commands: []string{"echo hey"},
		Outputs:  []string{"file"},
	}

	err := CopyRuleCacheToOutputs(fs, rule, "cache")
	if err == nil {
		t.Error("expected error when rule cache doesn't exist")
	}

	if !strings.Contains(err.Error(), "cache/file") {
		t.Errorf("error message should mention file name not found: %w", err)
	}
}
