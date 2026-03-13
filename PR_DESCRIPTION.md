# Phase 1 + 2: Foundation and Repository Pattern - Hexagonal Architecture

This PR implements Phases 1 and 2 of the Pet architecture refactoring, establishing the foundation for hexagonal architecture (ports & adapters pattern) and adding a file-based repository implementation with full backward compatibility.

## Phase 1: Foundation - Domain Layer

**What's Included:**

### Core Domain Models
- `Snippet` struct with validation, tag matching
- `Parameter` struct for command placeholders
- Pure business logic with zero I/O dependencies

### Repository Interface (Port)
- `SnippetRepository` interface for storage abstraction
- `MockRepository` for testing
- Methods: `Load()`, `Save()`, `FindByDescription()`, `FindByID()`, `FilterByTags()`

### Selector Interface (Port)
- `Selector` interface for snippet selection abstraction
- `DirectSelector` - exact match (for CLI-first mode)
- `FirstMatchSelector` - auto-select first match
- `MockSelector` for testing

### ParameterInput Interface (Port)
- `ParameterInput` interface for getting parameter values
- `DefaultParameterInput` - uses defaults without prompting
- `FlagParameterInput` - from CLI flags (with fallback to defaults)
- `PositionalParameterInput` - from positional args
- `MockParameterInput` for testing

### SnippetService (Business Logic)
- `ExtractParameters()` - parses `<param>`, `<param=default>`, `<param=|_v1_||_v2_|>`
- `SubstituteParameters()` - replaces placeholders with values
- Repository operations: `GetAll()`, `GetByDescription()`, `FilterByTags()`, `Create()`, `Update()`, `Delete()`
- Full validation and error handling
- Optimized: regex compiled once at package level

## Phase 2: Repository Pattern - FileRepository

**What's Included:**

### FileRepository Implementation
- Implements `domain.SnippetRepository` using TOML files
- Does NOT depend on global `config.Conf` - receives paths via constructor for testability
- TOML format identical to existing format (verified by compatibility tests)
- Supports multi-file and directory configurations

### Path Validation
- Validates snippet filenames are in valid locations (main file or snippet directories)
- Prevents Save/Load drift where snippets can be saved but not loaded back
- Returns clear errors for invalid paths

### Error Handling
- Logs errors when files are inaccessible during directory walks
- Doesn't silently skip corrupt/unreadable files
- Proper error messages for all failure modes

### Compatibility
- Both old (`snippet/`) and new (`infrastructure/`) code can read/write the same files
- Compatibility tests verify round-trip in both directions
- Multiline commands preserved correctly

## Migration Documentation

**`MIGRATION.md`** documents:
- Full 7-phase migration plan with architecture diagrams
- How old and new code coexist during transition
- Which GitHub issues will be resolved in Phase 6
- What changed (and didn't change) in each phase

## Test Coverage

### Phase 1 Tests (100+ test cases)
- Core models: `Snippet`, `Parameter` validation and tag matching
- All interfaces: Repository, Selector, ParameterInput with multiple implementations
- SnippetService: parameter extraction, substitution, CRUD operations
- Integration tests showing full workflows

### Phase 2 Tests (20+ test cases)
- FileRepository: load, save, multi-file, directories, error cases
- Path validation: valid/invalid paths, drift prevention
- Compatibility: bidirectional with existing `snippet.Snippets`
- Round-trip: multiline commands, tags, all fields preserved

**All tests passing** ✅ - both new tests and all existing project tests.

## Files Added

```
domain/
├── snippet.go, snippet_test.go           # Core models
├── repository.go, repository_test.go     # Storage interface (port)
├── selector.go, selector_test.go         # Selection interface (port)
├── parameter_input.go, parameter_input_test.go  # Input interface (port)
└── snippet_service.go, snippet_service_test.go  # Business logic

infrastructure/
├── file_repository.go, file_repository_test.go  # TOML storage implementation
├── compat_test.go                        # Compatibility with existing code
└── validation_test.go                    # Path validation tests

MIGRATION.md                              # Full migration plan
```

## What Was NOT Changed

All existing code remains untouched:
- `snippet/` package - existing TOML loading/saving
- `cmd/` package - all commands continue to work
- `dialog/` package - gocui parameter dialog
- `config/` package - configuration handling

The `domain/` and `infrastructure/` packages are purely additive.

## Benefits

### 1. Testability ✅
- Core logic is pure functions - trivial to unit test
- Mock adapters eliminate need for actual gocui/fzf in tests
- Test coverage can be comprehensive

### 2. Flexibility ✅
- Interfaces allow swapping implementations (interactive CUI vs direct CLI)
- Easy to add new adapters (REST API, HTTP server) in future
- Business logic changes don't affect UI and vice versa

### 3. Performance ✅
- Regex compiled once at package level (not on every call)
- Efficient file operations with proper error handling

### 4. Maintainability ✅
- Clear separation of concerns
- Easy to understand and modify
- Well-documented with MIGRATION.md

## What's Next

### Phase 3: Parameter Input Abstraction (Planned)
Refactor `dialog/` package to wrap gocui behind `ParameterInput` interface

### Phase 4: Selector Abstraction (Planned)
Wrap fzf/peco selection behind `Selector` interface

### Phase 5: Refactor Commands (Planned)
Update all commands in `cmd/` to use new architecture

### Phase 6: CLI-First Features (Planned)
Add direct execution flags:
```bash
pet exec "docker ps"
pet exec -d "docker run" --param container=nginx --param port=8080
pet exec --first-match -q "docker"
pet list --json
pet new --command "ls -la" --tag list --description "List files"
```

**This will resolve GitHub issues:** #73, #84, #117, #147, #164, #352

### Phase 7: Documentation & Polish (Planned)
Update README, add architecture docs, prepare v2.0 release

## Testing

```bash
# Run all tests
go test ./...

# Run domain tests
go test ./domain/... -v

# Run infrastructure tests (includes compatibility tests)
go test ./infrastructure/... -v

# Run with coverage
go test ./... -cover
```

---

**This PR is safe to merge** - it adds new, well-tested code but doesn't modify any existing functionality. All existing tests continue to pass.
