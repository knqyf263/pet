package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/knqyf263/pet/config"
	"github.com/knqyf263/pet/snippet"
	"github.com/stretchr/testify/assert"
)

func setupTestConfig(t *testing.T) (string, func()) {
	// Setup temporary directory for config
	tempDir, err := os.MkdirTemp("", "testdata")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	tempSnippetFile := filepath.Join(tempDir, "snippet.toml")

	// Create dummy snippets in the main snippet file
	mainSnippets := snippet.Snippets{
		Snippets: []snippet.SnippetInfo{
			{Description: "main snippet 1", Command: "echo main"},
			{Description: "main snippet 2", Command: "echo something else"},
		},
	}
	saveSnippetsToFile(t, tempSnippetFile, mainSnippets)

	// Mock configuration
	config.Conf.General.SnippetFile = tempSnippetFile

	// Set SelectCmd to a valid command with piping
	config.Conf.General.SelectCmd = "fzf"

	// Return cleanup function
	return tempDir, func() {
		os.RemoveAll(tempDir)
	}
}

func TestExecute_EmptyString(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()

	var stdout bytes.Buffer
	stdin := &MockReadCloser{strings.NewReader("\n")}

	err := _execute(stdin, &stdout)
	assert.NoError(t, err)
	assert.Contains(t, stdout.String(), "> ")
}

// func TestExecute_FindsMatchingCommand(t *testing.T) {
// 	_, cleanup := setupTestConfig(t)
// 	var stdout bytes.Buffer
// 	stdin := &MockReadCloser{strings.NewReader("\n")}
// 	defer cleanup()

// 	config.Flag.Query = "main | head -n 1"

// 	err := _execute(stdin, &stdout)
// 	assert.NoError(t, err)
// 	assert.Contains(t, stdout.String(), "> echo main")
// }

func TestExecute_SilentFlag(t *testing.T) {
	_, cleanup := setupTestConfig(t)
	defer cleanup()
	var stdout bytes.Buffer
	stdin := &MockReadCloser{strings.NewReader("\n")}

	config.Flag.Silent = true

	err := _execute(stdin, &stdout)
	assert.NoError(t, err)
	assert.Equal(t, stdout.String(), "")
}
