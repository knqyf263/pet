# Phase 3: Parameter Input Abstraction

**Part of:** Architecture Refactoring (Main Proposal: #XXX)
**Depends on:** Phase 2 (#XXX)

## Objective

Abstract parameter input behind an interface to support multiple input methods (interactive gocui, CLI flags, defaults) and improve testability.

## Tasks

- [ ] Define `ParameterInput` interface in `domain/parameter_input.go`
  ```go
  type ParameterInput interface {
      GetParameters(params []Parameter, command string) (map[string]string, error)
  }
  ```

- [ ] Create adapter package structure
  - [ ] `adapters/parameterinput/` directory

- [ ] Implement `GocuiParameterInput` adapter
  - [ ] Refactor `dialog/params.go` into adapter
  - [ ] Remove global state (CurrentCommand, FinalCommand, views)
  - [ ] Encapsulate all gocui logic in adapter struct
  - [ ] Maintain exact same UX as current implementation
  - [ ] Keep keybindings (TAB, ENTER, Ctrl+K, arrows)
  - [ ] Support single and multi-choice parameters

- [ ] Implement `FlagParameterInput` adapter
  - [ ] Accept pre-provided parameter map
  - [ ] Return error if required parameter missing
  - [ ] Use default values when available

- [ ] Implement `DefaultParameterInput` adapter
  - [ ] Use default values without prompting
  - [ ] Error if no default available

- [ ] Implement `MockParameterInput` for testing
  - [ ] Configurable return values
  - [ ] Configurable errors

- [ ] Update parameter extraction logic
  - [ ] Move to `domain/parameter.go`
  - [ ] Pure functions (no I/O)
  - [ ] Support current regex patterns: `<param>`, `<param=default>`, `<param=|_val1_||_val2_|>`

- [ ] Add comprehensive tests
  - [ ] Test parameter extraction
  - [ ] Test parameter substitution
  - [ ] Test each adapter implementation
  - [ ] Test error cases

## Acceptance Criteria

- [ ] `ParameterInput` interface defined and documented
- [ ] All adapters implement the interface
- [ ] GocuiParameterInput works identically to current behavior
- [ ] No global state in dialog package
- [ ] Parameter logic is pure and testable
- [ ] Test coverage >80% for parameter logic
- [ ] Existing commands still work with gocui
- [ ] Tests pass: `make test`

## Implementation Notes

### Current Global State Issue
```go
// dialog/params.go - Current (BAD)
var CurrentCommand string  // ❌ Global
var FinalCommand string    // ❌ Global

// New approach (GOOD)
type GocuiParameterInput struct {
    gui         *gocui.Gui
    command     string              // ✅ Instance state
    paramValues map[string]string   // ✅ Instance state
}
```

### Usage Example
```go
// Interactive mode (gocui)
paramInput := adapters.NewGocuiParameterInput()
values, err := paramInput.GetParameters(params, command)

// CLI mode (flags)
paramInput := adapters.NewFlagParameterInput(map[string]string{
    "container": "nginx",
})
values, err := paramInput.GetParameters(params, command)

// Testing
paramInput := adapters.NewMockParameterInput(expectedValues)
values, err := paramInput.GetParameters(params, command)
```

### gocui Event Loop Encapsulation

The key insight is that even though gocui has a blocking event loop, we hide it:

```go
func (g *GocuiParameterInput) GetParameters(...) (map[string]string, error) {
    gui := gocui.NewGui(...)
    defer gui.Close()

    gui.SetManagerFunc(g.layout)
    g.setupKeybindings()

    // Blocks until user presses ENTER or Ctrl+C
    gui.MainLoop()

    // Returns collected values
    return g.paramValues, g.err
}
```

Callers don't need to know about the event loop!

## Dependencies

- Phase 2: Repository Pattern (#XXX)

## Estimated Effort

Medium (4-6 days)

## Labels

`enhancement`, `refactoring`, `architecture`, `phase-3`, `testing`, `gocui`
