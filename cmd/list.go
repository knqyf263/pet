package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/knqyf263/pet/config"
	"github.com/knqyf263/pet/snippet"
	runewidth "github.com/mattn/go-runewidth"
	"github.com/spf13/cobra"
)

const (
	column = 40
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show all snippets",
	Long:  `Show all snippets`,
	RunE:  list,
}

func list(cmd *cobra.Command, args []string) error {
	var snippets snippet.Snippets
	if err := snippets.Load(); err != nil {
		return err
	}

	col := config.Conf.General.Column
	if col == 0 {
		col = column
	}

	for _, snippet := range snippets.Snippets {
		description := runewidth.FillRight(runewidth.Truncate(snippet.Description, col, "..."), col)
		command := runewidth.Truncate(snippet.Command, 100-4-col, "...")

		fmt.Fprintf(color.Output, "%s : %s\n",
			color.GreenString(description), color.YellowString(command))
	}
	return nil
}

func init() {
	RootCmd.AddCommand(listCmd)
}
