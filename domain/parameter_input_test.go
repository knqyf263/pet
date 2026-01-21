package domain

import (
	"errors"
	"testing"
)

// MockParameterInput is a test implementation of ParameterInput
type MockParameterInput struct {
	values map[string]string
	err    error
}

func NewMockParameterInput(values map[string]string) *MockParameterInput {
	return &MockParameterInput{
		values: values,
	}
}

func (m *MockParameterInput) GetParameters(params []Parameter) (map[string]string, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.values, nil
}

// TestMockParameterInput tests the mock parameter input
func TestMockParameterInput(t *testing.T) {
	params := []Parameter{
		{Name: "container", DefaultValue: "nginx"},
		{Name: "port", DefaultValue: "80"},
	}

	t.Run("returns configured values", func(t *testing.T) {
		values := map[string]string{
			"container": "apache",
			"port":      "8080",
		}
		input := NewMockParameterInput(values)

		result, err := input.GetParameters(params)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if result["container"] != "apache" {
			t.Errorf("expected 'apache', got '%s'", result["container"])
		}
		if result["port"] != "8080" {
			t.Errorf("expected '8080', got '%s'", result["port"])
		}
	})

	t.Run("returns error when configured", func(t *testing.T) {
		input := NewMockParameterInput(nil)
		input.err = errors.New("mock error")

		_, err := input.GetParameters(params)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("handles empty parameters", func(t *testing.T) {
		input := NewMockParameterInput(map[string]string{})

		result, err := input.GetParameters([]Parameter{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(result) != 0 {
			t.Errorf("expected empty result, got %d entries", len(result))
		}
	})
}

// TestDefaultParameterInput tests using default values
func TestDefaultParameterInput(t *testing.T) {
	t.Run("uses default values for all parameters", func(t *testing.T) {
		params := []Parameter{
			{Name: "container", DefaultValue: "nginx"},
			{Name: "port", DefaultValue: "80"},
		}

		input := NewDefaultParameterInput()
		result, err := input.GetParameters(params)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if result["container"] != "nginx" {
			t.Errorf("expected 'nginx', got '%s'", result["container"])
		}
		if result["port"] != "80" {
			t.Errorf("expected '80', got '%s'", result["port"])
		}
	})

	t.Run("returns error for parameter without default", func(t *testing.T) {
		params := []Parameter{
			{Name: "required", DefaultValue: ""},
		}

		input := NewDefaultParameterInput()
		_, err := input.GetParameters(params)

		if !errors.Is(err, ErrMissingParameter) {
			t.Errorf("expected ErrMissingParameter, got %v", err)
		}
	})

	t.Run("handles parameters with multiple defaults", func(t *testing.T) {
		params := []Parameter{
			{Name: "stage", MultipleDefaults: []string{"dev", "staging", "prod"}},
		}

		input := NewDefaultParameterInput()
		result, err := input.GetParameters(params)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Should use first default
		if result["stage"] != "dev" {
			t.Errorf("expected 'dev', got '%s'", result["stage"])
		}
	})

	t.Run("handles empty parameter list", func(t *testing.T) {
		input := NewDefaultParameterInput()
		result, err := input.GetParameters([]Parameter{})

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(result) != 0 {
			t.Errorf("expected empty result, got %d entries", len(result))
		}
	})
}

// TestFlagParameterInput tests CLI flag-based parameter input
func TestFlagParameterInput(t *testing.T) {
	t.Run("uses flag values for all parameters", func(t *testing.T) {
		params := []Parameter{
			{Name: "container", DefaultValue: "nginx"},
			{Name: "port", DefaultValue: "80"},
		}

		flagValues := map[string]string{
			"container": "apache",
			"port":      "8080",
		}

		input := NewFlagParameterInput(flagValues)
		result, err := input.GetParameters(params)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if result["container"] != "apache" {
			t.Errorf("expected 'apache', got '%s'", result["container"])
		}
		if result["port"] != "8080" {
			t.Errorf("expected '8080', got '%s'", result["port"])
		}
	})

	t.Run("falls back to defaults for missing flags", func(t *testing.T) {
		params := []Parameter{
			{Name: "container", DefaultValue: "nginx"},
			{Name: "port", DefaultValue: "80"},
		}

		flagValues := map[string]string{
			"container": "apache",
			// port not provided
		}

		input := NewFlagParameterInput(flagValues)
		result, err := input.GetParameters(params)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if result["container"] != "apache" {
			t.Errorf("expected 'apache', got '%s'", result["container"])
		}
		if result["port"] != "80" {
			t.Errorf("expected default '80', got '%s'", result["port"])
		}
	})

	t.Run("returns error for missing required parameter", func(t *testing.T) {
		params := []Parameter{
			{Name: "required", DefaultValue: ""},
		}

		flagValues := map[string]string{}

		input := NewFlagParameterInput(flagValues)
		_, err := input.GetParameters(params)

		if !errors.Is(err, ErrMissingParameter) {
			t.Errorf("expected ErrMissingParameter, got %v", err)
		}
	})

	t.Run("handles empty parameter list", func(t *testing.T) {
		input := NewFlagParameterInput(map[string]string{})
		result, err := input.GetParameters([]Parameter{})

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(result) != 0 {
			t.Errorf("expected empty result, got %d entries", len(result))
		}
	})

	t.Run("ignores extra flags not in parameters", func(t *testing.T) {
		params := []Parameter{
			{Name: "container", DefaultValue: "nginx"},
		}

		flagValues := map[string]string{
			"container": "apache",
			"extra":     "ignored",
		}

		input := NewFlagParameterInput(flagValues)
		result, err := input.GetParameters(params)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(result) != 1 {
			t.Errorf("expected 1 result, got %d", len(result))
		}
		if result["container"] != "apache" {
			t.Errorf("expected 'apache', got '%s'", result["container"])
		}
	})
}

// TestPositionalParameterInput tests positional argument-based input
func TestPositionalParameterInput(t *testing.T) {
	t.Run("maps positional args to parameters in order", func(t *testing.T) {
		params := []Parameter{
			{Name: "input"},
			{Name: "output"},
			{Name: "fps", DefaultValue: "10"},
		}

		args := []string{"video.mp4", "output.gif", "15"}

		input := NewPositionalParameterInput(args)
		result, err := input.GetParameters(params)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if result["input"] != "video.mp4" {
			t.Errorf("expected 'video.mp4', got '%s'", result["input"])
		}
		if result["output"] != "output.gif" {
			t.Errorf("expected 'output.gif', got '%s'", result["output"])
		}
		if result["fps"] != "15" {
			t.Errorf("expected '15', got '%s'", result["fps"])
		}
	})

	t.Run("uses defaults for missing positional args", func(t *testing.T) {
		params := []Parameter{
			{Name: "input"},
			{Name: "output"},
			{Name: "fps", DefaultValue: "10"},
		}

		args := []string{"video.mp4", "output.gif"}

		input := NewPositionalParameterInput(args)
		result, err := input.GetParameters(params)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if result["fps"] != "10" {
			t.Errorf("expected default '10', got '%s'", result["fps"])
		}
	})

	t.Run("returns error for missing required positional arg", func(t *testing.T) {
		params := []Parameter{
			{Name: "input"},
			{Name: "output"},
		}

		args := []string{"video.mp4"}

		input := NewPositionalParameterInput(args)
		_, err := input.GetParameters(params)

		if !errors.Is(err, ErrMissingParameter) {
			t.Errorf("expected ErrMissingParameter, got %v", err)
		}
	})

	t.Run("handles empty args and empty params", func(t *testing.T) {
		input := NewPositionalParameterInput([]string{})
		result, err := input.GetParameters([]Parameter{})

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(result) != 0 {
			t.Errorf("expected empty result, got %d entries", len(result))
		}
	})

	t.Run("ignores extra positional args", func(t *testing.T) {
		params := []Parameter{
			{Name: "input"},
			{Name: "output"},
		}

		args := []string{"video.mp4", "output.gif", "extra", "ignored"}

		input := NewPositionalParameterInput(args)
		result, err := input.GetParameters(params)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(result) != 2 {
			t.Errorf("expected 2 results, got %d", len(result))
		}
	})
}
