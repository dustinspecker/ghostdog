package graph

import (
	"os"
	"strings"
	"testing"

	"github.com/apex/log"
	"github.com/spf13/afero"

	"github.com/dustinspecker/ghostdog/internal/config"
)

func TestGetGraph(t *testing.T) {
	testConfig := config.NewTest()
	data := `
rule(name="test", sources=["build"], commands=["test"], outputs=[])
rule(name="build", sources=[], commands=["build"], outputs=[])
`

	if err := afero.WriteFile(testConfig.Config.Fs, "build.ghostdog", []byte(data), 0644); err != nil {
		t.Fatalf("unexpected error while writing build.ghostdog file: %s", err)
	}

	tempFile, err := afero.TempFile(testConfig.Config.Fs, "./", "temp")
	if err != nil {
		t.Fatalf("unexpected error while getting tempFile: %s", err)
	}

	if err := GetGraph(testConfig.Config, ".:all", tempFile); err != nil {
		t.Fatalf("unexpected error while getting graph: %s", err)
	}

	if !testConfig.HasLogEntry(log.InfoLevel, log.Fields{"buildFile": "build.ghostdog", "targetRule": "all"}, "build info") {
		t.Error("expected an info message saying build info")
	}

	tempFileContent, err := afero.ReadFile(testConfig.Config.Fs, tempFile.Name())
	if err != nil {
		t.Fatalf("unexpected error while reading tempFile: %s", err)
	}

	expectedGraph := "digraph g {\n  \"test\" -> \"build\";\n}"
	if string(tempFileContent) != expectedGraph {
		t.Errorf("expected tempFile content to be %s, but got: %s", expectedGraph, tempFileContent)
	}
}

func TestGetGraphOnlyBuildsDependenciesForGivenTarget(t *testing.T) {
	testConfig := config.NewTest()
	data := `
rule(name="test", sources=["build"], commands=["test"], outputs=[])
rule(name="build", sources=[], commands=["build"], outputs=[])
rule(name="publish", sources=["build"], commands=["publish"], outputs=[])
`

	if err := afero.WriteFile(testConfig.Config.Fs, "build.ghostdog", []byte(data), 0644); err != nil {
		t.Fatalf("unexpected error while writing build.ghostdog file: %s", err)
	}

	tempFile, err := afero.TempFile(testConfig.Config.Fs, "./", "temp")
	if err != nil {
		t.Fatalf("unexpected error while getting tempFile: %s", err)
	}

	if err := GetGraph(testConfig.Config, ".:publish", tempFile); err != nil {
		t.Fatalf("unexpected error while getting graph: %s", err)
	}

	tempFileContent, err := afero.ReadFile(testConfig.Config.Fs, tempFile.Name())
	if err != nil {
		t.Fatalf("unexpected error while reading tempFile: %s", err)
	}

	expectedGraph := "digraph g {\n  \"publish\" -> \"build\";\n}"
	if string(tempFileContent) != expectedGraph {
		t.Errorf("expected tempFile content to be %s, but got: %s", expectedGraph, tempFileContent)
	}
}

func TestGetGraphReturnsReturnsEmptyDigraphWhenNoRules(t *testing.T) {
	testConfig := config.NewTest()

	if err := afero.WriteFile(testConfig.Config.Fs, "build.ghostdog", []byte(""), 0644); err != nil {
		t.Fatalf("unexpected error while writing build.ghostdog file: %s", err)
	}

	tempFile, err := afero.TempFile(testConfig.Config.Fs, "./", "temp")
	if err != nil {
		t.Fatalf("unexpected error while getting tempFile: %s", err)
	}

	if err := GetGraph(testConfig.Config, ".:all", tempFile); err != nil {
		t.Fatalf("unexpected error while getting graph: %s", err)
	}

	tempFileContent, err := afero.ReadFile(testConfig.Config.Fs, tempFile.Name())
	if err != nil {
		t.Fatalf("unexpected error while reading tempFile: %s", err)
	}

	expectedGraph := "digraph g {\n}"
	if string(tempFileContent) != expectedGraph {
		t.Errorf("expected tempFile content to be %s, but got: %s", expectedGraph, tempFileContent)
	}
}

func TestGetGraphReturnsErrorWhenBuildFileDoesntExist(t *testing.T) {
	testConfig := config.NewTest()
	if err := GetGraph(testConfig.Config, ".:all", &os.File{}); err == nil {
		t.Error("expected an error when build.ghostdog file doesn't exist")
	}

	if !testConfig.HasLogEntry(log.ErrorLevel, log.Fields{"error": ""}, "getting build info") {
		t.Error("expected error log with error message")
	}
}

func TestGetGraphReturnsReturnsErrorWhenTargetDoesntExist(t *testing.T) {
	testConfig := config.NewTest()

	if err := afero.WriteFile(testConfig.Config.Fs, "build.ghostdog", []byte(""), 0644); err != nil {
		t.Fatalf("unexpected error while writing build.ghostdog file: %s", err)
	}

	err := GetGraph(testConfig.Config, ".:build", &os.File{})
	if err == nil {
		t.Fatal("expected an error when target doesn't exist")
	}

	if !strings.Contains(err.Error(), "target build not found") {
		t.Errorf("expected error message to container target build not found, but got: %s", err.Error())
	}
}

func TestGetGraphReturnsErrorWhenRulesDagCantBeBuilt(t *testing.T) {
	testConfig := config.NewTest()

	if err := afero.WriteFile(testConfig.Config.Fs, "build.ghostdog", []byte("not valid"), 0644); err != nil {
		t.Fatal("unexpected error while writing build.ghostdog file: %w", err)
	}

	if err := GetGraph(testConfig.Config, ".:all", &os.File{}); err == nil {
		t.Error("expected an error when rules dag couldn't be built")
	}
}
