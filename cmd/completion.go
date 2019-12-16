package cmd

import (
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	completionShells = map[string]func(out io.Writer, cmd *cobra.Command) error{
		"bash": runCompletionBash,
		"zsh":  runCompletionZsh,
	}
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion SHELL",
	Short: "Generate shell completions",
	Long: `Generate shell completions for Pet for the specified shell (bash or zsh).

This command can generate shell autocompletions. e.g.

	$ pet completion bash

Can be sourced as such

	$ source <(pet completion bash)
`,
	RunE: runCompletion,
}

func runCompletion(cmd *cobra.Command, args []string) (err error) {
	if len(args) == 0 {
		return errors.New("Shell not specified.")
	}

	if len(args) > 1 {
		return errors.New("Too many arguments. Expected only the shell type.")
	}

	run, found := completionShells[args[0]]
	if !found {
		return errors.New("Unsupported shell type.")
	}

	return run(os.Stdout, cmd.Root())
}

func runCompletionBash(out io.Writer, cmd *cobra.Command) error {
	return cmd.GenBashCompletion(out)
}

func runCompletionZsh(out io.Writer, cmd *cobra.Command) error {
	return cmd.GenZshCompletion(out)
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
