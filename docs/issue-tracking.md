# GitHub Issue Resolution Tracking

This document tracks which GitHub issues will be resolved by each phase of the architecture refactoring.

## Issue Categories

### 🎯 Direct Resolution (Will be Fully Resolved)
Issues that will be completely satisfied by the refactoring.

### 📊 Partial Resolution (Significant Progress)
Issues that will be significantly improved but may need follow-up work.

### 🔗 Related (Benefits from Refactoring)
Issues that aren't directly resolved but benefit from the new architecture.

---

## Phase 1: Foundation (v1.1.0)
**Goal**: Create domain layer with core types and interfaces

### Issues Resolved: None (Foundation work)
This phase sets up the architecture but doesn't close any user-facing issues yet.

---

## Phase 2: Repository Pattern (v1.2.0)
**Goal**: Abstract storage operations, improve testability

### Issues Resolved: None (Infrastructure work)
Improves test coverage but no user-facing features yet.

---

## Phase 3: Parameter Input Abstraction (v1.3.0)
**Goal**: Abstract parameter input (gocui, flags, defaults)

### Issues Resolved: None yet
Creates the foundation for CLI-first parameter passing, but doesn't expose it to users yet.

### Preparation for:
- #117 - Arguments for parameters
- #154 - Paste to prompt functionality
- #178 - Transfer command to prompt

---

## Phase 4: Selector Abstraction (v1.4.0)
**Goal**: Abstract snippet selection (fzf, peco, direct match)

### Issues Resolved: None yet
Creates the foundation for direct selection, but doesn't expose CLI flags yet.

### Preparation for:
- #147 - Exec without additional search
- #164 - Enter first match
- #180 - Search with query argument
- #73 - Search with positional args

---

## Phase 5: Refactor Commands (v1.5.0)
**Goal**: Update all commands to use new architecture

### Issues Resolved: None yet
Makes the architecture work across all commands, but still no new user-facing features.

---

## Phase 6: CLI-First Features (v1.6.0)
**Goal**: Expose CLI-first flags and direct execution

This is where we start closing issues! 🎉

### 🎯 Direct Resolution

#### **#147** - pet exec without additional search and enter
- **Status**: Open
- **Link**: https://github.com/knqyf263/pet/issues/147
- **Resolution**: Add `--first-match` flag
- **Syntax**:
  ```bash
  pet exec -q "private" --first-match
  ```
- **Implementation**: Use `FirstMatchSelector` adapter
- **Can Close**: ✅ Yes, after implementing `--first-match` flag

#### **#164** - Enter first match with shortcut key
- **Status**: Open
- **Link**: https://github.com/knqyf263/pet/issues/164
- **Resolution**: Same as #147, use `--first-match` flag
- **Syntax**:
  ```bash
  pet exec --first-match checkexp
  ```
- **Implementation**: Use `FirstMatchSelector` adapter
- **Can Close**: ✅ Yes, after implementing `--first-match` flag

#### **#180** - Add support for `pet search <search term>`
- **Status**: Closed (already completed)
- **Link**: https://github.com/knqyf263/pet/issues/180
- **Resolution**: This was already implemented! But we can enhance it:
  - Already supports: `pet search -q <term>`
  - Can add: `pet search <term>` (positional arg)
  - Can add: `--exec` flag for auto-execution
  - Can add: `--list` flag for non-interactive listing
- **New Syntax**:
  ```bash
  pet search docker              # Positional arg
  pet search docker --exec       # Auto-execute first match
  pet search docker --list       # List without fzf
  ```
- **Can Close**: Already closed, but we can add enhancements

#### **#73** - pet search [flags] <string> to query with an initial value
- **Status**: Open
- **Link**: https://github.com/knqyf263/pet/issues/73
- **Resolution**: Add positional argument support
- **Syntax**:
  ```bash
  pet search docker       # Instead of: pet search -q docker
  pet exec docker-ps      # Instead of interactive selection
  ```
- **Implementation**: Parse `args[0]` as query if present
- **Can Close**: ✅ Yes, after adding positional arg support

#### **#117** - [Feature] arguments for parameter
- **Status**: Closed
- **Link**: https://github.com/knqyf263/pet/issues/117
- **Resolution**: Add `--param` flag for named parameters
- **Syntax**:
  ```bash
  # Named parameters
  pet exec ffmpeg-to-gif --param input=./video.mkv --param output=./out.gif --param fps=15

  # Positional arguments (map to params in order)
  pet exec ffmpeg-to-gif ./video.mkv ./out.gif 15 720
  ```
- **Implementation**:
  - `FlagParameterInput` adapter for `--param` flags
  - `PositionalParameterInput` adapter for positional args
- **Can Close**: ✅ Yes, after implementing parameter passing

#### **#84** - Option to output all metadata from search
- **Status**: Open
- **Link**: https://github.com/knqyf263/pet/issues/84
- **Resolution**: Add `--json` flag to output structured data
- **Syntax**:
  ```bash
  pet list --json
  pet search docker --json

  # Use case: automation
  pet list --json | jq '.[] | select(.tags[] == "docker") | .command'
  ```
- **Implementation**: Create `JSONPresenter` adapter
- **Can Close**: ✅ Yes, after implementing `--json` output

### 📊 Partial Resolution

#### **#352** - Allow pet new to accept arguments via flags
- **Status**: Open
- **Link**: https://github.com/knqyf263/pet/issues/352
- **Resolution**: Add flags to `pet new` command
- **Syntax**:
  ```bash
  pet new --command "ls -la" --tag "list" --tag "files" --description "List all files"

  # Or stdin support
  echo "ls -la" | pet new --tag "list" --description "List all files"
  ```
