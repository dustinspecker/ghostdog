package rule

import (
	"os"
	"strings"
	"testing"

	"github.com/spf13/afero"
)

func TestGetHashDirectorySplitsHashesIntoMultipleDirectories(t *testing.T) {
	rule := Rule{
		Commands: []string{"echo hey"},
	}

	cacheDirectory := "ghostdog-cache"
	hashDirectory, err := rule.GetHashDirectory(afero.NewMemMapFs(), cacheDirectory)
	if err != nil {
		t.Fatalf("unexpected error from GetHashDirectory: %s", err)
	}

	hashDirectories := strings.Split(hashDirectory, string(os.PathSeparator))

	if len(hashDirectories) != 3 {
		t.Fatalf("expected rule hash directory to be comprised of 3 directories, but got %s", hashDirectories)
	}

	if hashDirectories[0] != cacheDirectory {
		t.Errorf("expected hashDirectory to start with %s, but got %s", cacheDirectory, hashDirectory)
	}

	if hashDirectories[1] != hashDirectories[2][0:2] {
		t.Errorf("expected hashDirectory to contain 2 character prefix from hash for grouping, but got %s", hashDirectories)
	}
}

func TestGetHashDirectoryUsesRuleCommandsToMakeHash(t *testing.T) {
	rule := Rule{
		Commands: []string{"echo"},
	}

	hashDirectoryEcho, err := rule.GetHashDirectory(afero.NewMemMapFs(), "cache")
	if err != nil {
		t.Fatalf("unexpected error from GetHashDirectory for echo: %s", err)
	}

	rule.Commands = []string{"./script"}
	hashDirectoryScript, err := rule.GetHashDirectory(afero.NewMemMapFs(), "cache")
	if err != nil {
		t.Fatalf("unexpected error from GetHashDirectory script: %s", err)
	}

	if hashDirectoryEcho == hashDirectoryScript {
		t.Errorf("expected GetHashDirectory to use rule command, but got same hash directory: %s", hashDirectoryEcho)
	}
}

func TestGetHashDirectoryUsesRuleOutputsNamesToMakeHash(t *testing.T) {
	rule := Rule{
		Outputs: []string{"out"},
	}

	hashDirectoryOut, err := rule.GetHashDirectory(afero.NewMemMapFs(), "cache")
	if err != nil {
		t.Fatalf("unexpected error from GetHashDirectory out: %s", err)
	}

	rule.Outputs = []string{"exe"}
	hashDirectoryExe, err := rule.GetHashDirectory(afero.NewMemMapFs(), "cache")
	if err != nil {
		t.Fatalf("unexpected error from GetHashDirectory for exe: %s", err)
	}

	if hashDirectoryOut == hashDirectoryExe {
		t.Errorf("expected GetHashDirectory to use rule outputs, but got same hash directory: %s", hashDirectoryOut)
	}
}

func TestGetHashDirectoryUsesChildrensOutputsContentToMakeHash(t *testing.T) {
	memFs := afero.NewMemMapFs()
	if err := afero.WriteFile(memFs, "file", []byte("hey"), 0644); err != nil {
		t.Fatalf("got error while writing file: %s", err)
	}
	if err := afero.WriteFile(memFs, "another_file", []byte("hey"), 0644); err != nil {
		t.Fatalf("got error while writing file: %s", err)
	}

	childRule := Rule{
		Outputs: []string{"file"},
	}

	rule := Rule{
		Children: []*Rule{&childRule},
	}

	hashDirectoryFile, err := rule.GetHashDirectory(memFs, "cache")
	if err != nil {
		t.Fatalf("unexpected error from GetHashDirectory for file: %s", err)
	}

	childRule.Outputs = []string{"another_file"}
	hashDirectoryAnotherFile, err := rule.GetHashDirectory(memFs, "cache")
	if err != nil {
		t.Fatalf("unexpected error from GetHashDirectory for another_file: %s", err)
	}

	if hashDirectoryFile == hashDirectoryAnotherFile {
		t.Errorf("expected GetHashDirectory to use rule's childrens's output's content, but got same hash directory: %s", hashDirectoryAnotherFile)
	}
}

func TestGetHashDirectoryReturnsErrorWhenFailToGetHashForChildrensOutputs(t *testing.T) {
	childRule := Rule{
		Outputs: []string{"file"},
	}

	rule := Rule{
		Children: []*Rule{&childRule},
	}

	_, err := rule.GetHashDirectory(afero.NewMemMapFs(), "cache")
	if err == nil {
		t.Fatalf("expected to get an error when output file is missing")
	}
}

func TestRunCommand(t *testing.T) {
	rule := Rule{
		Commands: []string{"echo hey"},
	}

	if rule.HasRan {
		t.Error("expected rule to not be marked as ran")
	}

	if err := rule.RunCommand(); err != nil {
		t.Error("expected rule to run command successfully")
	}

	if !rule.HasRan {
		t.Error("expected rule to be marked as ran")
	}
}

func TestRunCommandDoesNothingWhenNoCommandDefined(t *testing.T) {
	rule := Rule{}

	if rule.HasRan {
		t.Error("expected rule to not be marked as ran")
	}

	if err := rule.RunCommand(); err != nil {
		t.Error("expected RunCommand to do nothing if Command empty")
	}

	if !rule.HasRan {
		t.Error("expected rule to be marked as ran even when no commands")
	}
}

func TestRunCommandReturnsErrorWhenCommandFailsToBeParsed(t *testing.T) {
	rule := Rule{
		Commands: []string{"echo \"hey"},
	}

	if err := rule.RunCommand(); err == nil {
		t.Error("expected command to fail to parse")
	}
}

func TestRunCommandReturnsErrorWhenCommandReturnsNonZeroExitCode(t *testing.T) {
	rule := Rule{
		Commands: []string{"grep"},
	}

	err := rule.RunCommand()
	if err == nil {
		t.Fatalf("expected command to fail")
	}

	expectedMessage := "command \"grep\" failed:"
	if !strings.Contains(err.Error(), expectedMessage) {
		t.Errorf("expected error message to have %s, but got: %s", expectedMessage, err.Error())
	}

	expectedOutputMessage := "output: usage"
	if !strings.Contains(strings.ToLower(err.Error()), expectedOutputMessage) {
		t.Errorf("expected error message to contain %s, but got: %s", expectedOutputMessage, err.Error())
	}
}
