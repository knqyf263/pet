package domain

import "errors"

// Selector errors
var (
	ErrNoSnippetsAvailable = errors.New("no snippets available")
	ErrSelectionCancelled  = errors.New("selection cancelled by user")
)

// MatchType defines how DirectSelector matches snippets
type MatchType int

const (
	DirectMatchByDescription MatchType = iota
	DirectMatchByID
)

// Selector defines the interface for snippet selection
// Different implementations handle different selection modes:
// - Interactive (fzf/peco) - user selects from list
// - Direct (exact match) - programmatic selection
// - First match - auto-select first snippet
// - Mock - testing
type Selector interface {
	// Select chooses a snippet from the provided list
	// Returns error if selection fails or is cancelled
	Select(snippets []Snippet) (*Snippet, error)
}

// DirectSelector selects a snippet by exact match (description or ID)
type DirectSelector struct {
	query     string
	matchType MatchType
}

// NewDirectSelector creates a selector that finds exact matches
func NewDirectSelector(query string, matchType MatchType) *DirectSelector {
	return &DirectSelector{
		query:     query,
		matchType: matchType,
	}
}

func (d *DirectSelector) Select(snippets []Snippet) (*Snippet, error) {
	for i := range snippets {
		switch d.matchType {
		case DirectMatchByDescription:
			if snippets[i].Description == d.query {
				return &snippets[i], nil
			}
		case DirectMatchByID:
			// In Phase 1, ID is the same as description
			// Real ID support will come in later phases
			if snippets[i].Description == d.query {
				return &snippets[i], nil
			}
		}
	}
	return nil, ErrSnippetNotFound
}

// FirstMatchSelector automatically selects the first snippet
// Useful for automation and scripting
type FirstMatchSelector struct{}

// NewFirstMatchSelector creates a selector that picks the first snippet
func NewFirstMatchSelector() *FirstMatchSelector {
	return &FirstMatchSelector{}
}

func (f *FirstMatchSelector) Select(snippets []Snippet) (*Snippet, error) {
	if len(snippets) == 0 {
		return nil, ErrNoSnippetsAvailable
	}
	return &snippets[0], nil
}
