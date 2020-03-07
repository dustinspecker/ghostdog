package analyze

import (
	"strings"
	"testing"
)

func TestGetRules(t *testing.T) {
	data := `
rule(name="test", sources=[], commands=["echo hey"], outputs=[])
rule(name="publish", sources=["test"], commands=["echo bye"], outputs=[])
`
	rules, err := GetRules("BUILD", strings.NewReader(data), "all")
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
}

func TestGetRulesReturnsTargetRuleWhenNotAll(t *testing.T) {
	data := `
rule(name="build", sources=[], commands=["true"], outputs=[])
rule(name="test", sources=["build"], commands=["echo hey"], outputs=[])
rule(name="publish", sources=["build"], commands=["echo bye"], outputs=[])
`
	rules, err := GetRules("BUILD", strings.NewReader(data), "publish")
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
	data := `
rule(invalid_args=1)
`
	_, err := GetRules("BUILD", strings.NewReader(data), "all")
	if err == nil {
		t.Error("should return error if failed to run BUILD file")
	}
}

func TestGetRulesReturnsErrorWhenDuplicateRuleNameFound(t *testing.T) {
	data := `
rule(name="test", sources=[], commands=["echo hey"], outputs=[])
rule(name="test", sources=[], commands=["echo hey"], outputs=[])
`

	_, err := GetRules("BUILD", strings.NewReader(data), "all")
	if err == nil {
		t.Error("should return error when duplicate rule name is found")
	}

	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("error should mention already exists: %v", err)
	}
}

func TestGetRulesReturnsErrorWhenSourceDoesntExist(t *testing.T) {
	data := `
rule(name="test", sources=["build"], commands=["echo hey"], outputs=[])
`

	_, err := GetRules("BUILD", strings.NewReader(data), "all")
	if err == nil {
		t.Error("should return error when rule name is not found")
	}

	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("error should mention not found: %v", err)
	}
}

func TestGetRulesReturnsErrorWhenTargetDoesntExist(t *testing.T) {
	data := `
rule(name="build", sources=[], commands=["true"], outputs=[])
rule(name="test", sources=["build"], commands=["echo hey"], outputs=[])
rule(name="publish", sources=["build"], commands=["echo bye"], outputs=[])
`
	_, err := GetRules("BUILD", strings.NewReader(data), "deploy")
	if err == nil {
		t.Fatal("expected an error when target doesn't exist")
	}

	if !strings.Contains(err.Error(), "target deploy not found") {
		t.Errorf("expected error message to container target deploy not found, but got: %s", err.Error())
	}
}
