package infrastructure

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/knqyf263/pet/domain"
)

func TestFileRepository_Load(t *testing.T) {
	tempDir := t.TempDir()
	snippetFile := filepath.Join(tempDir, "snippets.toml")

	// Write a TOML file with snippets (same format as existing pet)
	tomlContent := `
[[Snippets]]
  Description = "List docker containers"
  Output = ""
  Tag = ["docker"]
  command = "docker ps -a"

[[Snippets]]
  Description = "Git status"
  Output = ""
  Tag = ["git"]
  command = "git status"
`
	err := os.WriteFile(snippetFile, []byte(tomlContent), 0644)
	if err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	repo := NewFileRepository(snippetFile, nil)

	t.Run("loads snippets from TOML file", func(t *testing.T) {
		snippets, err := repo.Load()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(snippets) != 2 {
			t.Fatalf("expected 2 snippets, got %d", len(snippets))
		}

		if snippets[0].Description != "List docker containers" {
			t.Errorf("expected 'List docker containers', got '%s'", snippets[0].Description)
		}
		if snippets[0].Command != "docker ps -a" {
			t.Errorf("expected 'docker ps -a', got '%s'", snippets[0].Command)
		}
		if len(snippets[0].Tag) != 1 || snippets[0].Tag[0] != "docker" {
			t.Errorf("expected tags ['docker'], got %v", snippets[0].Tag)
		}
	})

	t.Run("sets filename on loaded snippets", func(t *testing.T) {
		snippets, _ := repo.Load()

		if snippets[0].Filename != snippetFile {
			t.Errorf("expected filename '%s', got '%s'", snippetFile, snippets[0].Filename)
		}
	})
}

func TestFileRepository_LoadWithDirs(t *testing.T) {
	tempDir := t.TempDir()
	snippetFile := filepath.Join(tempDir, "snippets.toml")
	includeDir := filepath.Join(tempDir, "include")

	err := os.MkdirAll(includeDir, 0755)
	if err != nil {
		t.Fatalf("failed to create include dir: %v", err)
	}

	// Main snippet file
	mainContent := `
[[Snippets]]
  Description = "Main snippet"
  Output = ""
  Tag = ["main"]
  command = "echo main"
`
	err = os.WriteFile(snippetFile, []byte(mainContent), 0644)
	if err != nil {
		t.Fatalf("failed to write main file: %v", err)
	}

	// Included snippet file
	includeContent := `
[[Snippets]]
  Description = "Included snippet"
  Output = ""
  Tag = ["included"]
  command = "echo included"
`
	includeFile := filepath.Join(includeDir, "extra.toml")
	err = os.WriteFile(includeFile, []byte(includeContent), 0644)
	if err != nil {
		t.Fatalf("failed to write include file: %v", err)
	}

	repo := NewFileRepository(snippetFile, []string{includeDir})

	t.Run("loads snippets from main file and directories", func(t *testing.T) {
		snippets, err := repo.Load()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(snippets) != 2 {
			t.Fatalf("expected 2 snippets, got %d", len(snippets))
		}

		// Main file snippets come first
		if snippets[0].Description != "Main snippet" {
			t.Errorf("expected 'Main snippet', got '%s'", snippets[0].Description)
		}
		if snippets[0].Filename != snippetFile {
			t.Errorf("expected filename '%s', got '%s'", snippetFile, snippets[0].Filename)
		}

		// Included file snippets come second
		if snippets[1].Description != "Included snippet" {
			t.Errorf("expected 'Included snippet', got '%s'", snippets[1].Description)
		}
		if snippets[1].Filename != includeFile {
			t.Errorf("expected filename '%s', got '%s'", includeFile, snippets[1].Filename)
		}
	})
}

func TestFileRepository_LoadNonExistent(t *testing.T) {
	repo := NewFileRepository("/nonexistent/path/snippets.toml", nil)

	_, err := repo.Load()
	if err == nil {
		t.Error("expected error for non-existent file, got nil")
	}
}

