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

	err = editFile(editor, snippetFile)
	if err == nil && config.Conf.Gist.AutoSync {
		return autoSync(snippetFile)
	}
	return err
}

func init() {
	RootCmd.AddCommand(editCmd)
}
