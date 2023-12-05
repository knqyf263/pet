package snippet

import (
	"bytes"
	"fmt"
	"os"
	"slices"
	"sort"

	"github.com/BurntSushi/toml"
	"github.com/knqyf263/pet/config"
)

type Snippets struct {
	Snippets []SnippetInfo `toml:"snippets"`
}

type SnippetInfo struct {
	Description string   `toml:"description"`
	Commands    []string `toml:"commands"`
	Tag         []string `toml:"tag"`
	Output      string   `toml:"output"`
}

// Load reads toml file.
func (snippets *Snippets) Load() error {
	snippetFile := config.Conf.General.SnippetFile
	if _, err := os.Stat(snippetFile); os.IsNotExist(err) {
		return nil
	}
	if _, err := toml.DecodeFile(snippetFile, snippets); err != nil {
		return fmt.Errorf("Failed to load snippet file. %v", err)
	}
	snippets.Order()
	return nil
}

// Save saves the snippets to toml file.
func (snippets *Snippets) Save() error {
	snippetFile := config.Conf.General.SnippetFile
	f, err := os.Create(snippetFile)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("Failed to save snippet file. err: %s", err)
	}
	return toml.NewEncoder(f).Encode(snippets)
}

// Merge combines two slices of snippets together. If there is a duplicate, their commands and tags are merged as a set.
func (snippets *Snippets) Merge(moreSnippets *Snippets) *Snippets {
	// key our snippets by desc
	groups := map[string]SnippetInfo{}

	// add first group to map
	for _, snippet := range snippets.Snippets {
		groups[snippet.Description] = snippet
	}

	// add second group to map, handling duplicates
	for _, snippet := range moreSnippets.Snippets {
		// handle duplicate
		if existing, ok := groups[snippet.Description]; ok {
			// merge commands
			for _, command := range snippet.Commands {
				if !slices.Contains(existing.Commands, command) {
					existing.Commands = append(existing.Commands, command)
				}
			}
			// merge tags
			for _, tag := range snippet.Tag {
				if !slices.Contains(existing.Tag, tag) {
					existing.Tag = append(existing.Tag, tag)
				}
			}
		} else {
			// snippet is new, merge cleanly
			groups[snippet.Description] = snippet
		}
	}

	var result Snippets
	for _, s := range groups {
		result.Snippets = append(result.Snippets, s)
	}

	return &result
}

// ToString returns the contents of toml file.
func (snippets *Snippets) ToString() (string, error) {
	var buffer bytes.Buffer
	err := toml.NewEncoder(&buffer).Encode(snippets)
	if err != nil {
		return "", fmt.Errorf("Failed to convert struct to TOML string: %v", err)
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
func (a ByCommand) Less(i, j int) bool { return len(a[i].Commands) > len(a[j].Commands) }

type ByDescription []SnippetInfo

func (a ByDescription) Len() int           { return len(a) }
func (a ByDescription) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDescription) Less(i, j int) bool { return a[i].Description > a[j].Description }

type ByOutput []SnippetInfo

func (a ByOutput) Len() int           { return len(a) }
func (a ByOutput) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByOutput) Less(i, j int) bool { return a[i].Output > a[j].Output }
