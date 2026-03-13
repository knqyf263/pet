package infrastructure

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/knqyf263/pet/domain"
	"github.com/pelletier/go-toml"
)

// tomlSnippetInfo mirrors the TOML structure used by the existing snippet package
// This ensures backward compatibility with existing snippet files
type tomlSnippetInfo struct {
	Description string   `toml:"Description"`
	Command     string   `toml:"command,multiline"`
	Tag         []string `toml:"Tag"`
	Output      string   `toml:"Output"`
}

type tomlSnippets struct {
	Snippets []tomlSnippetInfo `toml:"Snippets"`
}

var tomlFileRegex = regexp.MustCompile(`^.+\.(toml)$`)

// FileRepository implements domain.SnippetRepository using TOML files
// This is an adapter in hexagonal architecture - it bridges the domain
// with file-based storage, maintaining compatibility with existing pet files
type FileRepository struct {
	snippetFile string   // Main snippet file path
	snippetDirs []string // Additional directories to scan for snippet files
}

// NewFileRepository creates a new file-based repository
// snippetFile: path to the main snippet TOML file
// snippetDirs: optional additional directories containing snippet TOML files
func NewFileRepository(snippetFile string, snippetDirs []string) *FileRepository {
	return &FileRepository{
		snippetFile: snippetFile,
		snippetDirs: snippetDirs,
	}
}

// Load reads all snippets from the main file and any additional directories
func (r *FileRepository) Load() ([]domain.Snippet, error) {
	var snippetFiles []string

	// Main snippet file
	if _, err := os.Stat(r.snippetFile); err == nil {
		snippetFiles = append(snippetFiles, r.snippetFile)
	} else if os.IsNotExist(err) {
		return nil, fmt.Errorf("snippet file not found: %s", r.snippetFile)
	} else {
		return nil, fmt.Errorf("failed to stat snippet file: %v", err)
	}

	// Additional directories
	for _, dir := range r.snippetDirs {
		if _, err := os.Stat(dir); err != nil {
			if os.IsNotExist(err) {
				return nil, fmt.Errorf("snippet directory not found: %s", dir)
			}
			return nil, fmt.Errorf("failed to stat snippet directory: %v", err)
		}
		snippetFiles = append(snippetFiles, getTomlFiles(dir)...)
	}

	// Read and parse all files
	var allSnippets []domain.Snippet
	for _, file := range snippetFiles {
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read snippet file %s: %v", file, err)
		}

		var parsed tomlSnippets
		if err := toml.Unmarshal(data, &parsed); err != nil {
			return nil, fmt.Errorf("failed to parse snippet file %s: %v", file, err)
		}

		for _, s := range parsed.Snippets {
			allSnippets = append(allSnippets, domain.Snippet{
				Description: s.Description,
				Command:     s.Command,
				Tag:         s.Tag,
				Output:      s.Output,
				Filename:    file,
			})
		}
	}

	return allSnippets, nil
}

// Save writes snippets to their respective files
// Snippets are grouped by Filename and written to the appropriate file
// Snippets with no Filename are written to the main snippet file
func (r *FileRepository) Save(snippets []domain.Snippet) error {
	// Group snippets by filename
	fileGroups := make(map[string][]domain.Snippet)
	for _, s := range snippets {
		filename := s.Filename
		if filename == "" {
			filename = r.snippetFile
		}
		fileGroups[filename] = append(fileGroups[filename], s)
	}

	// Write each group to its file
	for file, group := range fileGroups {
		if err := r.writeTomlFile(file, group); err != nil {
			return err
		}
	}

	return nil
}

// FindByDescription finds a snippet by exact description match
func (r *FileRepository) FindByDescription(description string) (*domain.Snippet, error) {
	snippets, err := r.Load()
	if err != nil {
		return nil, err
	}

	for i := range snippets {
		if snippets[i].Description == description {
			return &snippets[i], nil
		}
	}

	return nil, domain.ErrSnippetNotFound
}

// FindByID finds a snippet by ID
// In the current implementation, ID is the same as Description
func (r *FileRepository) FindByID(id string) (*domain.Snippet, error) {
	return r.FindByDescription(id)
}

// FilterByTags returns snippets matching any of the provided tags
// Empty tag list returns all snippets
func (r *FileRepository) FilterByTags(tags []string) ([]domain.Snippet, error) {
	snippets, err := r.Load()
	if err != nil {
		return nil, err
	}

	if len(tags) == 0 {
		return snippets, nil
	}

	var filtered []domain.Snippet
	for _, s := range snippets {
		if s.MatchesAnyTag(tags) {
			filtered = append(filtered, s)
		}
	}

	return filtered, nil
}

// writeTomlFile writes a list of snippets to a TOML file
func (r *FileRepository) writeTomlFile(filepath string, snippets []domain.Snippet) error {
	f, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create snippet file %s: %v", filepath, err)
	}
	defer f.Close()

	// Convert to TOML structs
	tomlData := tomlSnippets{}
	for _, s := range snippets {
		tomlData.Snippets = append(tomlData.Snippets, tomlSnippetInfo{
			Description: s.Description,
			Command:     s.Command,
			Tag:         s.Tag,
			Output:      s.Output,
		})
	}

	if err := toml.NewEncoder(f).Encode(tomlData); err != nil {
		return fmt.Errorf("failed to encode snippets to %s: %v", filepath, err)
	}

	return nil
}

// getTomlFiles returns all .toml files in a directory (recursive)
func getTomlFiles(dir string) []string {
	var files []string
	err := filepath.Walk(dir, func(p string, f os.FileInfo, err error) error {
		if err == nil && tomlFileRegex.MatchString(f.Name()) {
			files = append(files, p)
		}
		return nil
	})
	if err != nil {
		log.Printf("warning: failed to walk directory %s: %v", dir, err)
	}
	return files
}
