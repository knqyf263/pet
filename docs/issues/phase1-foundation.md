# Phase 1: Foundation - Create Core Domain Layer

**Part of:** Architecture Refactoring (Main Proposal: #XXX)

## Objective

Create the foundational packages and interfaces for the hexagonal architecture without breaking existing functionality.

## Tasks

- [ ] Create `domain/` package structure
  - [ ] `domain/snippet.go` - Core Snippet type and business logic
  - [ ] `domain/interfaces.go` - Define core interfaces (Repository, Service)
  - [ ] `domain/parameter.go` - Parameter extraction/substitution logic
  - [ ] `domain/errors.go` - Domain-specific error types

- [ ] Create `infrastructure/` package structure
  - [ ] `infrastructure/repository/` - Repository implementations

- [ ] Define core interfaces:
  ```go
  type SnippetRepository interface
  type SnippetService interface
  type ParameterExtractor interface
  ```

- [ ] Move snippet types from `snippet/` to `domain/` (maintain backward compatibility)

- [ ] Add initial unit tests for domain logic

## Acceptance Criteria

- [ ] New packages created with clear documentation
- [ ] Core interfaces defined and documented
- [ ] Existing code still compiles and runs
- [ ] No breaking changes to public API
- [ ] Tests pass: `make test`
- [ ] Code coverage maintained or improved

## Implementation Notes

- Keep existing `snippet/` package temporarily - we'll migrate incrementally
- Add type aliases in `snippet/` to point to new `domain/` types for backward compatibility
- Focus on structure, not full implementation yet

## Dependencies

None - this is the first step

## Estimated Effort

Small (1-2 days)

## Labels

`enhancement`, `refactoring`, `architecture`, `phase-1`
