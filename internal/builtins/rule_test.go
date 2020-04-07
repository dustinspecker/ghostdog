package builtins

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"go.starlark.net/starlark"

	"github.com/dustinspecker/ghostdog/internal/rule"
	rulesUtils "github.com/dustinspecker/ghostdog/internal/utils/rules"
)

var (
	simpleAddRule = func(rule rule.Rule) error {
		return nil
	}
)

func TestRule(t *testing.T) {
	thread := &starlark.Thread{}
	builtin := &starlark.Builtin{}
	args := []starlark.Value{
		starlark.String("rule_name"),
		starlark.NewList([]starlark.Value{
			starlark.String("source_code"),
		}),
		starlark.NewList([]starlark.Value{
			starlark.String("make build"),
		}),
		starlark.NewList([]starlark.Value{
			starlark.String("build.exe"),
		}),
	}
	kwargs := []starlark.Tuple{}

	addRuleCalled := false

	addRule := func(rule rule.Rule) error {
		addRuleCalled = true

		if rule.Name != "rule_name" {
			t.Errorf("expected rule.Name to be rule_name, but was %s", rule.Name)
		}

		if !reflect.DeepEqual(rule.Sources, []string{"source_code"}) {
			t.Errorf("expected rule to depend on source_code, but got %s", rule.Sources)
		}

		if !reflect.DeepEqual(rule.Commands, []string{"make build"}) {
			t.Errorf("expected Command to be make build, but was %s", rule.Commands)
		}

		if !reflect.DeepEqual(rule.Outputs, []string{"build.exe"}) {
			t.Errorf("expected rule to have build.exe as outputs, but got %s", rule.Outputs)
		}

		return nil
	}
	value, err := Rule(addRule)(thread, builtin, args, kwargs)
	if err != nil {
		t.Fatalf("Files failed: %w", err)
	}

	if value != starlark.None {
		t.Errorf("expected value to be None, but got: %v", value)
	}

	if !addRuleCalled {
		t.Error("expected addRule to be called")
	}
}

func TestRuleReturnsErrorWhenInvalidArgs(t *testing.T) {
	thread := &starlark.Thread{}
	builtin := &starlark.Builtin{}
	args := []starlark.Value{}
	kwargs := []starlark.Tuple{}

	_, err := Rule(simpleAddRule)(thread, builtin, args, kwargs)
	if err == nil {
		t.Fatal("Rule should have returned an error")
	}

	if !strings.Contains(err.Error(), "missing argument for") {
		t.Errorf("expected error message to contain \"missing argument for\", but got %s", err.Error())
	}
}

func TestRuleReturnsErrorWhenNameIsInvalid(t *testing.T) {
	thread := &starlark.Thread{}
	builtin := &starlark.Builtin{}
	args := []starlark.Value{
		starlark.String("Rule_name"),
		starlark.NewList([]starlark.Value{}),
		starlark.NewList([]starlark.Value{}),
		starlark.NewList([]starlark.Value{}),
	}
	kwargs := []starlark.Tuple{}

	_, err := Rule(simpleAddRule)(thread, builtin, args, kwargs)
	if err == nil {
		t.Fatal("expected Rule to return an error")
	}

	if !errors.Is(err, rulesUtils.ErrInvalidName) {
		t.Errorf("expected error to be ErrInvalidName, but got %v", err)
	}
}

func TestRuleReturnsErrorWhenCommandsIsEmpty(t *testing.T) {
	thread := &starlark.Thread{}
	builtin := &starlark.Builtin{}
	args := []starlark.Value{
		starlark.String("rule_name"),
		starlark.NewList([]starlark.Value{}),
		starlark.NewList([]starlark.Value{}),
		starlark.NewList([]starlark.Value{}),
	}
	kwargs := []starlark.Tuple{}

	_, err := Rule(simpleAddRule)(thread, builtin, args, kwargs)
	if err == nil {
		t.Fatal("expected Rule to return an error")
	}

	if !errors.Is(err, ErrNoCommands) {
		t.Errorf("expected error to be ErrNoCommands, but got %v", err)
	}
}

func TestRuleReturnsErrorWhenAddRuleErrs(t *testing.T) {
	thread := &starlark.Thread{}
	builtin := &starlark.Builtin{}
	args := []starlark.Value{
		starlark.String("rule_name"),
		starlark.NewList([]starlark.Value{}),
		starlark.NewList([]starlark.Value{
			starlark.String("make build"),
		}),
		starlark.NewList([]starlark.Value{}),
	}
	kwargs := []starlark.Tuple{}

	addRule := func(rule rule.Rule) error {
		return errors.New("bad stuff")
	}

	_, err := Rule(addRule)(thread, builtin, args, kwargs)
	if err == nil {
		t.Fatal("Files should have failed to call addRule")
	}

	if !strings.Contains(err.Error(), "bad stuff") {
		t.Errorf("expected error message to contain \"bad stuff\", but got %s", err.Error())
	}
}
