package rule

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/shlex"
	"github.com/spf13/afero"

	"github.com/dustinspecker/ghostdog/internal/hashing"
)

type Rule struct {
	Name             string
	Commands         []string
	Sources          []string
	Outputs          []string
	Parents          []*Rule
	Children         []*Rule
	WorkingDirectory string
}

func (rule Rule) GetHashDirectory(fs afero.Fs, cacheDirectory string) (string, error) {
	childrensOutputFilepaths := []string{}
	for _, child := range rule.Children {
		childrensOutputFilepaths = append(childrensOutputFilepaths, child.Outputs...)
	}

	childrensOutputFilepathsHash, err := hashing.GetHashForFiles(fs, childrensOutputFilepaths)
	if err != nil {
		return "", err
	}

	ruleHash := hashing.GetHashForStrings([]string{
		childrensOutputFilepathsHash,
		strings.Join(rule.Commands, ""),
		strings.Join(rule.Outputs, ""),
	})

	return filepath.Join(cacheDirectory, ruleHash[0:2], ruleHash), nil
}

func (rule Rule) RunCommand() error {
	if !rule.shouldRunCommand() {
		return nil
	}

	for _, command := range rule.Commands {
		splitCommand, err := shlex.Split(command)
		if err != nil {
			return err
		}

		cmd := exec.Command(splitCommand[0], splitCommand[1:]...)

		if err := cmd.Start(); err != nil {
			return err
		}

		if err := cmd.Wait(); err != nil {
			return fmt.Errorf("command \"%s\" failed: %w", command, err)
		}
	}

	return nil
}

func (rule Rule) shouldRunCommand() bool {
	// only File rules should have no command
	if len(rule.Commands) == 0 {
		return false
	}

	return true
}
