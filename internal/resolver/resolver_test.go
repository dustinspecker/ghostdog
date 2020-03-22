package resolver

import (
	"strings"
	"testing"

	"github.com/spf13/afero"
)

func TestGetBuildFileForPackage(t *testing.T) {
	fs := afero.NewMemMapFs()

	if err := fs.MkdirAll("/home/ghostdog/foo", 0755); err != nil {
		t.Fatalf("unexpected error while creating /home/ghostdog/foo directory: %w", err)
	}

	if err := afero.WriteFile(fs, "/home/ghostdog/foo/BUILD", []byte(""), 0644); err != nil {
		t.Fatalf("unexpected error while creating /home/ghostdog/foo/BUILD: %w", err)
	}

	buildFilePath, err := GetBuildFileForPackage(fs, "/home/ghostdog", "foo")
	if err != nil {
		t.Fatalf("unexpected error getting build file: %w", err)
	}

	if buildFilePath != "/home/ghostdog/foo/BUILD" {
		t.Errorf("expected buildFilePath to append package path to cwd, but got: %s", buildFilePath)
	}
}

func TestGetBuildFileForPackageReturnsErrorWhenNoBuildFileFound(t *testing.T) {
	_, err := GetBuildFileForPackage(afero.NewMemMapFs(), "/cool/project", "nope")
	if err == nil {
		t.Fatal("expected an error for a build file that doesn't exit")
	}

	expectedMessage := "no BUILD file found in /cool/project/nope"
	if !strings.Contains(err.Error(), expectedMessage) {
		t.Errorf("expected error messsage to contain %s, but got: %s", expectedMessage, err.Error())
	}
}
