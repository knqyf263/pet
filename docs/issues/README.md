# Architecture Refactoring Issues

This directory contains issue templates for the Pet architecture refactoring project.

## Overview

The refactoring is organized into 7 phases, each building on the previous:

1. **Phase 1: Foundation** - Create core domain layer and interfaces
2. **Phase 2: Repository Pattern** - Abstract storage operations
3. **Phase 3: Parameter Input Abstraction** - Abstract parameter input (gocui, flags, etc.)
4. **Phase 4: Selector Abstraction** - Abstract snippet selection (fzf, peco, direct)
5. **Phase 5: Refactor Commands** - Update all commands to use new architecture
6. **Phase 6: CLI-First Features** - Add direct execution and parameter flags
7. **Phase 7: Documentation & Polish** - Complete documentation and release v2.0

## Creating Issues

To create the issues on GitHub:

1. Read through `architecture-refactoring-proposal.md` first
2. Create a main tracking issue using the proposal content
3. Create individual issues for each phase using these templates
4. Link phase issues to the main tracking issue
5. Create a GitHub milestone for v2.0

## Issue Dependencies

```
Main Proposal Issue
  ├─ Phase 1: Foundation
  │    └─ Phase 2: Repository Pattern
  │         └─ Phase 3: Parameter Input Abstraction
  │              └─ Phase 4: Selector Abstraction
  │                   └─ Phase 5: Refactor Commands
  │                        └─ Phase 6: CLI-First Features
  │                             └─ Phase 7: Documentation & Polish (v2.0.0)
```

## Labels to Create

Suggested labels for organizing these issues:

- `architecture` - Architecture-related changes
- `refactoring` - Code refactoring
- `enhancement` - New features
- `testing` - Testing improvements
- `documentation` - Documentation updates
- `phase-1` through `phase-7` - Phase tracking
- `v2.0` - Target v2.0 release
- `cli` - CLI-related features
- `gocui` - gocui-related changes

## Milestones

Consider creating these milestones:

- `v1.1.0` - Phase 1
- `v1.2.0` - Phase 2
- `v1.3.0` - Phase 3
- `v1.4.0` - Phase 4
- `v1.5.0` - Phase 5
- `v1.6.0` - Phase 6
- `v2.0.0` - Phase 7

## File List

- `phase1-foundation.md` - Foundation phase tasks
- `phase2-repository-pattern.md` - Repository pattern tasks
- `phase3-parameter-input-abstraction.md` - Parameter input abstraction tasks
- `phase4-selector-abstraction.md` - Selector abstraction tasks
- `phase5-refactor-commands.md` - Command refactoring tasks
- `phase6-cli-first-features.md` - CLI-first feature tasks
- `phase7-documentation-polish.md` - Documentation and polish tasks

## Usage

Copy the content of each markdown file to create a new issue on GitHub. Update the placeholder issue numbers (#XXX) with actual issue numbers once created.

## Questions or Feedback

For questions about the architecture or refactoring plan, please discuss in the main proposal issue.
