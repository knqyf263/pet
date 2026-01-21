package domain

import (
	"testing"
)

// TestSnippetService_ExtractParameters tests parameter extraction from commands
func TestSnippetService_ExtractParameters(t *testing.T) {
	service := NewSnippetService(nil) // No repo needed for parameter extraction

	t.Run("extracts single parameter without default", func(t *testing.T) {
		command := "docker run <container>"
		params := service.ExtractParameters(command)

		if len(params) != 1 {
			t.Fatalf("expected 1 parameter, got %d", len(params))
		}
		if params[0].Name != "container" {
			t.Errorf("expected name 'container', got '%s'", params[0].Name)
		}
		if params[0].DefaultValue != "" {
			t.Errorf("expected empty default, got '%s'", params[0].DefaultValue)
		}
	})

	t.Run("extracts parameter with single default value", func(t *testing.T) {
		command := "docker run -p <port=8080> nginx"
		params := service.ExtractParameters(command)

		if len(params) != 1 {
			t.Fatalf("expected 1 parameter, got %d", len(params))
		}
		if params[0].Name != "port" {
			t.Errorf("expected name 'port', got '%s'", params[0].Name)
		}
		if params[0].DefaultValue != "8080" {
			t.Errorf("expected default '8080', got '%s'", params[0].DefaultValue)
		}
	})

	t.Run("extracts parameter with multiple default values", func(t *testing.T) {
		command := "deploy to <stage=|_dev_||_staging_||_prod_|>"
		params := service.ExtractParameters(command)

		if len(params) != 1 {
			t.Fatalf("expected 1 parameter, got %d", len(params))
		}
		if params[0].Name != "stage" {
			t.Errorf("expected name 'stage', got '%s'", params[0].Name)
		}
		if len(params[0].MultipleDefaults) != 3 {
			t.Fatalf("expected 3 defaults, got %d", len(params[0].MultipleDefaults))
		}
		if params[0].MultipleDefaults[0] != "dev" {
			t.Errorf("expected first default 'dev', got '%s'", params[0].MultipleDefaults[0])
		}
	})

	t.Run("extracts multiple parameters", func(t *testing.T) {
		command := "ffmpeg -i <input> -vf fps=<fps=10> <output>.gif"
		params := service.ExtractParameters(command)

		if len(params) != 3 {
			t.Fatalf("expected 3 parameters, got %d", len(params))
		}

		// Check names in order
		if params[0].Name != "input" {
			t.Errorf("expected first param 'input', got '%s'", params[0].Name)
		}
		if params[1].Name != "fps" {
			t.Errorf("expected second param 'fps', got '%s'", params[1].Name)
		}
		if params[2].Name != "output" {
			t.Errorf("expected third param 'output', got '%s'", params[2].Name)
		}

		// Check default for fps
		if params[1].DefaultValue != "10" {
			t.Errorf("expected fps default '10', got '%s'", params[1].DefaultValue)
		}
	})

	t.Run("deduplicates repeated parameters", func(t *testing.T) {
		command := "echo <var> and <var> again <var>"
		params := service.ExtractParameters(command)

		if len(params) != 1 {
			t.Errorf("expected 1 unique parameter, got %d", len(params))
		}
		if params[0].Name != "var" {
			t.Errorf("expected name 'var', got '%s'", params[0].Name)
		}
	})

	t.Run("returns empty for command without parameters", func(t *testing.T) {
		command := "git status"
		params := service.ExtractParameters(command)

		if len(params) != 0 {
			t.Errorf("expected 0 parameters, got %d", len(params))
		}
	})

	t.Run("ignores incomplete angle brackets", func(t *testing.T) {
		command := "echo < file.txt > output.txt"
		params := service.ExtractParameters(command)

		// Should not extract these as parameters
		if len(params) != 0 {
			t.Errorf("expected 0 parameters for redirect syntax, got %d", len(params))
		}
	})

	t.Run("handles parameters with spaces in default values", func(t *testing.T) {
		command := "git commit -m <message=Initial commit>"
		params := service.ExtractParameters(command)

		if len(params) != 1 {
			t.Fatalf("expected 1 parameter, got %d", len(params))
		}
		if params[0].DefaultValue != "Initial commit" {
			t.Errorf("expected default 'Initial commit', got '%s'", params[0].DefaultValue)
		}
	})
}

