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
	"github.com/ramiawar/superpet/config"
	"github.com/ramiawar/superpet/envvar"
	"github.com/ramiawar/superpet/snippet"
	petSync "github.com/ramiawar/superpet/sync"
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new COMMAND",
	Short: "Create a new snippet",
	Long:  `Create a new snippet (default: $HOME/.config/pet/snippet.toml)`,
	RunE:  new,
}

var newEnvCmd = &cobra.Command{
	Use:   "newenv",
	Short: "Create a new env var namespace",
	Long:  `Create a new env var namespace (default: $HOME/.config/pet/snippet.toml)`,
	RunE:  newenv,
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
		if t, err = scan(color.CyanString("Tag> ")); err != nil {
			return err
		}
		tags = strings.Fields(t)
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

func newenv(cmd *cobra.Command, args []string) (err error) {
	var variables []string
	var description string
	var tags []string

	var envvars envvar.EnvVar
	if err := envvars.Load(); err != nil {
		return err
	}

	description, err = scan(color.GreenString("Description> "))
	if err != nil {
		return err
	}

	fmt.Println("\n(variables format (space delimited): SDK_KEY=123 HOST=123 ENV=abc PORT=999)")
	v, err := scan(color.YellowString("Variables> "))
	if err != nil {
		return err
	}
	variables = strings.Fields(v)

	if config.Flag.Tag {
		var t string
		if t, err = scan(color.CyanString("Tag> ")); err != nil {
			return err
		}
		tags = strings.Fields(t)
	}

	for _, s := range envvars.EnvVars {
		if s.Description == description {
			return fmt.Errorf("env [%s] already exists", description)
		}
	}

	newEnv := envvar.EnvVarInfo{
		Description: description,
		Variables:   variables,
		Tag:         tags,
	}
	envvars.EnvVars = append(envvars.EnvVars, newEnv)
	if err = envvars.Save(); err != nil {
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
	RootCmd.AddCommand(newEnvCmd)
	newEnvCmd.Flags().BoolVarP(&config.Flag.Tag, "tag", "t", false, `Display tag prompt (delimiter: space)`)
}