func TestFileRepository_Save(t *testing.T) {
	tempDir := t.TempDir()
	snippetFile := filepath.Join(tempDir, "snippets.toml")

	// Create initial empty file
	err := os.WriteFile(snippetFile, []byte(""), 0644)
	if err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	repo := NewFileRepository(snippetFile, nil)

	t.Run("saves snippets to TOML file", func(t *testing.T) {
		snippets := []domain.Snippet{
			{
				Description: "Test snippet",
				Command:     "echo 'Hello, World!'",
				Tag:         []string{"test"},
				Output:      "Hello, World!",
				Filename:    snippetFile,
			},
		}

		err := repo.Save(snippets)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify file was written
		data, err := os.ReadFile(snippetFile)
		if err != nil {
			t.Fatalf("failed to read saved file: %v", err)
		}

		content := string(data)
		if content == "" {
			t.Error("expected non-empty file")
		}

		// Verify we can load the saved file back
		loaded, err := repo.Load()
		if err != nil {
			t.Fatalf("failed to load saved file: %v", err)
		}

		if len(loaded) != 1 {
			t.Fatalf("expected 1 snippet, got %d", len(loaded))
		}
		if loaded[0].Description != "Test snippet" {
			t.Errorf("expected 'Test snippet', got '%s'", loaded[0].Description)
		}
		if loaded[0].Command != "echo 'Hello, World!'" {
			t.Errorf("expected command 'echo 'Hello, World!'', got '%s'", loaded[0].Command)
		}
		if loaded[0].Output != "Hello, World!" {
			t.Errorf("expected output 'Hello, World!', got '%s'", loaded[0].Output)
		}
	})

	t.Run("saves snippets without filename to main file", func(t *testing.T) {
		snippets := []domain.Snippet{
			{
				Description: "No filename",
				Command:     "echo test",
				// No Filename set - should default to main snippet file
			},
		}

		err := repo.Save(snippets)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		loaded, err := repo.Load()
		if err != nil {
			t.Fatalf("failed to load: %v", err)
		}

		if len(loaded) != 1 {
			t.Fatalf("expected 1 snippet, got %d", len(loaded))
		}
		if loaded[0].Description != "No filename" {
			t.Errorf("expected 'No filename', got '%s'", loaded[0].Description)
		}
	})
}

func TestFileRepository_SaveMultipleFiles(t *testing.T) {
	tempDir := t.TempDir()
	snippetFile := filepath.Join(tempDir, "snippets.toml")
	includeDir := filepath.Join(tempDir, "include")
	includeFile := filepath.Join(includeDir, "extra.toml")

	err := os.MkdirAll(includeDir, 0755)
	if err != nil {
		t.Fatalf("failed to create dir: %v", err)
	}

	// Create initial files
	os.WriteFile(snippetFile, []byte(""), 0644)
	os.WriteFile(includeFile, []byte(""), 0644)

	repo := NewFileRepository(snippetFile, []string{includeDir})

	snippets := []domain.Snippet{
		{
			Description: "Main snippet",
			Command:     "echo main",
			Filename:    snippetFile,
		},
		{
			Description: "Included snippet",
			Command:     "echo included",
			Filename:    includeFile,
		},
	}

	err = repo.Save(snippets)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Load and verify
	loaded, err := repo.Load()
	if err != nil {
		t.Fatalf("failed to load: %v", err)
	}

	if len(loaded) != 2 {
		t.Fatalf("expected 2 snippets, got %d", len(loaded))
	}
}

