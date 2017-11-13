package snippet

import (
	"bytes"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/knqyf263/pet/config"
)

type Snippets struct {
	Snippets []SnippetInfo `toml:"snippets"`
}

type SnippetInfo struct {
	Description string   `toml:"description"`
	Command     string   `toml:"command"`
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

// ToString returns the contents of toml file.
func (snippets *Snippets) ToString() (string, error) {
	var buffer bytes.Buffer
	err := toml.NewEncoder(&buffer).Encode(snippets)
	if err != nil {
		return "", fmt.Errorf("Failed to convert struct to TOML string: %v", err)
	}
	return buffer.String(), nil
}
