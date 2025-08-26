package snippet

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/knqyf263/pet/config"
	"github.com/stretchr/testify/assert"
)

func TestLoadWithEnvironmentVariable(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()
	
	// Create two snippet files
	mainFile := filepath.Join(tempDir, "main-snippets.toml")
	envFile := filepath.Join(tempDir, "env-snippets.toml")
	
	// Create main snippets
	mainSnippets := &Snippets{
		Snippets: []SnippetInfo{
			{
				Description: "Main snippet 1",
				Command:     "echo main1",
				Tag:         []string{"main"},
				Output:      "",
			},
		},
	}
	createSnippetFile(t, mainFile, mainSnippets)
	
	// Create env snippets
	envSnippets := &Snippets{
		Snippets: []SnippetInfo{
			{
				Description: "Env snippet 1",
				Command:     "echo env1",
				Tag:         []string{"env"},
				Output:      "",
			},
		},
	}
	createSnippetFile(t, envFile, envSnippets)
	
	// Test 1: Load without environment variable
	config.Conf.General.SnippetFile = mainFile
	config.Conf.General.SnippetDirs = []string{}
	
	snippets1 := Snippets{}
	err := snippets1.Load(false)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(snippets1.Snippets))
	assert.Equal(t, "Main snippet 1", snippets1.Snippets[0].Description)
	
	// Test 2: Load with environment variable set
	os.Setenv("PET_SNIPPET_FILE", envFile)
	defer os.Unsetenv("PET_SNIPPET_FILE")
	
	snippets2 := Snippets{}
	err = snippets2.Load(false)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(snippets2.Snippets))
	assert.Equal(t, "Env snippet 1", snippets2.Snippets[0].Description)
}

func TestSaveWithEnvironmentVariable(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()
	
	// Create file paths
	mainFile := filepath.Join(tempDir, "main-snippets.toml")
	envFile := filepath.Join(tempDir, "env-snippets.toml")
	
	// Create empty files
	os.Create(mainFile)
	os.Create(envFile)
	
	// Configure main file
	config.Conf.General.SnippetFile = mainFile
	config.Conf.General.SnippetDirs = []string{}
	
	// Test 1: Save without environment variable
	snippet1 := SnippetInfo{
		Description: "Test snippet for main",
		Command:     "echo test main",
		Tag:         []string{"test"},
		Output:      "",
		Filename:    "", // Will be set to default
	}
	
	snippets1 := Snippets{
		Snippets: []SnippetInfo{snippet1},
	}
	
	err := snippets1.Save()
	assert.NoError(t, err)
	
	// Verify saved to main file
	loadedSnippets1 := Snippets{}
	err = loadedSnippets1.Load(false)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(loadedSnippets1.Snippets))
	assert.Equal(t, "Test snippet for main", loadedSnippets1.Snippets[0].Description)
	
	// Test 2: Save with environment variable set
	os.Setenv("PET_SNIPPET_FILE", envFile)
	defer os.Unsetenv("PET_SNIPPET_FILE")
	
	snippet2 := SnippetInfo{
		Description: "Test snippet for env",
		Command:     "echo test env",
		Tag:         []string{"test"},
		Output:      "",
		Filename:    "", // Will be set based on env var
	}
	
	snippets2 := Snippets{
		Snippets: []SnippetInfo{snippet2},
	}
	
	err = snippets2.Save()
	assert.NoError(t, err)
	
	// Verify saved to env file
	loadedSnippets2 := Snippets{}
	err = loadedSnippets2.Load(false)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(loadedSnippets2.Snippets))
	assert.Equal(t, "Test snippet for env", loadedSnippets2.Snippets[0].Description)
}

func TestEnvironmentVariableOverridesConfig(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()
	
	// Create two snippet files
	configFile := filepath.Join(tempDir, "config-snippets.toml")
	envFile := filepath.Join(tempDir, "override-snippets.toml")
	
	// Create config snippets
	configSnippets := &Snippets{
		Snippets: []SnippetInfo{
			{
				Description: "Config snippet",
				Command:     "echo config",
				Tag:         []string{"config"},
				Output:      "",
			},
		},
	}
	createSnippetFile(t, configFile, configSnippets)
	
	// Create env override snippets
	overrideSnippets := &Snippets{
		Snippets: []SnippetInfo{
			{
				Description: "Override snippet",
				Command:     "echo override",
				Tag:         []string{"override"},
				Output:      "",
			},
		},
	}
	createSnippetFile(t, envFile, overrideSnippets)
	
	// Set config to use configFile
	config.Conf.General.SnippetFile = configFile
	config.Conf.General.SnippetDirs = []string{}
	
	// Set environment variable to override
	os.Setenv("PET_SNIPPET_FILE", envFile)
	defer os.Unsetenv("PET_SNIPPET_FILE")
	
	// Load snippets - should use env file
	snippets := Snippets{}
	err := snippets.Load(false)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(snippets.Snippets))
	assert.Equal(t, "Override snippet", snippets.Snippets[0].Description)
	assert.Contains(t, snippets.Snippets[0].Tag, "override")
}