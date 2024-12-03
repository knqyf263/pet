package snippet

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/knqyf263/pet/config"
	"github.com/pelletier/go-toml"
	"github.com/stretchr/testify/assert"
)

func createSnippets1(absFilePath string) *Snippets {
	snippets := &Snippets{}
	snippets.Snippets = []SnippetInfo{
		{
			Description: "Test snippet",
			Command:     "echo 'Hello, World!'",
			Tag:         []string{"test"},
			Output:      "Hello, World!",
			Filename:    absFilePath,
		},
		{
			Description: "Test snippet 2",
			Command:     "echo 'Hello, World 2!'",
			Tag:         []string{"test"},
			Output:      "Hello, World 2!",
			Filename:    absFilePath,
		},
	}
	return snippets
}

func createSnippets2(absFilePath string) *Snippets {
	snippets := &Snippets{}
	snippets.Snippets = []SnippetInfo{
		{
			Description: "Test snippet 3",
			Command:     "echo 'Hello, World 3!'",
			Tag:         []string{"test"},
			Output:      "Hello, World 3!",
			Filename:    absFilePath,
		},
		{
			Description: "Test snippet 4",
			Command:     "echo 'Hello, World 4!'",
			Tag:         []string{"test"},
			Output:      "Hello, World 4!",
			Filename:    absFilePath,
		},
	}
	return snippets
}

func createSnippets3(absFilePath string) *Snippets {
	snippets := &Snippets{}
	snippets.Snippets = []SnippetInfo{
		{
			Description: "Test snippet 5",
			Command:     "echo 'Hello, World 5!'",
			Tag:         []string{"test"},
			Output:      "Hello, World 5!",
			Filename:    absFilePath,
		},
		{
			Description: "Test snippet 6",
			Command:     "echo 'Hello, World 6!'",
			Tag:         []string{"test"},
			Output:      "Hello, World 6!",
			Filename:    absFilePath,
		},
	}
	return snippets
}

func createSnippetFile(t *testing.T, filename string, snippets *Snippets) {
	// Encode the snippets to TOML without using save and write to the file
	snippetsFile, err := os.Create(filename)
	assert.NoError(t, err)
	err = toml.NewEncoder(snippetsFile).Encode(snippets)
	assert.NoError(t, err)

	err = snippetsFile.Close()
	assert.NoError(t, err)
}

func TestLoad(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "testdata")
	defer os.RemoveAll(tempDir)
	// Mock configuration
	config.Conf.General.SnippetFile = filepath.Join(tempDir, "snippets.toml")

	// Create file with snippets
	testSnippets := createSnippets1(config.Conf.General.SnippetFile)
	createSnippetFile(t, config.Conf.General.SnippetFile, testSnippets)

	// Create Snippets instance and load the snippets
	snippets := &Snippets{}
	err := snippets.Load(false)
	assert.NoError(t, err)

	// Verify the snippets were loaded
	assert.Len(t, snippets.Snippets, 2)
	assert.Equal(t, "Test snippet", snippets.Snippets[0].Description)
	assert.Equal(t, "echo 'Hello, World!'", snippets.Snippets[0].Command)
	assert.Equal(t, []string{"test"}, snippets.Snippets[0].Tag)
	assert.Equal(t, "Hello, World!", snippets.Snippets[0].Output)
	assert.Equal(t, "Test snippet 2", snippets.Snippets[1].Description)
	assert.Equal(t, "echo 'Hello, World 2!'", snippets.Snippets[1].Command)
	assert.Equal(t, []string{"test"}, snippets.Snippets[1].Tag)
	assert.Equal(t, "Hello, World 2!", snippets.Snippets[1].Output)

	// Make sure filename is loaded properly (not stored, just loaded on the fly based on file)
	assert.Equal(t, config.Conf.General.SnippetFile, snippets.Snippets[0].Filename)
	assert.Equal(t, config.Conf.General.SnippetFile, snippets.Snippets[1].Filename)
}

func TestLoadWithIncludeDirectories(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "testdata")
	defer os.RemoveAll(tempDir)
	includeDir := filepath.Join(tempDir, "include1")

	// Mock configuration
	config.Conf.General.SnippetFile = filepath.Join(tempDir, "snippets.toml")
	config.Conf.General.SnippetDirs = []string{includeDir}

	// Create directory
	err := os.MkdirAll(includeDir, 0755)
	assert.NoError(t, err)

	// Add multiple snippet files, 2 in the directory and 1 as the main file
	testSnippets1 := createSnippets1(config.Conf.General.SnippetFile)
	createSnippetFile(t, config.Conf.General.SnippetFile, testSnippets1)

	snippetFile2 := filepath.Join(includeDir, "snippets1.toml")
	testSnippets2 := createSnippets2(snippetFile2)
	createSnippetFile(t, snippetFile2, testSnippets2)

	snippetFile3 := filepath.Join(includeDir, "snippets2.toml")
	testSnippets3 := createSnippets3(snippetFile3)
	createSnippetFile(t, snippetFile3, testSnippets3)

	// Create Snippets instance and load the snippets
	snippets := &Snippets{}
	err = snippets.Load(true)
	assert.NoError(t, err)

	// Verify the snippets were loaded in the correct order (recency by default)
	assert.Len(t, snippets.Snippets, 6)
	assert.Equal(t, testSnippets1.Snippets[0].Description, snippets.Snippets[0].Description)
	assert.Equal(t, testSnippets1.Snippets[1].Description, snippets.Snippets[1].Description)
	assert.Equal(t, testSnippets2.Snippets[0].Description, snippets.Snippets[2].Description)
	assert.Equal(t, testSnippets2.Snippets[1].Description, snippets.Snippets[3].Description)
	assert.Equal(t, testSnippets3.Snippets[0].Description, snippets.Snippets[4].Description)
	assert.Equal(t, testSnippets3.Snippets[1].Description, snippets.Snippets[5].Description)

	// Check that filenames are loaded properly
	assert.Equal(t, config.Conf.General.SnippetFile, snippets.Snippets[0].Filename)
	assert.Equal(t, config.Conf.General.SnippetFile, snippets.Snippets[1].Filename)
	assert.Equal(t, filepath.Join(includeDir, "snippets1.toml"), snippets.Snippets[2].Filename)
	assert.Equal(t, filepath.Join(includeDir, "snippets1.toml"), snippets.Snippets[3].Filename)
	assert.Equal(t, filepath.Join(includeDir, "snippets2.toml"), snippets.Snippets[4].Filename)
	assert.Equal(t, filepath.Join(includeDir, "snippets2.toml"), snippets.Snippets[5].Filename)
}