func TestFileRepository_FindByDescription(t *testing.T) {
	tempDir := t.TempDir()
	snippetFile := filepath.Join(tempDir, "snippets.toml")

	tomlContent := `
[[Snippets]]
  Description = "docker ps"
  Output = ""
  Tag = ["docker"]
  command = "docker ps -a"

[[Snippets]]
  Description = "git status"
  Output = ""
  Tag = ["git"]
  command = "git status"
`
	os.WriteFile(snippetFile, []byte(tomlContent), 0644)

	repo := NewFileRepository(snippetFile, nil)

	t.Run("finds existing snippet", func(t *testing.T) {
		snippet, err := repo.FindByDescription("git status")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if snippet.Command != "git status" {
			t.Errorf("expected 'git status', got '%s'", snippet.Command)
		}
	})

	t.Run("returns error for non-existent snippet", func(t *testing.T) {
		_, err := repo.FindByDescription("non-existent")
		if err != domain.ErrSnippetNotFound {
			t.Errorf("expected ErrSnippetNotFound, got %v", err)
		}
	})
}

func TestFileRepository_FilterByTags(t *testing.T) {
	tempDir := t.TempDir()
	snippetFile := filepath.Join(tempDir, "snippets.toml")

	tomlContent := `
[[Snippets]]
  Description = "docker ps"
  Output = ""
  Tag = ["docker"]
  command = "docker ps -a"

[[Snippets]]
  Description = "git status"
  Output = ""
  Tag = ["git"]
  command = "git status"

[[Snippets]]
  Description = "kubectl get pods"
  Output = ""
  Tag = ["k8s", "kubectl"]
  command = "kubectl get pods"
`
	os.WriteFile(snippetFile, []byte(tomlContent), 0644)

	repo := NewFileRepository(snippetFile, nil)

	t.Run("filters by single tag", func(t *testing.T) {
		snippets, err := repo.FilterByTags([]string{"docker"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(snippets) != 1 {
			t.Fatalf("expected 1 snippet, got %d", len(snippets))
		}
		if snippets[0].Description != "docker ps" {
			t.Errorf("expected 'docker ps', got '%s'", snippets[0].Description)
		}
	})

	t.Run("filters by multiple tags", func(t *testing.T) {
		snippets, err := repo.FilterByTags([]string{"docker", "git"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(snippets) != 2 {
			t.Fatalf("expected 2 snippets, got %d", len(snippets))
		}
	})

	t.Run("empty tags returns all snippets", func(t *testing.T) {
		snippets, err := repo.FilterByTags([]string{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(snippets) != 3 {
			t.Fatalf("expected 3 snippets, got %d", len(snippets))
		}
	})
}

func TestFileRepository_RoundTrip(t *testing.T) {
	tempDir := t.TempDir()
	snippetFile := filepath.Join(tempDir, "snippets.toml")

	// Create initial file
	os.WriteFile(snippetFile, []byte(""), 0644)

	repo := NewFileRepository(snippetFile, nil)

	// Save some snippets
	original := []domain.Snippet{
		{
			Description: "Multi-tag snippet",
			Command:     "docker compose up -d",
			Tag:         []string{"docker", "compose"},
			Output:      "containers started",
		},
		{
			Description: "No tags",
			Command:     "whoami",
			Tag:         []string{},
			Output:      "",
		},
		{
			Description: "Multiline command",
			Command:     "docker run -d \\\n  --name nginx \\\n  -p 80:80 \\\n  nginx:latest",
			Tag:         []string{"docker"},
			Output:      "",
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

	if len(loaded) != 3 {
		t.Fatalf("expected 3 snippets, got %d", len(loaded))
	}

	// Verify all fields survive round-trip
	for i, orig := range original {
		if loaded[i].Description != orig.Description {
			t.Errorf("[%d] description mismatch: '%s' vs '%s'", i, orig.Description, loaded[i].Description)
		}
		if loaded[i].Command != orig.Command {
			t.Errorf("[%d] command mismatch: '%s' vs '%s'", i, orig.Command, loaded[i].Command)
		}
		if loaded[i].Output != orig.Output {
			t.Errorf("[%d] output mismatch: '%s' vs '%s'", i, orig.Output, loaded[i].Output)
		}
		if len(loaded[i].Tag) != len(orig.Tag) {
			t.Errorf("[%d] tag count mismatch: %d vs %d", i, len(orig.Tag), len(loaded[i].Tag))
		}
	}
}
