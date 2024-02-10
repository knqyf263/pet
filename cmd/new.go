package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
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

func CanceledError() error {
	return errors.New("canceled")
}

func scan(message string, out io.Writer, in io.ReadCloser, allowEmpty bool) (string, error) {
	f, err := os.CreateTemp("", "pet-")
	if err != nil {
		return "", err
	}
	defer os.Remove(f.Name()) // clean up temp file
	tempFile := f.Name()

	l, err := readline.NewEx(&readline.Config{
		Stdout:          out,
		Stdin:           in,
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

		// If empty string, just ignore tags
		line = strings.TrimSpace(line)
		if line == "" && !allowEmpty {
			continue
		} else if line == "" {
			return "", nil
		}
		return line, nil
	}
	return "", CanceledError()
}

func new(cmd *cobra.Command, args []string) (err error) {
	var command string
	var description string
	var tags []string

	var snippets snippet.Snippets
	if err := snippets.Load(); err != nil {
		return err
	}

	if len(args) > 0 {
		command = strings.Join(args, " ")
		fmt.Fprintf(color.Output, "%s %s\n", color.HiYellowString("Command>"), command)
	} else {
		command, err = scan(color.HiYellowString("Command> "), os.Stdout, os.Stdin, false)
		if err != nil {
			return err
		}
	}
	description, err = scan(color.HiGreenString("Description> "), os.Stdout, os.Stdin, false)
	if err != nil {
		return err
	}

	if config.Flag.Tag {
		var t string
		if t, err = scan(color.HiCyanString("Tag> "), os.Stdout, os.Stdin, true); err != nil {
			return err
		}

		if t != "" {
			tags = strings.Fields(t)
		}
	}

	for _, s := range snippets.Snippets {
		if s.Description == description {
			return fmt.Errorf("snippet [%s] already exists", description)
		}
	}

	newSnippet := snippet.SnippetInfo{
		Description: description,
		Command:     command,
		Tag:         tags,
	}
	snippets.Snippets = append(snippets.Snippets, newSnippet)
	if err = snippets.Save(); err != nil {
		return err
	}

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
}
