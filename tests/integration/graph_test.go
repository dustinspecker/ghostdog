// +build integration

package integration

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/sebdah/goldie/v2"
)

func TestGraphExamples(t *testing.T) {
	ghostdogBinaryPath := os.Getenv("GHOSTDOG_BINARY")
	examplesDirectory := os.Getenv("EXAMPLES_DIRECTORY")

	tests := []struct {
		exampleDirectory string
	}{
		{"hello-world"},
		{"simple-go"},
		{"simple-go-with-libs"},
	}

	for _, tt := range tests {
		cmd := exec.Cmd{
			Path: ghostdogBinaryPath,
			Args: []string{ghostdogBinaryPath, "graph", tt.exampleDirectory + ":all"},
			Dir:  examplesDirectory,
		}

		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Errorf("%s failed with: %w %s", tt.exampleDirectory, err, string(output))
		}

		g := goldie.New(t, goldie.WithTestNameForDir(true))
		g.Assert(t, strings.ReplaceAll(tt.exampleDirectory, "/", "_"), output)
	}
}
