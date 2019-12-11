package cmd

import (
	"io/ioutil"

	"github.com/spf13/viper"
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

func edit(cmd *cobra.Command, args []string) (err error) {
	editor := viper.GetString("general.editor")
	snippetFile := viper.GetString("general.snippetfile")

	// file content before editing
	before := fileContent(snippetFile)

	err = editFile(editor, snippetFile)
	if err != nil {
		return
	}

	// file content after editing
	after := fileContent(snippetFile)

	// return if same file content
	if before == after {
		return nil
	}

	if viper.GetBool("gist.auto_sync") {
		return petSync.AutoSync(snippetFile)
	}

	return nil
}

func fileContent(fname string) string {
	data, _ := ioutil.ReadFile(fname)
	return string(data)
}

func init() {
	rootCmd.AddCommand(editCmd)
}
