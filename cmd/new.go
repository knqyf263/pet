package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/fatih/color"
	"github.com/knqyf263/pet/snippet"
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

	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		return "", err
	}
	defer terminal.Restore(0, oldState)

	t := terminal.NewTerminal(os.Stdin, message)
	s, err := t.ReadLine()

	if s == "" {
		return "", errors.New("canceled")
	}

	return s, err
}

func new(cmd *cobra.Command, args []string) (err error) {
	var command string
	var description string

	var snippets snippet.Snippets
	if err := snippets.Load(); err != nil {
		return err
	}

	if len(args) > 0 {
		command = strings.Join(args, " ")
		fmt.Printf("%s %s\n", color.YellowString("Command:"), command)
	} else {
		command, err = scan(color.YellowString("Command: "))
		if err != nil {
			return err
		}
	}
	description, err = scan(color.GreenString("Description: "))
	if err != nil {
		return err
	}

	for _, s := range snippets.Snippets {
		if s.Description == description {
			return fmt.Errorf("Snippet [%s] already exists", description)
		}
	}

	newSnippet := snippet.SnippetInfo{
		Description: description,
		Command:     command,
	}
	snippets.Snippets = append(snippets.Snippets, newSnippet)
	if err = snippets.Save(); err != nil {
		return err
	}

	return nil
}

func init() {
	RootCmd.AddCommand(newCmd)
}
