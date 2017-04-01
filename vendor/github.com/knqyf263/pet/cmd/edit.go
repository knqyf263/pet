package cmd

import (
	"github.com/knqyf263/pet/config"
	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit snippet file",
	Long:  `Edit snippet file (default: opened by vim)`,
	RunE:  edit,
}

func edit(cmd *cobra.Command, args []string) (err error) {
	editor := config.Conf.General.Editor
	snippetFile := config.Conf.General.SnippetFile

	return editFile(editor, snippetFile)
}

func init() {
	RootCmd.AddCommand(editCmd)
}
