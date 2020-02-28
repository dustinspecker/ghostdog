package graph

import (
	"os"
	"testing"

	"github.com/spf13/afero"
)

func TestGetGraph(t *testing.T) {
	fs := afero.NewMemMapFs()
	data := `
rule(name="test", sources=["build"], commands=["test"], outputs=[])
rule(name="build", sources=[], commands=["build"], outputs=[])
`

	if err := afero.WriteFile(fs, "BUILD", []byte(data), 0644); err != nil {
		t.Fatalf("unexpected error while writing BUILD file: %w", err)
	}

	tempFile, err := afero.TempFile(fs, "./", "temp")
	if err != nil {
		t.Fatalf("unexpected error while getting tempFile: %w", err)
	}

	if err := GetGraph(fs, "BUILD", tempFile); err != nil {
		t.Fatalf("unexpected error while getting graph: %w", err)
	}

	tempFileContent, err := afero.ReadFile(fs, tempFile.Name())
	if err != nil {
		t.Fatalf("unexpected error while reading tempFile: %w", err)
	}

	expectedGraph := "digraph g {\n  \"test\" -> \"build\";\n}"
	if string(tempFileContent) != expectedGraph {
		t.Errorf("expected tempFile content to be %s, but got: %s", expectedGraph, tempFileContent)
	}
}

func TestGetGraphReturnsReturnsEmptyDigraphWhenNoRules(t *testing.T) {
	fs := afero.NewMemMapFs()

	if err := afero.WriteFile(fs, "BUILD", []byte(""), 0644); err != nil {
		t.Fatalf("unexpected error while writing BUILD file: %w", err)
	}

	tempFile, err := afero.TempFile(fs, "./", "temp")
	if err != nil {
		t.Fatalf("unexpected error while getting tempFile: %w", err)
	}

	if err := GetGraph(fs, "BUILD", tempFile); err != nil {
		t.Fatalf("unexpected error while getting graph: %w", err)
	}

	tempFileContent, err := afero.ReadFile(fs, tempFile.Name())
	if err != nil {
		t.Fatalf("unexpected error while reading tempFile: %w", err)
	}

	expectedGraph := "digraph g {\n}"
	if string(tempFileContent) != expectedGraph {
		t.Errorf("expected tempFile content to be %s, but got: %s", expectedGraph, tempFileContent)
	}
}

func TestGetGraphReturnsErrorWhenBuildFileDoesntExist(t *testing.T) {
	if err := GetGraph(afero.NewMemMapFs(), "BUILD", &os.File{}); err == nil {
		t.Error("expected an error when BUILD file doesn't exist")
	}
}

func TestGetGraphReturnsErrorWhenRulesDagCantBeBuilt(t *testing.T) {
	fs := afero.NewMemMapFs()

	if err := afero.WriteFile(fs, "BUILD", []byte("not valid"), 0644); err != nil {
		t.Fatal("unexpected error while writing BUILD file: %w", err)
	}

	if err := GetGraph(fs, "BUILD", &os.File{}); err == nil {
		t.Error("expected an error when rules dag couldn't be built")
	}
}
