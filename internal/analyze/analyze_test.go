package analyze

import (
	"reflect"
	"strings"
	"testing"

	"github.com/dustinspecker/ghostdog/internal/rule"
)

func TestGetRulesDag(t *testing.T) {
	data := `
rule(name="test", sources=[], commands=["echo hey"], outputs=[])
rule(name="publish", sources=["test"], commands=["echo bye"], outputs=[])
`
	rulesDag, err := GetRulesDag("BUILD", strings.NewReader(data))
	if err != nil {
		t.Fatalf("expected `rule` function to work: %w", err)
	}

	testRule, ok := rulesDag.Rules["test"]
	if !ok {
		t.Fatal("expected rulesDag to have a rule with test id")
	}

	if len(testRule.Children) != 0 {
		t.Errorf("expected test's children to be empty, but got %v", testRule.Children)
	}

	publishRule, ok := rulesDag.Rules["publish"]
	if !ok {
		t.Fatal("expected rulesDag to have a rule with publish id")
	}

	if !reflect.DeepEqual(publishRule.Children, []*rule.Rule{testRule}) {
		t.Errorf("expected publish's children to have test, but had %v", publishRule.Children)
	}
}

func TestGetRulesDagReturnsErrorWhenItFailsToRunBuildFile(t *testing.T) {
	data := `
rule(invalid_args=1)
`
	_, err := GetRulesDag("BUILD", strings.NewReader(data))
	if err == nil {
		t.Error("should return error if failed to run BUILD file")
	}
}

func TestGetRulesDagReturnsErrorWhenDuplicateRuleNameFound(t *testing.T) {
	data := `
rule(name="test", sources=[], commands=["echo hey"], outputs=[])
rule(name="test", sources=[], commands=["echo hey"], outputs=[])
`

	_, err := GetRulesDag("BUILD", strings.NewReader(data))
	if err == nil {
		t.Error("should return error when duplicate rule name is found")
	}

	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("error should mention already exists: %v", err)
	}
}

func TestGetRulesDagReturnsErrorWhenSourceDoesntExist(t *testing.T) {
	data := `
rule(name="test", sources=["build"], commands=["echo hey"], outputs=[])
`

	_, err := GetRulesDag("BUILD", strings.NewReader(data))
	if err == nil {
		t.Error("should return error when rule name is not found")
	}

	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("error should mention not found: %v", err)
	}
}
