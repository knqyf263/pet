# Architecture Refactoring: Hexagonal Architecture for Testability and CLI-First Design

## Problem Statement

Pet currently faces two major challenges:

### 1. **Testability Issues**
- The application is CUI-first with tightly coupled I/O operations
- Features chain together in a pipeline that involves blocking user input
- Testing requires mocking complex interactions with fzf/peco and gocui
- The `filter()` function in `cmd/util.go` is monolithic: loads snippets → filters tags → calls external selector → shows parameter dialog → substitutes values
- Global state in `dialog/params.go` (CurrentCommand, FinalCommand, views) makes testing harder

### 2. **Limited CLI Usage**
- Pet is designed CUI-first, making it hard to use in non-interactive contexts
- No way to execute commands directly: `pet exec "docker ps"` doesn't exist
- Can't pass parameters via flags: `pet exec -d "docker run" --param container=nginx`
- Users want one-line commands but are forced through interactive selection
- This limits Pet's potential for scripting and automation

## Proposed Solution: Hexagonal Architecture (Ports & Adapters)

Refactor Pet to use a hexagonal/ports-and-adapters architecture where:
- **Core business logic** is pure, testable, and has no I/O dependencies
- **Adapters** handle I/O (gocui, fzf/peco, file system, etc.)
- **Ports** are interfaces that allow swapping implementations
- **CLI commands** become thin layers that compose core logic with adapters

### Architecture Diagram

```
┌─────────────────────────────────────────────────────┐
│                   Adapters (UI Layer)               │
│  ┌──────────────┐  ┌─────────────┐  ┌────────────┐ │
│  │ CLI Adapter  │  │ CUI Adapter │  │Future:REST │ │
│  │ (Direct)     │  │ (Interactive)│  │  API       │ │
│  └──────┬───────┘  └──────┬──────┘  └─────┬──────┘ │
└─────────┼──────────────────┼───────────────┼────────┘
          │                  │               │
┌─────────▼──────────────────▼───────────────▼────────┐
│              Application Ports (Interfaces)          │
│  - SnippetRepository                                 │
│  - ParameterInput                                    │
│  - Selector                                          │
│  - CommandExecutor                                   │
└──────────────────────┬───────────────────────────────┘
                       │
┌──────────────────────▼───────────────────────────────┐
│            Core Business Logic (Pure)                │
│  - Snippet CRUD operations                           │
│  - Parameter extraction/substitution                 │
│  - Tag filtering                                     │
│  - Command validation                                │
│  - NO I/O, NO UI dependencies                        │
└──────────────────────┬───────────────────────────────┘
                       │
┌──────────────────────▼───────────────────────────────┐
│            Repository Implementations                │
│  - FileRepository (TOML)                             │
│  - GistRepository                                    │
│  - GitLabRepository                                  │
└──────────────────────────────────────────────────────┘
```

### Key Interfaces (Ports)

```go
// Core service
type SnippetService interface {
    GetAll() ([]Snippet, error)
    GetByID(id string) (*Snippet, error)
    GetByDescription(desc string) (*Snippet, error)
    FilterByTags(tags []string) ([]Snippet, error)
    Create(snippet Snippet) error
    Update(id string, snippet Snippet) error
    Delete(id string) error
    ExtractParameters(command string) []Parameter
    SubstituteParameters(command string, params map[string]string) string
}

// Storage abstraction
type SnippetRepository interface {
    Load() ([]Snippet, error)
    Save(snippets []Snippet) error
    FindByID(id string) (*Snippet, error)
    FindByDescription(desc string) (*Snippet, error)
}

// User interaction abstraction
type ParameterInput interface {
    GetParameters(params []Parameter, command string) (map[string]string, error)
}

type Selector interface {
    Select(items []SelectableItem) (*SelectableItem, error)
}
```

### Example: gocui as an Adapter

The key insight is that gocui's event loop can be hidden behind a simple interface:

```go
// Adapter that implements ParameterInput
type GocuiParameterInput struct {
    gui         *gocui.Gui
    paramValues map[string]string  // Instance state (not global!)
}

func (g *GocuiParameterInput) GetParameters(params []Parameter, command string) (map[string]string, error) {
    // Setup gocui, create views, setup keybindings
    gui := gocui.NewGui(...)
    defer gui.Close()

    gui.SetManagerFunc(...)
    g.setupKeybindings()

    // MainLoop blocks until user presses ENTER
    gui.MainLoop()

    // Return collected values
    return g.paramValues, nil
}
```

Callers don't know about the event loop:
```go
paramInput := adapters.NewGocuiParameterInput()
values, err := paramInput.GetParameters(params, command) // Blocks but looks synchronous
```

## Benefits

### 1. **Testability**
- Core logic is pure functions - easy to unit test
- Mock adapters for testing: `MockParameterInput`, `MockSelector`
- No need to spawn actual gocui or fzf processes in tests
- Test coverage can be comprehensive

### 2. **CLI-First Capability**
```bash
# Interactive (current behavior, backward compatible)
pet exec

# Direct execution by description
pet exec -d "docker ps"

# Direct with parameters
pet exec -d "docker run" --param container=nginx --param port=8080

# Direct by ID
pet exec -i "abc123"

# Use defaults without prompting
pet exec -d "docker ps" --use-defaults
```

### 3. **Maintainability**
- Clear separation of concerns
- Easy to add new adapters (REST API, HTTP server, etc.)
- Business logic changes don't affect UI
- UI changes don't affect business logic

