package rules

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	ErrAbsolutePath = errors.New("paths cannot contain an absolute file path")
	ErrInvalidName  = errors.New("name must only contain lowercase letters or underscores")
	ErrParentPath   = errors.New("paths cannot depend on parent path (..)")
)

func ValidatePaths(paths []string) error {
	for _, file := range paths {
		if strings.HasPrefix(file, "/") {
			return fmt.Errorf("%s is invalid: %w", file, ErrAbsolutePath)
		}

		if strings.Contains(file, "..") {
			return fmt.Errorf("%s is invalid: %w", file, ErrParentPath)
		}
	}

	return nil
}

func ValidateName(name string) error {
	re := regexp.MustCompile(`^[a-z_]*$`)

	if re.MatchString(name) {
		return nil
	}

	return ErrInvalidName
}
