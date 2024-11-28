package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/knqyf263/pet/config"
	"github.com/knqyf263/pet/snippet"
	"github.com/pelletier/go-toml"
)

// MockReadCloser is a mock implementation of io.ReadCloser
type MockReadCloser struct {
	*strings.Reader
}

// Close does nothing for this mock implementation
func (m *MockReadCloser) Close() error {
	return nil
}

func TestScan(t *testing.T) {
	message := "Enter something: "

	input := "test\n" // Simulated user input
	want := "test"    // Expected output
	expectedError := error(nil)

	// Create a buffer for output
	var outputBuffer bytes.Buffer
	// Create a mock ReadCloser for input
	inputReader := &MockReadCloser{strings.NewReader(input)}

	result, err := scan(message, &outputBuffer, inputReader, false)

	// Check if the input was printed
	got := result

	// Check if the result matches the expected result
	if want != got {
		t.Errorf("Expected result %q, but got %q", want, got)
	}

	// Check if the error matches the expected error
	if err != expectedError {
		t.Errorf("Expected error %v, but got %v", expectedError, err)
	}
}

func TestScan_EmptyStringWithAllowEmpty(t *testing.T) {
	message := "Enter something: "

	input := "\n"               // Simulated user input
	want := ""                  // Expected output
	expectedError := error(nil) // Should not error

	// Create a buffer for output
	var outputBuffer bytes.Buffer
	// Create a mock ReadCloser for input
	inputReader := &MockReadCloser{strings.NewReader(input)}

	result, err := scan(message, &outputBuffer, inputReader, true)

	// Check if the input was printed
	got := result

	// Check if the result is empty
	if want != got {
		t.Errorf("Expected result %q, but got %q", want, got)
	}

	// Check if the error matches the expected error
	if err != expectedError {
		t.Errorf("Expected error %v, but got %v", expectedError, err)
	}
}

func TestScan_EmptyStringWithoutAllowEmpty(t *testing.T) {
	message := "Enter something: "

	input := "\n" // Simulated user input
	want := ""    // Expected output
	expectedError := CanceledError()

	// Create a buffer for output
	var outputBuffer bytes.Buffer
	// Create a mock ReadCloser for input
	inputReader := &MockReadCloser{strings.NewReader(input)}

	result, err := scan(message, &outputBuffer, inputReader, false)

	// Check if the input was printed
	got := result
	// Check if the result matches the expected result
	if want != got {
		t.Errorf("Expected result %q, but got %q", want, got)
	}

	// Check if the error matches the expected error
	if err.Error() != expectedError.Error() {
		t.Errorf("Expected error %v, but got %v", expectedError, err)
	}
}

func TestScanMultiLine_ExitsOnTwoEmptyLines(t *testing.T) {
	prompt := "Enter something: "
	secondPrompt := "whatever"

	input := "test\nnewline here\nand another;\n\n\n" // Simulated user input
	want := "test\nnewline here\nand another;"        // Expected output
	expectedError := error(nil)

	// Create a buffer for output
	var outputBuffer bytes.Buffer
	// Create a mock ReadCloser for input
	inputReader := &MockReadCloser{strings.NewReader(input)}

	result, err := scanMultiLine(prompt, secondPrompt, &outputBuffer, inputReader)

	// Check if the input was printed
	got := result

	// Check if the result matches the expected result
	if want != got {
		t.Errorf("Expected result %q, but got %q", want, got)
	}

	// Check if the error matches the expected error
	if err != expectedError {
		t.Errorf("Expected error %v, but got %v", expectedError, err)
	}
}

