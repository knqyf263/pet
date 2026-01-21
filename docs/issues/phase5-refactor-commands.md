# Phase 5: Refactor Commands to Use New Architecture

**Part of:** Architecture Refactoring (Main Proposal: #XXX)
**Depends on:** Phase 4 (#XXX)

## Objective

Refactor all CLI commands to use the new hexagonal architecture, removing the monolithic `filter()` function and making commands thin composition layers.

## Tasks

### Command Refactoring

- [ ] Refactor `exec` command (`cmd/exec.go`)
  - [ ] Use SnippetService for snippet operations
  - [ ] Inject Selector for snippet selection
  - [ ] Inject ParameterInput for parameter collection
  - [ ] Add integration tests
  - [ ] Maintain backward compatibility

- [ ] Refactor `search` command (`cmd/search.go`)
  - [ ] Use SnippetService
  - [ ] Inject Selector
  - [ ] Inject ParameterInput
  - [ ] Add integration tests

- [ ] Refactor `clip` command (`cmd/clip.go`)
  - [ ] Use SnippetService
  - [ ] Inject Selector
  - [ ] Inject ParameterInput
  - [ ] Add integration tests

- [ ] Refactor `new` command (`cmd/new.go`)
  - [ ] Use SnippetService for creation
  - [ ] Add validation tests

- [ ] Refactor `edit` command (`cmd/edit.go`)
  - [ ] Use SnippetService
  - [ ] Inject Selector
  - [ ] Add tests

- [ ] Refactor `list` command (`cmd/list.go`)
  - [ ] Use SnippetService
  - [ ] Add filtering tests

- [ ] Refactor `sync` command (`cmd/sync.go`)
  - [ ] Use SnippetService
  - [ ] Consider sync as repository operation
  - [ ] Add sync tests

- [ ] Refactor `configure` command (`cmd/configure.go`)
  - [ ] Minimal changes
  - [ ] Add tests if needed

### Cleanup

- [ ] Remove old `filter()` function from `cmd/util.go`
- [ ] Remove other monolithic utility functions that are now in domain
- [ ] Update any remaining uses of old patterns
- [ ] Clean up deprecated code

### Testing

- [ ] Add integration tests for each command
  - [ ] Test with MockRepository
  - [ ] Test with MockSelector
  - [ ] Test with MockParameterInput
  - [ ] Test error handling
  - [ ] Test edge cases

- [ ] Add end-to-end tests
  - [ ] Test full command workflows
  - [ ] Test with real file repository (temp files)
  - [ ] Test command chaining if applicable

### Documentation

- [ ] Update command help text if needed
- [ ] Add inline code documentation
- [ ] Document adapter injection patterns

## Acceptance Criteria

- [ ] All commands refactored to use new architecture
- [ ] No more direct calls to dialog, fzf/peco from commands
- [ ] Old `filter()` function removed
- [ ] Test coverage for commands >70%
- [ ] Integration tests passing
- [ ] All existing functionality preserved
- [ ] No breaking changes to CLI interface
- [ ] Tests pass: `make test`
- [ ] No performance regression

## Implementation Notes

### Before (Monolithic)

```go
// cmd/exec.go - Old
var execCmd = &cobra.Command{
    Use: "exec",
    RunE: func(cmd *cobra.Command, args []string) error {
        snippet, err := filter()  // ❌ Monolithic, hard to test
        if err != nil {
            return err
        }
        return run(snippet.Command)
    },
}
```

### After (Hexagonal)

```go
// cmd/exec.go - New
var execCmd = &cobra.Command{
    Use: "exec",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Setup dependencies
        repo := repository.NewFileRepository(config.Conf.Snippetfile)
        service := domain.NewSnippetService(repo)
        selector := adapters.NewFzfSelector()
        paramInput := adapters.NewGocuiParameterInput()

        // Business logic composition
        snippets, err := service.FilterByTags(tags)
        if err != nil {
            return err
        }

        items := toSelectableItems(snippets)
        selected, err := selector.Select(items)
        if err != nil {
            return err
        }

        snippet := selected.Value.(domain.Snippet)
        params := service.ExtractParameters(snippet.Command)

        if len(params) > 0 {
            values, err := paramInput.GetParameters(params, snippet.Command)
            if err != nil {
                return err
            }
            snippet.Command = service.SubstituteParameters(snippet.Command, values)
        }

        return run(snippet.Command)
    },
}
```

### Testing Example

```go
func TestExecCommand(t *testing.T) {
    // Setup mocks
    mockRepo := repository.NewMockRepository()
    mockRepo.AddSnippet(testSnippet)

    service := domain.NewSnippetService(mockRepo)
    selector := &adapters.MockSelector{ReturnIndex: 0}
    paramInput := &adapters.MockParameterInput{
        Values: map[string]string{"name": "test"},
    }

    // Test command logic without actual I/O
    snippets, _ := service.GetAll()
    items := toSelectableItems(snippets)
    selected, _ := selector.Select(items)

    assert.Equal(t, "test-snippet", selected.ID)
}
```

### Dependency Injection Pattern

Consider creating a `CommandContext` struct to simplify injection:

```go
type CommandContext struct {
    Service      domain.SnippetService
    Selector     domain.Selector
    ParameterInput domain.ParameterInput
    Executor     CommandExecutor
}

func NewCommandContext() *CommandContext {
    repo := repository.NewFileRepository(config.Conf.Snippetfile)
    return &CommandContext{
        Service:      domain.NewSnippetService(repo),
        Selector:     adapters.NewFzfSelector(),
        ParameterInput: adapters.NewGocuiParameterInput(),
        Executor:     &ShellExecutor{},
    }
}

// In tests
func NewTestCommandContext() *CommandContext {
    return &CommandContext{
        Service:      domain.NewSnippetService(repository.NewMockRepository()),
        Selector:     &adapters.MockSelector{ReturnIndex: 0},
        ParameterInput: &adapters.MockParameterInput{Values: testValues},
        Executor:     &MockExecutor{},
    }
}
```

## Dependencies

- Phase 4: Selector Abstraction (#XXX)

## Estimated Effort

Large (7-10 days)

## Labels

`enhancement`, `refactoring`, `architecture`, `phase-5`, `testing`
