package domain

import (
	"testing"
)

// TestSnippet tests the Snippet domain model
func TestSnippet(t *testing.T) {
	t.Run("creates snippet with all fields", func(t *testing.T) {
		snippet := Snippet{
			Description: "List all docker containers",
			Command:     "docker ps -a",
			Tag:         []string{"docker", "containers"},
			Output:      "container list",
			Filename:    "/path/to/snippets.toml",
		}

		if snippet.Description != "List all docker containers" {
			t.Errorf("expected description 'List all docker containers', got '%s'", snippet.Description)
		}
		if snippet.Command != "docker ps -a" {
			t.Errorf("expected command 'docker ps -a', got '%s'", snippet.Command)
		}
		if len(snippet.Tag) != 2 {
			t.Errorf("expected 2 tags, got %d", len(snippet.Tag))
		}
	})

	t.Run("snippet with empty tags", func(t *testing.T) {
		snippet := Snippet{
			Description: "Test",
			Command:     "echo test",
			Tag:         []string{},
		}

		if snippet.Tag == nil {
			t.Error("expected empty slice, got nil")
		}
		if len(snippet.Tag) != 0 {
			t.Errorf("expected 0 tags, got %d", len(snippet.Tag))
		}
	})

	t.Run("snippet with multiline command", func(t *testing.T) {
		multilineCmd := `docker run -d \
  --name nginx \
  -p 80:80 \
  nginx:latest`

		snippet := Snippet{
			Description: "Run nginx",
			Command:     multilineCmd,
		}

		if snippet.Command != multilineCmd {
			t.Error("multiline command not preserved correctly")
		}
	})
}

// TestParameter tests the Parameter domain model
func TestParameter(t *testing.T) {
	t.Run("creates parameter with name only", func(t *testing.T) {
		param := Parameter{
			Name: "container",
		}

		if param.Name != "container" {
			t.Errorf("expected name 'container', got '%s'", param.Name)
		}
		if param.DefaultValue != "" {
			t.Errorf("expected empty default value, got '%s'", param.DefaultValue)
		}
		if len(param.MultipleDefaults) != 0 {
			t.Errorf("expected 0 multiple defaults, got %d", len(param.MultipleDefaults))
		}
	})

	t.Run("creates parameter with single default value", func(t *testing.T) {
		param := Parameter{
			Name:         "port",
			DefaultValue: "8080",
		}

		if param.DefaultValue != "8080" {
			t.Errorf("expected default value '8080', got '%s'", param.DefaultValue)
		}
	})

	t.Run("creates parameter with multiple default values", func(t *testing.T) {
		param := Parameter{
			Name:             "stage",
			MultipleDefaults: []string{"dev", "staging", "prod"},
		}

		if len(param.MultipleDefaults) != 3 {
			t.Errorf("expected 3 multiple defaults, got %d", len(param.MultipleDefaults))
		}
		if param.MultipleDefaults[0] != "dev" {
			t.Errorf("expected first default 'dev', got '%s'", param.MultipleDefaults[0])
		}
	})

	t.Run("parameter equality by name", func(t *testing.T) {
		param1 := Parameter{Name: "container", DefaultValue: "nginx"}
		param2 := Parameter{Name: "container", DefaultValue: "apache"}

		// Parameters with same name are considered the same parameter
		// even if defaults differ (first occurrence wins in deduplication)
		if param1.Name != param2.Name {
			t.Error("expected parameters with same name to have equal names")
		}
	})
}

// TestSnippetValidation tests snippet validation logic
func TestSnippetValidation(t *testing.T) {
	t.Run("valid snippet", func(t *testing.T) {
		snippet := Snippet{
			Description: "Valid snippet",
			Command:     "echo hello",
		}

		if err := snippet.Validate(); err != nil {
			t.Errorf("expected valid snippet, got error: %v", err)
		}
	})

	t.Run("snippet without description", func(t *testing.T) {
		snippet := Snippet{
			Description: "",
			Command:     "echo hello",
		}

		if err := snippet.Validate(); err == nil {
			t.Error("expected validation error for empty description")
		}
	})

	t.Run("snippet without command", func(t *testing.T) {
		snippet := Snippet{
			Description: "Test",
			Command:     "",
		}

		if err := snippet.Validate(); err == nil {
			t.Error("expected validation error for empty command")
		}
	})

	t.Run("snippet with whitespace-only fields", func(t *testing.T) {
		snippet := Snippet{
			Description: "   ",
			Command:     "  \n  ",
		}

		if err := snippet.Validate(); err == nil {
			t.Error("expected validation error for whitespace-only fields")
		}
	})
}

// TestSnippetHasTag tests tag checking
func TestSnippetHasTag(t *testing.T) {
	snippet := Snippet{
		Description: "Test",
		Command:     "echo test",
		Tag:         []string{"docker", "containers", "k8s"},
	}

	t.Run("has tag that exists", func(t *testing.T) {
		if !snippet.HasTag("docker") {
			t.Error("expected snippet to have tag 'docker'")
		}
	})

	t.Run("does not have tag that doesn't exist", func(t *testing.T) {
		if snippet.HasTag("nginx") {
			t.Error("expected snippet to not have tag 'nginx'")
		}
	})

	t.Run("empty tag list", func(t *testing.T) {
		emptySnippet := Snippet{
			Description: "Test",
			Command:     "echo test",
			Tag:         []string{},
		}

		if emptySnippet.HasTag("anything") {
			t.Error("expected snippet with no tags to not have any tag")
		}
	})
}

// TestSnippetMatchesAnyTag tests multiple tag matching
func TestSnippetMatchesAnyTag(t *testing.T) {
	snippet := Snippet{
		Description: "Test",
		Command:     "echo test",
		Tag:         []string{"docker", "containers"},
	}

	t.Run("matches one of multiple tags", func(t *testing.T) {
		if !snippet.MatchesAnyTag([]string{"k8s", "docker", "nginx"}) {
			t.Error("expected snippet to match at least one tag")
		}
	})

	t.Run("matches none of the tags", func(t *testing.T) {
		if snippet.MatchesAnyTag([]string{"k8s", "nginx", "apache"}) {
			t.Error("expected snippet to not match any tags")
		}
	})

	t.Run("empty tag filter matches nothing", func(t *testing.T) {
		if snippet.MatchesAnyTag([]string{}) {
			t.Error("expected empty tag filter to match nothing")
		}
	})

	t.Run("snippet with no tags matches nothing", func(t *testing.T) {
		emptySnippet := Snippet{
			Description: "Test",
			Command:     "echo test",
			Tag:         []string{},
		}

		if emptySnippet.MatchesAnyTag([]string{"docker"}) {
			t.Error("expected snippet with no tags to match nothing")
		}
	})
}
