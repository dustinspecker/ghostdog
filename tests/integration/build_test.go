// +build integration

package integration

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sebdah/goldie/v2"
)

func TestBuildExamples(t *testing.T) {
	ghostdogBinaryPath := os.Getenv("GHOSTDOG_BINARY")
	examplesDirectory := os.Getenv("EXAMPLES_DIRECTORY")

	tests := []struct {
		exampleDirectory string
	}{
		{"hello-world"},
		// might be flakey as the go compiler isn't being pinnned...
		// TODO: pin go compiler
		{"simple-go"},
		{"simple-go-with-libs"},
	}

	for _, tt := range tests {
		tempDir, err := ioutil.TempDir("", "ghostdog-cache")
		if err != nil {
			t.Fatalf("unexpected error while creating tempDir: %w", err)
		}

		t.Logf("using %s for cache directory", tempDir)

		cleanOutput, err := run(ghostdogBinaryPath, tempDir, tt.exampleDirectory, examplesDirectory)
		if err != nil {
			t.Errorf("%s failed with: %w %s", tt.exampleDirectory, err, cleanOutput)
		}

		parsedCleanOutput := string(cleanOutput)

		if strings.Contains(parsedCleanOutput, "skipping") {
			t.Errorf("expected no targets to be skipped on a clean run, but got %s", parsedCleanOutput)
		}

		g := goldie.New(t, goldie.WithTestNameForDir(true))
		g.Assert(t, strings.ReplaceAll(tt.exampleDirectory, "/", "_")+"-clean", cleanOutput)

		baseCacheDirectory := filepath.Join(tempDir, "ghostdog")
		cleanCacheFileTree, err := getFileTree(baseCacheDirectory)
		if err != nil {
			t.Fatalf("unexpected error while walking after clean run: %w", err)
		}
		g.Assert(t, strings.ReplaceAll(tt.exampleDirectory, "/", "_")+"-cache-tree", []byte(strings.Join(cleanCacheFileTree, "\n")))

		// verify cache is used on sequential runs
		// verify working directory and packagepath change doesn't use different cache
		cacheOutput, err := run(ghostdogBinaryPath, tempDir, ".", filepath.Join(examplesDirectory, tt.exampleDirectory))
		if err != nil {
			t.Errorf("%s failed with: %w %s", tt.exampleDirectory, err, cacheOutput)
		}

		parsedCacheOutput := string(cacheOutput)
		if !strings.Contains(parsedCacheOutput, "skipping") {
			t.Errorf("expected some targets to be skipped on a cache run, but got %s", parsedCacheOutput)
		}

		g.Assert(t, strings.ReplaceAll(tt.exampleDirectory, "/", "_")+"-cache", cacheOutput)

		cacheCacheFileTree, err := getFileTree(baseCacheDirectory)
		if err != nil {
			t.Fatalf("unexpected error while walking after cache run: %w", err)
		}
		g.Assert(t, strings.ReplaceAll(tt.exampleDirectory, "/", "_")+"-cache-tree", []byte(strings.Join(cacheCacheFileTree, "\n")))
	}
}

func run(ghostdogBinaryPath, cacheDirectory, packagePath, workingDirectory string) ([]byte, error) {
	cmd := exec.Cmd{
		Path: ghostdogBinaryPath,
		Args: []string{ghostdogBinaryPath, "build", "--cache-directory", cacheDirectory, packagePath + ":all"},
		Dir:  workingDirectory,
	}

	output, err := cmd.CombinedOutput()

	return output, err
}

func getFileTree(baseDir string) ([]string, error) {
	files := []string{}

	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == baseDir {
			return nil
		}

		files = append(files, strings.Replace(path, baseDir, "", 1))

		return nil
	})

	return files, err
}
