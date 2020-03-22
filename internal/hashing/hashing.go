package hashing

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	fpath "path/filepath"
	"sort"
	"strings"

	"github.com/spf13/afero"
)

func getHashForFile(fs afero.Fs, basepath, filepath string) (string, error) {
	file, err := fs.Open(fpath.Join(basepath, filepath))
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err = io.Copy(hash, file); err != nil {
		return "", err
	}
	if _, err = hash.Write([]byte(filepath)); err != nil {
		return "", err
	}

	return string(hash.Sum(nil)), nil
}

func GetHashForFiles(fs afero.Fs, basepath string, filepaths []string) (string, error) {
	sort.Strings(filepaths)

	hashes := []string{}
	for _, filepath := range filepaths {
		hash, err := getHashForFile(fs, basepath, filepath)
		if err != nil {
			return "", err
		}

		hashes = append(hashes, hash)
	}

	hash := sha256.New()
	hash.Write([]byte(strings.Join(hashes, "")))

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func GetHashForStrings(stringsToHash []string) string {
	sort.Strings(stringsToHash)

	hash := sha256.New()

	hash.Write([]byte(strings.Join(stringsToHash, "")))

	return hex.EncodeToString(hash.Sum(nil))
}
