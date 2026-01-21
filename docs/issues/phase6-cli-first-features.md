# Phase 6: CLI-First Features

**Part of:** Architecture Refactoring (Main Proposal: #XXX)
**Depends on:** Phase 5 (#XXX)

## Objective

Add CLI-first capabilities to Pet, allowing users to execute commands directly without interactive selection. This unlocks Pet's potential for scripting and automation.

## GitHub Issues Resolved by This Phase

This phase will allow us to close the following issues:

### ✅ Can Close After Implementation
- **#73** - pet search [flags] <string> (positional args)
- **#84** - Output all metadata from search (JSON output)
- **#117** - Arguments for parameters (--param and positional args)
- **#147** - pet exec without additional search (--first-match)
- **#164** - Enter first match with shortcut (--first-match)
- **#352** - Allow pet new to accept arguments via flags

### ⚡ Enhances (Already Closed)
- **#180** - Add support for `pet search <search term>` (add --exec, --list flags)

### ❌ Out of Scope (Cannot Close)
- **#146** - pet exec --last-command (needs history tracking - separate feature)
- **#154** - Paste to prompt (needs terminal integration)
- **#163** - Set cursor point (needs terminal integration)
- **#178** - Transfer to prompt (needs terminal integration)

## Tasks

### New Command Flags

- [ ] Add `-d/--description` flag for direct selection by description
  - [ ] Add to `exec` command
  - [ ] Add to `search` command
  - [ ] Add to `clip` command
  - [ ] Support exact match or regex (configurable)

- [ ] Add `-i/--id` flag for direct selection by snippet ID
  - [ ] Add to `exec` command
  - [ ] Add to `search` command
  - [ ] Add to `clip` command

- [ ] Add `--param key=value` flag for parameter passing
  - [ ] Support multiple `--param` flags
  - [ ] Override defaults with flag values
  - [ ] Error clearly if required param missing

- [ ] Add `--use-defaults` flag to skip parameter prompts
  - [ ] Use default values from snippet
  - [ ] Error if no default available
  - [ ] Works with `-d` and `-i` flags

- [ ] Add `--raw` flag to skip parameter substitution entirely
  - [ ] Output raw command with `<param>` placeholders
  - [ ] Useful for further processing

- [ ] Add `--json` flag for machine-readable output
  - [ ] Output snippet details as JSON
  - [ ] Useful for scripting and integration

### Implementation

- [ ] Wire up DirectSelector when `-d` or `-i` flags used
- [ ] Wire up FlagParameterInput when `--param` flags provided
- [ ] Wire up DefaultParameterInput when `--use-defaults` used
- [ ] Keep backward compatibility (no flags = interactive mode)

### Command Behavior Matrix

| Flags | Selector | ParameterInput | Behavior |
|-------|----------|----------------|----------|
| None | Fzf/Peco | Gocui | Current interactive |
| `-d "desc"` | Direct | Gocui | Direct select, prompt params |
| `-d "desc" --use-defaults` | Direct | Default | Direct select, use defaults |
| `-d "desc" --param x=1` | Direct | Flag | Direct select, use flag params |
| `-i "abc"` | Direct | Gocui | Direct by ID, prompt params |
| `--raw` | Fzf/Peco | None | Interactive select, no substitution |

### Documentation

- [ ] Update README with CLI-first examples
- [ ] Add USAGE.md with comprehensive examples
- [ ] Update command help text
- [ ] Add man pages (optional)

### Shell Completion

- [ ] Add completion for `--param` names
  - [ ] Parse snippet file for parameter names
  - [ ] Generate dynamic completion

- [ ] Add completion for `-d` descriptions
  - [ ] Complete from snippet descriptions

- [ ] Add completion for `-i` IDs
  - [ ] Complete from snippet IDs

- [ ] Update existing completion scripts
  - [ ] Bash completion
  - [ ] Zsh completion
  - [ ] Fish completion

### Examples to Document

```bash
# Interactive (current behavior)
pet exec

# Direct execution by description
pet exec -d "docker ps"

# Direct with parameters via flags
pet exec -d "docker run" --param container=nginx --param port=8080

# Direct using defaults
pet exec -d "git commit" --use-defaults

# Direct by ID
pet exec -i "abc123"

# Get command without executing
pet search -d "docker ps" --raw

# JSON output for scripting
pet search -d "docker" --json | jq '.command'

# Clipboard without interaction
pet clip -d "kubectl get pods" --use-defaults
```

### Testing

- [ ] Test each flag combination
- [ ] Test error cases (missing params, no match)
- [ ] Test with mock adapters
- [ ] Test actual CLI invocation (integration tests)
- [ ] Test shell completion scripts

## Acceptance Criteria

- [ ] All new flags implemented and working
- [ ] Direct selection works (by description and ID)
- [ ] Parameter passing via flags works
- [ ] Default parameter mode works
- [ ] Raw output mode works
- [ ] JSON output mode works
- [ ] Backward compatibility maintained (no flags = interactive)
- [ ] Help text updated for all commands
- [ ] README updated with CLI-first examples
- [ ] Shell completion updated
- [ ] Tests pass: `make test`
- [ ] Integration tests for CLI mode pass

## Implementation Notes

### Flag Parsing

```go
var (
    description  string
    snippetID    string
    paramFlags   []string
    useDefaults  bool
    rawOutput    bool
    jsonOutput   bool
)

execCmd.Flags().StringVarP(&description, "description", "d", "", "Select snippet by description (exact match)")
execCmd.Flags().StringVarP(&snippetID, "id", "i", "", "Select snippet by ID")
execCmd.Flags().StringSliceVar(&paramFlags, "param", []string{}, "Parameter value (can be repeated): --param name=value")
execCmd.Flags().BoolVar(&useDefaults, "use-defaults", false, "Use default parameter values without prompting")
execCmd.Flags().BoolVar(&rawOutput, "raw", false, "Output raw command without parameter substitution")
execCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output snippet as JSON")
```

### Adapter Selection Logic

```go
func selectSnippetSelector(cmd *cobra.Command) domain.Selector {
    if description := cmd.Flag("description").Value.String(); description != "" {
        return adapters.NewDirectSelector(description, DirectMatchByDescription)
    }
    if snippetID := cmd.Flag("id").Value.String(); snippetID != "" {
        return adapters.NewDirectSelector(snippetID, DirectMatchByID)
    }
    // Default to interactive
    return adapters.NewFzfSelector()
}

func selectParameterInput(cmd *cobra.Command) domain.ParameterInput {
    if cmd.Flag("use-defaults").Changed {
        return adapters.NewDefaultParameterInput()
    }
    if paramFlags, _ := cmd.Flags().GetStringSlice("param"); len(paramFlags) > 0 {
        params := parseParamFlags(paramFlags)
        return adapters.NewFlagParameterInput(params)
    }
    // Default to interactive
    return adapters.NewGocuiParameterInput()
}
```

### Shell Completion Example

```bash
# Bash completion for --param
_pet_exec_complete_param() {
    local cur="${COMP_WORDS[COMP_CWORD]}"
    local params=$(pet list --json | jq -r '.[] | .command' | grep -o '<[^>]*>' | tr -d '<>' | sort -u)
    COMPREPLY=( $(compgen -W "$params" -- "$cur") )
}

complete -F _pet_exec_complete_param pet exec --param
```

## User Benefits

1. **Scripting**: Use Pet snippets in shell scripts without interaction
2. **Automation**: Integrate Pet into CI/CD pipelines
3. **Speed**: Power users can bypass interactive selection
4. **Composability**: Chain Pet commands with other tools via pipes
5. **Discoverability**: JSON output enables building tools on top of Pet

## Dependencies

- Phase 5: Refactor Commands (#XXX)

## Estimated Effort

Medium (5-7 days)

## Labels

`enhancement`, `feature`, `cli`, `architecture`, `phase-6`
