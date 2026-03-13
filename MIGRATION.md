# Architecture Migration Plan

Pet is being refactored from a CUI-first architecture to a CLI-first, testable
architecture using the **Ports & Adapters** (Hexagonal Architecture) pattern.

This document describes the migration plan, current progress, and how the old
and new code coexist during the transition.

## Goals

1. **Testability**: Core business logic should be testable without I/O (no
   gocui, no fzf/peco, no file system).
2. **CLI-First**: Users should be able to execute snippets directly without
   interactive prompts (e.g. `pet exec "docker ps" --param port=8080`).
3. **Backward Compatibility**: The existing interactive CUI workflow must
   continue to work exactly as before.

## Architecture Overview

```
┌───────────────────────────────────────────────┐
│            Adapters (UI Layer)                 │
│  ┌────────────┐  ┌────────────┐  ┌──────────┐ │
│  │ CLI Direct │  │ CUI (fzf/  │  │  Future  │ │
│  │ (flags,    │  │ gocui)     │  │  (API,   │ │
│  │ positional)│  │            │  │   HTTP)  │ │
│  └─────┬──────┘  └─────┬──────┘  └────┬─────┘ │
└────────┼───────────────┼──────────────┼───────┘
         │               │              │
┌────────▼───────────────▼──────────────▼───────┐
│         Application Ports (Interfaces)        │
│  SnippetRepository · Selector · ParameterInput│
└────────────────────┬──────────────────────────┘
                     │
┌────────────────────▼──────────────────────────┐
│       Core Business Logic (domain/)           │
│  SnippetService · ExtractParameters ·         │
│  SubstituteParameters · Validate              │
│  Pure functions — zero I/O dependencies       │
└────────────────────┬──────────────────────────┘
                     │
┌────────────────────▼──────────────────────────┐
│       Storage (infrastructure/)               │
│  FileRepository (TOML — compatible with       │
│  existing snippet files)                      │
└───────────────────────────────────────────────┘
```

## Migration Phases

### Phase 1: Foundation ✅

Create the `domain/` package with pure business logic and interfaces.

**What was added:**
- `domain/snippet.go` — `Snippet` and `Parameter` domain models
- `domain/repository.go` — `SnippetRepository` interface (port)
- `domain/selector.go` — `Selector` interface + `DirectSelector`,
  `FirstMatchSelector` implementations
- `domain/parameter_input.go` — `ParameterInput` interface +
  `DefaultParameterInput`, `FlagParameterInput`, `PositionalParameterInput`
- `domain/snippet_service.go` — `SnippetService` with parameter extraction,
  substitution, CRUD operations

**What was NOT changed:**
- No existing code was modified. The `domain/` package is entirely new.

**Test coverage:** 100+ test cases, all passing.

---

### Phase 2: Repository Pattern ✅

Create `infrastructure/FileRepository` that implements `SnippetRepository`
using the same TOML format as the existing `snippet/` package.

**What was added:**
- `infrastructure/file_repository.go` — `FileRepository` implementing
  `domain.SnippetRepository`
- Compatibility tests verifying round-trip between `FileRepository` and the
  existing `snippet.Snippets` code

**Key design decisions:**
- `FileRepository` does NOT depend on the global `config.Conf` — it receives
  the snippet file path and directories via constructor.
- TOML format is identical to the existing format, ensuring full backward
  compatibility with existing snippet files.
- Both the old and new code can read/write the same files.

**What was NOT changed:**
- No existing code was modified. The `infrastructure/` package is new.
- The existing `snippet/` package continues to work as-is.

---

### Phase 3: Parameter Input Abstraction (Planned)

Refactor the `dialog/` package to wrap gocui behind the `ParameterInput`
interface. This removes global state (`CurrentCommand`, `FinalCommand`, `views`)
and encapsulates it in a `GocuiParameterInput` adapter.

**What will change:**
- `dialog/params.go` global variables will move into the adapter struct
- The gocui event loop will be hidden behind `GetParameters()`
- No user-visible behavior changes

---

### Phase 4: Selector Abstraction (Planned)

Wrap fzf/peco selection behind the `Selector` interface. Add `FzfSelector`
and `PecoSelector` adapters.

**What will change:**
- The `filter()` function in `cmd/util.go` will be decomposed
- Selection logic will move into adapter implementations
- No user-visible behavior changes

---

### Phase 5: Refactor Commands (Planned)

Update all commands in `cmd/` to use `SnippetService` + adapters instead of
the monolithic `filter()` function.

**What will change:**
- Commands compose `SnippetService` + `Selector` + `ParameterInput`
- The `filter()` function will be removed
- Default behavior (no flags) remains interactive (backward compatible)

---

### Phase 6: CLI-First Features (Planned)

Add new CLI flags to enable non-interactive usage:

```bash
# Direct execution by description
pet exec "docker ps"
pet exec -d "docker run" --param container=nginx --param port=8080

# First-match auto-execution
pet exec --first-match -q "docker"

# Use default values without prompting
pet exec "docker ps" --use-defaults

# JSON output for scripting
pet list --json

# Non-interactive snippet creation
pet new --command "ls -la" --tag list --description "List files"
```

**GitHub issues that will be resolved:**
- #73 — Positional args for search
- #84 — JSON output for automation
- #117 — Pass parameters as arguments
- #147 — Exec without additional search
- #164 — First match shortcut
- #352 — Non-interactive `pet new`

---

### Phase 7: Documentation & Polish (Planned)

Update README with CLI-first examples, add architecture documentation, prepare
for release.

---

## How Old and New Code Coexist

During migration, both systems work side by side:

```
cmd/                    ← Uses old snippet/ package (unchanged)
  exec.go
  search.go
  ...

snippet/                ← Existing code (unchanged)
  snippet.go            ← SnippetInfo, Load(), Save(), FilterByTags()
  util.go               ← getFiles()

domain/                 ← NEW: Pure business logic
  snippet.go            ← Snippet, Parameter models
  repository.go         ← SnippetRepository interface
  selector.go           ← Selector interface + implementations
  parameter_input.go    ← ParameterInput interface + implementations
  snippet_service.go    ← ExtractParameters, SubstituteParameters

infrastructure/         ← NEW: File-based repository
  file_repository.go    ← FileRepository (reads same TOML files)
```

The `domain/` and `infrastructure/` packages are additive — they don't modify
or replace any existing code. Commands in `cmd/` continue to use the old
`snippet/` package until Phase 5, when they are migrated one at a time.

Both `snippet.Snippets.Load()/Save()` and `FileRepository.Load()/Save()`
read and write the **same TOML format**, so files created by either system
are fully compatible with the other.

## Running Tests

```bash
# Run all tests (old + new)
go test ./...

# Run only domain tests
go test ./domain/... -v

# Run only infrastructure tests (includes compatibility tests)
go test ./infrastructure/... -v
```