- **Implementation**: Add command-line flags to `pet new`
- **Can Close**: ✅ Yes, after implementing non-interactive `pet new`
- **Note**: Skip HEREDOC syntax (user doesn't want it)

#### **#146** - pet exec --last-command
- **Status**: Closed
- **Link**: https://github.com/knqyf263/pet/issues/146
- **Resolution**: Would require history tracking (separate feature)
- **Can Close**: ❌ Not directly resolved by architecture refactoring
- **Follow-up**: Could be added later with a `~/.pet/history.json` file

#### **#154** - Paste command to terminal prompt instead of instant execution
- **Status**: Closed
- **Link**: https://github.com/knqyf263/pet/issues/154
- **Resolution**: Would require terminal integration (out of scope)
- **Can Close**: ❌ Not resolved by architecture refactoring
- **Note**: This requires ANSI escape sequences or shell integration

#### **#178** - Transfer command to prompt instead of interactive box
- **Status**: Closed
- **Link**: https://github.com/knqyf263/pet/issues/178
- **Resolution**: Same as #154, requires terminal integration
- **Can Close**: ❌ Not resolved by architecture refactoring

#### **#163** - Set cursor point
- **Status**: Open
- **Link**: https://github.com/knqyf263/pet/issues/163
- **Resolution**: Partially solved by better parameter handling
- **Can Close**: ⚠️ Partially - parameter flags help, but cursor positioning needs terminal integration

### 🔗 Related Issues (Benefit from Refactoring)

#### **#367** - Architecture refactoring proposal
- **Status**: Open
- **Link**: https://github.com/knqyf263/pet/issues/367
- **Resolution**: This IS the implementation of #367
- **Can Close**: ✅ Yes, after completing all 7 phases (v2.0.0)

---

## Phase 7: Documentation & Polish (v2.0.0)
**Goal**: Complete documentation and release

### Issues Resolved

#### **#367** - Architecture refactoring (the main issue)
- **Can Close**: ✅ Yes, when all phases are complete

---

## Summary: Issues We Can Close

### After Phase 6 (v1.6.0):
1. ✅ **#147** - Exec without search (add `--first-match`)
2. ✅ **#164** - First match shortcut (same as #147)
3. ✅ **#73** - Positional args for search
4. ✅ **#117** - Arguments for parameters (add `--param` and positional args)
5. ✅ **#84** - JSON output (add `--json`)
6. ✅ **#352** - Non-interactive `pet new` (add flags)

### After Phase 7 (v2.0.0):
7. ✅ **#367** - Architecture refactoring complete

### Cannot Close (Out of Scope):
- ❌ **#146** - Last command (needs history tracking - separate feature)
- ❌ **#154** - Paste to prompt (needs terminal integration)
- ❌ **#178** - Transfer to prompt (needs terminal integration)
- ⚠️ **#163** - Cursor positioning (partially solved, full solution needs terminal integration)

---

## Implementation Checklist for Phase 6

When implementing Phase 6, add these features and track which issues they close:

### Feature 1: Direct Selection by Description/ID
- [ ] Add positional arg support: `pet exec <description>`
- [ ] Add `--id` flag: `pet exec --id <id>`
- [ ] Add `--description` / `-d` flag: `pet exec -d <description>`
- **Closes**: #73

### Feature 2: First Match Auto-Execution
- [ ] Add `--first-match` flag to `pet exec`
- [ ] Add `--first-match` flag to `pet search`
- **Closes**: #147, #164

### Feature 3: Parameter Passing
- [ ] Add `--param key=value` flag (repeatable)
- [ ] Add positional argument mapping to parameters
- [ ] Add `--use-defaults` flag to skip parameter prompts
- **Closes**: #117

### Feature 4: Structured Output
- [ ] Add `--json` flag to `pet list`
- [ ] Add `--json` flag to `pet search`
- [ ] Add `--json` flag to `pet exec` (dry-run mode)
- **Closes**: #84

### Feature 5: Non-Interactive Creation
- [ ] Add `--command` flag to `pet new`
- [ ] Add `--description` flag to `pet new`
- [ ] Add `--tag` flag to `pet new` (repeatable)
- [ ] Add stdin support to `pet new`
- **Closes**: #352

### Feature 6: Enhanced Search
- [ ] Support positional args: `pet search <query>`
- [ ] Add `--exec` flag to auto-execute
- [ ] Add `--list` flag for non-interactive listing
- **Enhances**: #180 (already closed, but improves it)

---

## Testing Requirements

For each feature above, ensure:
- [ ] Unit tests for core logic
- [ ] Integration tests for CLI commands
- [ ] Manual testing for user workflows
- [ ] Update documentation
- [ ] Add examples to README
- [ ] Update help text

---

## Communication Plan

When closing issues:
1. Comment on the issue explaining what was implemented
2. Show example usage with new syntax
3. Link to documentation
4. Link to the PR that implemented it
5. Thank the user for the feature request
6. Ask for feedback on the implementation

Example comment template:
```markdown
This has been implemented in v1.6.0! 🎉

You can now use:
```bash
[show example syntax]
```

Implementation details:
- [Link to PR]
- [Link to documentation]

Thanks for the feature request! Please let us know if this solves your use case or if you'd like any adjustments.
```

---

## Notes

- Skip HEREDOC syntax for `pet new` (user feedback: "weird in this day and age")
- Maintain 100% backward compatibility - all new features are opt-in via flags
- Default behavior (no flags) remains unchanged
- Focus on Unix philosophy: composable, scriptable, automatable
