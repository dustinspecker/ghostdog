package graph

import (
	"os"
	"strings"
	"testing"

	"github.com/spf13/afero"

	"github.com/dustinspecker/ghostdog/internal/config"
)

func TestGetGraph(t *testing.T) {
	config := config.NewTest()
	data := `
rule(name="test", sources=["build"], commands=["test"], outputs=[])
rule(name="build", sources=[], commands=["build"], outputs=[])
`

	if err := afero.WriteFile(config.Fs, "build.ghostdog", []byte(data), 0644); err != nil {
		t.Fatalf("unexpected error while writing build.ghostdog file: %w", err)
	}

	tempFile, err := afero.TempFile(config.Fs, "./", "temp")
	if err != nil {
		t.Fatalf("unexpected error while getting tempFile: %w", err)
	}

	if err := GetGraph(config, ".:all", tempFile); err != nil {
		t.Fatalf("unexpected error while getting graph: %w", err)
	}

	tempFileContent, err := afero.ReadFile(config.Fs, tempFile.Name())
	if err != nil {
		t.Fatalf("unexpected error while reading tempFile: %w", err)
	}

	expectedGraph := "digraph g {\n  \"test\" -> \"build\";\n}"
	if string(tempFileContent) != expectedGraph {
		t.Errorf("expected tempFile content to be %s, but got: %s", expectedGraph, tempFileContent)
	}
}

func TestGetGraphOnlyBuildsDependenciesForGivenTarget(t *testing.T) {
	config := config.NewTest()
	data := `
rule(name="test", sources=["build"], commands=["test"], outputs=[])
rule(name="build", sources=[], commands=["build"], outputs=[])
rule(name="publish", sources=["build"], commands=["publish"], outputs=[])
`

	if err := afero.WriteFile(config.Fs, "build.ghostdog", []byte(data), 0644); err != nil {
		t.Fatalf("unexpected error while writing build.ghostdog file: %w", err)
	}

	tempFile, err := afero.TempFile(config.Fs, "./", "temp")
	if err != nil {
		t.Fatalf("unexpected error while getting tempFile: %w", err)
	}

	if err := GetGraph(config, ".:publish", tempFile); err != nil {
		t.Fatalf("unexpected error while getting graph: %w", err)
	}

	tempFileContent, err := afero.ReadFile(config.Fs, tempFile.Name())
	if err != nil {
		t.Fatalf("unexpected error while reading tempFile: %w", err)
	}

	expectedGraph := "digraph g {\n  \"publish\" -> \"build\";\n}"
	if string(tempFileContent) != expectedGraph {
		t.Errorf("expected tempFile content to be %s, but got: %s", expectedGraph, tempFileContent)
	}
}

func TestGetGraphReturnsReturnsEmptyDigraphWhenNoRules(t *testing.T) {
	config := config.NewTest()

	if err := afero.WriteFile(config.Fs, "build.ghostdog", []byte(""), 0644); err != nil {
		t.Fatalf("unexpected error while writing build.ghostdog file: %w", err)
	}

	tempFile, err := afero.TempFile(config.Fs, "./", "temp")
	if err != nil {
		t.Fatalf("unexpected error while getting tempFile: %w", err)
	}

	if err := GetGraph(config, ".:all", tempFile); err != nil {
		t.Fatalf("unexpected error while getting graph: %w", err)
	}

	tempFileContent, err := afero.ReadFile(config.Fs, tempFile.Name())
	if err != nil {
		t.Fatalf("unexpected error while reading tempFile: %w", err)
	}

	expectedGraph := "digraph g {\n}"
	if string(tempFileContent) != expectedGraph {
		t.Errorf("expected tempFile content to be %s, but got: %s", expectedGraph, tempFileContent)
	}
}

func TestGetGraphReturnsErrorWhenBuildFileDoesntExist(t *testing.T) {
	if err := GetGraph(config.NewTest(), ".:all", &os.File{}); err == nil {
		t.Error("expected an error when build.ghostdog file doesn't exist")
	}
}

func TestGetGraphReturnsReturnsErrorWhenTargetDoesntExist(t *testing.T) {
	config := config.NewTest()

	if err := afero.WriteFile(config.Fs, "build.ghostdog", []byte(""), 0644); err != nil {
		t.Fatalf("unexpected error while writing build.ghostdog file: %w", err)
	}

	err := GetGraph(config, ".:build", &os.File{})
	if err == nil {
		t.Fatal("expected an error when target doesn't exist")
	}

	if !strings.Contains(err.Error(), "target build not found") {
		t.Errorf("expected error message to container target build not found, but got: %s", err.Error())
	}
}

func TestGetGraphReturnsErrorWhenRulesDagCantBeBuilt(t *testing.T) {
	config := config.NewTest()

	if err := afero.WriteFile(config.Fs, "build.ghostdog", []byte("not valid"), 0644); err != nil {
		t.Fatal("unexpected error while writing build.ghostdog file: %w", err)
	}

	if err := GetGraph(config, ".:all", &os.File{}); err == nil {
		t.Error("expected an error when rules dag couldn't be built")
	}
}
