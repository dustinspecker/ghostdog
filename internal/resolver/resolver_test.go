package resolver

import (
	"strings"
	"testing"

	"github.com/spf13/afero"
)

func TestGetBuildInfoForPackage(t *testing.T) {
	fs := afero.NewMemMapFs()

	if err := fs.MkdirAll("/home/ghostdog/foo", 0755); err != nil {
		t.Fatalf("unexpected error while creating /home/ghostdog/foo directory: %w", err)
	}

	if err := afero.WriteFile(fs, "/home/ghostdog/foo/BUILD", []byte(""), 0644); err != nil {
		t.Fatalf("unexpected error while creating /home/ghostdog/foo/BUILD: %w", err)
	}

	tests := []struct {
		workingDirectory      string
		buildTarget           string
		expectedBuildFilePath string
		expectedTargetRule    string
	}{
		{"/home/ghostdog", "foo", "/home/ghostdog/foo/BUILD", "all"},
		{"/home/ghostdog", "foo:bar", "/home/ghostdog/foo/BUILD", "bar"},
		{"/home/ghostdog/foo", ":bar", "/home/ghostdog/foo/BUILD", "bar"},
		{"/home/ghostdog/foo", "", "/home/ghostdog/foo/BUILD", "all"},
	}

	for _, tt := range tests {
		buildFilePath, targetRule, err := GetBuildInfoForPackage(fs, tt.workingDirectory, tt.buildTarget)
		if err != nil {
			t.Fatalf("unexpected error getting build file: %w", err)
		}

		if buildFilePath != tt.expectedBuildFilePath {
			t.Errorf("expected buildFilePath to append package path to cwd, but got: %s", buildFilePath)
		}

		if targetRule != tt.expectedTargetRule {
			t.Errorf("expected default target rule to be %s, but got %s", tt.expectedTargetRule, targetRule)
		}
	}
}

func TestGetBuildInfoForPackageReturnsErrorWhenNoBuildFileFound(t *testing.T) {
	_, _, err := GetBuildInfoForPackage(afero.NewMemMapFs(), "/cool/project", "nope")
	if err == nil {
		t.Fatal("expected an error for a build file that doesn't exit")
	}

	expectedMessage := "no BUILD file found in /cool/project/nope"
	if !strings.Contains(err.Error(), expectedMessage) {
		t.Errorf("expected error messsage to contain %s, but got: %s", expectedMessage, err.Error())
	}
}

func TestGetBuildInfoForPackageReturnsErrorWhenInvalidTarget(t *testing.T) {
	_, _, err := GetBuildInfoForPackage(afero.NewMemMapFs(), "/cool/project", "nope:hey:bye")
	if err == nil {
		t.Fatal("expected an error for a build file that doesn't exit")
	}

	expectedMessage := "nope:hey:bye is an invalid target"
	if !strings.Contains(err.Error(), expectedMessage) {
		t.Errorf("expected error messsage to contain %s, but got: %s", expectedMessage, err.Error())
	}
}
