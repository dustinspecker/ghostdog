package analyze

import (
	"strings"
	"testing"

	"github.com/apex/log"
	"github.com/spf13/afero"
)

var testLogCtx = log.WithFields(log.Fields{
	"testPath": "internal/analyze/analyze_test.go",
})

func TestGetRules(t *testing.T) {
	fs := afero.NewMemMapFs()

	data := `
rule(name="test", sources=[], commands=["echo hey"], outputs=[])
rule(name="publish", sources=["test"], commands=["echo bye"], outputs=[])
`
	if err := afero.WriteFile(fs, "build.ghostdog", []byte(data), 0644); err != nil {
		t.Fatalf("unexpected error writing build.ghostdog file: %s", err)
	}

	rules, err := GetRules(testLogCtx, fs, "build.ghostdog", "all")
	if err != nil {
		t.Fatalf("expected `rule` function to work: %s", err)
	}

	if len(rules) != 1 {
		t.Fatalf("expected rules to only be source rules, but got: %v", rules)
	}

	publishRule := rules[0]
	if publishRule.Name != "publish" {
		t.Fatal("expected rulesDag to have a rule with publish id")
	}

	if len(publishRule.Children) != 1 {
		t.Fatalf("expected publishRule to only have 1 child, but got: %v", publishRule.Children)
	}

	if publishRule.Children[0].Name != "test" {
		t.Errorf("expected publish's children to have test, but had %v", publishRule.Children)
	}

	if publishRule.WorkingDirectory != "." {
		t.Errorf("expected publish's WorkingDirectory to be ., but got %s", publishRule.WorkingDirectory)
	}
}

func TestGetRulesSupportsStarlarksLoad(t *testing.T) {
	fs := afero.NewMemMapFs()

	if err := fs.MkdirAll("libs/function", 0755); err != nil {
		t.Fatalf("unexpected while creating libs/function directory: %s", err)
	}

	constantsData := `
MAKEFILE_NAME = "Makefile"
`
	if err := afero.WriteFile(fs, "libs/constants.ghostdog", []byte(constantsData), 0644); err != nil {
		t.Fatalf("unexpected error writing libs/constants.ghostdog file: %s", err)
	}

	libData := `
load("constants.ghostdog", "MAKEFILE_NAME")

def test():
  files(name="make", paths=[MAKEFILE_NAME])
  rule(name="test", sources=[], commands=["echo hey"], outputs=[])
`
	if err := afero.WriteFile(fs, "libs/test.ghostdog", []byte(libData), 0644); err != nil {
		t.Fatalf("unexpected error writing libs/test.ghostdog file: %s", err)
	}

	buildData := `
load("libs/test.ghostdog", "test")
test()
`
	if err := afero.WriteFile(fs, "build.ghostdog", []byte(buildData), 0644); err != nil {
		t.Fatalf("unexpected error writing build.ghostdog file: %s", err)
	}

	rules, err := GetRules(testLogCtx, fs, "build.ghostdog", "all")
	if err != nil {
		t.Fatalf("expected `load` function to work and load modules relative to moduel calling load: %s", err)
	}

	if len(rules) != 2 {
		t.Fatalf("expected only 2 rules to be created from libs/test.ghostdog's test(): %v", rules)
	}
}

func TestGetRulesReturnsErrorWhenLoadFunctionCalledOnNonexistentFile(t *testing.T) {
	fs := afero.NewMemMapFs()

	buildData := `
load("libs/test.ghostdog", "test")
`
	if err := afero.WriteFile(fs, "build.ghostdog", []byte(buildData), 0644); err != nil {
		t.Fatalf("unexpected error writing build.ghostdog file: %s", err)
	}

	_, err := GetRules(testLogCtx, fs, "build.ghostdog", "all")
	if err == nil {
		t.Fatal("expected `load` function to return error when file not found")
	}

	if !strings.Contains(err.Error(), "cannot load libs/test.ghostdog") {
		t.Error("expected error message to contain, but got: %w", err)
	}
}

