package domain

import (
	"fmt"
	"strings"
)

// Snippet represents a command snippet with metadata
// This is the core domain model, independent of storage format
type Snippet struct {
	Description string   // Human-readable description
	Command     string   // The actual command (supports multiline)
	Tag         []string // Tags for categorization
	Output      string   // Optional command output example
	Filename    string   // Source file path (for multi-file support)
}

// Parameter represents a placeholder parameter in a command
// Format: <paramName> or <paramName=defaultValue> or <paramName=|_val1_||_val2_|>
type Parameter struct {
	Name             string   // Parameter name (e.g., "container")
	DefaultValue     string   // Single default value (e.g., "nginx")
	MultipleDefaults []string // Multiple default values for selection (e.g., ["dev", "staging", "prod"])
}

// Validate checks if the snippet has required fields
func (s *Snippet) Validate() error {
	if strings.TrimSpace(s.Description) == "" {
		return fmt.Errorf("snippet description cannot be empty")
	}
	if strings.TrimSpace(s.Command) == "" {
		return fmt.Errorf("snippet command cannot be empty")
	}
	return nil
}

// HasTag checks if the snippet has a specific tag
func (s *Snippet) HasTag(tag string) bool {
	for _, t := range s.Tag {
		if t == tag {
			return true
		}
	}
	return false
}

// MatchesAnyTag checks if the snippet has any of the provided tags
func (s *Snippet) MatchesAnyTag(tags []string) bool {
	if len(tags) == 0 || len(s.Tag) == 0 {
		return false
	}

	for _, filterTag := range tags {
		if s.HasTag(filterTag) {
			return true
		}
	}
	return false
}
