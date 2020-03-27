// +build integration

package integration

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestBuildExamples(t *testing.T) {
	ghostdogBinaryPath := os.Getenv("GHOSTDOG_BINARY")
	examplesDirectory := os.Getenv("EXAMPLES_DIRECTORY")

	tempDir, err := ioutil.TempDir("", "ghostdog-cache")
	if err != nil {
		t.Fatalf("unexpected error while creating tempDir: %w", err)
	}

	t.Logf("using %s for cache directory", tempDir)

	tests := []struct {
		exampleDirectory         string
		expectedCacheDirectories []string
	}{
		{"single", []string{
			filepath.Join("0f", "0f0963f849e1bd3b7197c29b66d59a531c1806e22f9d86478142d2917cd46038"),
			filepath.Join("30", "308a5d09381b31966c398432ae1542f0b35564243c411dd2f4116b84452d5ef6"),
		}},
	}

	for _, tt := range tests {
		cleanOutput, err := run(ghostdogBinaryPath, tempDir, tt.exampleDirectory, examplesDirectory)
		if err != nil {
			t.Errorf("%s failed with: %w %s", tt.exampleDirectory, err, cleanOutput)
		}

		if strings.Contains(cleanOutput, "skipping") {
			t.Errorf("expected no targetrs to be skipped on a clean run, but got %s", cleanOutput)
		}

		// verify cache is used on sequential runs
		// verify working directory and packagepath change doesn't use different cache
		cacheOutput, err := run(ghostdogBinaryPath, tempDir, ".", filepath.Join(examplesDirectory, tt.exampleDirectory))
		if err != nil {
			t.Errorf("%s failed with: %w %s", tt.exampleDirectory, err, cacheOutput)
		}

		if !strings.Contains(cacheOutput, "skipping") {
			t.Errorf("expected some targets to be skipped on a cache run, but got %s", cacheOutput)
		}

		baseCacheDirectory := filepath.Join(tempDir, "ghostdog")
		dirNames := []string{}
		if err := filepath.Walk(baseCacheDirectory, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			dirName := strings.Replace(path, baseCacheDirectory+string(filepath.Separator), "", 1)
			if len(dirName) == 2 || path == baseCacheDirectory || !info.IsDir() {
				return nil
			}

			dirNames = append(dirNames, dirName)
			return filepath.SkipDir
		}); err != nil {
			t.Fatalf("unexpected error while walking: %w", err)
		}

		if !reflect.DeepEqual(tt.expectedCacheDirectories, dirNames) {
			t.Errorf("expected cache directories to be %s, but got %s", tt.expectedCacheDirectories, dirNames)
		}
	}
}

func run(ghostdogBinaryPath, cacheDirectory, packagePath, workingDirectory string) (string, error) {
	cmd := exec.Cmd{
		Path: ghostdogBinaryPath,
		Args: []string{ghostdogBinaryPath, "build", "--cache-directory", cacheDirectory, packagePath + ":all"},
		Dir:  workingDirectory,
	}

	output, err := cmd.CombinedOutput()

	return string(output), err
}
