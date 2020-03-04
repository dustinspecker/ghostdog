package build

import (
	"testing"

	"github.com/spf13/afero"
)

func TestRunBuildFileDefinesRuleFunction(t *testing.T) {
	fs := afero.NewMemMapFs()

	data := `
print('hello')
rule(name="test", sources=[], commands=["echo hey"], outputs=[])
rule(name="publish", sources=["test"], commands=["echo bye"], outputs=[])
`

	if err := afero.WriteFile(fs, "BUILD", []byte(data), 0644); err != nil {
		t.Fatalf("unexpected error while writing BUILD file: %w", err)
	}

	if err := RunBuildFile(fs, "BUILD", "cache-dir"); err != nil {
		t.Fatalf("expected `rule` function to work: %w", err)
	}
}

func TestRunBuildFileReturnsErrorWhenBuildFileDoesntExist(t *testing.T) {
	if err := RunBuildFile(afero.NewMemMapFs(), "BUILD", "cache-dir"); err == nil {
		t.Fatal("expected an error when BUILD file didn't exist")
	}
}

func TestRunBuildFileReturnsErrorWhenFailToBuildRulesDag(t *testing.T) {
	fs := afero.NewMemMapFs()

	data := `
doesnt_exist()
`

	if err := afero.WriteFile(fs, "BUILD", []byte(data), 0644); err != nil {
		t.Fatalf("unexpected error while writing BUILD file: %w", err)
	}

	err := RunBuildFile(fs, "BUILD", "cache-dir")
	if err == nil {
		t.Fatal("expected to fail to build dag")
	}
}

func TestRunBuildFileReturnsErrorWhenACommandReturnsNonZeroExit(t *testing.T) {
	fs := afero.NewMemMapFs()

	data := `
rule(name="test", sources=[], commands=["false"], outputs=[])
`

	if err := afero.WriteFile(fs, "BUILD", []byte(data), 0644); err != nil {
		t.Fatalf("unexpected error while writing BUILD file: %w", err)
	}

	err := RunBuildFile(fs, "BUILD", "cache-dir")
	if err == nil {
		t.Fatal("expected test command to fail")
	}
}
