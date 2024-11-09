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

	got := Expand(test_path)

	if got != want {
		t.Errorf("Expected result to be %s, but got %s", want, got)
	}
}

func TestExpandAbsolutePath(t *testing.T) {
	test_path := "/var/tmp/"
	want := "/var/tmp/"

	got := Expand(test_path)

	if got != want {
		t.Errorf("Expected result to be %s, but got %s", want, got)
	}
}

func TestExpandPathWithEmptyInput(t *testing.T) {
	test_path := ""
	want := ""

	got := Expand(test_path)

	if got != want {
		t.Errorf("Expected result to be empty, but got %s", got)
	}
}

func TestValidateRelativePathWithAbsPath(t *testing.T) {
	test_path := "/var/tmp"
	want := false
	got := isRelativePath(test_path)

	if want != got {
		t.Errorf("Expected %s to be %t, but got %t", test_path, want, got)
	}
}

func TestValidateRelativePathWithRelPath(t *testing.T) {
	test_path := "~/.config"
	want := true
	got := isRelativePath(test_path)

	if want != got {
		t.Errorf("Expected %s to be %t, but got %t", test_path, want, got)
	}
}

func TestValidateRelativePathWithInvalidPath(t *testing.T) {
	test_path := "~"
	want := true
	got := isRelativePath(test_path)

	if want != got {
		t.Errorf("Expected %s to be %t, but got %t", test_path, want, got)
	}
}
