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
