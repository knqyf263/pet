package cmd

import (
	petSync "github.com/knqyf263/pet/sync"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync snippets",
	Long:  `Sync snippets with gist/gitlab`,
	RunE:  sync,
}

func sync(cmd *cobra.Command, args []string) (err error) {
	return petSync.AutoSync(viper.GetString("general.snippetFile"))
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