### 4. **Composability**
- Different commands can compose the same core logic differently
- Easy to add new workflows without duplicating code

## Migration Strategy

We have two options:

### Option A: Gradual Migration (Recommended)
Refactor incrementally on `main` branch:
- Maintains backward compatibility at each step
- Users benefit from improvements gradually
- Less risky, easier to review
- Can be released as incremental minor versions

### Option B: Big Bang 2.0
Create a `v2-architecture` branch:
- Complete rewrite on separate branch
- Merge as major 2.0 release
- Higher risk but cleaner end result
- Harder to review, longer feedback cycle

**Recommendation: Option A (Gradual)** - Less magical, easier to understand each step, maintains project momentum.

## Implementation Plan

### Phase 1: Foundation (v1.1.0)
- [ ] #XXX: Create `domain/` package with core types and interfaces
- [ ] #XXX: Create `infrastructure/` package for repository implementations
- [ ] #XXX: Move snippet data structures to `domain/snippet.go`
- [ ] #XXX: Define core interfaces (SnippetRepository, ParameterInput, Selector)
- [ ] #XXX: No breaking changes - just organizing code

**Deliverable**: New packages with interfaces, existing code still works

### Phase 2: Repository Pattern (v1.2.0)
- [ ] #XXX: Implement FileRepository wrapping current TOML logic
- [ ] #XXX: Implement SnippetService with repository injection
- [ ] #XXX: Add comprehensive tests for repository layer
- [ ] #XXX: Update `snippet/` package to use repository pattern

**Deliverable**: Testable storage layer, backward compatible

### Phase 3: Parameter Input Abstraction (v1.3.0)
- [ ] #XXX: Create GocuiParameterInput adapter (refactor current dialog/params.go)
- [ ] #XXX: Remove global state from dialog package
- [ ] #XXX: Create FlagParameterInput adapter for CLI mode
- [ ] #XXX: Create MockParameterInput for testing
- [ ] #XXX: Add tests for parameter extraction and substitution

**Deliverable**: Multiple parameter input methods, testable

### Phase 4: Selector Abstraction (v1.4.0)
- [ ] #XXX: Create Selector interface
- [ ] #XXX: Implement FzfSelector adapter (wrap current fzf logic)
- [ ] #XXX: Implement PecoSelector adapter
- [ ] #XXX: Implement DirectSelector for exact match (CLI mode)
- [ ] #XXX: Implement MockSelector for testing
- [ ] #XXX: Add tests for selection logic

**Deliverable**: Multiple selection methods, testable

### Phase 5: Refactor Commands (v1.5.0)
- [ ] #XXX: Refactor `exec` command to use new architecture
- [ ] #XXX: Refactor `search` command to use new architecture
- [ ] #XXX: Refactor `clip` command to use new architecture
- [ ] #XXX: Refactor `new` command to use new architecture
- [ ] #XXX: Refactor remaining commands
- [ ] #XXX: Remove old `cmd/util.go` filter() function
- [ ] #XXX: Add integration tests for commands

**Deliverable**: All commands using clean architecture, backward compatible

### Phase 6: CLI-First Features (v1.6.0)
- [ ] #XXX: Add direct execution flags to `exec` command
- [ ] #XXX: Add direct execution flags to `search` command
- [ ] #XXX: Add direct execution flags to `clip` command
- [ ] #XXX: Add `--use-defaults` flag to skip parameter prompts
- [ ] #XXX: Add `--param key=value` flag for parameter passing
- [ ] #XXX: Update documentation with CLI-first examples
- [ ] #XXX: Add shell completion for parameter names

**Deliverable**: Full CLI-first capability, new features

### Phase 7: Documentation & Polish (v2.0.0)
- [ ] #XXX: Update README with architecture overview
- [ ] #XXX: Add ARCHITECTURE.md explaining design
- [ ] #XXX: Add developer guide for contributing
- [ ] #XXX: Update all command help text
- [ ] #XXX: Add examples for CLI and CUI modes
- [ ] #XXX: Deprecation notices for old patterns (if any)
- [ ] #XXX: Performance benchmarks

**Deliverable**: Complete documentation, ready for 2.0 release

## Success Metrics

- [ ] Test coverage increases from ~X% to >80%
- [ ] All commands testable without mocking I/O
- [ ] CLI-first usage documented and working
- [ ] No breaking changes to existing CUI workflows
- [ ] Performance equivalent or better than current

## Open Questions

1. Should we add a `-d/--description` flag or use positional args for direct selection?
2. How should parameter defaults work in CLI mode? Error if missing or use defaults?
3. Should we support regex matching in DirectSelector or only exact matches?
4. Do we want to expose the core as a Go library for other tools to use?
5. Should repository pattern support multiple snippet files or continue with single file?

## Related Issues

This proposal addresses:
- Testing difficulties (no existing issue, but mentioned in discussions)
- CLI-first usage requests (need to search for related issues)
- Code maintainability concerns

## References

- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Ports and Adapters Pattern](https://herbertograca.com/2017/09/14/ports-adapters-architecture/)
- [Testing with Interfaces](https://go.dev/blog/laws-of-reflection)

---

**Next Steps:**
1. Gather feedback on this proposal in GitHub Discussions
2. Refine approach based on community input
3. Create milestone and individual task issues
4. Begin Phase 1 implementation

**Estimated Timeline:** 2-3 months for gradual migration (depending on review cycles)

**Breaking Changes:** None until 2.0 (if needed)
