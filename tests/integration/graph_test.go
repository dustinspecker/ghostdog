// +build integration

package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestGraphExamples(t *testing.T) {
	ghostdogBinaryPath := os.Getenv("GHOSTDOG_BINARY")
	examplesDirectory := os.Getenv("EXAMPLES_DIRECTORY")

	tests := []struct {
		exampleDirectory string
	}{
		{"single"},
	}

	for _, tt := range tests {
		cmd := exec.Cmd{
			Path: ghostdogBinaryPath,
			Args: []string{ghostdogBinaryPath, "graph", "BUILD", "all"},
			Dir:  filepath.Join(examplesDirectory, tt.exampleDirectory),
		}

		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Errorf("%s failed with: %w %s", tt.exampleDirectory, err, string(output))
		}
		t.Log(string(output))
	}
}