package build

import (
	"strings"
	"testing"

	"github.com/apex/log"
	"github.com/spf13/afero"
)

var testLogCtx = log.WithFields(log.Fields{
	"testPath": "internal/build/build_test.go",
})

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

	if err := RunBuildFile(testLogCtx, fs, ".", ".:all", "cache-dir"); err != nil {
		t.Fatalf("expected `rule` function to work: %w", err)
	}
}

func TestRunBuildFileRunsSpecificTargetWhenNotAll(t *testing.T) {
	fs := afero.NewMemMapFs()

	data := `
rule(name="pass", sources=[], commands=["true"], outputs=[])
rule(name="fail", sources=[], commands=["false"], outputs=[])
`

	if err := afero.WriteFile(fs, "BUILD", []byte(data), 0644); err != nil {
		t.Fatalf("unexpected error while writing BUILD file: %w", err)
	}

	err := RunBuildFile(testLogCtx, fs, ".", ".:pass", "cache-dir")
	if err != nil {
		t.Fatalf("expected build to only run pass rule, but failed: %w", err)
	}
}

func TestRunBuildFileReturnsErrorWhenBuildFileDoesntExist(t *testing.T) {
	if err := RunBuildFile(testLogCtx, afero.NewMemMapFs(), ".", ".:all", "cache-dir"); err == nil {
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

	err := RunBuildFile(testLogCtx, fs, ".", ".:all", "cache-dir")
	if err == nil {
		t.Fatal("expected to fail to build dag")
	}
}

func TestRunBuildFileReturnsErrorWhenATargetDoesntExist(t *testing.T) {
	fs := afero.NewMemMapFs()

	if err := afero.WriteFile(fs, "BUILD", []byte(""), 0644); err != nil {
		t.Fatalf("unexpected error while writing BUILD file: %w", err)
	}

	err := RunBuildFile(testLogCtx, fs, ".", ".:pass", "cache-dir")
	if err == nil {
		t.Fatal("expected an error when target not found")
	}

	if !strings.Contains(err.Error(), "target pass not found") {
		t.Errorf("expected error message to container target pass not found, but got: %s", err.Error())
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

	err := RunBuildFile(testLogCtx, fs, ".", ".:all", "cache-dir")
	if err == nil {
		t.Fatal("expected test command to fail")
	}
}
