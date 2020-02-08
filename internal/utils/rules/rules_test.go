package rules

import (
	"errors"
	"testing"
)

func TestValidatePathsReturnNilWhenNoErrors(t *testing.T) {
	err := ValidatePaths([]string{"hello", "BYE", "./some/dir/some/file"})
	if err != nil {
		t.Error("expected nil, but got: %w", err)
	}
}

func TestValidatePathsReturnsErrorWhenAbsolutePath(t *testing.T) {
	err := ValidatePaths([]string{"/hello"})
	if err == nil {
		t.Fatal("ValidatePaths should return error for an absolute path")
	}

	if !errors.Is(err, ErrAbsolutePath) {
		t.Errorf("expected error to be %w, but got %w", ErrAbsolutePath, err)
	}
}

func TestValidatePathsReturnsErrorWhenParentPathIsFound(t *testing.T) {
	err := ValidatePaths([]string{"./hello/.."})
	if err == nil {
		t.Fatal("ValidatePaths should return error when a parent path is used")
	}

	if !errors.Is(err, ErrParentPath) {
		t.Errorf("expected error to be %w, but got %w", ErrParentPath, err)
	}
}

func TestValidateNameReturnNilWhenNoErrors(t *testing.T) {
	err := ValidateName("go_src_code")
	if err != nil {
		t.Error("expected nil, but got: %w", err)
	}
}

func TestValiidateNameReturnsErrorWhenNotLowercaseOrUnderscore(t *testing.T) {
	err := ValidateName("python37")
	if err == nil {
		t.Fatal("ValidateName should return an error when invalid name")
	}

	if !errors.Is(err, ErrInvalidName) {
		t.Errorf("expected error to be %w, but got %w", ErrInvalidName, err)
	}
}
