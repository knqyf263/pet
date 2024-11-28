package path

import (
	"fmt"
	"os"
	"path/filepath"
)

// AbsolutePath represents an interface for an absolute file path.
type AbsolutePath interface {
	Get() string
	Set(newPath string) error
}

// absolutePath is an unexported struct that guarantees the file path is always absolute.
// The `path` field is also unexported to prevent direct access.
type absolutePath struct {
	path string
}

// expandPath expands a relative path string to an absolute path string
func expandPath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}

	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error resolving relative path - unknown home dir:\n%v", err)
	}

	if path == "~" {
		return homedir, nil
	} else if len(path) > 1 {
		relativePath := path[2:]
		fullPath := filepath.Join(homedir, relativePath)
		return fullPath, nil
	}

	return "", fmt.Errorf("error resolving relative path - unknown path: %s", path)
}

// NewabsolutePath creates a new absolutePath instance from an absolute path or relative path string.
func NewAbsolutePath(path string) (AbsolutePath, error) {
	// Expand path if relative
	path, err := expandPath(path)
	if err != nil {
		return nil, err
	}

	return &absolutePath{path: path}, nil
}

// Get returns the absolute file path as a string.
func (ap *absolutePath) Get() string {
	return ap.path
}

// Set updates the absolute file path after validating it is absolute.
func (ap *absolutePath) Set(newPath string) error {
	path, err := expandPath(newPath)
	if err != nil {
		return nil
	}

	ap.path = path
	return nil
}
