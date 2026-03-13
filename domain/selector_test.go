package domain

import (
	"errors"
	"testing"
)

// MockSelector is a test implementation of Selector
type MockSelector struct {
	returnSnippet *Snippet
	returnErr     error
	returnIndex   int
}

func NewMockSelector(snippet *Snippet) *MockSelector {
	return &MockSelector{
		returnSnippet: snippet,
	}
}

func NewMockSelectorByIndex(index int) *MockSelector {
	return &MockSelector{
		returnIndex: index,
	}
}

func (m *MockSelector) Select(snippets []Snippet) (*Snippet, error) {
	if m.returnErr != nil {
		return nil, m.returnErr
	}

	if m.returnSnippet != nil {
		return m.returnSnippet, nil
	}

	if m.returnIndex >= 0 && m.returnIndex < len(snippets) {
		return &snippets[m.returnIndex], nil
	}

	if len(snippets) == 0 {
		return nil, ErrNoSnippetsAvailable
	}

	// Default: return first snippet
	return &snippets[0], nil
}

// TestMockSelector tests the mock selector implementation
func TestMockSelector(t *testing.T) {
	testSnippets := []Snippet{
		{Description: "First", Command: "echo first"},
		{Description: "Second", Command: "echo second"},
		{Description: "Third", Command: "echo third"},
	}

	t.Run("returns configured snippet", func(t *testing.T) {
		expectedSnippet := &Snippet{Description: "Custom", Command: "echo custom"}
		selector := NewMockSelector(expectedSnippet)

		snippet, err := selector.Select(testSnippets)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if snippet.Description != "Custom" {
			t.Errorf("expected 'Custom', got '%s'", snippet.Description)
		}
	})

	t.Run("returns snippet by index", func(t *testing.T) {
		selector := NewMockSelectorByIndex(1)

		snippet, err := selector.Select(testSnippets)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if snippet.Description != "Second" {
			t.Errorf("expected 'Second', got '%s'", snippet.Description)
		}
	})

	t.Run("returns error when configured", func(t *testing.T) {
		selector := NewMockSelector(nil)
		selector.returnErr = errors.New("mock error")

		_, err := selector.Select(testSnippets)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("returns error for empty snippet list", func(t *testing.T) {
		selector := NewMockSelectorByIndex(0)

		_, err := selector.Select([]Snippet{})
		if err != ErrNoSnippetsAvailable {
			t.Errorf("expected ErrNoSnippetsAvailable, got %v", err)
		}
	})

	t.Run("returns first snippet by default", func(t *testing.T) {
		selector := &MockSelector{returnIndex: -1} // Invalid index, should default

		snippet, err := selector.Select(testSnippets)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if snippet.Description != "First" {
			t.Errorf("expected 'First', got '%s'", snippet.Description)
		}
	})
}

// TestDirectSelector tests direct selection by description
func TestDirectSelector(t *testing.T) {
	testSnippets := []Snippet{
		{Description: "docker ps", Command: "docker ps -a"},
		{Description: "git status", Command: "git status"},
		{Description: "kubectl get pods", Command: "kubectl get pods"},
	}

	t.Run("selects by exact description match", func(t *testing.T) {
		selector := NewDirectSelector("git status", DirectMatchByDescription)

		snippet, err := selector.Select(testSnippets)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if snippet.Command != "git status" {
			t.Errorf("expected 'git status', got '%s'", snippet.Command)
		}
	})

	t.Run("returns error for non-existent description", func(t *testing.T) {
		selector := NewDirectSelector("non-existent", DirectMatchByDescription)

		_, err := selector.Select(testSnippets)
		if err != ErrSnippetNotFound {
			t.Errorf("expected ErrSnippetNotFound, got %v", err)
		}
	})

	t.Run("case-sensitive match", func(t *testing.T) {
		selector := NewDirectSelector("Docker Ps", DirectMatchByDescription)

		_, err := selector.Select(testSnippets)
		if err != ErrSnippetNotFound {
			t.Errorf("expected ErrSnippetNotFound for case mismatch, got %v", err)
		}
	})
}

// TestFirstMatchSelector tests auto-selection of first match
func TestFirstMatchSelector(t *testing.T) {
	testSnippets := []Snippet{
		{Description: "docker ps", Command: "docker ps -a", Tag: []string{"docker"}},
		{Description: "docker images", Command: "docker images", Tag: []string{"docker"}},
		{Description: "git status", Command: "git status", Tag: []string{"git"}},
	}

	t.Run("returns first snippet", func(t *testing.T) {
		selector := NewFirstMatchSelector()

		snippet, err := selector.Select(testSnippets)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if snippet.Description != "docker ps" {
			t.Errorf("expected first snippet, got '%s'", snippet.Description)
		}
	})

	t.Run("returns error for empty list", func(t *testing.T) {
		selector := NewFirstMatchSelector()

		_, err := selector.Select([]Snippet{})
		if err != ErrNoSnippetsAvailable {
			t.Errorf("expected ErrNoSnippetsAvailable, got %v", err)
		}
	})

	t.Run("works with single snippet", func(t *testing.T) {
		selector := NewFirstMatchSelector()
		singleSnippet := []Snippet{
			{Description: "only one", Command: "echo one"},
		}

		snippet, err := selector.Select(singleSnippet)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if snippet.Description != "only one" {
			t.Errorf("expected 'only one', got '%s'", snippet.Description)
		}
	})
}
