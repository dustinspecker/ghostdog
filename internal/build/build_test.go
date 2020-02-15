package build

import (
	"strings"
	"testing"
)

func TestRunBuildFileParsesStarlark(t *testing.T) {
	if err := RunBuildFile("BUILD", strings.NewReader("print('dustin!')")); err != nil {
		t.Fatalf("expected the BUILD file to run: %w", err)
	}
}

func TestRunBuildFileDefinesFilesFunction(t *testing.T) {
	data := `
print('hello')
files(name="test", paths=[])
`
	if err := RunBuildFile("BUILD", strings.NewReader(data)); err != nil {
		t.Fatalf("expected `files` function to work: %w", err)
	}
}

func TestRunBuildFileDefinesRuleFunction(t *testing.T) {
	data := `
print('hello')
rule(name="test", sources=[], commands=["make build"], outputs=[])
`
	if err := RunBuildFile("BUILD", strings.NewReader(data)); err != nil {
		t.Fatalf("expected `rule` function to work: %w", err)
	}
}
