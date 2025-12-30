package fs

import (
	"os"
	"path/filepath"
)

// EnsureDir ensures that the directory exists, creating it if necessary.
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// WriteFile writes data to a file, creating the directory if it doesn't exist.
func WriteFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := EnsureDir(dir); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// ReadFile reads data from a file.
func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// Exists checks if a file exists.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
