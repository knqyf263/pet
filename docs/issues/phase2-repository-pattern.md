# Phase 2: Repository Pattern Implementation

**Part of:** Architecture Refactoring (Main Proposal: #XXX)
**Depends on:** Phase 1 (#XXX)

## Objective

Implement the Repository pattern to abstract storage operations, making the codebase more testable and maintainable.

## Tasks

- [ ] Implement `FileRepository` in `infrastructure/repository/file_repository.go`
  - [ ] Wrap existing TOML loading logic from `snippet/`
  - [ ] Implement `Load()` method
  - [ ] Implement `Save()` method
  - [ ] Implement `FindByID()` method
  - [ ] Implement `FindByDescription()` method
  - [ ] Add path validation and error handling

- [ ] Implement `SnippetService` in `domain/service.go`
  - [ ] Constructor with repository injection
  - [ ] Implement CRUD operations
  - [ ] Implement filtering logic (tags, search)
  - [ ] Implement ordering logic

- [ ] Create mock repository for testing
  - [ ] `infrastructure/repository/mock_repository.go`
  - [ ] In-memory storage for tests

- [ ] Add comprehensive tests
  - [ ] Unit tests for FileRepository
  - [ ] Unit tests for SnippetService
  - [ ] Test error cases and edge cases
  - [ ] Test with mock repository

- [ ] Update existing `snippet/` package to delegate to repository
  - [ ] Maintain backward compatibility
  - [ ] Add deprecation notices (optional)

## Acceptance Criteria

- [ ] FileRepository fully implements SnippetRepository interface
- [ ] SnippetService works with any repository implementation
- [ ] Can swap FileRepository for MockRepository in tests
- [ ] Test coverage for repository layer >80%
- [ ] Existing commands still work unchanged
- [ ] Tests pass: `make test`
- [ ] No performance regression

## Implementation Notes

### Example Repository Usage

```go
// Old way (direct)
snippets, err := snippet.Load()

// New way (with repository)
repo := repository.NewFileRepository(config.Conf.Snippetfile)
service := domain.NewSnippetService(repo)
snippets, err := service.GetAll()
```

### Testing Example

```go
func TestSnippetService_FilterByTags(t *testing.T) {
    mockRepo := repository.NewMockRepository()
    mockRepo.AddSnippet(snippet1)
    mockRepo.AddSnippet(snippet2)

    service := domain.NewSnippetService(mockRepo)
    result := service.FilterByTags([]string{"docker"})

    // Easy to test - no file I/O!
    assert.Len(t, result, 1)
}
```

## Dependencies

- Phase 1: Foundation (#XXX)

## Estimated Effort

Medium (3-5 days)

## Labels

`enhancement`, `refactoring`, `architecture`, `phase-2`, `testing`
