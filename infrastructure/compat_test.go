package infrastructure

import (
	"path/filepath"
	"testing"

	"github.com/knqyf263/pet/config"
	"github.com/knqyf263/pet/domain"
	"github.com/knqyf263/pet/snippet"
)

// TestCompat_FileRepositorySave_SnippetLoad verifies that files saved by
// FileRepository can be loaded by the existing snippet.Snippets.Load()
func TestCompat_FileRepositorySave_SnippetLoad(t *testing.T) {
	tempDir := t.TempDir()
	snippetFile := filepath.Join(tempDir, "snippets.toml")

	// Save via FileRepository
	repo := NewFileRepository(snippetFile, nil)
	err := repo.Save([]domain.Snippet{
		{
			Description: "Test compat",
			Command:     "echo 'compat test'",
			Tag:         []string{"test", "compat"},
			Output:      "compat test",
		},
	})
	if err != nil {
		t.Fatalf("FileRepository.Save failed: %v", err)
	}

	// Load via existing snippet package
	config.Conf.General.SnippetFile = snippetFile
	config.Conf.General.SnippetDirs = nil
	config.Conf.General.SortBy = ""

	var snippets snippet.Snippets
	err = snippets.Load(false)
	if err != nil {
		t.Fatalf("snippet.Snippets.Load failed: %v", err)
	}

	if len(snippets.Snippets) != 1 {
		t.Fatalf("expected 1 snippet, got %d", len(snippets.Snippets))
	}

	s := snippets.Snippets[0]
	if s.Description != "Test compat" {
		t.Errorf("description mismatch: got '%s'", s.Description)
	}
	if s.Command != "echo 'compat test'" {
		t.Errorf("command mismatch: got '%s'", s.Command)
	}
	if len(s.Tag) != 2 || s.Tag[0] != "test" || s.Tag[1] != "compat" {
		t.Errorf("tag mismatch: got %v", s.Tag)
	}
	if s.Output != "compat test" {
		t.Errorf("output mismatch: got '%s'", s.Output)
	}
}

// TestCompat_SnippetSave_FileRepositoryLoad verifies that files saved by
// the existing snippet.Snippets.Save() can be loaded by FileRepository
func TestCompat_SnippetSave_FileRepositoryLoad(t *testing.T) {
	tempDir := t.TempDir()
	snippetFile := filepath.Join(tempDir, "snippets.toml")

	// Save via existing snippet package
	config.Conf.General.SnippetFile = snippetFile

	snippets := &snippet.Snippets{
		Snippets: []snippet.SnippetInfo{
			{
				Description: "Existing snippet",
				Command:     "echo 'existing'",
				Tag:         []string{"existing"},
				Output:      "existing output",
				Filename:    snippetFile,
			},
		},
	}
	err := snippets.Save()
	if err != nil {
		t.Fatalf("snippet.Snippets.Save failed: %v", err)
	}

	// Load via FileRepository
	repo := NewFileRepository(snippetFile, nil)
	loaded, err := repo.Load()
	if err != nil {
		t.Fatalf("FileRepository.Load failed: %v", err)
	}

	if len(loaded) != 1 {
		t.Fatalf("expected 1 snippet, got %d", len(loaded))
	}

	s := loaded[0]
	if s.Description != "Existing snippet" {
		t.Errorf("description mismatch: got '%s'", s.Description)
	}
	if s.Command != "echo 'existing'" {
		t.Errorf("command mismatch: got '%s'", s.Command)
	}
	if len(s.Tag) != 1 || s.Tag[0] != "existing" {
		t.Errorf("tag mismatch: got %v", s.Tag)
	}
	if s.Output != "existing output" {
		t.Errorf("output mismatch: got '%s'", s.Output)
	}
}

// TestCompat_MultilineCommand verifies multiline commands survive round-trip
// between both systems
func TestCompat_MultilineCommand(t *testing.T) {
	tempDir := t.TempDir()
	snippetFile := filepath.Join(tempDir, "snippets.toml")

	multilineCmd := "docker run -d \\\n  --name nginx \\\n  -p 80:80 \\\n  nginx:latest"

	// Save via FileRepository
	repo := NewFileRepository(snippetFile, nil)
	err := repo.Save([]domain.Snippet{
		{
			Description: "Multiline test",
			Command:     multilineCmd,
			Tag:         []string{"docker"},
		},
	})
	if err != nil {
		t.Fatalf("FileRepository.Save failed: %v", err)
	}

	// Load via existing snippet package
	config.Conf.General.SnippetFile = snippetFile
	config.Conf.General.SnippetDirs = nil
	config.Conf.General.SortBy = ""

	var snippets snippet.Snippets
	err = snippets.Load(false)
	if err != nil {
		t.Fatalf("snippet.Snippets.Load failed: %v", err)
	}

	if len(snippets.Snippets) != 1 {
		t.Fatalf("expected 1 snippet, got %d", len(snippets.Snippets))
	}

	if snippets.Snippets[0].Command != multilineCmd {
		t.Errorf("multiline command not preserved.\nExpected:\n%s\nGot:\n%s", multilineCmd, snippets.Snippets[0].Command)
	}
}
