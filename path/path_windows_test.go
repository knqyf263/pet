//go:build windows

package path

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExpandPathWithTilde(t *testing.T) {
	test_path := "~/.config/pet"

	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Errorf("Error occurred: %s", err)
	}
	want := filepath.Join(homeDir, ".config", "pet")

	got, err := expandPath(test_path)
	if err != nil {
		t.Errorf("Error occurred: %s", err)
	}

	if got != want {
		t.Errorf("Expected result to be %s, but got %s", want, got)
	}
}

func TestNewAbsolutePathIsAbsolute(t *testing.T) {
	test_path := "~/relative/path"
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Errorf("Error occurred: %s", err)
	}
	want := filepath.Join(homeDir, "relative", "path")

	absPath, err := NewAbsolutePath(test_path)
	if err != nil {
		t.Errorf("Error occurred: %s", err)
	}

	got := absPath.Get()
	if got != want {
		t.Errorf("Expected result to be %s, but got %s", want, got)
	}
}

func TestSetAbsolutePathIsAbsolute(t *testing.T) {
	test_path := "~/relative/path"
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Errorf("Error occurred: %s", err)
	}
	want := filepath.Join(homeDir, "relative", "path")

	absPath, err := NewAbsolutePath("/whatever")
	if err != nil {
		t.Errorf("Error occurred: %s", err)
	}

	err = absPath.Set(test_path)
	if err != nil {
		t.Errorf("Error occurred: %s", err)
	}

	got := absPath.Get()
	if got != want {
		t.Errorf("Expected result to be %s, but got %s", want, got)
	}
}

func TestExpandAbsolutePathDoesNothing(t *testing.T) {
	test_path := "C:\\absolute\\path"
	want := test_path

	got, err := expandPath(test_path)
	if err != nil {
		t.Errorf("Error occurred: %s", err)
	}

	if got != want {
		t.Errorf("Expected result to be %s, but got %s", want, got)
	}
}
