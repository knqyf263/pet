# Phase 4: Selector Abstraction

**Part of:** Architecture Refactoring (Main Proposal: #XXX)
**Depends on:** Phase 3 (#XXX)

## Objective

Abstract snippet selection behind an interface to support multiple selection methods (fzf, peco, direct match, fuzzy search) and improve testability.

## Tasks

- [ ] Define `Selector` interface in `domain/selector.go`
  ```go
  type Selector interface {
      Select(items []SelectableItem) (*SelectableItem, error)
      SelectMultiple(items []SelectableItem) ([]SelectableItem, error)
  }

  type SelectableItem struct {
      ID          string
      DisplayText string
      Value       interface{}
  }
  ```

- [ ] Implement `FzfSelector` adapter
  - [ ] Wrap current fzf piping logic
  - [ ] Respect `config.General.SelectCmd` setting
  - [ ] Maintain current formatting: `[$description]: $command $tags`
  - [ ] Handle empty selections gracefully

- [ ] Implement `PecoSelector` adapter
  - [ ] Similar to FzfSelector but for peco
  - [ ] Support peco-specific options

- [ ] Implement `DirectSelector` adapter (for CLI mode)
  - [ ] Match by exact description
  - [ ] Match by ID
  - [ ] Match by regex pattern
  - [ ] Return error if no match or multiple matches

- [ ] Implement `FuzzySelector` adapter (for CLI mode)
  - [ ] Fuzzy matching without external tool
  - [ ] Return best match or interactive prompt if ambiguous

- [ ] Implement `MockSelector` for testing
  - [ ] Return item at specified index
  - [ ] Configurable errors

- [ ] Add comprehensive tests
  - [ ] Test each selector implementation
  - [ ] Test error cases (no selection, empty list)
  - [ ] Test formatting of display text

## Acceptance Criteria

- [ ] `Selector` interface defined and documented
- [ ] All adapters implement the interface
- [ ] FzfSelector/PecoSelector work identically to current behavior
- [ ] DirectSelector supports exact match by description/ID
- [ ] Test coverage >80% for selector logic
- [ ] Existing commands still work with fzf/peco
- [ ] Tests pass: `make test`

## Implementation Notes

### Current Filter Function Issue

```go
// cmd/util.go - Current (monolithic)
func filter() (snippet.SnippetInfo, error) {
    // Load snippets
    // Filter by tags
    // Pipe to fzf/peco        // ❌ Hard to test
    // Parse selection
    // Show parameter dialog    // ❌ Chained I/O
    // Return result
}
```

### New Approach (composable)

```go
// Separate concerns
snippets := service.GetAll()
filtered := service.FilterByTags(tags)
items := toSelectableItems(filtered)

// Inject selector
selector := adapters.NewFzfSelector()
selected, err := selector.Select(items)

// Inject parameter input
snippet := selected.Value.(domain.Snippet)
params := service.ExtractParameters(snippet.Command)
paramInput := adapters.NewGocuiParameterInput()
values, err := paramInput.GetParameters(params, snippet.Command)
```

### DirectSelector Usage (CLI mode)

```bash
# User runs: pet exec -d "docker ps"

# Code:
selector := adapters.NewDirectSelector("docker ps", DirectMatchByDescription)
selected, err := selector.Select(items) // Returns immediately, no user input!
```

### Formatting

Current format: `[$description]: $command $tags`

```go
func toSelectableItems(snippets []domain.Snippet) []domain.SelectableItem {
    items := make([]domain.SelectableItem, len(snippets))
    for i, s := range snippets {
        items[i] = domain.SelectableItem{
            ID:          s.ID,
            DisplayText: fmt.Sprintf("[%s]: %s %s", s.Description, s.Command, formatTags(s.Tag)),
            Value:       s,
        }
    }
    return items
}
```

## Dependencies

- Phase 3: Parameter Input Abstraction (#XXX)

## Estimated Effort

Medium (3-5 days)

## Labels

`enhancement`, `refactoring`, `architecture`, `phase-4`, `testing`