func TestSave(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "testdata")
	defer os.RemoveAll(tempDir)

	// Mock configuration
	config.Conf.General.SnippetFile = filepath.Join(tempDir, "snippets.toml")

	// Create a snippet
	snippet := SnippetInfo{
		Description: "Test snippet",
		Command:     "echo 'Hello, World!'",
		Tag:         []string{"test"},
		Output:      "Hello, World!",
		Filename:    config.Conf.General.SnippetFile,
	}

	// Create Snippets instance and add the snippet
	snippets := &Snippets{
		Snippets: []SnippetInfo{snippet},
	}

	// Call Save method
	err := snippets.Save()
	assert.NoError(t, err)

	// Verify the file was created
	assert.NoError(t, err)
	_, err = os.Stat(config.Conf.General.SnippetFile)
	assert.NoError(t, err)

	// Verify the contents of the file
	data, err := os.ReadFile(config.Conf.General.SnippetFile)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	// Filename is not stored in the file, so it should not be present in the file
	want := `
[[Snippets]]
  Description = "Test snippet"
  Output = "Hello, World!"
  Tag = ["test"]
  command = "echo 'Hello, World!'"
`
	assert.Equal(t, want, string(data))
}

func TestSaveWithMultipleSnippetFiles(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "testdata")
	defer os.RemoveAll(tempDir)

	includeDir := filepath.Join(tempDir, "include1")

	// Mock configuration
	config.Conf.General.SnippetFile = filepath.Join(tempDir, "snippets.toml")
	config.Conf.General.SnippetDirs = []string{includeDir}

	// Create directory
	err := os.MkdirAll(includeDir, 0755)
	assert.NoError(t, err)

	// Create snippets but do not create the files - only empty files
	// are needed for this test - Save method should write the snippets
	testSnippets1 := createSnippets1(config.Conf.General.SnippetFile)
	os.Create(config.Conf.General.SnippetFile)

	snippetFile2 := filepath.Join(includeDir, "snippets1.toml")
	testSnippets2 := createSnippets2(snippetFile2)
	os.Create(snippetFile2)

	snippetFile3 := filepath.Join(includeDir, "snippets2.toml")
	testSnippets3 := createSnippets3(snippetFile3)
	os.Create(snippetFile3)

	// Create snippets instance and add all snippets
	// with filenames pointing to the respective snippet files
	snippets := &Snippets{}
	snippets.Snippets = append(snippets.Snippets, testSnippets1.Snippets...)
	snippets.Snippets = append(snippets.Snippets, testSnippets2.Snippets...)
	snippets.Snippets = append(snippets.Snippets, testSnippets3.Snippets...)

	// Call Save method
	err = snippets.Save()
	assert.NoError(t, err)

	// Verify the contents of the main file
	data, err := os.ReadFile(config.Conf.General.SnippetFile)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)
	want := `
[[Snippets]]
  Description = "Test snippet"
  Output = "Hello, World!"
  Tag = ["test"]
  command = "echo 'Hello, World!'"

[[Snippets]]
  Description = "Test snippet 2"
  Output = "Hello, World 2!"
  Tag = ["test"]
  command = "echo 'Hello, World 2!'"
`
	assert.Equal(t, want, string(data))

	// Verify the contents of the included files
	data, err = os.ReadFile(filepath.Join(includeDir, "snippets1.toml"))
	assert.NoError(t, err)
	assert.NotEmpty(t, data)
	want = `
[[Snippets]]
  Description = "Test snippet 3"
  Output = "Hello, World 3!"
  Tag = ["test"]
  command = "echo 'Hello, World 3!'"

[[Snippets]]
  Description = "Test snippet 4"
  Output = "Hello, World 4!"
  Tag = ["test"]
  command = "echo 'Hello, World 4!'"
`
	assert.Equal(t, want, string(data))

	data, err = os.ReadFile(filepath.Join(includeDir, "snippets2.toml"))
	assert.NoError(t, err)
	assert.NotEmpty(t, data)
	want = `
[[Snippets]]
  Description = "Test snippet 5"
  Output = "Hello, World 5!"
  Tag = ["test"]
  command = "echo 'Hello, World 5!'"

[[Snippets]]
  Description = "Test snippet 6"
  Output = "Hello, World 6!"
  Tag = ["test"]
  command = "echo 'Hello, World 6!'"
`
	assert.Equal(t, want, string(data))
}
