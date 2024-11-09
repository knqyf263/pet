package config

import (
	"os"
	"path/filepath"
	"strings"
)

// Given a path to either a file or directory, returns its absolute path format.
// expandPath resolves "~/" (UNIX) and "~\" (Windows) prefix in a given system path.
func Expand(path string) string {
	if isRelativePath(path) {
		homedir, err := os.UserHomeDir()
		if err != nil {
            panic(err)
		}

		if path == "~" {
			return homedir
		} else if len(path) >= 2 {
			relativePath := path[2:]
			return filepath.Join(homedir, relativePath)
		}
	}

	return path
}

// Determine if a given path is in relative format containing user home directory notation (~)
func isRelativePath(path string) bool {
    if path == "~" {
        return true
    }

    if len(path) >= 2 && strings.HasPrefix(path, "~") && os.IsPathSeparator(path[1]) {
        return true
    }

    return false
}
