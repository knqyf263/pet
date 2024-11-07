package config

import (
	"os"
	"path/filepath"
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
	want := filepath.Join(getUserHomeDir(), "/.config/pet")

	got, err := ExpandPath(test_path)

	if err != nil {
		t.Errorf("Expected no error, but %v was raised", err)
	}

	if got != want {
		t.Errorf("Expected result to be %s, but got %s", want, got)
	}
}

func TestExpandAbsolutePath(t *testing.T) {
	test_path := "/var/tmp/"
	want := "/var/tmp/"

	got, err := ExpandPath(test_path)

	if err != nil {
		t.Errorf("Expected no error, but %v was raised", err)
	}

	if got != want {
		t.Errorf("Expected result to be %s, but got %s", want, got)
	}
}

func TestExpandPathWithEmptyInput(t *testing.T) {
	test_path := ""
	want := ""

	got, err := ExpandPath(test_path)

	if err == nil {
		t.Errorf("Expected error to be raised, but got nil")
	}

	if got != want {
		t.Errorf("Expected result to be empty, but got %s", got)
	}
}
