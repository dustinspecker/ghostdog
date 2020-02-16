package rule

import (
	"testing"
)

func TestRunCommand(t *testing.T) {
	rule := Rule{
		Commands: []string{"echo hey"},
	}

	if err := rule.RunCommand(); err != nil {
		t.Error("expected rule to run command successfully")
	}
}

func TestRunCommandDoesNothingWhenNoCommandDefined(t *testing.T) {
	rule := Rule{}

	if err := rule.RunCommand(); err != nil {
		t.Error("expected RunCommand to do nothing if Command empty")
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
		Commands: []string{"false"},
	}

	err := rule.RunCommand()
	if err == nil {
		t.Fatalf("expected command to fail")
	}
}
