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

func TestValidateNameReturnsErrorWhenNameIsAReservedName(t *testing.T) {
	err := ValidateName("all")
	if err == nil {
		t.Fatal("ValidateName should return an error when using a reserved name")
	}

	if !errors.Is(err, ErrReservedName) {
		t.Errorf("expected error to be %w, but got %w", ErrReservedName, err)
	}
}

func TestValiidateNameReturnsErrorWhenInvalidNameFormat(t *testing.T) {
	invalidNames := []string{"python37", "some name", "Name"}
	for _, invalidName := range invalidNames {
		err := ValidateName(invalidName)
		if err == nil {
			t.Fatalf("ValidateName should return an error when invalid name like %s", invalidName)
		}

		if !errors.Is(err, ErrInvalidName) {
			t.Errorf("expected error to be %w, but got %w", ErrInvalidName, err)
		}
	}

	validNames := []string{"python", "some_name"}
	for _, validName := range validNames {
		err := ValidateName(validName)
		if err != nil {
			t.Fatalf("ValidateName should not return an error when valid name like %s, but got: %w", validName, err)
		}
	}
}
