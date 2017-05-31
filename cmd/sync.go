package cmd

import (
	"fmt"

	"github.com/knqyf263/pet/config"
	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync snippets",
	Long:  `Sync snippets with gist`,
	RunE:  sync,
}

func sync(cmd *cobra.Command, args []string) (err error) {
	if config.Conf.Gist.AccessToken == "" {
		return fmt.Errorf(`access_token is empty.
Go https://github.com/settings/tokens/new and create access_token (only need "gist" scope).
Write access_token in config file (pet configure).
		`)
	}

	return autoSync(config.Conf.General.SnippetFile)
}

func init() {
	RootCmd.AddCommand(syncCmd)
}
