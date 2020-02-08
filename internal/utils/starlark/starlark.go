package starlark

import (
	"go.starlark.net/starlark"
)

func GetStringSlice(starlarkList starlark.List) []string {
	var strings []string

	iter := starlarkList.Iterate()
	defer iter.Done()

	var value starlark.Value
	for iter.Next(&value) {
		// remove surrounding quotes from string value
		str := value.String()
		strippedStr := str[1 : len(str)-1]
		strings = append(strings, strippedStr)
	}

	return strings
}
