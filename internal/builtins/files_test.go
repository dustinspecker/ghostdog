package builtins

import (
	"testing"

	"go.starlark.net/starlark"
)

func TestFilesReturnsNone(t *testing.T) {
	thread := &starlark.Thread{}
	builtin := &starlark.Builtin{}
	args := []starlark.Value{
		starlark.String("rule_name"),
		starlark.NewList([]starlark.Value{}),
	}
	kwargs := []starlark.Tuple{}

	value, err := Files(thread, builtin, args, kwargs)
	if err != nil {
		t.Fatalf("Files failed: %w", err)
	}

	if value != starlark.None {
		t.Errorf("expected value to be None, but got: %v", value)
	}
}

func TestFilesReturnsErrorWhenInvalidArgs(t *testing.T) {
	thread := &starlark.Thread{}
	builtin := &starlark.Builtin{}
	args := []starlark.Value{}
	kwargs := []starlark.Tuple{}

	_, err := Files(thread, builtin, args, kwargs)
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

	_, err := Files(thread, builtin, args, kwargs)
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

	_, err := Files(thread, builtin, args, kwargs)
	if err == nil {
		t.Fatal("expected Files to return an error")
	}
}
