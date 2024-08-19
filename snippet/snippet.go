package snippet

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/kennygrant/sanitize"
	"github.com/knqyf263/pet/config"
	"github.com/pelletier/go-toml"
)

type Snippets struct {
	Snippets []SnippetInfo
}

type SnippetInfo struct {
	Filename    string
	Description string
	Command     string `toml:"command,multiline"`
	Tag         []string
	Output      string
}

// Load reads toml file.
func (snippets *Snippets) Load() error {
	var snippetFiles []string

	snippetFile := config.Conf.General.SnippetFile
	if snippetFile != "" {
		if _, err := os.Stat(snippetFile); err == nil {
			snippetFiles = append(snippetFiles, snippetFile)
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("failed to load snippet file. %v", err)
		}
	}

	for _, dir := range config.Conf.General.SnippetDirs {
		if _, err := os.Stat(dir); err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("snippet directory not found. %s", dir)
			}
			return fmt.Errorf("failed to load snippet directory. %v", err)
		}
		snippetFiles = append(snippetFiles, getFiles(dir)...)
	}

	// Read files and load snippets
	for _, file := range snippetFiles {
		tmp := Snippets{}
		f, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to load snippet file. %v", err)
		}

		err = toml.Unmarshal(f, &tmp)
		if err != nil {
			return fmt.Errorf("failed to parse snippet file. %v", err)
		}

		for _, snippet := range tmp.Snippets {
			snippet.Filename = file
			snippets.Snippets = append(snippets.Snippets, snippet)
		}
	}

	snippets.Order()
	return nil
}

// Save saves the snippets to toml file.
func (snippets *Snippets) Save() error {
	var snippetFile string
	var newSnippets Snippets
	for _, snippet := range snippets.Snippets {
		if snippet.Filename == "" {
			snippetFile = config.Conf.General.SnippetDirs[0] + fmt.Sprintf("%s.toml", strings.ToLower(sanitize.BaseName(snippet.Description)))
			newSnippets.Snippets = append(newSnippets.Snippets, snippet)
		} else if snippet.Filename == config.Conf.General.SnippetFile {
			snippetFile = config.Conf.General.SnippetFile
			newSnippets.Snippets = append(newSnippets.Snippets, snippet)
		}
	}
	f, err := os.Create(snippetFile)
	if err != nil {
		return fmt.Errorf("failed to save snippet file. err: %s", err)
	}

	defer f.Close()
	return toml.NewEncoder(f).Encode(snippets)
}

// ToString returns the contents of toml file.
func (snippets *Snippets) ToString() (string, error) {
	var buffer bytes.Buffer
	err := toml.NewEncoder(&buffer).Encode(snippets)
	if err != nil {
		return "", fmt.Errorf("failed to convert struct to TOML string: %v", err)
	}
	return buffer.String(), nil
}

// Order snippets regarding SortBy option defined in config toml
// Prefix "-" reverses the order, default is "recency", "+<expressions>" is the same as "<expression>"
func (snippets *Snippets) Order() {
	sortBy := config.Conf.General.SortBy
	switch {
	case sortBy == "command" || sortBy == "+command":
		sort.Sort(ByCommand(snippets.Snippets))
	case sortBy == "-command":
		sort.Sort(sort.Reverse(ByCommand(snippets.Snippets)))

	case sortBy == "description" || sortBy == "+description":
		sort.Sort(ByDescription(snippets.Snippets))
	case sortBy == "-description":
		sort.Sort(sort.Reverse(ByDescription(snippets.Snippets)))

	case sortBy == "output" || sortBy == "+output":
		sort.Sort(ByOutput(snippets.Snippets))
	case sortBy == "-output":
		sort.Sort(sort.Reverse(ByOutput(snippets.Snippets)))

	case sortBy == "-recency":
		snippets.reverse()
	}
}

func (snippets *Snippets) reverse() {
	for i, j := 0, len(snippets.Snippets)-1; i < j; i, j = i+1, j-1 {
		snippets.Snippets[i], snippets.Snippets[j] = snippets.Snippets[j], snippets.Snippets[i]
	}
}

type ByCommand []SnippetInfo

func (a ByCommand) Len() int           { return len(a) }
func (a ByCommand) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCommand) Less(i, j int) bool { return a[i].Command > a[j].Command }

type ByDescription []SnippetInfo

func (a ByDescription) Len() int           { return len(a) }
func (a ByDescription) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDescription) Less(i, j int) bool { return a[i].Description > a[j].Description }

type ByOutput []SnippetInfo

func (a ByOutput) Len() int           { return len(a) }
func (a ByOutput) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByOutput) Less(i, j int) bool { return a[i].Output > a[j].Output }
