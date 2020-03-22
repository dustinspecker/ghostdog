package analyze

import (
	"strings"
	"testing"

	"github.com/spf13/afero"
)

func TestGetRules(t *testing.T) {
	fs := afero.NewMemMapFs()

	data := `
rule(name="test", sources=[], commands=["echo hey"], outputs=[])
rule(name="publish", sources=["test"], commands=["echo bye"], outputs=[])
`
	if err := afero.WriteFile(fs, "BUILD", []byte(data), 0644); err != nil {
		t.Fatalf("unexpected error writing BUILD file: %w", err)
	}

	rules, err := GetRules(fs, "BUILD", "all")
	if err != nil {
		t.Fatalf("expected `rule` function to work: %w", err)
	}

	if len(rules) != 1 {
		t.Fatalf("expected rules to only be source rules, but got: %v", rules)
	}

	publishRule, ok := rules["publish"]
	if !ok {
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

func TestGetRulesReturnsErrorWhenBuildFileDoesntExist(t *testing.T) {
	_, err := GetRules(afero.NewMemMapFs(), "BUILD", "all")
	if err == nil {
		t.Fatal("expected an error when BUILD file doesn't exist")
	}

	if !strings.Contains(err.Error(), "BUILD") {
		t.Errorf("expected message to contain BUILD, but got: %s", err.Error())
	}
}

func TestGetRulesReturnsTargetRuleWhenNotAll(t *testing.T) {
	fs := afero.NewMemMapFs()

	data := `
rule(name="build", sources=[], commands=["true"], outputs=[])
rule(name="test", sources=["build"], commands=["echo hey"], outputs=[])
rule(name="publish", sources=["build"], commands=["echo bye"], outputs=[])
`
	if err := afero.WriteFile(fs, "BUILD", []byte(data), 0644); err != nil {
		t.Fatalf("unexpected error while writing BUILD file: %v", err)
	}

	rules, err := GetRules(fs, "BUILD", "publish")
	if err != nil {
		t.Fatalf("expected `rule` function to work: %w", err)
	}

	if _, ok := rules["publish"]; !ok {
		t.Fatalf("expected rules to have publish rule, but got: %v", rules)
	}

	if len(rules) != 1 {
		t.Errorf("expected rules to only contain target rule, but got: %v", rules)
	}
}
func TestGetRulesReturnsErrorWhenItFailsToRunBuildFile(t *testing.T) {
	fs := afero.NewMemMapFs()

	data := `
rule(invalid_args=1)
`
	if err := afero.WriteFile(fs, "BUILD", []byte(data), 0644); err != nil {
		t.Fatalf("unexpected error while writing BUILD file: %v", err)
	}

	_, err := GetRules(fs, "BUILD", "all")
	if err == nil {
		t.Error("should return error if failed to run BUILD file")
	}
}

func TestGetRulesReturnsErrorWhenDuplicateRuleNameFound(t *testing.T) {
	fs := afero.NewMemMapFs()

	data := `
rule(name="test", sources=[], commands=["echo hey"], outputs=[])
rule(name="test", sources=[], commands=["echo hey"], outputs=[])
`

	if err := afero.WriteFile(fs, "BUILD", []byte(data), 0644); err != nil {
		t.Fatalf("unexpected error while writing BUILD file: %v", err)
	}

	_, err := GetRules(fs, "BUILD", "all")
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

	if err := afero.WriteFile(fs, "BUILD", []byte(data), 0644); err != nil {
		t.Fatalf("unexpected error while writing BUILD file: %v", err)
	}

	_, err := GetRules(fs, "BUILD", "all")
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

	if err := afero.WriteFile(fs, "BUILD", []byte(data), 0644); err != nil {
		t.Fatalf("unexpected error while writing BUILD file: %v", err)
	}

	_, err := GetRules(fs, "BUILD", "deploy")
	if err == nil {
		t.Fatal("expected an error when target doesn't exist")
	}

	if !strings.Contains(err.Error(), "target deploy not found") {
		t.Errorf("expected error message to container target deploy not found, but got: %s", err.Error())
	}
}
