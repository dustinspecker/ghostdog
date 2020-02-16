package build

import (
	"strings"
	"testing"
)

func TestRunBuildFileDefinesRuleFunction(t *testing.T) {
	data := `
print('hello')
rule(name="test", sources=[], commands=["echo hey"], outputs=[])
rule(name="publish", sources=["test"], commands=["echo bye"], outputs=[])
`
	if err := RunBuildFile("BUILD", strings.NewReader(data)); err != nil {
		t.Fatalf("expected `rule` function to work: %w", err)
	}
}
