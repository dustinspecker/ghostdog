package builtins

import (
	"errors"
	"reflect"
	"testing"

	"go.starlark.net/starlark"

	"github.com/dustinspecker/ghostdog/internal/rule"
)

var (
	simpleAddFileRule = func(rule rule.Rule) error {
		return nil
	}
)

func TestFiles(t *testing.T) {
	thread := &starlark.Thread{}
	builtin := &starlark.Builtin{}
	args := []starlark.Value{
		starlark.String("rule_name"),
		starlark.NewList([]starlark.Value{
			starlark.String("Makefile"),
		}),
	}
	kwargs := []starlark.Tuple{}

	addRuleCalled := false

	addRule := func(rule rule.Rule) error {
		addRuleCalled = true

		if rule.Name != "rule_name" {
			t.Errorf("expeceted rule.Name to be rule_name, but was %s", rule.Name)
		}

		if !reflect.DeepEqual(rule.Outputs, []string{"Makefile"}) {
			t.Errorf("expected outputs to contain Makefile, but got %s", rule.Outputs)
		}

		return nil
	}

	value, err := Files(addRule)(thread, builtin, args, kwargs)
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

func TestFilesReturnsErrorWhenInvalidArgs(t *testing.T) {
	thread := &starlark.Thread{}
	builtin := &starlark.Builtin{}
	args := []starlark.Value{}
	kwargs := []starlark.Tuple{}

	_, err := Files(simpleAddFileRule)(thread, builtin, args, kwargs)
	if err == nil {
		t.Fatal("Files should have returned an error")
	}
}

func TestFilesReturnsErrorWhenNameIsInvalid(t *testing.T) {
	thread := &starlark.Thread{}
	builtin := &starlark.Builtin{}
	args := []starlark.Value{
		starlark.String("123rule_name"),
		starlark.NewList([]starlark.Value{}),
	}
	kwargs := []starlark.Tuple{}

	_, err := Files(simpleAddFileRule)(thread, builtin, args, kwargs)
	if err == nil {
		t.Fatal("expected Files to return an error")
	}
}

func TestFilesReturnsErrorWhenAPathIsInvalid(t *testing.T) {
	thread := &starlark.Thread{}
	builtin := &starlark.Builtin{}
	args := []starlark.Value{
		starlark.String("absolute_path"),
		starlark.NewList([]starlark.Value{starlark.String("/etc/test")}),
	}
	kwargs := []starlark.Tuple{}

	_, err := Files(simpleAddFileRule)(thread, builtin, args, kwargs)
	if err == nil {
		t.Fatal("expected Files to return an error")
	}
}

func TestFilesReturnsErrorWhenAddRuleErrs(t *testing.T) {
	thread := &starlark.Thread{}
	builtin := &starlark.Builtin{}
	args := []starlark.Value{
		starlark.String("rule_name"),
		starlark.NewList([]starlark.Value{}),
	}
	kwargs := []starlark.Tuple{}

	addRule := func(rule rule.Rule) error {
		return errors.New("bad stuff")
	}

	_, err := Files(addRule)(thread, builtin, args, kwargs)
	if err == nil {
		t.Fatal("Files should have failed to call addRule")
	}
}
