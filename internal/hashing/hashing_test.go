package hashing

import (
	"testing"

	"github.com/spf13/afero"
)

func TestGetHashForFilesSortsFilesForConsistentHash(t *testing.T) {
	fs := afero.NewMemMapFs()
	if err := fs.MkdirAll("some/dir", 0755); err != nil {
		t.Fatalf("got error while creating some/dir: %v", err)
	}

	filesWithContent := map[string][]byte{
		"file":          []byte("a file"),
		"some/file":     []byte("another file"),
		"some/dir/file": []byte("yet another file"),
	}

	for filepath, fileContent := range filesWithContent {
		if err := afero.WriteFile(fs, filepath, fileContent, 0755); err != nil {
			t.Fatalf("received error while creating files: %w", err)
		}
	}

	hash, err := GetHashForFiles(fs, []string{"file", "some/file", "some/dir/file"})
	if err != nil {
		t.Fatalf("expected no error retrieving hash, but got: %w", err)
	}

	reorderedHash, err := GetHashForFiles(fs, []string{"some/dir/file", "some/file", "file"})
	if err != nil {
		t.Fatalf("expected no error retrieving reorderedHash, but got: %w", err)
	}

	if hash != reorderedHash {
		t.Errorf("should return the same hash for a list of files, regardless or order")
	}
}

func TestGetHashForFilesAccountsForFileNameWhenHashing(t *testing.T) {
	fs := afero.NewMemMapFs()

	filesWithContent := map[string][]byte{
		"file_says_hello":      []byte("hello"),
		"file_also_says_hello": []byte("hello"),
	}

	for filepath, fileContent := range filesWithContent {
		if err := afero.WriteFile(fs, filepath, fileContent, 0755); err != nil {
			t.Fatalf("received error while creating files: %w", err)
		}
	}

	hash1, err := GetHashForFiles(fs, []string{"file_says_hello"})
	if err != nil {
		t.Fatalf("expected no error retrieving hash1, but got: %w", err)
	}

	hash2, err := GetHashForFiles(fs, []string{"file_also_says_hello"})
	if err != nil {
		t.Fatalf("expected no error retrieving hash2, but got: %w", err)
	}

	if hash1 == hash2 {
		t.Error("files with the same content, but different name should have different hashes")
	}
}

func TestGetHashForFilesReturnsErrorWhenFileCannotOpen(t *testing.T) {
	fs := afero.NewMemMapFs()

	_, err := GetHashForFiles(fs, []string{"Makefile"})
	if err == nil {
		t.Error("expected an error to be returned")
	}
}

func TestGetHashForStrings(t *testing.T) {
	hash1 := GetHashForStrings([]string{"hey", "bye"})
	hash2 := GetHashForStrings([]string{"bye", "hey"})

	if hash1 != "f56c790bfecb3de08ae7cb47eb4dc33cc90feede8e01f69ecc569098c072aabf" {
		t.Fatalf("expected GetHashForStrings to return desired sha256, but got %s", hash1)
	}

	if hash1 != hash2 {
		t.Error("expected GetHashForStrings to sort strings before hashing")
	}
}
