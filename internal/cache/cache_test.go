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
		t.Fatalf("got error while creating cache: %s", err)
	}

	if err := fs.MkdirAll(filepath.Join("build", "output"), 0755); err != nil {
		t.Fatalf("got error while creating build/output/dir: %s", err)
	}

	if err := afero.WriteFile(fs, filepath.Join("build", "output", "file"), []byte("hey"), 0644); err != nil {
		t.Fatalf("got error writing build/output/file: %s", err)
	}

	if err := CopyFileToDestination(fs, filepath.Join("build", "output", "file"), filepath.Join("cache", "build", "output", "file")); err != nil {
		t.Fatalf("failed to copy files: %s", err)
	}

	if _, err := fs.Stat(filepath.Join("cache/build/output/file")); err != nil {
		t.Errorf("expected build/output/file to be copied: %s", err)
	}
}

func TestCopyFileToDestinationReturnsErrorWhenCreatingDestinationDirectoryFails(t *testing.T) {
	fs := afero.NewMemMapFs()

	if err := fs.MkdirAll("cache", 0755); err != nil {
		t.Fatalf("got error while creating cache: %s", err)
	}

	if err := fs.MkdirAll(filepath.Join("build", "output"), 0755); err != nil {
		t.Fatalf("got error while creating build/output/dir: %s", err)
	}

	if err := afero.WriteFile(fs, filepath.Join("build", "output", "file"), []byte("hey"), 0644); err != nil {
		t.Fatalf("got error writing build/output/file: %s", err)
	}

	if err := CopyFileToDestination(afero.NewReadOnlyFs(fs), filepath.Join("build", "output", "file"), filepath.Join("cache", "build", "output", "file")); err == nil {
		t.Error("expected an error creating destination directory")
	}
}

func TestCopyFileToDestinationReturnsErrorWhenFilepathDoesntExist(t *testing.T) {
	fs := afero.NewMemMapFs()

	err := CopyFileToDestination(fs, "some_file", "cache")
	if err == nil {
		t.Fatalf("expected CopyFileToDestination to return an error when file doesn't exist")
	}

	if !strings.Contains(err.Error(), "some_file") {
		t.Errorf("error message should contain which file didn't exist: %s", err)
	}
}

func TestCopyFilePathsToDirectoryReturnsErrorWhenFilepathIsNotRegularFile(t *testing.T) {
	fs := afero.NewMemMapFs()

	if err := fs.MkdirAll("build", 0755); err != nil {
		t.Fatalf("got an error while creating build: %s", err)
	}

	err := CopyFileToDestination(fs, "build", "cache")
	if err == nil {
		t.Fatalf("expected an error when trying to copy a non-regular file")
	}

	if !strings.Contains(err.Error(), "build") {
		t.Errorf("error message should contain filepath that isn't regular file: %s", err)
	}
}

func TestCopyOutputsToRuleCache(t *testing.T) {
	fs := afero.NewMemMapFs()

	if err := afero.WriteFile(fs, filepath.Join("working", "some", "file"), []byte("hey"), 0644); err != nil {
		t.Fatalf("expected no error creating some/file, but got: %s", err)
	}

	rule := rule.Rule{
		Commands:         []string{"echo hey"},
		Outputs:          []string{"some/file"},
		WorkingDirectory: "working",
	}

	if err := CopyOutputsToRuleCache(fs, rule, "rule-cache"); err != nil {
		t.Fatalf("expected no error when calling CopyOutputsToRuleCache, but got: %s", err)
	}

	if _, err := fs.Stat(filepath.Join("rule-cache", "some", "file")); err != nil {
		t.Errorf("expected rule-cache/some/file to exist: %s", err)
	}
}

func TestCopyOutputsToRuleCacheDoesNothingWhenNoCommands(t *testing.T) {
	if err := CopyOutputsToRuleCache(afero.NewMemMapFs(), rule.Rule{}, "cache"); err != nil {
		t.Errorf("expected no error when rule had zero commands, but got: %s", err)
	}
}

func TestCopyOutputsToRuleCacheReturnsErrorWhenCreatingDirectoryFails(t *testing.T) {
	rule := rule.Rule{
		Commands: []string{"echo hey"},
	}
	err := CopyOutputsToRuleCache(afero.NewReadOnlyFs(afero.NewMemMapFs()), rule, "cache")
	if err == nil {
		t.Error("expected error when rule cache directory is failed to be created")
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
		t.Errorf("expected no error when rule had zero outputs, but got: %s", err)
	}

	if _, err := fs.Stat("cache/rule-cache"); err != nil {
		t.Errorf("expected cache/rule-cache to be created even when no outputs: %s", err)
	}
}

func TestCopyRuleCacheToOutputs(t *testing.T) {
	fs := afero.NewMemMapFs()

	rule := rule.Rule{
		Commands:         []string{"echo hey"},
		Outputs:          []string{"dist/file"},
		WorkingDirectory: "workingdir",
	}

	if err := fs.MkdirAll(filepath.Join("rule-cache", "dist"), 0755); err != nil {
		t.Fatalf("expected no error creating rule-cache/dist: %s", err)
	}

	if err := afero.WriteFile(fs, filepath.Join("rule-cache", "dist", "file"), []byte("hey"), 0644); err != nil {
		t.Fatalf("expected no error creating rule-cache/dist/file: %s", err)
	}

	if err := CopyRuleCacheToOutputs(fs, rule, "rule-cache"); err != nil {
		t.Fatalf("expected no error calling CopyRuleCacheToOutputs: %s", err)
	}

	if _, err := fs.Stat(filepath.Join("workingdir", "dist", "file")); err != nil {
		t.Errorf("expected dist/file to be created: %s", err)
	}
}

func TestCopyRuleCacheToOutputsDoesNothingWhenNoCommands(t *testing.T) {
	if err := CopyRuleCacheToOutputs(afero.NewMemMapFs(), rule.Rule{}, "cache"); err != nil {
		t.Errorf("expected no error when rule had zero commands, but got: %s", err)
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
		t.Errorf("error message should mention file name not found: %s", err)
	}
}
