package domain

import (
	"fmt"
	"testing"
)

// MockRepository is a test implementation of SnippetRepository
type MockRepository struct {
	snippets []Snippet
	saveErr  error
	loadErr  error
}

func NewMockRepository(snippets []Snippet) *MockRepository {
	return &MockRepository{
		snippets: snippets,
	}
}

func (m *MockRepository) Load() ([]Snippet, error) {
	if m.loadErr != nil {
		return nil, m.loadErr
	}
	return m.snippets, nil
}

func (m *MockRepository) Save(snippets []Snippet) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.snippets = snippets
	return nil
}

func (m *MockRepository) FindByDescription(desc string) (*Snippet, error) {
	for i := range m.snippets {
		if m.snippets[i].Description == desc {
			return &m.snippets[i], nil
		}
	}
	return nil, ErrSnippetNotFound
}

func (m *MockRepository) FindByID(id string) (*Snippet, error) {
	// For now, use description as ID (Phase 1)
	// Real ID support will come later
	return m.FindByDescription(id)
}

func (m *MockRepository) FilterByTags(tags []string) ([]Snippet, error) {
	if len(tags) == 0 {
		return m.snippets, nil
	}

	var result []Snippet
	for _, snippet := range m.snippets {
		if snippet.MatchesAnyTag(tags) {
			result = append(result, snippet)
		}
	}
	return result, nil
}

// TestMockRepository tests the mock repository implementation
func TestMockRepository(t *testing.T) {
	testSnippets := []Snippet{
		{
			Description: "List docker containers",
			Command:     "docker ps -a",
			Tag:         []string{"docker"},
		},
		{
			Description: "List kubernetes pods",
			Command:     "kubectl get pods",
			Tag:         []string{"k8s", "kubectl"},
		},
		{
			Description: "Git status",
			Command:     "git status",
			Tag:         []string{"git"},
		},
	}

	t.Run("Load returns all snippets", func(t *testing.T) {
		repo := NewMockRepository(testSnippets)

		snippets, err := repo.Load()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(snippets) != 3 {
			t.Errorf("expected 3 snippets, got %d", len(snippets))
		}
	})

	t.Run("Load returns error when configured", func(t *testing.T) {
		repo := NewMockRepository(testSnippets)
		repo.loadErr = fmt.Errorf("mock load error")

		_, err := repo.Load()
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("Save updates snippets", func(t *testing.T) {
		repo := NewMockRepository(testSnippets)

		newSnippets := []Snippet{
			{Description: "New snippet", Command: "echo new"},
		}

		err := repo.Save(newSnippets)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		loaded, _ := repo.Load()
		if len(loaded) != 1 {
			t.Errorf("expected 1 snippet after save, got %d", len(loaded))
		}
		if loaded[0].Description != "New snippet" {
			t.Errorf("expected 'New snippet', got '%s'", loaded[0].Description)
		}
	})

	t.Run("Save returns error when configured", func(t *testing.T) {
		repo := NewMockRepository(testSnippets)
		repo.saveErr = fmt.Errorf("mock save error")

		err := repo.Save(testSnippets)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("FindByDescription finds exact match", func(t *testing.T) {
		repo := NewMockRepository(testSnippets)

		snippet, err := repo.FindByDescription("Git status")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if snippet.Command != "git status" {
			t.Errorf("expected command 'git status', got '%s'", snippet.Command)
		}
	})

	t.Run("FindByDescription returns error for non-existent snippet", func(t *testing.T) {
		repo := NewMockRepository(testSnippets)

		_, err := repo.FindByDescription("Non-existent")
		if err != ErrSnippetNotFound {
			t.Errorf("expected ErrSnippetNotFound, got %v", err)
		}
	})

	t.Run("FilterByTags returns matching snippets", func(t *testing.T) {
		repo := NewMockRepository(testSnippets)

		snippets, err := repo.FilterByTags([]string{"docker"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(snippets) != 1 {
			t.Errorf("expected 1 snippet, got %d", len(snippets))
		}
		if snippets[0].Description != "List docker containers" {
			t.Errorf("unexpected snippet: %s", snippets[0].Description)
		}
	})

	t.Run("FilterByTags with multiple tags returns all matches", func(t *testing.T) {
		repo := NewMockRepository(testSnippets)

		snippets, err := repo.FilterByTags([]string{"docker", "git"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(snippets) != 2 {
			t.Errorf("expected 2 snippets, got %d", len(snippets))
		}
	})

	t.Run("FilterByTags with kubectl tag returns pod listing", func(t *testing.T) {
		repo := NewMockRepository(testSnippets)

		snippets, err := repo.FilterByTags([]string{"kubectl"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(snippets) != 1 {
			t.Errorf("expected 1 snippet, got %d", len(snippets))
		}
		if snippets[0].Description != "List kubernetes pods" {
			t.Errorf("unexpected snippet: %s", snippets[0].Description)
		}
	})

	t.Run("FilterByTags with no tags returns all snippets", func(t *testing.T) {
		repo := NewMockRepository(testSnippets)

		snippets, err := repo.FilterByTags([]string{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(snippets) != 3 {
			t.Errorf("expected 3 snippets, got %d", len(snippets))
		}
	})

	t.Run("FilterByTags with non-matching tag returns empty", func(t *testing.T) {
		repo := NewMockRepository(testSnippets)

		snippets, err := repo.FilterByTags([]string{"nonexistent"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(snippets) != 0 {
			t.Errorf("expected 0 snippets, got %d", len(snippets))
		}
	})
}
