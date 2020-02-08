package starlark

import (
	"reflect"
	"testing"

	"go.starlark.net/starlark"
)

func TestGetStringSlice(t *testing.T) {
	starlarkStrings := starlark.NewList([]starlark.Value{
		starlark.String("ghost"),
		starlark.String("dog"),
	})

	strings := GetStringSlice(*starlarkStrings)
	expectedStrings := []string{"ghost", "dog"}

	if !reflect.DeepEqual(strings, expectedStrings) {
		t.Error("bad")
	}
}
