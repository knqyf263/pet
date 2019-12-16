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
	"github.com/knqyf263/pet/snippet"
	petSync "github.com/knqyf263/pet/sync"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new COMMAND",
	Short: "Create a new snippet",
	Long:  `Create a new snippet (default: $HOME/.config/pet/snippet.toml)`,
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("tag", cmd.Flags().Lookup("tag"))
	},
	RunE: new,
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
	var filename string = ""
	var command string
	var description string
	var tags []string

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

	if viper.GetBool("tag") {
		var t string
		if t, err = scan(color.CyanString("Tag> ")); err != nil {
			return err
		}
		tags = strings.Fields(t)
	}

	for _, s := range snippets.Snippets {
		if s.Description == description {
			return fmt.Errorf("Snippet [%s] already exists", description)
		}
	}

	if viper.GetString("general.snippetfile") != "" {
		filename = viper.GetString("general.snippetfile")
	}

	newSnippet := snippet.SnippetInfo{
		Filename:    filename,
		Description: description,
		Command:     command,
		Tag:         tags,
	}
	snippets.Snippets = append(snippets.Snippets, newSnippet)
	if err = snippets.Save(); err != nil {
		return err
	}

	if viper.GetBool("gist.auto_sync") {
		return petSync.AutoSync()
	}

	return nil
}

func init() {
	rootCmd.AddCommand(newCmd)
	newCmd.Flags().BoolP("tag", "t", false, `Display tag prompt (delimiter: space)`)
}
