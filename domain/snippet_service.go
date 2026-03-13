package domain

import (
	"regexp"
	"strings"
)

// parameterRegex matches parameter placeholders: <name> or <name=default> or <name=|_val1_||_val2_|>
// This matches most encountered patterns
// Skips match if there is a whitespace at the end ex. <param='my >
// Ignores <, > characters since they're used to match the pattern
// Compiled once at package initialization for performance
var parameterRegex = regexp.MustCompile(`<([^<>]*[^\s])>`)

// SnippetService provides business logic for managing snippets
// This is pure business logic with no I/O dependencies
type SnippetService struct {
	repo SnippetRepository
}

// NewSnippetService creates a new snippet service
func NewSnippetService(repo SnippetRepository) *SnippetService {
	return &SnippetService{
		repo: repo,
	}
}

// ExtractParameters extracts parameter placeholders from a command string
// Returns parameters in the order they appear (deduplicated)
// Format: <paramName> or <paramName=defaultValue> or <paramName=|_val1_||_val2_||_val3_|>
func (s *SnippetService) ExtractParameters(command string) []Parameter {
	matches := parameterRegex.FindAllStringSubmatch(command, -1)

	if len(matches) == 0 {
		return []Parameter{}
	}

	// Use map to deduplicate, slice to preserve order
	seen := make(map[string]bool)
	var params []Parameter

	for _, match := range matches {
		// match[0] is full match with brackets: <param=value>
		// match[1] is captured group without brackets: param=value
		matchedGroup := match[1]

		// Split on '=' to separate name from default value(s)
		paramName, defaultPart, hasDefault := strings.Cut(matchedGroup, "=")

		// Skip if we've already seen this parameter
		if seen[paramName] {
			continue
		}
		seen[paramName] = true

		param := Parameter{Name: paramName}

		if hasDefault {
			// Check if it's multiple defaults format: |_val1_||_val2_||_val3_|
			if strings.HasPrefix(defaultPart, "|_") && strings.HasSuffix(defaultPart, "_|") {
				// Extract multiple defaults
				// Remove leading "|_" and trailing "_|"
				stripped := strings.TrimPrefix(defaultPart, "|_")
				stripped = strings.TrimSuffix(stripped, "_|")

				// Split on "_||_"
				values := strings.Split(stripped, "_||_")
				param.MultipleDefaults = values
			} else {
				// Single default value
				param.DefaultValue = defaultPart
			}
		}

		params = append(params, param)
	}

	return params
}

// SubstituteParameters replaces parameter placeholders with actual values
// Format: <paramName> or <paramName=defaultValue> -> actual value
// Returns the command with all placeholders replaced
func (s *SnippetService) SubstituteParameters(command string, values map[string]string) string {
	matches := parameterRegex.FindAllStringSubmatch(command, -1)

	if len(matches) == 0 {
		return command
	}

	result := command

	// Replace each parameter placeholder
	for _, match := range matches {
		fullMatch := match[0]      // <param=value>
		matchedGroup := match[1]   // param=value

		// Extract parameter name (before '=')
		paramName, _, _ := strings.Cut(matchedGroup, "=")

		// Get the value to substitute
		if value, ok := values[paramName]; ok {
			result = strings.ReplaceAll(result, fullMatch, value)
		}
	}

	return result
}

// Repository operations - these delegate to the repository

// GetAll retrieves all snippets from the repository
func (s *SnippetService) GetAll() ([]Snippet, error) {
	return s.repo.Load()
}

// GetByDescription finds a snippet by exact description match
func (s *SnippetService) GetByDescription(description string) (*Snippet, error) {
	return s.repo.FindByDescription(description)
}

// GetByID finds a snippet by its ID
func (s *SnippetService) GetByID(id string) (*Snippet, error) {
	return s.repo.FindByID(id)
}

// FilterByTags returns snippets that match any of the provided tags
func (s *SnippetService) FilterByTags(tags []string) ([]Snippet, error) {
	return s.repo.FilterByTags(tags)
}

// Create adds a new snippet to the repository
// Validates the snippet before adding
func (s *SnippetService) Create(snippet Snippet) error {
	if err := snippet.Validate(); err != nil {
		return err
	}

	// Load all snippets
	snippets, err := s.repo.Load()
	if err != nil {
		return err
	}

	// Add new snippet
	snippets = append(snippets, snippet)

	// Save back
	return s.repo.Save(snippets)
}

// Update modifies an existing snippet
// In Phase 1, this finds by description (no proper ID yet)
func (s *SnippetService) Update(description string, updated Snippet) error {
	if err := updated.Validate(); err != nil {
		return err
	}

	snippets, err := s.repo.Load()
	if err != nil {
		return err
	}

	// Find and update
	found := false
	for i := range snippets {
		if snippets[i].Description == description {
			snippets[i] = updated
			found = true
			break
		}
	}

	if !found {
		return ErrSnippetNotFound
	}

	return s.repo.Save(snippets)
}

// Delete removes a snippet by description
// In Phase 1, this finds by description (no proper ID yet)
func (s *SnippetService) Delete(description string) error {
	snippets, err := s.repo.Load()
	if err != nil {
		return err
	}

	// Find and remove
	found := false
	var filtered []Snippet
	for i := range snippets {
		if snippets[i].Description == description {
			found = true
			continue // Skip this one (delete it)
		}
		filtered = append(filtered, snippets[i])
	}

	if !found {
		return ErrSnippetNotFound
	}

	return s.repo.Save(filtered)
}
