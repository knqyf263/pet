package cmd

import (
	"os"
	
	"github.com/knqyf263/pet/config"
	"github.com/knqyf263/pet/path"
	petSync "github.com/knqyf263/pet/sync"
	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync snippets",
	Long:  `Sync snippets with gist/gitlab`,
	RunE:  sync,
}

func sync(cmd *cobra.Command, args []string) (err error) {
	// Check for environment variable override first
	snippetFile := config.Conf.General.SnippetFile
	if envFile := os.Getenv("PET_SNIPPET_FILE"); envFile != "" {
		snippetFile = envFile
	}
	filePath, err := path.NewAbsolutePath(snippetFile)
	return petSync.AutoSync(filePath)
}

func init() {
	RootCmd.AddCommand(syncCmd)
}
