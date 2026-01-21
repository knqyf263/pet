# Architecture Refactoring: Hexagonal Architecture for Testability and CLI-First Design

## Problem Statement

Pet currently faces two major architectural challenges:

### 1. Testing is Extremely Difficult
- The application is CUI-first with tightly coupled I/O operations
- Features chain together in a pipeline requiring user interaction (load → filter → fzf/peco → gocui → substitute)
- The monolithic `filter()` function in `cmd/util.go` is nearly impossible to unit test
- Global state in `dialog/params.go` makes testing even harder

### 2. CLI Potential is Wasted
- Pet is designed CUI-first, requiring interactive selection for everything
- No way to execute commands directly: `pet exec "docker ps"` doesn't exist
- Can't pass parameters via flags for scripting/automation
- Users want one-line commands but are forced through interactive prompts

## Proposed Solution

Refactor Pet to use **Hexagonal Architecture (Ports & Adapters)** where:
- **Core business logic** is pure, testable, with zero I/O dependencies
- **Adapters** handle all I/O (gocui, fzf/peco, files, etc.)
- **Interfaces** allow swapping implementations (interactive vs direct, real vs mock)
- **CLI commands** become thin composition layers

### Architecture Overview

```
┌─────────────────────────────────────────────┐
│         Adapters (UI Layer)                 │
│  ┌──────────┐  ┌──────────┐  ┌───────────┐ │
│  │CLI Direct│  │CUI Inter.│  │Future:API │ │
│  └────┬─────┘  └────┬─────┘  └─────┬─────┘ │
└───────┼─────────────┼──────────────┼───────┘
        │             │              │
┌───────▼─────────────▼──────────────▼───────┐
│     Application Ports (Interfaces)         │
│  - SnippetRepository                       │
│  - ParameterInput                          │
│  - Selector                                │
└──────────────────┬─────────────────────────┘
                   │
┌──────────────────▼─────────────────────────┐
│    Core Business Logic (Pure, No I/O)     │
│  - Snippet CRUD                            │
│  - Parameter extraction/substitution       │
│  - Tag filtering, validation               │
└──────────────────┬─────────────────────────┘
                   │
┌──────────────────▼─────────────────────────┐
│    Repository Implementations              │
│  - FileRepository (TOML)                   │
│  - GistRepository, GitLabRepository        │
└────────────────────────────────────────────┘
```

### Key Insight: gocui as an Adapter

Even though gocui has a blocking event loop, we hide it behind a simple interface:

```go
// Interface
type ParameterInput interface {
    GetParameters(params []Parameter, command string) (map[string]string, error)
}

// Implementations
- GocuiParameterInput (interactive, current behavior)
- FlagParameterInput (CLI flags, new feature)
- MockParameterInput (testing, new capability)
```

Callers don't know about the event loop - they just call `GetParameters()` and get values back.

## Benefits

### Testability ✅
- Core logic is pure functions - trivial to unit test
- Mock adapters eliminate need for actual gocui/fzf in tests
- Test coverage can reach >80%

### CLI-First Capability ✅
```bash
# Interactive (current, backward compatible)
pet exec

# Direct execution (NEW)
pet exec -d "docker ps"
pet exec -d "docker run" --param container=nginx --param port=8080
pet exec -i "abc123" --use-defaults
```

### Maintainability ✅
- Clear separation of concerns
- Easy to add new adapters (REST API, HTTP server)
- Business logic changes don't affect UI and vice versa

## Implementation Roadmap

We'll do a **gradual migration** (not a big-bang rewrite) to keep changes reviewable and maintain backward compatibility.

### Phase 1: Foundation (v1.1.0)
Create `domain/` package with core types and interfaces. No breaking changes.

### Phase 2: Repository Pattern (v1.2.0)
Abstract storage behind `SnippetRepository` interface. Add tests.

### Phase 3: Parameter Input Abstraction (v1.3.0)
Refactor gocui into `GocuiParameterInput` adapter. Remove global state. Add `FlagParameterInput` and `MockParameterInput`.

### Phase 4: Selector Abstraction (v1.4.0)
Abstract fzf/peco behind `Selector` interface. Add `DirectSelector` for CLI mode.

### Phase 5: Refactor Commands (v1.5.0)
Update all commands to use new architecture. Remove monolithic `filter()`. Add integration tests.

### Phase 6: CLI-First Features (v1.6.0)
Add new flags: `-d/--description`, `-i/--id`, `--param`, `--use-defaults`, `--json`. Enable scripting!

### Phase 7: Documentation & Polish (v2.0.0)
Complete ARCHITECTURE.md, update README, add examples, prepare release.

**Timeline:** 2-3 months for gradual migration
**Breaking Changes:** None until 2.0 (if needed)

## Success Metrics

- [ ] Test coverage increases to >80%
- [ ] All commands testable without mocking I/O
- [ ] CLI-first usage works: `pet exec -d "name" --param x=1`
- [ ] No breaking changes to existing CUI workflows
- [ ] Performance equivalent or better

## Migration Strategy

**Gradual Migration (Recommended)** - Each phase:
- Is independently reviewable and releasable
- Maintains 100% backward compatibility
- Adds value incrementally
- Reduces risk

**NOT Big Bang** - Avoids:
- Long-lived feature branches
- Massive PRs that are hard to review
- Extended periods without releases
- High risk of breaking changes

## Next Steps

1. **Gather feedback** on this proposal (Discussion or comments)
2. **Refine approach** based on community input
3. **Create phase issues** and milestone when ready
4. **Begin Phase 1** implementation

## Questions for Discussion

1. Should we use `-d/--description` flag or positional args for direct selection?
2. How should missing parameters work in CLI mode - error or use defaults?
3. Do we want to expose core logic as a Go library for other tools?
4. Should we support regex matching in DirectSelector or only exact matches?

## References

- Full detailed proposal: `docs/architecture-refactoring-proposal.md`
- Phase-specific task breakdowns: `docs/issues/phase*.md`
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Ports and Adapters Pattern](https://herbertograca.com/2017/09/14/ports-adapters-architecture/)

---

**This is a proposal for discussion** - no code changes yet. Feedback welcome!