func TestNewSnippetCreationWithSnippetDirectory(t *testing.T) {
	// Setup temporary directory for config
	tempDir := t.TempDir()
	tempSnippetFile := filepath.Join(tempDir, "snippet.toml")
	tempSnippetDir1 := filepath.Join(tempDir, "snippets1")
	tempSnippetDir2 := filepath.Join(tempDir, "snippets2")

	// Clean up temp dirs, needed for windows
	// https://github.com/golang/go/issues/51442
	defer os.RemoveAll(tempSnippetFile)
	defer os.RemoveAll(tempSnippetDir1)
	defer os.RemoveAll(tempSnippetDir2)

	// Create snippet directories
	if err := os.Mkdir(tempSnippetDir1, 0755); err != nil {
		t.Fatalf("Failed to create temp snippet directory: %v", err)
	}
	if err := os.Mkdir(tempSnippetDir2, 0755); err != nil {
		t.Fatalf("Failed to create temp snippet directory: %v", err)
	}

	// Create dummy snippets in the main snippet file
	mainSnippets := snippet.Snippets{
		Snippets: []snippet.SnippetInfo{
			{Description: "main snippet 1", Command: "echo main1"},
			{Description: "main snippet 2", Command: "echo main2"},
		},
	}
	saveSnippetsToFile(t, tempSnippetFile, mainSnippets)

	// Create dummy snippets in the snippet directories
	dirSnippets1 := snippet.Snippets{
		Snippets: []snippet.SnippetInfo{
			{Description: "dir1 snippet 1", Command: "echo dir1-1"},
		},
	}
	dirSnippets2 := snippet.Snippets{
		Snippets: []snippet.SnippetInfo{
			{Description: "dir2 snippet 1", Command: "echo dir2-1"},
		},
	}
	saveSnippetsToFile(t, filepath.Join(tempSnippetDir1, "snippets1.toml"), dirSnippets1)
	saveSnippetsToFile(t, filepath.Join(tempSnippetDir2, "snippets2.toml"), dirSnippets2)

	// Mock configuration
	config.Conf.General.SnippetFile = tempSnippetFile
	config.Conf.General.SnippetDirs = []string{tempSnippetDir1, tempSnippetDir2}

	// Simulate creating a new snippet
	args := []string{"echo new command"}

	// Create a buffer for output
	var outputBuffer bytes.Buffer
	// Create a mock ReadCloser for input
	inputReader := &MockReadCloser{strings.NewReader("test\ntest")}

	err := _new(inputReader, &outputBuffer, args)
	if err != nil {
		t.Fatalf("Failed to create new snippet: %v", err)
	}

	// Load the main snippet file and check:
	// 1 - if the new snippet is added
	// 2 - if the number of snippets is correct (to avoid bugs like overwriting with dir snippets)
	var updatedMainSnippets snippet.Snippets
	loadSnippetsFromFile(t, tempSnippetFile, &updatedMainSnippets)

	if len(updatedMainSnippets.Snippets) != 3 {
		t.Fatalf("Expected 3 snippets in main snippet file, got %d", len(updatedMainSnippets.Snippets))
	}

	newSnippet := updatedMainSnippets.Snippets[2]
	if newSnippet.Command != "echo new command" {
		t.Errorf("Expected new command to be 'echo new command', got '%s'", newSnippet.Command)
	}

	// Ensure the snippet files in the directories remain unchanged
	var unchangedDirSnippets1, unchangedDirSnippets2 snippet.Snippets
	loadSnippetsFromFile(t, filepath.Join(tempSnippetDir1, "snippets1.toml"), &unchangedDirSnippets1)
	loadSnippetsFromFile(t, filepath.Join(tempSnippetDir2, "snippets2.toml"), &unchangedDirSnippets2)

	if !compareSnippets(dirSnippets1, unchangedDirSnippets1) {
		t.Errorf("Snippets in directory 1 have changed")
	}
	if !compareSnippets(dirSnippets2, unchangedDirSnippets2) {
		t.Errorf("Snippets in directory 2 have changed")
	}
}

func saveSnippetsToFile(t *testing.T, filename string, snippets snippet.Snippets) {
	f, err := os.Create(filename)
	if err != nil {
		t.Fatalf("Failed to create snippet file: %v", err)
	}
	defer f.Close()

	if err := toml.NewEncoder(f).Encode(snippets); err != nil {
		t.Fatalf("Failed to encode snippets to file: %v", err)
	}
}

func loadSnippetsFromFile(t *testing.T, filename string, snippets *snippet.Snippets) {
	f, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read snippet file: %v", err)
	}

	if err := toml.Unmarshal(f, snippets); err != nil {
		t.Fatalf("Failed to unmarshal snippets from file: %v", err)
	}
}

func compareSnippets(a, b snippet.Snippets) bool {
	if len(a.Snippets) != len(b.Snippets) {
		return false
	}
	for i := range a.Snippets {
		if a.Snippets[i].Description != b.Snippets[i].Description || a.Snippets[i].Command != b.Snippets[i].Command {
			return false
		}
	}
	return true
}
