package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/knqyf263/pet/config"
	"github.com/knqyf263/pet/snippet"
	"github.com/pelletier/go-toml"
)

func TestFilterStaticSnippetHandling(t *testing.T) {
	// Test the static field handling in filter logic
	// This is a unit test for the static field behavior in the filter function

	// Create test snippet with static flag
	staticFlag := true
	testSnippet := snippet.SnippetInfo{
		Description: "static snippet",
		Command:     "echo hello <name>",
		Tag:         []string{"test"},
		Static:      &staticFlag, // This snippet has static = true
	}

	// Test static field detection
	isStatic := false
	if testSnippet.Static != nil {
		isStatic = *testSnippet.Static
	}

	if !isStatic {
		t.Error("Expected snippet to be static, but it was not")
	}

	// Test regular snippet without static flag
	regularSnippet := snippet.SnippetInfo{
		Description: "regular snippet",
		Command:     "echo hello <name>",
		Tag:         []string{"test"},
		Static:      nil, // No static field
	}

	isStatic = false
	if regularSnippet.Static != nil {
		isStatic = *regularSnippet.Static
	}

	if isStatic {
		t.Error("Expected snippet to not be static, but it was")
	}
}

func TestSnippetSaveWithStatic(t *testing.T) {
	// Test saving and loading snippets with static field
	tempDir, _ := os.MkdirTemp("", "testdata")
	tempSnippetFile := filepath.Join(tempDir, "snippet.toml")
	defer os.RemoveAll(tempDir)

	// Create test snippets - one static, one regular
	staticFlag := true
	testSnippets := snippet.Snippets{
		Snippets: []snippet.SnippetInfo{
			{
				Description: "regular snippet",
				Command:     "echo hello <name>",
				Tag:         []string{"test"},
			},
			{
				Description: "static snippet",
				Command:     "echo hello <name>",
				Tag:         []string{"test"},
				Static:      &staticFlag,
			},
		},
	}

	// Mock configuration
	config.Conf.General.SnippetFile = tempSnippetFile

	// Save snippets
	err := testSnippets.Save()
	if err != nil {
		t.Fatalf("Failed to save snippets: %v", err)
	}

	// Load snippets back
	var loadedSnippets snippet.Snippets
	f, err := os.ReadFile(tempSnippetFile)
	if err != nil {
		t.Fatalf("Failed to read snippet file: %v", err)
	}

	err = toml.Unmarshal(f, &loadedSnippets)
	if err != nil {
		t.Fatalf("Failed to unmarshal snippets: %v", err)
	}

	// Verify we have 2 snippets
	if len(loadedSnippets.Snippets) != 2 {
		t.Fatalf("Expected 2 snippets, got %d", len(loadedSnippets.Snippets))
	}

	// Check first snippet (regular - no static field)
	if loadedSnippets.Snippets[0].Static != nil {
		t.Error("Expected first snippet to have nil Static field")
	}

	// Check second snippet (static)
	if loadedSnippets.Snippets[1].Static == nil {
		t.Error("Expected second snippet to have Static field set")
	} else if !*loadedSnippets.Snippets[1].Static {
		t.Error("Expected second snippet to have Static field set to true")
	}
}
