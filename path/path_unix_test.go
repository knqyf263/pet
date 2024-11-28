//go:build !windows

package path

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExpandPathWithTilde(t *testing.T) {
	test_path := "~/.config/pet"
	want := filepath.Join(os.Getenv("HOME"), ".config", "pet")

	got, err := expandPath(test_path)
	if err != nil {
		t.Errorf("Error occured: %s", err)
	}

	if got != want {
		t.Errorf("Expected result to be %s, but got %s", want, got)
	}
}

func TestNewAbsolutePathIsAbsolute(t *testing.T) {
	test_path := "~/relative/path"
	want := filepath.Join(os.Getenv("HOME"), "relative", "path")

	absPath, err := NewAbsolutePath(test_path)
	if err != nil {
		t.Errorf("Error occured: %s", err)
	}

	got := absPath.Get()
	if got != want {
		t.Errorf("Expected result to be %s, but got %s", want, got)
	}
}

func TestExpandAbsolutePathDoesNothing(t *testing.T) {
	test_path := "/var/tmp"
	want := test_path

	got, err := expandPath(test_path)
	if err != nil {
		t.Errorf("Error occured: %s", err)
	}

	if got != want {
		t.Errorf("Expected result to be %s, but got %s", want, got)
	}
}
