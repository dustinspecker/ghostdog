package build

import (
	"strings"
	"testing"

	"github.com/apex/log"
	"github.com/spf13/afero"

	"github.com/dustinspecker/ghostdog/internal/config"
)

func TestRunBuildFileDefinesRuleFunction(t *testing.T) {
	testConfig := config.NewTest()

	data := `
print('hello')
rule(name="test", sources=[], commands=["echo hey"], outputs=[])
rule(name="publish", sources=["test"], commands=["echo bye"], outputs=[])
`

	if err := afero.WriteFile(testConfig.Config.Fs, "build.ghostdog", []byte(data), 0644); err != nil {
		t.Fatalf("unexpected error while writing build.ghostdog file: %w", err)
	}

	if err := RunBuildFile(testConfig.Config, ".:all", "cache-dir"); err != nil {
		t.Fatalf("expected `rule` function to work: %w", err)
	}
}

func TestRunBuildFileRunsSpecificTargetWhenNotAll(t *testing.T) {
	testConfig := config.NewTest()

	data := `
rule(name="pass", sources=[], commands=["true"], outputs=[])
rule(name="fail", sources=[], commands=["false"], outputs=[])
`

	if err := afero.WriteFile(testConfig.Config.Fs, "build.ghostdog", []byte(data), 0644); err != nil {
		t.Fatalf("unexpected error while writing build.ghostdog file: %w", err)
	}

	err := RunBuildFile(testConfig.Config, ".:pass", "cache-dir")
	if err != nil {
		t.Fatalf("expected build to only run pass rule, but failed: %w", err)
	}
}

func TestRunBuildFileReturnsErrorWhenBuildFileDoesntExist(t *testing.T) {
	testConfig := config.NewTest()
	if err := RunBuildFile(testConfig.Config, ".:all", "cache-dir"); err == nil {
		t.Fatal("expected an error when build.ghostdog file didn't exist")
	}

	if !testConfig.HasLogEntry(log.ErrorLevel, log.Fields{"error": ""}, "getting build info") {
		t.Error("expected error log with error message")
	}
}

func TestRunBuildFileReturnsErrorWhenFailToBuildRulesDag(t *testing.T) {
	testConfig := config.NewTest()

	data := `
doesnt_exist()
`

	if err := afero.WriteFile(testConfig.Config.Fs, "build.ghostdog", []byte(data), 0644); err != nil {
		t.Fatalf("unexpected error while writing build.ghostdog file: %w", err)
	}

	err := RunBuildFile(testConfig.Config, ".:all", "cache-dir")
	if err == nil {
		t.Fatal("expected to fail to build dag")
	}
}

func TestRunBuildFileReturnsErrorWhenATargetDoesntExist(t *testing.T) {
	testConfig := config.NewTest()

	if err := afero.WriteFile(testConfig.Config.Fs, "build.ghostdog", []byte(""), 0644); err != nil {
		t.Fatalf("unexpected error while writing build.ghostdog file: %w", err)
	}

	err := RunBuildFile(testConfig.Config, ".:pass", "cache-dir")
	if err == nil {
		t.Fatal("expected an error when target not found")
	}

	if !strings.Contains(err.Error(), "target pass not found") {
		t.Errorf("expected error message to container target pass not found, but got: %s", err.Error())
	}
}

func TestRunBuildFileReturnsErrorWhenACommandReturnsNonZeroExit(t *testing.T) {
	testConfig := config.NewTest()

	data := `
rule(name="test", sources=[], commands=["false"], outputs=[])
`

	if err := afero.WriteFile(testConfig.Config.Fs, "build.ghostdog", []byte(data), 0644); err != nil {
		t.Fatalf("unexpected error while writing build.ghostdog file: %w", err)
	}

	err := RunBuildFile(testConfig.Config, ".:all", "cache-dir")
	if err == nil {
		t.Fatal("expected test command to fail")
	}
}
