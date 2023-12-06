package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/knqyf263/pet/config"
	"github.com/knqyf263/pet/snippet"
	petSync "github.com/knqyf263/pet/sync"
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new COMMAND",
	Short: "Create a new snippet",
	Long:  `Create a new snippet (default: $HOME/.config/pet/snippet.toml)`,
	RunE:  new,
}

func scan(message string) (string, error) {
	tempFile := "/tmp/pet.tmp"
	if runtime.GOOS == "windows" {
		tempDir := os.Getenv("TEMP")
		tempFile = filepath.Join(tempDir, "pet.tmp")
	}
	l, err := readline.NewEx(&readline.Config{
		Prompt:          message,
		HistoryFile:     tempFile,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold: true,
	})
	if err != nil {
		return "", err
	}
	defer l.Close()

	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		return line, nil
	}
	return "", errors.New("canceled")
}

func new(cmd *cobra.Command, args []string) (err error) {
	var command string
	var description string
	var tags []string
	var output string

	var snippets snippet.Snippets
	if err := snippets.Load(); err != nil {
		return err
	}

	if len(args) > 0 {
		command = strings.Join(args, " ")
		fmt.Fprintf(color.Output, "%s %s\n", color.YellowString("Command>"), command)
	} else {
		command, err = scan(color.YellowString("Command> "))
		if err != nil {
			return err
		}
	}
	description, err = scan(color.GreenString("Description> "))
	if err != nil {
		return err
	}

	if config.Flag.Tag {
		var t string
		if t, err = scan(color.CyanString("Tags (space-delimited)> ")); err != nil {
			return err
		}
		tags = strings.Fields(t)
	}

	if config.Flag.Output {
		var o string
		if o, err = scan(color.RedString("Output> ")); err != nil {
			return err
		}
		output = o
	}

	snippetExists := false
	var newSnippets snippet.Snippets
	for _, s := range snippets.Snippets {
		// if the description matches an existing one, merge it into the existing one
		if s.Description == description {
			snippetExists = true
			fmt.Fprintf(color.Output, "%s\"%s\"\n", color.HiGreenString("Merging with existing snippet: "), description)
			// if the description and command already exist, then return an error
			if slices.Contains(s.Commands, command) {
				return fmt.Errorf("snippet already exists with that description and command")
			}
			s.Commands = append(s.Commands, command)

			// add any new tags to the set of tags
			for _, tag := range tags {
				if !slices.Contains(s.Tags, tag) {
					s.Tags = append(s.Tags, tag)
				}
			}

			// overwrite the output if there isn't one
			if s.Output == "" {
				s.Output = output
			} else {
				fmt.Fprintf(color.Output, "%s\n", color.HiRedString("'Output' is already set on the snippet, ignoring."))
			}
		}

		newSnippets.Snippets = append(newSnippets.Snippets, s)
	}

	// if we didnt match an existing snippet, then create a new one
	if !snippetExists {
		fmt.Fprintf(color.Output, "%s\n", color.HiGreenString("Creating new snippet..."))
		newSnippet := snippet.SnippetInfo{
			Description: description,
			Commands:    []string{command},
			Tags:        tags,
			Output:      output,
		}
		newSnippets.Snippets = append(newSnippets.Snippets, newSnippet)
	}

	if err = newSnippets.Save(); err != nil {
		return err
	}
	fmt.Fprintf(color.Output, "%s\n", color.HiGreenString("Success!"))

	snippetFile := config.Conf.General.SnippetFile
	if config.Conf.Gist.AutoSync {
		return petSync.AutoSync(snippetFile)
	}

	return nil
}

func init() {
	RootCmd.AddCommand(newCmd)
	newCmd.Flags().BoolVarP(&config.Flag.Tag, "tag", "t", false,
		`Display tag prompt (delimiter: space)`)
	newCmd.Flags().BoolVarP(&config.Flag.Output, "output", "o", false,
		`Display output prompt`)
}
