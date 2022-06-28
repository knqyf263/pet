package cmd

import (
	"io/ioutil"

	"github.com/knqyf263/pet/config"
	petSync "github.com/knqyf263/pet/sync"
	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit snippet file",
	Long:  `Edit snippet file (default: opened by vim)`,
	RunE:  edit,
}

var editEnvCmd = &cobra.Command{
	Use:   "editenv",
	Short: "Edit env file",
	Long:  `Edit env file (default: opened by vim)`,
	RunE:  editenv,
}

func edit(cmd *cobra.Command, args []string) (err error) {
	return handleEdit(config.Conf.General.SnippetFile)
}

func editenv(cmd *cobra.Command, args []string) (err error) {
	return handleEdit(config.Conf.General.EnvFile)
}

func handleEdit(file string) (err error) {
	editor := config.Conf.General.Editor

	// file content before editing
	before := fileContent(file)

	err = editFile(editor, file)
	if err != nil {
		return
	}

	// file content after editing
	after := fileContent(file)

	// return if same file content
	if before == after {
		return nil
	}

	if config.Conf.Gist.AutoSync {
		return petSync.AutoSync(file)
	}

	return nil
}

func fileContent(fname string) string {
	data, _ := ioutil.ReadFile(fname)
	return string(data)
}

func init() {
	RootCmd.AddCommand(editCmd)
	RootCmd.AddCommand(editEnvCmd)
}
