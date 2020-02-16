package rule

import (
	"fmt"
	"os/exec"

	"github.com/google/shlex"
)

type Rule struct {
	Name     string
	Commands []string
	Sources  []string
	Outputs  []string
	Parents  []*Rule
	Children []*Rule
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
