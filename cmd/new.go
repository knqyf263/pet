package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
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

func scan(prompt string, out io.Writer, in io.ReadCloser, allowEmpty bool) (string, error) {
	f, err := os.CreateTemp("", "pet-")
	if err != nil {
		return "", err
	}
	defer os.Remove(f.Name()) // clean up temp file
	tempFile := f.Name()

	l, err := readline.NewEx(&readline.Config{
		Stdout:            out,
		Stdin:             in,
		Prompt:            prompt,
		HistoryFile:       tempFile,
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
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

// States of scanMultiLine state machine
const (
	start = iota
	lastLineNotEmpty
	lastLineEmpty
)

func scanMultiLine(prompt string, secondMessage string, out io.Writer, in io.ReadCloser) (string, error) {
	tempFile := "/tmp/pet.tmp"
	if runtime.GOOS == "windows" {
		tempDir := os.Getenv("TEMP")
		tempFile = filepath.Join(tempDir, "pet.tmp")
	}
	l, err := readline.NewEx(&readline.Config{
		Stdout:            out,
		Stdin:             in,
		Prompt:            prompt,
		HistoryFile:       tempFile,
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		VimMode:           false,
		HistorySearchFold: true,
	})
	if err != nil {
		return "", err
	}
	defer l.Close()

	state := start
	multiline := ""
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
		switch state {
		case start:
			if line == "" {
				continue
			}
			multiline += line
			state = lastLineNotEmpty
			l.SetPrompt(secondMessage)
		case lastLineNotEmpty:
			if line == "" {
				state = lastLineEmpty
				continue
			}
			multiline += "\n" + line
		case lastLineEmpty:
			if line == "" {
				return multiline, nil
			}
			multiline += "\n" + line
			state = lastLineNotEmpty
		}
	}
	return "", errors.New("canceled")
}

// createAndEditSnippet creates and saves a given snippet, then opens the
// configured editor to edit the snippet file at startLine.
func createAndEditSnippet(newSnippet snippet.SnippetInfo, snippets snippet.Snippets, startLine int) error {
	snippets.Snippets = append(snippets.Snippets, newSnippet)
	if err := snippets.Save(); err != nil {
		return err
	}

	// Open snippet for editing
	snippetFile := config.Conf.General.SnippetFile
	editor := config.Conf.General.Editor
	err := editFile(editor, snippetFile, startLine)
	if err != nil {
		return err
	}

	if config.Conf.Gist.AutoSync {
		return petSync.AutoSync(snippetFile)
	}

	return nil
}

func countSnippetLines() int {
	// Count lines in snippet file
	f, err := os.Open(config.Conf.General.SnippetFile)
	if err != nil {
		panic("Error reading snippet file")
	}
	lineCount, err := CountLines(f)
	if err != nil {
		panic("Error counting lines in snippet file")
	}

	return lineCount
}

func new(cmd *cobra.Command, args []string) (err error) {
	var command string
	var description string
	var tags []string

	var snippets snippet.Snippets
	if err := snippets.Load(); err != nil {
		return err
	}

	lineCount := countSnippetLines()

	// Get the command from the user
	if len(args) > 0 {
		command = strings.Join(args, " ")
		fmt.Fprintf(color.Output, "%s %s\n", color.HiYellowString("Command>"), command)
	} else {
		if config.Flag.UseMultiLine {
			command, err = scanMultiLine(
				color.YellowString("Command> "),
				color.YellowString(".......> "),
				os.Stdout, os.Stdin,
			)
		} else if config.Flag.UseEditor {
			// Create and save empty snippet
			newSnippet := snippet.SnippetInfo{
				Description: description,
				Command:     command,
				Tag:         tags,
			}

			return createAndEditSnippet(newSnippet, snippets, lineCount+3)

		} else {
			command, err = scan(color.HiYellowString("Command> "), os.Stdout, os.Stdin, false)
		}
		if err != nil {
			return err
		}
	}
	// Get the description from the user
	description, err = scan(color.HiGreenString("Description> "), os.Stdout, os.Stdin, false)
	if err != nil {
		return err
	}
	// Get the tags from the user
	if config.Flag.Tag {
		var t string
		if t, err = scan(color.HiCyanString("Tag> "), os.Stdout, os.Stdin, true); err != nil {
			return err
		}

		if t != "" {
			tags = strings.Fields(t)
		}
	}

	snippetFile := config.Conf.General.SnippetFile

	// Sync snippets beforehand to get the latest version if the remote is newer
	if config.Conf.Gist.AutoSync {
		if err := petSync.AutoSync(snippetFile); err != nil {
			return err
		}
	}

	for _, s := range snippets.Snippets {
		if s.Description == description {
			return fmt.Errorf("snippet [%s] already exists", description)
		}
	}
	// Save the command, description, and tags as a new snippet
	newSnippet := snippet.SnippetInfo{
		Description: description,
		Command:     command,
		Tag:         tags,
	}
	snippets.Snippets = append(snippets.Snippets, newSnippet)
	if err = snippets.Save(); err != nil {
		return err
	}
	// Sync snippets after update to keep the remote up to date
	if config.Conf.Gist.AutoSync {
		return petSync.AutoSync(snippetFile)
	}

	return nil
}

func init() {
	RootCmd.AddCommand(newCmd)
	newCmd.Flags().BoolVarP(&config.Flag.Tag, "tag", "t", false,
		`Display tag prompt (delimiter: space)`)
	newCmd.Flags().BoolVarP(&config.Flag.UseMultiLine, "multiline", "m", false,
		`Can enter multiline snippet (Double \n to quit)`)
	newCmd.Flags().BoolVarP(&config.Flag.UseEditor, "editor", "e", false,
		`Use editor to create snippet`)
}