func TestGetRulesReturnsErrorWhenBuildFileDoesntExist(t *testing.T) {
	_, err := GetRules(testLogCtx, afero.NewMemMapFs(), "build.ghostdog", "all")
	if err == nil {
		t.Fatal("expected an error when build.ghostdog file doesn't exist")
	}

	if !strings.Contains(err.Error(), "build.ghostdog") {
		t.Errorf("expected message to contain build.ghostdog, but got: %s", err.Error())
	}
}

func TestGetRulesReturnsTargetRuleWhenNotAll(t *testing.T) {
	fs := afero.NewMemMapFs()

	data := `
rule(name="build", sources=[], commands=["true"], outputs=[])
rule(name="test", sources=["build"], commands=["echo hey"], outputs=[])
rule(name="publish", sources=["build"], commands=["echo bye"], outputs=[])
`
	if err := afero.WriteFile(fs, "build.ghostdog", []byte(data), 0644); err != nil {
		t.Fatalf("unexpected error while writing build.ghostdog file: %v", err)
	}

	rules, err := GetRules(testLogCtx, fs, "build.ghostdog", "publish")
	if err != nil {
		t.Fatalf("expected `rule` function to work: %s", err)
	}

	if len(rules) != 1 {
		t.Errorf("expected rules to only contain target rule, but got: %v", rules)
	}

	publishRule := rules[0]
	if publishRule.Name != "publish" {
		t.Fatalf("expected rules to have publish rule, but got: %v", rules)
	}

}
func TestGetRulesReturnsErrorWhenItFailsToRunBuildFile(t *testing.T) {
	fs := afero.NewMemMapFs()

	data := `
rule(invalid_args=1)
`
	if err := afero.WriteFile(fs, "build.ghostdog", []byte(data), 0644); err != nil {
		t.Fatalf("unexpected error while writing build.ghostdog file: %v", err)
	}

	_, err := GetRules(testLogCtx, fs, "build.ghostdog", "all")
	if err == nil {
		t.Error("should return error if failed to run build.ghostdog file")
	}
}

func TestGetRulesReturnsErrorWhenDuplicateRuleNameFound(t *testing.T) {
	fs := afero.NewMemMapFs()

	data := `
rule(name="test", sources=[], commands=["echo hey"], outputs=[])
rule(name="test", sources=[], commands=["echo hey"], outputs=[])
`

	if err := afero.WriteFile(fs, "build.ghostdog", []byte(data), 0644); err != nil {
		t.Fatalf("unexpected error while writing build.ghostdog file: %v", err)
	}

	_, err := GetRules(testLogCtx, fs, "build.ghostdog", "all")
	if err == nil {
		t.Error("should return error when duplicate rule name is found")
	}

	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("error should mention already exists: %v", err)
	}
}

func TestGetRulesReturnsErrorWhenSourceDoesntExist(t *testing.T) {
	fs := afero.NewMemMapFs()

	data := `
rule(name="test", sources=["build"], commands=["echo hey"], outputs=[])
`

	if err := afero.WriteFile(fs, "build.ghostdog", []byte(data), 0644); err != nil {
		t.Fatalf("unexpected error while writing build.ghostdog file: %v", err)
	}

	_, err := GetRules(testLogCtx, fs, "build.ghostdog", "all")
	if err == nil {
		t.Error("should return error when rule name is not found")
	}

	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("error should mention not found: %v", err)
	}
}

func TestGetRulesReturnsErrorWhenTargetDoesntExist(t *testing.T) {
	fs := afero.NewMemMapFs()

	data := `
rule(name="build", sources=[], commands=["true"], outputs=[])
rule(name="test", sources=["build"], commands=["echo hey"], outputs=[])
rule(name="publish", sources=["build"], commands=["echo bye"], outputs=[])
`

	if err := afero.WriteFile(fs, "build.ghostdog", []byte(data), 0644); err != nil {
		t.Fatalf("unexpected error while writing build.ghostdog file: %v", err)
	}

	_, err := GetRules(testLogCtx, fs, "build.ghostdog", "deploy")
	if err == nil {
		t.Fatal("expected an error when target doesn't exist")
	}

	if !strings.Contains(err.Error(), "target deploy not found") {
		t.Errorf("expected error message to container target deploy not found, but got: %s", err.Error())
	}
}
