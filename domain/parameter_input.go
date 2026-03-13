package domain

import (
	"errors"
	"fmt"
)

// ParameterInput errors
var (
	ErrMissingParameter      = errors.New("required parameter missing")
	ErrParameterInputCancelled = errors.New("parameter input cancelled by user")
)

// ParameterInput defines the interface for getting parameter values
// Different implementations handle different input modes:
// - Interactive (gocui) - user fills in dialog
// - Flag - values from CLI flags
// - Default - use default values without prompting
// - Positional - values from positional arguments
// - Mock - testing
type ParameterInput interface {
	// GetParameters retrieves values for the provided parameters
	// Returns a map of parameter name -> value
	// Returns error if required parameters are missing or input is cancelled
	GetParameters(params []Parameter) (map[string]string, error)
}

// DefaultParameterInput uses default values without prompting
// Returns error if a parameter has no default value
type DefaultParameterInput struct{}

func NewDefaultParameterInput() *DefaultParameterInput {
	return &DefaultParameterInput{}
}

func (d *DefaultParameterInput) GetParameters(params []Parameter) (map[string]string, error) {
	result := make(map[string]string)

	for _, param := range params {
		// Use single default value
		if param.DefaultValue != "" {
			result[param.Name] = param.DefaultValue
			continue
		}

		// Use first of multiple defaults
		if len(param.MultipleDefaults) > 0 {
			result[param.Name] = param.MultipleDefaults[0]
			continue
		}

		// No default available
		return nil, fmt.Errorf("%w: %s", ErrMissingParameter, param.Name)
	}

	return result, nil
}

// FlagParameterInput uses values from CLI flags
// Falls back to defaults if flag not provided
type FlagParameterInput struct {
	flagValues map[string]string
}

func NewFlagParameterInput(flagValues map[string]string) *FlagParameterInput {
	return &FlagParameterInput{
		flagValues: flagValues,
	}
}

func (f *FlagParameterInput) GetParameters(params []Parameter) (map[string]string, error) {
	result := make(map[string]string)

	for _, param := range params {
		// Try flag value first
		if val, ok := f.flagValues[param.Name]; ok {
			result[param.Name] = val
			continue
		}

		// Fall back to single default
		if param.DefaultValue != "" {
			result[param.Name] = param.DefaultValue
			continue
		}

		// Fall back to first of multiple defaults
		if len(param.MultipleDefaults) > 0 {
			result[param.Name] = param.MultipleDefaults[0]
			continue
		}

		// No value and no default
		return nil, fmt.Errorf("%w: %s", ErrMissingParameter, param.Name)
	}

	return result, nil
}

// PositionalParameterInput uses positional arguments
// Maps args to parameters in order: args[0] -> params[0], args[1] -> params[1], etc.
// Falls back to defaults if args don't cover all parameters
type PositionalParameterInput struct {
	args []string
}

func NewPositionalParameterInput(args []string) *PositionalParameterInput {
	return &PositionalParameterInput{
		args: args,
	}
}

func (p *PositionalParameterInput) GetParameters(params []Parameter) (map[string]string, error) {
	result := make(map[string]string)

	for i, param := range params {
		// Use positional arg if available
		if i < len(p.args) {
			result[param.Name] = p.args[i]
			continue
		}

		// Fall back to single default
		if param.DefaultValue != "" {
			result[param.Name] = param.DefaultValue
			continue
		}

		// Fall back to first of multiple defaults
		if len(param.MultipleDefaults) > 0 {
			result[param.Name] = param.MultipleDefaults[0]
			continue
		}

		// No arg and no default
		return nil, fmt.Errorf("%w: %s (position %d)", ErrMissingParameter, param.Name, i)
	}

	return result, nil
}