// TestSnippetService_SubstituteParameters tests parameter substitution
func TestSnippetService_SubstituteParameters(t *testing.T) {
	service := NewSnippetService(nil)

	t.Run("substitutes single parameter", func(t *testing.T) {
		command := "docker run <container>"
		values := map[string]string{"container": "nginx"}

		result := service.SubstituteParameters(command, values)

		expected := "docker run nginx"
		if result != expected {
			t.Errorf("expected '%s', got '%s'", expected, result)
		}
	})

	t.Run("substitutes multiple parameters", func(t *testing.T) {
		command := "ffmpeg -i <input> -vf fps=<fps> <output>.gif"
		values := map[string]string{
			"input":  "video.mp4",
			"fps":    "15",
			"output": "animation",
		}

		result := service.SubstituteParameters(command, values)

		expected := "ffmpeg -i video.mp4 -vf fps=15 animation.gif"
		if result != expected {
			t.Errorf("expected '%s', got '%s'", expected, result)
		}
	})

	t.Run("substitutes same parameter multiple times", func(t *testing.T) {
		command := "echo <var> and <var> again <var>"
		values := map[string]string{"var": "test"}

		result := service.SubstituteParameters(command, values)

		expected := "echo test and test again test"
		if result != expected {
			t.Errorf("expected '%s', got '%s'", expected, result)
		}
	})

	t.Run("returns unchanged command if no parameters", func(t *testing.T) {
		command := "git status"
		values := map[string]string{}

		result := service.SubstituteParameters(command, values)

		if result != command {
			t.Errorf("expected unchanged command, got '%s'", result)
		}
	})

	t.Run("substitutes parameter with default value syntax", func(t *testing.T) {
		command := "docker run -p <port=8080>"
		values := map[string]string{"port": "3000"}

		result := service.SubstituteParameters(command, values)

		expected := "docker run -p 3000"
		if result != expected {
			t.Errorf("expected '%s', got '%s'", expected, result)
		}
	})

	t.Run("substitutes parameter with spaces in value", func(t *testing.T) {
		command := "git commit -m <message>"
		values := map[string]string{"message": "Fix: resolve bug"}

		result := service.SubstituteParameters(command, values)

		expected := "git commit -m Fix: resolve bug"
		if result != expected {
			t.Errorf("expected '%s', got '%s'", expected, result)
		}
	})

	t.Run("handles empty values", func(t *testing.T) {
		command := "echo <var>"
		values := map[string]string{"var": ""}

		result := service.SubstituteParameters(command, values)

		expected := "echo "
		if result != expected {
			t.Errorf("expected '%s', got '%s'", expected, result)
		}
	})
}

// TestSnippetService_Repository tests repository integration
func TestSnippetService_Repository(t *testing.T) {
	testSnippets := []Snippet{
		{Description: "docker ps", Command: "docker ps -a", Tag: []string{"docker"}},
		{Description: "git status", Command: "git status", Tag: []string{"git"}},
	}
	repo := NewMockRepository(testSnippets)
	service := NewSnippetService(repo)

	t.Run("GetAll returns all snippets from repository", func(t *testing.T) {
		snippets, err := service.GetAll()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(snippets) != 2 {
			t.Errorf("expected 2 snippets, got %d", len(snippets))
		}
	})

	t.Run("GetByDescription finds snippet", func(t *testing.T) {
		snippet, err := service.GetByDescription("git status")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if snippet.Command != "git status" {
			t.Errorf("expected command 'git status', got '%s'", snippet.Command)
		}
	})

	t.Run("FilterByTags returns matching snippets", func(t *testing.T) {
		snippets, err := service.FilterByTags([]string{"docker"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(snippets) != 1 {
			t.Errorf("expected 1 snippet, got %d", len(snippets))
		}
		if snippets[0].Description != "docker ps" {
			t.Errorf("unexpected snippet: %s", snippets[0].Description)
		}
	})

	t.Run("Create adds snippet to repository", func(t *testing.T) {
		newSnippet := Snippet{
			Description: "new snippet",
			Command:     "echo new",
			Tag:         []string{"test"},
		}

		err := service.Create(newSnippet)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify it was added
		all, _ := service.GetAll()
		if len(all) != 3 {
			t.Errorf("expected 3 snippets after create, got %d", len(all))
		}
	})

	t.Run("Create validates snippet before adding", func(t *testing.T) {
		invalidSnippet := Snippet{
			Description: "", // Invalid - empty description
			Command:     "echo test",
		}

		err := service.Create(invalidSnippet)
		if err == nil {
			t.Error("expected validation error, got nil")
		}
	})
}

// TestSnippetService_Integration tests full workflow
func TestSnippetService_Integration(t *testing.T) {
	testSnippets := []Snippet{
		{
			Description: "Run docker container",
			Command:     "docker run -p <port=8080> <container=nginx>",
			Tag:         []string{"docker"},
		},
	}
	repo := NewMockRepository(testSnippets)
	service := NewSnippetService(repo)

	t.Run("full workflow: find snippet, extract params, substitute", func(t *testing.T) {
		// Find snippet
		snippet, err := service.GetByDescription("Run docker container")
		if err != nil {
			t.Fatalf("failed to get snippet: %v", err)
		}

		// Extract parameters
		params := service.ExtractParameters(snippet.Command)
		if len(params) != 2 {
			t.Fatalf("expected 2 parameters, got %d", len(params))
		}

		// Provide values (simulating user input)
		values := map[string]string{
			"port":      "3000",
			"container": "apache",
		}

		// Substitute
		finalCommand := service.SubstituteParameters(snippet.Command, values)

		expected := "docker run -p 3000 apache"
		if finalCommand != expected {
			t.Errorf("expected '%s', got '%s'", expected, finalCommand)
		}
	})

	t.Run("workflow with default values", func(t *testing.T) {
		snippet, _ := service.GetByDescription("Run docker container")
		params := service.ExtractParameters(snippet.Command)

		// Use defaults via DefaultParameterInput
		paramInput := NewDefaultParameterInput()
		values, err := paramInput.GetParameters(params)
		if err != nil {
			t.Fatalf("failed to get parameters: %v", err)
		}

		finalCommand := service.SubstituteParameters(snippet.Command, values)

		expected := "docker run -p 8080 nginx"
		if finalCommand != expected {
			t.Errorf("expected '%s', got '%s'", expected, finalCommand)
		}
	})
}
