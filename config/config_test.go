package config

import (
	"os"
	"runtime"
	"testing"
)

func getUserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

func TestExpandPathWithTilde(t *testing.T) {
	test_path := "~/.config/pet"
	expectedPath := getUserHomeDir() + "/.config/pet"

	result, err := ExpandPath(test_path)

	if err != nil {
		t.Errorf("Expected no error, but %v was raised", err)
	}

	if result != expectedPath {
		t.Errorf("Expected result to be %s, but got %s", expectedPath, result)
	}
}

func TestExpandAbsolutePath(t *testing.T) {
	test_path := "/var/tmp/"
	expectedPath := "/var/tmp/"

	result, err := ExpandPath(test_path)

	if err != nil {
		t.Errorf("Expected no error, but %v was raised", err)
	}

	if result != expectedPath {
		t.Errorf("Expected result to be %s, but got %s", expectedPath, result)
	}
}

func TestExpandPathWithEmptyInput(t *testing.T) {
	test_path := ""
	expectedPath := ""

	result, err := ExpandPath(test_path)

	if err == nil {
		t.Errorf("Expected error to be raised, but got nil")
	}

	if result != expectedPath {
		t.Errorf("Expected result to be empty, but got %s", result)
	}
}
