package domain

import "errors"

// Common errors
var (
	ErrSnippetNotFound = errors.New("snippet not found")
)

// SnippetRepository defines the interface for snippet storage operations
// This is a port in hexagonal architecture - the domain doesn't know
// about the actual storage mechanism (TOML files, database, etc.)
type SnippetRepository interface {
	// Load retrieves all snippets from storage
	Load() ([]Snippet, error)

	// Save persists snippets to storage
	Save(snippets []Snippet) error

	// FindByDescription finds a snippet by exact description match
	FindByDescription(description string) (*Snippet, error)

	// FindByID finds a snippet by its ID
	// Note: In Phase 1, ID might be the same as description
	// Proper ID support will come in later phases
	FindByID(id string) (*Snippet, error)

	// FilterByTags returns snippets that match any of the provided tags
	// Empty tag list returns all snippets
	FilterByTags(tags []string) ([]Snippet, error)
}
