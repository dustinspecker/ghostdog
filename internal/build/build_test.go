package build

import (
	"strings"
	"testing"

	"github.com/spf13/afero"
)

func TestRunBuildFileDefinesRuleFunction(t *testing.T) {
	data := `
print('hello')
rule(name="test", sources=[], commands=["echo hey"], outputs=[])
rule(name="publish", sources=["test"], commands=["echo bye"], outputs=[])
`
	if err := RunBuildFile(afero.NewMemMapFs(), "BUILD", strings.NewReader(data), "cache-dir"); err != nil {
		t.Fatalf("expected `rule` function to work: %w", err)
	}
}
