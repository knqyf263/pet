package cmd

import (
	"fmt"
	"os"

	"github.com/knqyf263/pet/config"
	"github.com/knqyf263/pet/path"
	petSync "github.com/knqyf263/pet/sync"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/alessio/shellescape.v1"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit snippet file",
	Long:  `Edit snippet file (default: opened by vim)`,
	RunE:  edit,
}

func edit(cmd *cobra.Command, args []string) (err error) {
	flag := config.Flag
	editor := config.Conf.General.Editor
	snippetFilePath, err := path.NewAbsolutePath(config.Conf.General.SnippetFile)
	if err != nil {
		return err
	}

	var options []string
	if flag.Query != "" {
		options = append(options, fmt.Sprintf("--query %s", shellescape.Quote(flag.Query)))
	}

	// If we have multiple snippet directories, we need to find the right
	// snippet file to edit - so we need to prompt the user to select a snippet first
	if len(config.Conf.General.SnippetDirs) > 0 {
		snippetFilePath, err = selectFile(options, flag.FilterTag)
		if err != nil {
			return err
		}
	}
	if snippetFilePath.Get() == "" {
		return errors.New("No snippet file selected")
	}

	// only sync if content has changed
	contentBefore := fileContent(snippetFilePath)
	err = editFile(editor, snippetFilePath, 0)
	if err != nil {
		return err
	}
	contentAfter := fileContent(snippetFilePath)
	if contentBefore == contentAfter {
		return nil
	}

	// sync snippet file
	if config.Conf.Gist.AutoSync {
		return petSync.AutoSync(snippetFilePath)
	}

	return nil
}

func fileContent(filePath path.AbsolutePath) string {
	data, _ := os.ReadFile(filePath.Get())
	return string(data)
}

func init() {
	RootCmd.AddCommand(editCmd)
	editCmd.Flags().StringVarP(&config.Flag.Query, "query", "q", "",
		`Initial value for query`)
	editCmd.Flags().StringVarP(&config.Flag.FilterTag, "tag", "t", "",
		`Filter tag`)
}
