package infrastructure

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/knqyf263/pet/domain"
)

func TestFileRepository_ValidateSnippetPath(t *testing.T) {
	tempDir := t.TempDir()
	snippetFile := filepath.Join(tempDir, "snippets.toml")
	includeDir := filepath.Join(tempDir, "include")
	outsideDir := filepath.Join(tempDir, "outside")

	// Create directories
	os.MkdirAll(includeDir, 0755)
	os.MkdirAll(outsideDir, 0755)

	// Create initial files
	os.WriteFile(snippetFile, []byte(""), 0644)

	repo := NewFileRepository(snippetFile, []string{includeDir})

	t.Run("validates main snippet file", func(t *testing.T) {
		err := repo.validateSnippetPath(snippetFile)
		if err != nil {
			t.Errorf("expected main snippet file to be valid, got: %v", err)
		}
	})

	t.Run("validates file in snippet directory", func(t *testing.T) {
		includedFile := filepath.Join(includeDir, "extra.toml")
		err := repo.validateSnippetPath(includedFile)
		if err != nil {
			t.Errorf("expected file in snippet dir to be valid, got: %v", err)
		}
	})

	t.Run("rejects file outside valid locations", func(t *testing.T) {
		outsideFile := filepath.Join(outsideDir, "bad.toml")
		err := repo.validateSnippetPath(outsideFile)
		if err == nil {
			t.Error("expected error for file outside valid locations")
		}
	})

	t.Run("rejects file in parent directory", func(t *testing.T) {
		parentFile := filepath.Join(tempDir, "parent.toml")
		err := repo.validateSnippetPath(parentFile)
		if err == nil {
			t.Error("expected error for file in parent directory")
		}
	})
}

func TestFileRepository_SaveWithInvalidPath(t *testing.T) {
	tempDir := t.TempDir()
	snippetFile := filepath.Join(tempDir, "snippets.toml")
	outsideDir := filepath.Join(tempDir, "outside")

	os.MkdirAll(outsideDir, 0755)
	os.WriteFile(snippetFile, []byte(""), 0644)

	repo := NewFileRepository(snippetFile, nil)

	snippets := []domain.Snippet{
		{
			Description: "Valid snippet",
			Command:     "echo valid",
			Filename:    snippetFile,
		},
		{
			Description: "Invalid snippet",
			Command:     "echo invalid",
			Filename:    filepath.Join(outsideDir, "bad.toml"), // Outside valid locations
		},
	}

	err := repo.Save(snippets)
	if err == nil {
		t.Error("expected error when saving snippet with invalid path")
	}
	if err != nil && !strings.Contains(err.Error(), "invalid snippet path") {
		t.Errorf("expected 'invalid snippet path' error, got: %v", err)
	}
}

func TestFileRepository_SaveLoadConsistency(t *testing.T) {
	tempDir := t.TempDir()
	snippetFile := filepath.Join(tempDir, "snippets.toml")

	os.WriteFile(snippetFile, []byte(""), 0644)

	repo := NewFileRepository(snippetFile, nil)

	// Save snippets
	original := []domain.Snippet{
		{
			Description: "First",
			Command:     "echo first",
			Tag:         []string{"test"},
		},
		{
			Description: "Second",
			Command:     "echo second",
			Tag:         []string{"test"},
		},
	}

	err := repo.Save(original)
	if err != nil {
		t.Fatalf("save failed: %v", err)
	}

	// Load them back
	loaded, err := repo.Load()
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	// All snippets should have been loaded
	if len(loaded) != len(original) {
		t.Errorf("expected %d snippets, got %d", len(original), len(loaded))
	}

	// All loaded snippets should have proper filenames
	for i, s := range loaded {
		if s.Filename != snippetFile {
			t.Errorf("snippet[%d] has wrong filename: %s", i, s.Filename)
		}
	}
}
